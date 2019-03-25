/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package endpoints

import (
	"fmt"
	"net/http"
	gpath "path"
	"reflect"
	"sort"
	"strings"
	"time"
	"unicode"

	restful "github.com/emicklei/go-restful"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/endpoints/discovery"
	"k8s.io/apiserver/pkg/endpoints/handlers"
	"k8s.io/apiserver/pkg/endpoints/handlers/fieldmanager"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/apiserver/pkg/endpoints/metrics"
	"k8s.io/apiserver/pkg/features"
	"k8s.io/apiserver/pkg/registry/rest"
	genericfilters "k8s.io/apiserver/pkg/server/filters"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
)

const (
	ROUTE_META_GVK    = "x-kubernetes-group-version-kind"
	ROUTE_META_ACTION = "x-kubernetes-action"
)

type APIInstaller struct {
	group                        *APIGroupVersion
	prefix                       string // Path prefix where API resources are to be registered.
	minRequestTimeout            time.Duration
	enableAPIResponseCompression bool
}

// Struct capturing information about an action ("GET", "POST", "WATCH", "PROXY", etc).
type action struct {
	Verb          string               // Verb identifying the action ("GET", "POST", "WATCH", "PROXY", etc).
	Path          string               // The path of the action
	Params        []*restful.Parameter // List of parameters associated with the action.
	AllNamespaces bool                 // true iff the action is namespaced but works on aggregate result for all namespaces

	handler     restful.RouteFunction
	readSample  interface{} // optional. Sample object of request body
	writeSample interface{} // required. Sample object of response body

	// list of potential http status codes given action returns. The current
	// state and assumption is all status code responses are of writeSample
	// type
	returnStatusCodes []int
	consumedMIMETypes []string
	producedMIMETypes []string
}

func newAction(verb, path string, params []*restful.Parameter, allNamespaces bool, additionalReturnStatusCodes ...int) *action {
	a := &action{
		Verb:          verb,
		Path:          path,
		Params:        params,
		AllNamespaces: allNamespaces,
	}
	a.returnStatusCodes = []int{http.StatusOK}
	// TODO: in some cases, the API may return a v1.Status instead of the versioned object
	// but currently go-restful can't handle multiple different objects being returned.
	a.returnStatusCodes = append(a.returnStatusCodes, additionalReturnStatusCodes...)
	return a
}

// An interface to see if one storage supports override its default verb for monitoring
type StorageMetricsOverride interface {
	// OverrideMetricsVerb gives a storage object an opportunity to override the verb reported to the metrics endpoint
	OverrideMetricsVerb(oldVerb string) (newVerb string)
}

// An interface to see if an object supports swagger documentation as a method
type documentable interface {
	SwaggerDoc() map[string]string
}

// toDiscoveryKubeVerb maps an action.Verb to the logical kube verb, used for discovery
var toDiscoveryKubeVerb = map[string]string{
	"CONNECT":          "", // do not list in discovery.
	"DELETE":           "delete",
	"DELETECOLLECTION": "deletecollection",
	"GET":              "get",
	"LIST":             "list",
	"PATCH":            "patch",
	"POST":             "create",
	"PROXY":            "proxy",
	"PUT":              "update",
	"WATCH":            "watch",
	"WATCHLIST":        "watch",
}

// toRESTMethod maps an action.Verb to the RESTful method
var toRESTMethod = map[string]string{
	"DELETE":           "DELETE",
	"DELETECOLLECTION": "DELETE",
	"GET":              "GET",
	"LIST":             "GET",
	"PATCH":            "PATCH",
	"POST":             "POST",
	"PUT":              "PUT",
	"WATCH":            "GET",
	"WATCHLIST":        "GET",
}

// Install handlers for API resources.
func (a *APIInstaller) Install() ([]metav1.APIResource, *restful.WebService, []error) {
	var apiResources []metav1.APIResource
	var errors []error
	ws := a.newWebService()

	// Register the paths in a deterministic (sorted) order to get a deterministic swagger spec.
	paths := make([]string, len(a.group.Storage))
	var i int = 0
	for path := range a.group.Storage {
		paths[i] = path
		i++
	}
	sort.Strings(paths)
	for _, path := range paths {
		apiResource, err := a.registerResourceHandlers(path, a.group.Storage[path], ws)
		if err != nil {
			errors = append(errors, fmt.Errorf("error in registering resource: %s, %v", path, err))
		}
		if apiResource != nil {
			apiResources = append(apiResources, *apiResource)
		}
	}
	return apiResources, ws, errors
}

// newWebService creates a new restful webservice with the api installer's prefix and version.
func (a *APIInstaller) newWebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(a.prefix)
	// a.prefix contains "prefix/group/version"
	ws.Doc("API at " + a.prefix)
	// Backwards compatibility, we accepted objects with empty content-type at V1.
	// If we stop using go-restful, we can default empty content-type to application/json on an
	// endpoint by endpoint basis
	ws.Consumes("*/*")
	mediaTypes, streamMediaTypes := negotiation.MediaTypesForSerializer(a.group.Serializer)
	ws.Produces(append(mediaTypes, streamMediaTypes...)...)
	ws.ApiVersion(a.group.GroupVersion.String())

	return ws
}

// calculate the storage gvk, the gvk objects are converted to before persisted to the etcd.
func getStorageVersionKind(storageVersioner runtime.GroupVersioner, storage rest.Storage, typer runtime.ObjectTyper) (schema.GroupVersionKind, error) {
	object := storage.New()
	fqKinds, _, err := typer.ObjectKinds(object)
	if err != nil {
		return schema.GroupVersionKind{}, err
	}
	gvk, ok := storageVersioner.KindForGroupVersionKinds(fqKinds)
	if !ok {
		return schema.GroupVersionKind{}, fmt.Errorf("cannot find the storage version kind for %v", reflect.TypeOf(object))
	}
	return gvk, nil
}

// GetResourceKind returns the external group version kind registered for the given storage
// object. If the storage object is a subresource and has an override supplied for it, it returns
// the group version kind supplied in the override.
func GetResourceKind(groupVersion schema.GroupVersion, storage rest.Storage, typer runtime.ObjectTyper) (schema.GroupVersionKind, error) {
	// Let the storage tell us exactly what GVK it has
	if gvkProvider, ok := storage.(rest.GroupVersionKindProvider); ok {
		return gvkProvider.GroupVersionKind(groupVersion), nil
	}

	object := storage.New()
	fqKinds, _, err := typer.ObjectKinds(object)
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	// a given go type can have multiple potential fully qualified kinds.  Find the one that corresponds with the group
	// we're trying to register here
	fqKindToRegister := schema.GroupVersionKind{}
	for _, fqKind := range fqKinds {
		if fqKind.Group == groupVersion.Group {
			fqKindToRegister = groupVersion.WithKind(fqKind.Kind)
			break
		}
	}
	if fqKindToRegister.Empty() {
		return schema.GroupVersionKind{}, fmt.Errorf("unable to locate fully qualified kind for %v: found %v when registering for %v", reflect.TypeOf(object), fqKinds, groupVersion)
	}

	// group is guaranteed to match based on the check above
	return fqKindToRegister, nil
}

func (a *APIInstaller) registerResourceHandlers(path string, storage rest.Storage, ws *restful.WebService) (*metav1.APIResource, error) {
	admit := a.group.Admit

	optionsExternalVersion := a.group.GroupVersion
	if a.group.OptionsExternalVersion != nil {
		optionsExternalVersion = *a.group.OptionsExternalVersion
	}

	resource, subresource, err := splitSubresource(path)
	if err != nil {
		return nil, err
	}

	group, version := a.group.GroupVersion.Group, a.group.GroupVersion.Version

	fqKindToRegister, err := GetResourceKind(a.group.GroupVersion, storage, a.group.Typer)
	if err != nil {
		return nil, err
	}

	versionedPtr, err := a.group.Creater.New(fqKindToRegister)
	if err != nil {
		return nil, err
	}
	defaultVersionedObject := indirectArbitraryPointer(versionedPtr)
	kind := fqKindToRegister.Kind
	isSubresource := len(subresource) > 0

	// If there is a subresource, namespace scoping is defined by the parent resource
	namespaceScoped := true
	if isSubresource {
		parentStorage, ok := a.group.Storage[resource]
		if !ok {
			return nil, fmt.Errorf("missing parent storage: %q", resource)
		}
		scoper, ok := parentStorage.(rest.Scoper)
		if !ok {
			return nil, fmt.Errorf("%q must implement scoper", resource)
		}
		namespaceScoped = scoper.NamespaceScoped()

	} else {
		scoper, ok := storage.(rest.Scoper)
		if !ok {
			return nil, fmt.Errorf("%q must implement scoper", resource)
		}
		namespaceScoped = scoper.NamespaceScoped()
	}

	// what verbs are supported by the storage, used to know what verbs we support per path
	creater, isCreater := storage.(rest.Creater)
	namedCreater, isNamedCreater := storage.(rest.NamedCreater)
	lister, isLister := storage.(rest.Lister)
	getter, isGetter := storage.(rest.Getter)
	getterWithOptions, isGetterWithOptions := storage.(rest.GetterWithOptions)
	gracefulDeleter, isGracefulDeleter := storage.(rest.GracefulDeleter)
	collectionDeleter, isCollectionDeleter := storage.(rest.CollectionDeleter)
	updater, isUpdater := storage.(rest.Updater)
	patcher, isPatcher := storage.(rest.Patcher)
	watcher, isWatcher := storage.(rest.Watcher)
	connecter, isConnecter := storage.(rest.Connecter)
	storageMeta, isMetadata := storage.(rest.StorageMetadata)
	if !isMetadata {
		storageMeta = defaultStorageMetadata{}
	}
	exporter, isExporter := storage.(rest.Exporter)
	if !isExporter {
		exporter = nil
	}

	versionedExportOptions, err := a.group.Creater.New(optionsExternalVersion.WithKind("ExportOptions"))
	if err != nil {
		return nil, err
	}

	if isNamedCreater {
		isCreater = true
	}

	var versionedList interface{}
	if isLister {
		list := lister.NewList()
		listGVKs, _, err := a.group.Typer.ObjectKinds(list)
		if err != nil {
			return nil, err
		}
		versionedListPtr, err := a.group.Creater.New(a.group.GroupVersion.WithKind(listGVKs[0].Kind))
		if err != nil {
			return nil, err
		}
		versionedList = indirectArbitraryPointer(versionedListPtr)
	}

	versionedListOptions, err := a.group.Creater.New(optionsExternalVersion.WithKind("ListOptions"))
	if err != nil {
		return nil, err
	}
	versionedCreateOptions, err := a.group.Creater.New(optionsExternalVersion.WithKind("CreateOptions"))
	if err != nil {
		return nil, err
	}
	versionedPatchOptions, err := a.group.Creater.New(optionsExternalVersion.WithKind("PatchOptions"))
	if err != nil {
		return nil, err
	}
	versionedUpdateOptions, err := a.group.Creater.New(optionsExternalVersion.WithKind("UpdateOptions"))
	if err != nil {
		return nil, err
	}

	var versionedDeleteOptions runtime.Object
	var versionedDeleterObject interface{}
	if isGracefulDeleter {
		versionedDeleteOptions, err = a.group.Creater.New(optionsExternalVersion.WithKind("DeleteOptions"))
		if err != nil {
			return nil, err
		}
		versionedDeleterObject = indirectArbitraryPointer(versionedDeleteOptions)
	}

	versionedStatusPtr, err := a.group.Creater.New(optionsExternalVersion.WithKind("Status"))
	if err != nil {
		return nil, err
	}
	versionedStatus := indirectArbitraryPointer(versionedStatusPtr)
	var (
		getOptions             runtime.Object
		versionedGetOptions    runtime.Object
		getOptionsInternalKind schema.GroupVersionKind
		getSubpath             bool
	)
	if isGetterWithOptions {
		getOptions, getSubpath, _ = getterWithOptions.NewGetOptions()
		getOptionsInternalKinds, _, err := a.group.Typer.ObjectKinds(getOptions)
		if err != nil {
			return nil, err
		}
		getOptionsInternalKind = getOptionsInternalKinds[0]
		versionedGetOptions, err = a.group.Creater.New(a.group.GroupVersion.WithKind(getOptionsInternalKind.Kind))
		if err != nil {
			versionedGetOptions, err = a.group.Creater.New(optionsExternalVersion.WithKind(getOptionsInternalKind.Kind))
			if err != nil {
				return nil, err
			}
		}
		isGetter = true
	}

	var versionedWatchEvent interface{}
	if isWatcher {
		versionedWatchEventPtr, err := a.group.Creater.New(a.group.GroupVersion.WithKind("WatchEvent"))
		if err != nil {
			return nil, err
		}
		versionedWatchEvent = indirectArbitraryPointer(versionedWatchEventPtr)
	}

	var (
		connectOptions             runtime.Object
		versionedConnectOptions    runtime.Object
		connectOptionsInternalKind schema.GroupVersionKind
		connectSubpath             bool
	)
	if isConnecter {
		connectOptions, connectSubpath, _ = connecter.NewConnectOptions()
		if connectOptions != nil {
			connectOptionsInternalKinds, _, err := a.group.Typer.ObjectKinds(connectOptions)
			if err != nil {
				return nil, err
			}

			connectOptionsInternalKind = connectOptionsInternalKinds[0]
			versionedConnectOptions, err = a.group.Creater.New(a.group.GroupVersion.WithKind(connectOptionsInternalKind.Kind))
			if err != nil {
				versionedConnectOptions, err = a.group.Creater.New(optionsExternalVersion.WithKind(connectOptionsInternalKind.Kind))
				if err != nil {
					return nil, err
				}
			}
		}
	}

	allowWatchList := isWatcher && isLister // watching on lists is allowed only for kinds that support both watch and list.
	nameParam := ws.PathParameter("name", "name of the "+kind).DataType("string")
	pathParam := ws.PathParameter("path", "path to the resource").DataType("string")

	params := []*restful.Parameter{}
	actions := []*action{}

	var resourceKind string
	kindProvider, ok := storage.(rest.KindProvider)
	if ok {
		resourceKind = kindProvider.Kind()
	} else {
		resourceKind = kind
	}

	tableProvider, _ := storage.(rest.TableConvertor)

	var namer handlers.ContextBasedNaming
	var resourcePath, itemPath string
	var resourceParams, nameParams, proxyParams []*restful.Parameter
	// Get the list of actions for the given scope.
	switch {
	case !namespaceScoped:
		// Handle non-namespace scoped resources like nodes.
		resourcePath = resource
		resourceParams = params
		itemPath = resourcePath + "/{name}"
		nameParams = append(params, nameParam)
		proxyParams = append(nameParams, pathParam)
		suffix := ""
		if isSubresource {
			suffix = "/" + subresource
			itemPath = itemPath + suffix
			resourcePath = itemPath
			resourceParams = nameParams
		}
		namer = handlers.ContextBasedNaming{
			SelfLinker:         a.group.Linker,
			ClusterScoped:      true,
			SelfLinkPathPrefix: gpath.Join(a.prefix, resource) + "/",
			SelfLinkPathSuffix: suffix,
		}
	default:
		namespaceParamName := "namespaces"
		// Handler for standard REST verbs (GET, PUT, POST and DELETE).
		namespaceParam := ws.PathParameter("namespace", "object name and auth scope, such as for teams and projects").DataType("string")
		namespacedPath := namespaceParamName + "/{namespace}/" + resource
		namespaceParams := []*restful.Parameter{namespaceParam}

		resourcePath = namespacedPath
		resourceParams = namespaceParams
		itemPath = namespacedPath + "/{name}"
		nameParams = append(namespaceParams, nameParam)
		proxyParams = append(nameParams, pathParam)
		itemPathSuffix := ""
		if isSubresource {
			itemPathSuffix = "/" + subresource
			itemPath = itemPath + itemPathSuffix
			resourcePath = itemPath
			resourceParams = nameParams
		}
		namer = handlers.ContextBasedNaming{
			SelfLinker:         a.group.Linker,
			ClusterScoped:      false,
			SelfLinkPathPrefix: gpath.Join(a.prefix, namespaceParamName) + "/",
			SelfLinkPathSuffix: itemPathSuffix,
		}

		// list or post across namespace.
		// E.g. LIST all pods in all namespaces by sending a LIST request at /api/apiVersion/pods.
		// TODO: more strongly type whether a resource allows these actions on "all namespaces" (bulk delete)
		if !isSubresource {
			actions = appendIf(actions, newAction("LIST", resource, params, true), isLister)
			// DEPRECATED in 1.11
			actions = appendIf(actions, newAction("WATCHLIST", "watch/"+resource, params, true), allowWatchList)
		}
	}

	// Handler for standard REST verbs (GET, PUT, POST and DELETE).
	// Add actions at the resource path: /api/apiVersion/resource
	actions = appendIf(actions, newAction("LIST", resourcePath, resourceParams, false), isLister)
	actions = appendIf(actions, newAction("POST", resourcePath, resourceParams, false, http.StatusCreated, http.StatusAccepted), isCreater)
	actions = appendIf(actions, newAction("DELETECOLLECTION", resourcePath, resourceParams, false), isCollectionDeleter)

	// DEPRECATED in 1.11
	actions = appendIf(actions, newAction("WATCHLIST", "watch/"+resourcePath, resourceParams, false), allowWatchList)

	// Add actions at the item path: /api/apiVersion/resource/{name}
	actions = appendIf(actions, newAction("GET", itemPath, nameParams, false), isGetter)
	if getSubpath {
		actions = appendIf(actions, newAction("GET", itemPath+"/{path:*}", proxyParams, false), isGetter)
	}
	actions = appendIf(actions, newAction("PUT", itemPath, nameParams, false, http.StatusCreated), isUpdater)
	actions = appendIf(actions, newAction("PATCH", itemPath, nameParams, false), isPatcher)
	actions = appendIf(actions, newAction("DELETE", itemPath, nameParams, false, http.StatusAccepted), isGracefulDeleter)
	// DEPRECATED in 1.11
	actions = appendIf(actions, newAction("WATCH", "watch/"+itemPath, nameParams, false), isWatcher)
	actions = appendIf(actions, newAction("CONNECT", itemPath, nameParams, false), isConnecter)
	actions = appendIf(actions, newAction("CONNECT", itemPath+"/{path:*}", proxyParams, false), isConnecter && connectSubpath)

	// TODO: Add status documentation using Returns()
	// Errors (see api/errors/errors.go as well as go-restful router):
	// http.StatusNotFound, http.StatusMethodNotAllowed,
	// http.StatusUnsupportedMediaType, http.StatusNotAcceptable,
	// http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden,
	// http.StatusRequestTimeout, http.StatusConflict, http.StatusPreconditionFailed,
	// http.StatusUnprocessableEntity, http.StatusInternalServerError,
	// http.StatusServiceUnavailable
	// and api error codes
	// Note that if we specify a versioned Status object here, we may need to
	// create one for the tests, also
	// Success:
	// http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent
	//
	// test/integration/auth_test.go is currently the most comprehensive status code test

	for _, s := range a.group.Serializer.SupportedMediaTypes() {
		if len(s.MediaTypeSubType) == 0 || len(s.MediaTypeType) == 0 {
			return nil, fmt.Errorf("all serializers in the group Serializer must have MediaTypeType and MediaTypeSubType set: %s", s.MediaType)
		}
	}

	// construct handler for each action
	reqScope, err := a.constructRequestScope(fqKindToRegister, tableProvider, resource, subresource, namer)
	if err != nil {
		return nil, err
	}

	for _, action := range actions {
		var handler restful.RouteFunction
		verb := action.Verb
		switch action.Verb {
		case "GET":
			if isGetterWithOptions {
				handler = restfulGetResourceWithOptions(getterWithOptions, reqScope, isSubresource)
			} else {
				handler = restfulGetResource(getter, exporter, reqScope)
			}
			verbOverrider, needOverride := storage.(StorageMetricsOverride)
			if needOverride {
				// need change the reported verb
				verb = verbOverrider.OverrideMetricsVerb(action.Verb)
			}
		case "LIST":
			handler = restfulListResource(lister, watcher, reqScope, false, a.minRequestTimeout)
		case "PUT":
			handler = restfulUpdateResource(updater, reqScope, admit)
		case "PATCH":
			handler = restfulPatchResource(patcher, reqScope, admit, supportedPatchTypes())
		case "POST":
			if isNamedCreater {
				handler = restfulCreateNamedResource(namedCreater, reqScope, admit)
			} else {
				handler = restfulCreateResource(creater, reqScope, admit)
			}
		case "DELETE":
			handler = restfulDeleteResource(gracefulDeleter, isGracefulDeleter, reqScope, admit)
		case "DELETECOLLECTION":
			handler = restfulDeleteCollection(collectionDeleter, isCollectionDeleter, reqScope, admit)
		case "WATCH", "WATCHLIST":
			handler = restfulListResource(lister, watcher, reqScope, true, a.minRequestTimeout)
		case "CONNECT":
			handler = restfulConnectResource(connecter, reqScope, admit, path, isSubresource)
		}
		requestScope := "cluster"
		if namespaceScoped {
			requestScope = "namespace"
		}
		if strings.HasSuffix(action.Path, "/{path:*}") {
			requestScope = "resource"
		}
		if action.AllNamespaces {
			requestScope = "cluster"
		}
		// instrument handler
		handler = metrics.InstrumentRouteFunc(verb, group, version, resource, subresource, requestScope, metrics.APIServerComponent, handler)
		switch action.Verb {
		case "GET", "LIST":
			if a.enableAPIResponseCompression {
				handler = genericfilters.RestfulWithCompression(handler)
			}
		}
		action.handler = handler
	}

	mediaTypes, streamMediaTypes := negotiation.MediaTypesForSerializer(a.group.Serializer)
	allMediaTypes := append(mediaTypes, streamMediaTypes...)
	ws.Produces(allMediaTypes...)

	// complete actions
	for _, action := range actions {
		// assign producedMIMETypes
		producedMIMETypes := storageMeta.ProducesMIMETypes(action.Verb)
		action = action.assignProducedMIMETypes(producedMIMETypes, mediaTypes, allMediaTypes)

		// assign consumedMIMETypes
		action = action.assignConsumedMIMETypes()

		// assign writeSample
		writeSample := storageMeta.ProducesObject(action.Verb)
		if writeSample == nil {
			writeSample = defaultVersionedObject
		}
		action = action.assignWriteSample(writeSample, versionedList, versionedStatus, versionedWatchEvent)

		// assign readSample
		action = action.assignReadSample(defaultVersionedObject, versionedDeleterObject)
	}

	// add object params for each action
	for _, action := range actions {
		switch action.Verb {
		case "GET":
			if isGetterWithOptions {
				params, err := addObjectParams(ws, versionedGetOptions)
				if err != nil {
					return nil, err
				}
				action.Params = append(action.Params, params...)
			}
			if isExporter {
				params, err := addObjectParams(ws, versionedExportOptions)
				if err != nil {
					return nil, err
				}
				action.Params = append(action.Params, params...)
			}
		case "LIST", "DELETECOLLECTION", "WATCH", "WATCHLIST":
			params, err := addObjectParams(ws, versionedListOptions)
			if err != nil {
				return nil, err
			}
			action.Params = append(action.Params, params...)
		case "PUT":
			params, err := addObjectParams(ws, versionedUpdateOptions)
			if err != nil {
				return nil, err
			}
			action.Params = append(action.Params, params...)
		case "PATCH":
			params, err := addObjectParams(ws, versionedPatchOptions)
			if err != nil {
				return nil, err
			}
			action.Params = append(action.Params, params...)
		case "POST":
			params, err := addObjectParams(ws, versionedCreateOptions)
			if err != nil {
				return nil, err
			}
			action.Params = append(action.Params, params...)
		case "DELETE":
			if isGracefulDeleter {
				params, err := addObjectParams(ws, versionedDeleteOptions)
				if err != nil {
					return nil, err
				}
				action.Params = append(action.Params, params...)
			}
		case "CONNECT":
			if versionedConnectOptions != nil {
				params, err := addObjectParams(ws, versionedConnectOptions)
				if err != nil {
					return nil, err
				}
				action.Params = append(action.Params, params...)
			}
		}
	}

	// If there is a subresource, kind should be the parent's kind.
	if isSubresource {
		parentStorage, ok := a.group.Storage[resource]
		if !ok {
			return nil, fmt.Errorf("missing parent storage: %q", resource)
		}

		fqParentKind, err := GetResourceKind(a.group.GroupVersion, parentStorage, a.group.Typer)
		if err != nil {
			return nil, err
		}
		kind = fqParentKind.Kind
	}
	kubeVerbs := map[string]struct{}{}
	// Create Routes for the actions.
	for _, action := range actions {
		if kubeVerb, found := toDiscoveryKubeVerb[action.Verb]; found {
			if len(kubeVerb) != 0 {
				kubeVerbs[kubeVerb] = struct{}{}
			}
		} else {
			return nil, fmt.Errorf("unknown action verb for discovery: %s", action.Verb)
		}

		var routes []*restful.RouteBuilder
		var err error
		if action.Verb == "CONNECT" {
			routes, err = registerCONNECTToWebService(action, ws, namespaceScoped, isSubresource, isLister, isWatcher, isGracefulDeleter, kind, subresource, storageMeta, connecter, kubeVerbs)
		} else {
			doc, err := documentationAction(action, kind, subresource, isSubresource, isLister, isWatcher)
			if err != nil {
				return nil, err
			}
			routes, err = registerActionsToWebService(action, ws, namespaceScoped, isSubresource, isLister, isWatcher, isGracefulDeleter, kind, subresource, toRESTMethod[action.Verb], doc)
		}
		if err != nil {
			return nil, err
		}
		for _, route := range routes {
			route.Metadata(ROUTE_META_GVK, metav1.GroupVersionKind{
				Group:   reqScope.Kind.Group,
				Version: reqScope.Kind.Version,
				Kind:    reqScope.Kind.Kind,
			})
			route.Metadata(ROUTE_META_ACTION, strings.ToLower(action.Verb))
			ws.Route(route)
		}
		// Note: update GetAuthorizerAttributes() when adding a custom handler.
	}

	return a.constructAPIResource(storage, kubeVerbs, namespaceScoped, path, resourceKind)
}

func (a *APIInstaller) constructAPIResource(storage rest.Storage, kubeVerbs map[string]struct{}, namespaceScoped bool, path, resourceKind string) (*metav1.APIResource, error) {
	var apiResource metav1.APIResource
	apiResource.Name = path
	apiResource.Namespaced = namespaceScoped
	apiResource.Kind = resourceKind
	storageVersionProvider, isStorageVersionProvider := storage.(rest.StorageVersionProvider)
	if utilfeature.DefaultFeatureGate.Enabled(features.StorageVersionHash) &&
		isStorageVersionProvider &&
		storageVersionProvider.StorageVersion() != nil {
		versioner := storageVersionProvider.StorageVersion()
		gvk, err := getStorageVersionKind(versioner, storage, a.group.Typer)
		if err != nil {
			return nil, err
		}
		apiResource.StorageVersionHash = discovery.StorageVersionHash(gvk.Group, gvk.Version, gvk.Kind)
	}
	apiResource.Verbs = make([]string, 0, len(kubeVerbs))
	for kubeVerb := range kubeVerbs {
		apiResource.Verbs = append(apiResource.Verbs, kubeVerb)
	}
	sort.Strings(apiResource.Verbs)

	if shortNamesProvider, ok := storage.(rest.ShortNamesProvider); ok {
		apiResource.ShortNames = shortNamesProvider.ShortNames()
	}
	if categoriesProvider, ok := storage.(rest.CategoriesProvider); ok {
		apiResource.Categories = categoriesProvider.Categories()
	}
	if gvkProvider, ok := storage.(rest.GroupVersionKindProvider); ok {
		gvk := gvkProvider.GroupVersionKind(a.group.GroupVersion)
		apiResource.Group = gvk.Group
		apiResource.Version = gvk.Version
		apiResource.Kind = gvk.Kind
	}

	return &apiResource, nil
}

func supportedPatchTypes() []string {
	supportedTypes := []string{
		string(types.JSONPatchType),
		string(types.MergePatchType),
		string(types.StrategicMergePatchType),
	}
	if utilfeature.DefaultFeatureGate.Enabled(features.ServerSideApply) {
		supportedTypes = append(supportedTypes, string(types.ApplyPatchType))
	}
	return supportedTypes
}

func (a *APIInstaller) constructRequestScope(fqKindToRegister schema.GroupVersionKind, tableProvider rest.TableConvertor, resource, subresource string, namer handlers.ContextBasedNaming) (handlers.RequestScope, error) {
	reqScope := handlers.RequestScope{
		Serializer:      a.group.Serializer,
		ParameterCodec:  a.group.ParameterCodec,
		Creater:         a.group.Creater,
		Convertor:       a.group.Convertor,
		Defaulter:       a.group.Defaulter,
		Typer:           a.group.Typer,
		UnsafeConvertor: a.group.UnsafeConvertor,
		Authorizer:      a.group.Authorizer,

		// TODO: Check for the interface on storage
		TableConvertor: tableProvider,

		// TODO: This seems wrong for cross-group subresources. It makes an assumption that a subresource and its parent are in the same group version. Revisit this.
		Resource:    a.group.GroupVersion.WithResource(resource),
		Subresource: subresource,
		Kind:        fqKindToRegister,

		HubGroupVersion: schema.GroupVersion{Group: fqKindToRegister.Group, Version: runtime.APIVersionInternal},

		MetaGroupVersion: metav1.SchemeGroupVersion,

		MaxRequestBodyBytes: a.group.MaxRequestBodyBytes,
	}
	if a.group.MetaGroupVersion != nil {
		reqScope.MetaGroupVersion = *a.group.MetaGroupVersion
	}
	if a.group.OpenAPIModels != nil && utilfeature.DefaultFeatureGate.Enabled(features.ServerSideApply) {
		fm, err := fieldmanager.NewFieldManager(
			a.group.OpenAPIModels,
			a.group.UnsafeConvertor,
			a.group.Defaulter,
			fqKindToRegister.GroupVersion(),
			reqScope.HubGroupVersion,
		)
		if err != nil {
			return handlers.RequestScope{}, fmt.Errorf("failed to create field manager: %v", err)
		}
		reqScope.FieldManager = fm
	}
	reqScope.Namer = namer

	return reqScope, nil
}

func (a *action) assignConsumedMIMETypes() *action {
	switch a.Verb {
	case "PATCH": // Partially update a resource
		a.consumedMIMETypes = supportedPatchTypes()
	case "CONNECT":
		a.consumedMIMETypes = []string{"*/*"}
	}
	return a
}

func (a *action) assignProducedMIMETypes(producedMIMETypes, mediaTypes, allMediaTypes []string) *action {
	switch a.Verb {
	case "LIST":
		producedMIMETypes = append(producedMIMETypes, allMediaTypes...)
	case "WATCH", "WATCHLIST":
		producedMIMETypes = allMediaTypes
	case "CONNECT":
		producedMIMETypes = []string{"*/*"}
	default:
		producedMIMETypes = append(producedMIMETypes, mediaTypes...)
	}
	a.producedMIMETypes = producedMIMETypes
	return a
}

func (a *action) assignWriteSample(writeSample, versionedList, versionedStatus, versionedWatchEvent interface{}) *action {
	switch a.Verb {
	case "LIST":
		writeSample = versionedList
	case "DELETE", "DELETECOLLECTION":
		writeSample = versionedStatus
	case "WATCH", "WATCHLIST":
		writeSample = versionedWatchEvent
	}

	a.writeSample = writeSample
	return a
}

func (a *action) assignReadSample(defaultVersionedObject, versionedDeleterObject interface{}) *action {
	switch a.Verb {
	case "PUT", "POST":
		a.readSample = defaultVersionedObject
	case "DELETE":
		a.readSample = versionedDeleterObject
	case "PATCH":
		a.readSample = metav1.Patch{}
	}
	return a
}

func documentationAction(action *action, kind, subresource string, isSubresource, isLister, isWatcher bool) (string, error) {
	var doc string
	switch action.Verb {
	case "GET": // Get a resource.
		doc = "read the specified " + kind
		if isSubresource {
			doc = "read " + subresource + " of the specified " + kind
		}
	case "LIST": // List all resources of a kind.
		doc = "list objects of kind " + kind
		if isSubresource {
			doc = "list " + subresource + " of objects of kind " + kind
		}
		switch {
		case isLister && isWatcher:
			doc = "list or watch objects of kind " + kind
			if isSubresource {
				doc = "list or watch " + subresource + " of objects of kind " + kind
			}
		case isWatcher:
			doc = "watch objects of kind " + kind
			if isSubresource {
				doc = "watch " + subresource + "of objects of kind " + kind
			}
		}
	case "PUT": // Update a resource.
		doc = "replace the specified " + kind
		if isSubresource {
			doc = "replace " + subresource + " of the specified " + kind
		}
	case "PATCH": // Partially update a resource
		doc = "partially update the specified " + kind
		if isSubresource {
			doc = "partially update " + subresource + " of the specified " + kind
		}
	case "POST": // Create a resource.
		article := GetArticleForNoun(kind, " ")
		doc = "create" + article + kind
		if isSubresource {
			doc = "create " + subresource + " of" + article + kind
		}
	case "DELETE": // Delete a resource.
		article := GetArticleForNoun(kind, " ")
		doc = "delete" + article + kind
		if isSubresource {
			doc = "delete " + subresource + " of" + article + kind
		}
	case "DELETECOLLECTION":
		doc = "delete collection of " + kind
		if isSubresource {
			doc = "delete collection of " + subresource + " of a " + kind
		}
	// deprecated in 1.11
	case "WATCH": // Watch a resource.
		doc = "watch changes to an object of kind " + kind
		if isSubresource {
			doc = "watch changes to " + subresource + " of an object of kind " + kind
		}
		doc += ". deprecated: use the 'watch' parameter with a list operation instead, filtered to a single item with the 'fieldSelector' parameter."
	case "WATCHLIST": // Watch all resources of a kind.
		doc = "watch individual changes to a list of " + kind
		if isSubresource {
			doc = "watch individual changes to a list of " + subresource + " of " + kind
		}
		doc += ". deprecated: use the 'watch' parameter with a list operation instead."
	default:
		return "", fmt.Errorf("unrecognized action verb: %s", action.Verb)
	}
	return doc, nil
}

func registerActionsToWebService(action *action, ws *restful.WebService, isNamespaced, isSubresource, isLister, isWatcher, isGracefulDeleter bool, kind, subresource, method, doc string) ([]*restful.RouteBuilder, error) {
	var namespaced string
	var operationSuffix string
	if isNamespaced {
		namespaced = "Namespaced"
	}
	if strings.HasSuffix(action.Path, "/{path:*}") {
		operationSuffix = operationSuffix + "WithPath"
	}
	if action.AllNamespaces {
		operationSuffix = operationSuffix + "ForAllNamespaces"
		namespaced = ""
	}

	routes := []*restful.RouteBuilder{}
	var operation string
	var route *restful.RouteBuilder
	switch action.Verb {
	case "GET": // Get a resource.
		operation = "read"
	case "LIST": // List all resources of a kind.
		operation = "list"
	case "PUT": // Update a resource.
		operation = "replace"
	case "PATCH": // Partially update a resource
		operation = "patch"
	case "POST": // Create a resource.
		operation = "create"
	case "DELETE": // Delete a resource.
		operation = "delete"
	case "DELETECOLLECTION":
		operation = "deletecollection"
	// deprecated in 1.11
	case "WATCH": // Watch a resource.
		operation = "watch"
	// deprecated in 1.11
	case "WATCHLIST": // Watch all resources of a kind.
		operation = "watch"
	default:
		return nil, fmt.Errorf("unrecognized action verb: %s", action.Verb)
	}

	route = ws.Method(method).Path(action.Path).To(action.handler).
		Doc(doc).
		Param(ws.QueryParameter("pretty", "If 'true', then the output is pretty printed.")).
		Operation(operation + namespaced + kind + strings.Title(subresource) + operationSuffix).
		Produces(action.producedMIMETypes...).
		Writes(action.writeSample)

	// corner case
	if action.Verb == "WATCHLIST" {
		route = route.Operation("watch" + namespaced + kind + strings.Title(subresource) + "List" + operationSuffix)
	}

	for _, statusCode := range action.returnStatusCodes {
		route.Returns(statusCode, http.StatusText(statusCode), action.writeSample)
	}
	if action.readSample != nil {
		if action.Verb == "DELETE" {
			if isGracefulDeleter {
				route.Reads(action.readSample)
				route.ParameterNamed("body").Required(false)
			}
		} else {
			route.Reads(action.readSample)
		}
	}
	if action.consumedMIMETypes != nil {
		route.Consumes(action.consumedMIMETypes...)
	}
	addParams(route, action.Params)
	routes = append(routes, route)
	return routes, nil
}

func registerCONNECTToWebService(action *action, ws *restful.WebService, isNamespaced, isSubresource, isLister, isWatcher, isGracefulDeleter bool, kind, subresource string, storageMeta rest.StorageMetadata, connecter rest.Connecter, kubeVerbs map[string]struct{}) ([]*restful.RouteBuilder, error) {
	var namespaced string
	var operationSuffix string
	if isNamespaced {
		namespaced = "Namespaced"
	}
	if strings.HasSuffix(action.Path, "/{path:*}") {
		operationSuffix = operationSuffix + "WithPath"
	}
	if action.AllNamespaces {
		operationSuffix = operationSuffix + "ForAllNamespaces"
		namespaced = ""
	}

	routes := []*restful.RouteBuilder{}
	for _, method := range connecter.ConnectMethods() {
		connectProducedObject := storageMeta.ProducesObject(method)
		if connectProducedObject == nil {
			connectProducedObject = "string"
		}
		doc := "connect " + method + " requests to " + kind
		if isSubresource {
			doc = "connect " + method + " requests to " + subresource + " of " + kind
		}
		route := ws.Method(method).Path(action.Path).
			To(action.handler).
			Doc(doc).
			Operation("connect" + strings.Title(strings.ToLower(method)) + namespaced + kind + strings.Title(subresource) + operationSuffix).
			Produces(action.producedMIMETypes...).
			Consumes(action.consumedMIMETypes...).
			Writes(connectProducedObject)
		addParams(route, action.Params)
		routes = append(routes, route)

		// transform ConnectMethods to kube verbs
		if kubeVerb, found := toDiscoveryKubeVerb[method]; found {
			if len(kubeVerb) != 0 {
				kubeVerbs[kubeVerb] = struct{}{}
			}
		}
	}
	return routes, nil
}

// indirectArbitraryPointer returns *ptrToObject for an arbitrary pointer
func indirectArbitraryPointer(ptrToObject interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(ptrToObject)).Interface()
}

func appendIf(actions []*action, a *action, shouldAppend bool) []*action {
	if shouldAppend {
		actions = append(actions, a)
	}
	return actions
}

func addParams(route *restful.RouteBuilder, params []*restful.Parameter) {
	for _, param := range params {
		route.Param(param)
	}
}

// AddObjectParams converts a runtime.Object into a set of go-restful Param() definitions on the route.
// The object must be a pointer to a struct; only fields at the top level of the struct that are not
// themselves interfaces or structs are used; only fields with a json tag that is non empty (the standard
// Go JSON behavior for omitting a field) become query parameters. The name of the query parameter is
// the JSON field name. If a description struct tag is set on the field, that description is used on the
// query parameter. In essence, it converts a standard JSON top level object into a query param schema.
func AddObjectParams(ws *restful.WebService, route *restful.RouteBuilder, obj interface{}) error {
	sv, err := conversion.EnforcePtr(obj)
	if err != nil {
		return err
	}
	st := sv.Type()
	switch st.Kind() {
	case reflect.Struct:
		for i := 0; i < st.NumField(); i++ {
			name := st.Field(i).Name
			sf, ok := st.FieldByName(name)
			if !ok {
				continue
			}
			switch sf.Type.Kind() {
			case reflect.Interface, reflect.Struct:
			case reflect.Ptr:
				// TODO: This is a hack to let metav1.Time through. This needs to be fixed in a more generic way eventually. bug #36191
				if (sf.Type.Elem().Kind() == reflect.Interface || sf.Type.Elem().Kind() == reflect.Struct) && strings.TrimPrefix(sf.Type.String(), "*") != "metav1.Time" {
					continue
				}
				fallthrough
			default:
				jsonTag := sf.Tag.Get("json")
				if len(jsonTag) == 0 {
					continue
				}
				jsonName := strings.SplitN(jsonTag, ",", 2)[0]
				if len(jsonName) == 0 {
					continue
				}

				var desc string
				if docable, ok := obj.(documentable); ok {
					desc = docable.SwaggerDoc()[jsonName]
				}
				route.Param(ws.QueryParameter(jsonName, desc).DataType(typeToJSON(sf.Type.String())))
			}
		}
	}
	return nil
}

// addObjectParams converts a runtime.Object into a set of go-restful Param() definitions on the route.
// The object must be a pointer to a struct; only fields at the top level of the struct that are not
// themselves interfaces or structs are used; only fields with a json tag that is non empty (the standard
// Go JSON behavior for omitting a field) become query parameters. The name of the query parameter is
// the JSON field name. If a description struct tag is set on the field, that description is used on the
// query parameter. In essence, it converts a standard JSON top level object into a query param schema.
func addObjectParams(ws *restful.WebService, obj interface{}) ([]*restful.Parameter, error) {
	params := []*restful.Parameter{}
	sv, err := conversion.EnforcePtr(obj)
	if err != nil {
		return nil, err
	}
	st := sv.Type()
	switch st.Kind() {
	case reflect.Struct:
		for i := 0; i < st.NumField(); i++ {
			name := st.Field(i).Name
			sf, ok := st.FieldByName(name)
			if !ok {
				continue
			}
			switch sf.Type.Kind() {
			case reflect.Interface, reflect.Struct:
			case reflect.Ptr:
				// TODO: This is a hack to let metav1.Time through. This needs to be fixed in a more generic way eventually. bug #36191
				if (sf.Type.Elem().Kind() == reflect.Interface || sf.Type.Elem().Kind() == reflect.Struct) && strings.TrimPrefix(sf.Type.String(), "*") != "metav1.Time" {
					continue
				}
				fallthrough
			default:
				jsonTag := sf.Tag.Get("json")
				if len(jsonTag) == 0 {
					continue
				}
				jsonName := strings.SplitN(jsonTag, ",", 2)[0]
				if len(jsonName) == 0 {
					continue
				}

				var desc string
				if docable, ok := obj.(documentable); ok {
					desc = docable.SwaggerDoc()[jsonName]
				}
				params = append(params, ws.QueryParameter(jsonName, desc).DataType(typeToJSON(sf.Type.String())))
			}
		}
	}
	return params, nil
}

// TODO: this is incomplete, expand as needed.
// Convert the name of a golang type to the name of a JSON type
func typeToJSON(typeName string) string {
	switch typeName {
	case "bool", "*bool":
		return "boolean"
	case "uint8", "*uint8", "int", "*int", "int32", "*int32", "int64", "*int64", "uint32", "*uint32", "uint64", "*uint64":
		return "integer"
	case "float64", "*float64", "float32", "*float32":
		return "number"
	case "metav1.Time", "*metav1.Time":
		return "string"
	case "byte", "*byte":
		return "string"
	case "v1.DeletionPropagation", "*v1.DeletionPropagation":
		return "string"

	// TODO: Fix these when go-restful supports a way to specify an array query param:
	// https://github.com/emicklei/go-restful/issues/225
	case "[]string", "[]*string":
		return "string"
	case "[]int32", "[]*int32":
		return "integer"

	default:
		return typeName
	}
}

// defaultStorageMetadata provides default answers to rest.StorageMetadata.
type defaultStorageMetadata struct{}

// defaultStorageMetadata implements rest.StorageMetadata
var _ rest.StorageMetadata = defaultStorageMetadata{}

func (defaultStorageMetadata) ProducesMIMETypes(verb string) []string {
	return nil
}

func (defaultStorageMetadata) ProducesObject(verb string) interface{} {
	return nil
}

// splitSubresource checks if the given storage path is the path of a subresource and returns
// the resource and subresource components.
func splitSubresource(path string) (string, string, error) {
	var resource, subresource string
	switch parts := strings.Split(path, "/"); len(parts) {
	case 2:
		resource, subresource = parts[0], parts[1]
	case 1:
		resource = parts[0]
	default:
		// TODO: support deeper paths
		return "", "", fmt.Errorf("api_installer allows only one or two segment paths (resource or resource/subresource)")
	}
	return resource, subresource, nil
}

// GetArticleForNoun returns the article needed for the given noun.
func GetArticleForNoun(noun string, padding string) string {
	if noun[len(noun)-2:] != "ss" && noun[len(noun)-1:] == "s" {
		// Plurals don't have an article.
		// Don't catch words like class
		return fmt.Sprintf("%v", padding)
	}

	article := "a"
	if isVowel(rune(noun[0])) {
		article = "an"
	}

	return fmt.Sprintf("%s%s%s", padding, article, padding)
}

// isVowel returns true if the rune is a vowel (case insensitive).
func isVowel(c rune) bool {
	vowels := []rune{'a', 'e', 'i', 'o', 'u'}
	for _, value := range vowels {
		if value == unicode.ToLower(c) {
			return true
		}
	}
	return false
}

func restfulListResource(r rest.Lister, rw rest.Watcher, scope handlers.RequestScope, forceWatch bool, minRequestTimeout time.Duration) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.ListResource(r, rw, scope, forceWatch, minRequestTimeout)(res.ResponseWriter, req.Request)
	}
}

func restfulCreateNamedResource(r rest.NamedCreater, scope handlers.RequestScope, admit admission.Interface) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.CreateNamedResource(r, scope, admit)(res.ResponseWriter, req.Request)
	}
}

func restfulCreateResource(r rest.Creater, scope handlers.RequestScope, admit admission.Interface) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.CreateResource(r, scope, admit)(res.ResponseWriter, req.Request)
	}
}

func restfulDeleteResource(r rest.GracefulDeleter, allowsOptions bool, scope handlers.RequestScope, admit admission.Interface) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.DeleteResource(r, allowsOptions, scope, admit)(res.ResponseWriter, req.Request)
	}
}

func restfulDeleteCollection(r rest.CollectionDeleter, checkBody bool, scope handlers.RequestScope, admit admission.Interface) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.DeleteCollection(r, checkBody, scope, admit)(res.ResponseWriter, req.Request)
	}
}

func restfulUpdateResource(r rest.Updater, scope handlers.RequestScope, admit admission.Interface) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.UpdateResource(r, scope, admit)(res.ResponseWriter, req.Request)
	}
}

func restfulPatchResource(r rest.Patcher, scope handlers.RequestScope, admit admission.Interface, supportedTypes []string) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.PatchResource(r, scope, admit, supportedTypes)(res.ResponseWriter, req.Request)
	}
}

func restfulGetResource(r rest.Getter, e rest.Exporter, scope handlers.RequestScope) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.GetResource(r, e, scope)(res.ResponseWriter, req.Request)
	}
}

func restfulGetResourceWithOptions(r rest.GetterWithOptions, scope handlers.RequestScope, isSubresource bool) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.GetResourceWithOptions(r, scope, isSubresource)(res.ResponseWriter, req.Request)
	}
}

func restfulConnectResource(connecter rest.Connecter, scope handlers.RequestScope, admit admission.Interface, restPath string, isSubresource bool) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.ConnectResource(connecter, scope, admit, restPath, isSubresource)(res.ResponseWriter, req.Request)
	}
}
