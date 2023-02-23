/*
Copyright 2016 The Kubernetes Authors.

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

package etcd3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/attribute"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/audit"
	"k8s.io/apiserver/pkg/features"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/etcd3/metrics"
	"k8s.io/apiserver/pkg/storage/value"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/component-base/tracing"
	"k8s.io/klog/v2"
)

const (
	// maxLimit is a maximum page limit increase used when fetching objects from etcd.
	// This limit is used only for increasing page size by kube-apiserver. If request
	// specifies larger limit initially, it won't be changed.
	maxLimit = 10000
)

// authenticatedDataString satisfies the value.Context interface. It uses the key to
// authenticate the stored data. This does not defend against reuse of previously
// encrypted values under the same key, but will prevent an attacker from using an
// encrypted value from a different key. A stronger authenticated data segment would
// include the etcd3 Version field (which is incremented on each write to a key and
// reset when the key is deleted), but an attacker with write access to etcd can
// force deletion and recreation of keys to weaken that angle.
type authenticatedDataString string

// AuthenticatedData implements the value.Context interface.
func (d authenticatedDataString) AuthenticatedData() []byte {
	return []byte(string(d))
}

var _ value.Context = authenticatedDataString("")

type store struct {
	client              *clientv3.Client
	codec               runtime.Codec
	versioner           storage.Versioner
	transformer         value.Transformer
	pathPrefix          string
	groupResource       schema.GroupResource
	groupResourceString string
	watcher             *watcher
	pagingEnabled       bool
	leaseManager        *leaseManager
}

type objState struct {
	obj   runtime.Object
	meta  *storage.ResponseMeta
	rev   int64
	data  []byte
	stale bool
}

// New returns an etcd3 implementation of storage.Interface.
func New(c *clientv3.Client, codec runtime.Codec, newFunc func() runtime.Object, prefix string, groupResource schema.GroupResource, transformer value.Transformer, pagingEnabled bool, leaseManagerConfig LeaseManagerConfig) storage.Interface {
	return newStore(c, codec, newFunc, prefix, groupResource, transformer, pagingEnabled, leaseManagerConfig)
}

func newStore(c *clientv3.Client, codec runtime.Codec, newFunc func() runtime.Object, prefix string, groupResource schema.GroupResource, transformer value.Transformer, pagingEnabled bool, leaseManagerConfig LeaseManagerConfig) *store {
	versioner := storage.APIObjectVersioner{}
	// for compatibility with etcd2 impl.
	// no-op for default prefix of '/registry'.
	// keeps compatibility with etcd2 impl for custom prefixes that don't start with '/'
	pathPrefix := path.Join("/", prefix)
	if !strings.HasSuffix(pathPrefix, "/") {
		// Ensure the pathPrefix ends in "/" here to simplify key concatenation later.
		pathPrefix += "/"
	}
	result := &store{
		client:              c,
		codec:               codec,
		versioner:           versioner,
		transformer:         transformer,
		pagingEnabled:       pagingEnabled,
		pathPrefix:          pathPrefix,
		groupResource:       groupResource,
		groupResourceString: groupResource.String(),
		watcher:             newWatcher(c, codec, groupResource, newFunc, versioner),
		leaseManager:        newDefaultLeaseManager(c, leaseManagerConfig),
	}
	return result
}

// Versioner implements storage.Interface.Versioner.
func (s *store) Versioner() storage.Versioner {
	return s.versioner
}

// Get implements storage.Interface.Get.
func (s *store) Get(ctx context.Context, key string, opts storage.GetOptions, out runtime.Object) error {
	preparedKey, err := s.prepareKey(key)
	if err != nil {
		return err
	}
	startTime := time.Now()
	getResp, err := s.client.KV.Get(ctx, preparedKey)
	metrics.RecordEtcdRequestLatency("get", s.groupResourceString, startTime)
	if err != nil {
		return err
	}
	if err = s.validateMinimumResourceVersion(opts.ResourceVersion, uint64(getResp.Header.Revision)); err != nil {
		return err
	}

	if len(getResp.Kvs) == 0 {
		if opts.IgnoreNotFound {
			return runtime.SetZeroValue(out)
		}
		return storage.NewKeyNotFoundError(preparedKey, 0)
	}
	kv := getResp.Kvs[0]

	data, _, err := s.transformer.TransformFromStorage(ctx, kv.Value, authenticatedDataString(preparedKey))
	if err != nil {
		return storage.NewInternalError(err.Error())
	}

	err = decode(s.codec, s.versioner, data, out, kv.ModRevision)
	if err != nil {
		recordDecodeError(s.groupResourceString, preparedKey)
		return err
	}
	return nil
}

// Create implements storage.Interface.Create.
func (s *store) Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error {
	preparedKey, err := s.prepareKey(key)
	if err != nil {
		return err
	}
	ctx, span := tracing.Start(ctx, "Create etcd3",
		attribute.String("audit-id", audit.GetAuditIDTruncated(ctx)),
		attribute.String("key", key),
		attribute.String("type", getTypeName(obj)),
		attribute.String("resource", s.groupResourceString),
	)
	defer span.End(500 * time.Millisecond)
	if version, err := s.versioner.ObjectResourceVersion(obj); err == nil && version != 0 {
		return errors.New("resourceVersion should not be set on objects to be created")
	}
	if err := s.versioner.PrepareObjectForStorage(obj); err != nil {
		return fmt.Errorf("PrepareObjectForStorage failed: %v", err)
	}
	span.AddEvent("About to Encode")
	data, err := runtime.Encode(s.codec, obj)
	if err != nil {
		span.AddEvent("Encode failed", attribute.Int("len", len(data)), attribute.String("err", err.Error()))
		return err
	}
	span.AddEvent("Encode succeeded", attribute.Int("len", len(data)))

	opts, err := s.ttlOpts(ctx, int64(ttl))
	if err != nil {
		return err
	}

	newData, err := s.transformer.TransformToStorage(ctx, data, authenticatedDataString(preparedKey))
	if err != nil {
		span.AddEvent("TransformToStorage failed", attribute.String("err", err.Error()))
		return storage.NewInternalError(err.Error())
	}
	span.AddEvent("TransformToStorage succeeded")

	startTime := time.Now()
	txnResp, err := s.client.KV.Txn(ctx).If(
		notFound(preparedKey),
	).Then(
		clientv3.OpPut(preparedKey, string(newData), opts...),
	).Commit()
	metrics.RecordEtcdRequestLatency("create", s.groupResourceString, startTime)
	if err != nil {
		span.AddEvent("Txn call failed", attribute.String("err", err.Error()))
		return err
	}
	span.AddEvent("Txn call succeeded")

	if !txnResp.Succeeded {
		return storage.NewKeyExistsError(preparedKey, 0)
	}

	if out != nil {
		putResp := txnResp.Responses[0].GetResponsePut()
		err = decode(s.codec, s.versioner, data, out, putResp.Header.Revision)
		if err != nil {
			span.AddEvent("decode failed", attribute.Int("len", len(data)), attribute.String("err", err.Error()))
			recordDecodeError(s.groupResourceString, preparedKey)
			return err
		}
		span.AddEvent("decode succeeded", attribute.Int("len", len(data)))
	}
	return nil
}

// Delete implements storage.Interface.Delete.
func (s *store) Delete(
	ctx context.Context, key string, out runtime.Object, preconditions *storage.Preconditions,
	validateDeletion storage.ValidateObjectFunc, cachedExistingObject runtime.Object) error {
	preparedKey, err := s.prepareKey(key)
	if err != nil {
		return err
	}
	v, err := conversion.EnforcePtr(out)
	if err != nil {
		return fmt.Errorf("unable to convert output object to pointer: %v", err)
	}
	return s.conditionalDelete(ctx, preparedKey, out, v, preconditions, validateDeletion, cachedExistingObject)
}

func (s *store) conditionalDelete(
	ctx context.Context, key string, out runtime.Object, v reflect.Value, preconditions *storage.Preconditions,
	validateDeletion storage.ValidateObjectFunc, cachedExistingObject runtime.Object) error {
	getCurrentState := func() (*objState, error) {
		startTime := time.Now()
		getResp, err := s.client.KV.Get(ctx, key)
		metrics.RecordEtcdRequestLatency("get", s.groupResourceString, startTime)
		if err != nil {
			return nil, err
		}
		return s.getState(ctx, getResp, key, v, false)
	}

	var origState *objState
	var err error
	var origStateIsCurrent bool
	if cachedExistingObject != nil {
		origState, err = s.getStateFromObject(cachedExistingObject)
	} else {
		origState, err = getCurrentState()
		origStateIsCurrent = true
	}
	if err != nil {
		return err
	}

	for {
		if preconditions != nil {
			if err := preconditions.Check(key, origState.obj); err != nil {
				if origStateIsCurrent {
					return err
				}

				// It's possible we're working with stale data.
				// Remember the revision of the potentially stale data and the resulting update error
				cachedRev := origState.rev
				cachedUpdateErr := err

				// Actually fetch
				origState, err = getCurrentState()
				if err != nil {
					return err
				}
				origStateIsCurrent = true

				// it turns out our cached data was not stale, return the error
				if cachedRev == origState.rev {
					return cachedUpdateErr
				}

				// Retry
				continue
			}
		}
		if err := validateDeletion(ctx, origState.obj); err != nil {
			if origStateIsCurrent {
				return err
			}

			// It's possible we're working with stale data.
			// Remember the revision of the potentially stale data and the resulting update error
			cachedRev := origState.rev
			cachedUpdateErr := err

			// Actually fetch
			origState, err = getCurrentState()
			if err != nil {
				return err
			}
			origStateIsCurrent = true

			// it turns out our cached data was not stale, return the error
			if cachedRev == origState.rev {
				return cachedUpdateErr
			}

			// Retry
			continue
		}

		startTime := time.Now()
		txnResp, err := s.client.KV.Txn(ctx).If(
			clientv3.Compare(clientv3.ModRevision(key), "=", origState.rev),
		).Then(
			clientv3.OpDelete(key),
		).Else(
			clientv3.OpGet(key),
		).Commit()
		metrics.RecordEtcdRequestLatency("delete", s.groupResourceString, startTime)
		if err != nil {
			return err
		}
		if !txnResp.Succeeded {
			getResp := (*clientv3.GetResponse)(txnResp.Responses[0].GetResponseRange())
			klog.V(4).Infof("deletion of %s failed because of a conflict, going to retry", key)
			origState, err = s.getState(ctx, getResp, key, v, false)
			if err != nil {
				return err
			}
			origStateIsCurrent = true
			continue
		}

		if len(txnResp.Responses) == 0 || txnResp.Responses[0].GetResponseDeleteRange() == nil {
			return errors.New(fmt.Sprintf("invalid DeleteRange response: %v", txnResp.Responses))
		}
		deleteResp := txnResp.Responses[0].GetResponseDeleteRange()
		if deleteResp.Header == nil {
			return errors.New("invalid DeleteRange response - nil header")
		}
		err = decode(s.codec, s.versioner, origState.data, out, deleteResp.Header.Revision)
		if err != nil {
			recordDecodeError(s.groupResourceString, key)
			return err
		}
		return nil
	}
}

// GuaranteedUpdate implements storage.Interface.GuaranteedUpdate.
func (s *store) GuaranteedUpdate(
	ctx context.Context, key string, destination runtime.Object, ignoreNotFound bool,
	preconditions *storage.Preconditions, tryUpdate storage.UpdateFunc, cachedExistingObject runtime.Object) error {
	preparedKey, err := s.prepareKey(key)
	if err != nil {
		return err
	}
	ctx, span := tracing.Start(ctx, "GuaranteedUpdate etcd3",
		attribute.String("audit-id", audit.GetAuditIDTruncated(ctx)),
		attribute.String("key", key),
		attribute.String("type", getTypeName(destination)),
		attribute.String("resource", s.groupResourceString))
	defer span.End(500 * time.Millisecond)

	v, err := conversion.EnforcePtr(destination)
	if err != nil {
		return fmt.Errorf("unable to convert output object to pointer: %v", err)
	}

	getCurrentState := func() (*objState, error) {
		startTime := time.Now()
		getResp, err := s.client.KV.Get(ctx, preparedKey)
		metrics.RecordEtcdRequestLatency("get", s.groupResourceString, startTime)
		if err != nil {
			return nil, err
		}
		return s.getState(ctx, getResp, preparedKey, v, ignoreNotFound)
	}

	var origState *objState
	var origStateIsCurrent bool
	if cachedExistingObject != nil {
		origState, err = s.getStateFromObject(cachedExistingObject)
	} else {
		origState, err = getCurrentState()
		origStateIsCurrent = true
	}
	if err != nil {
		return err
	}
	span.AddEvent("initial value restored")

	transformContext := authenticatedDataString(preparedKey)
	for {
		if err := preconditions.Check(preparedKey, origState.obj); err != nil {
			// If our data is already up to date, return the error
			if origStateIsCurrent {
				return err
			}

			// It's possible we were working with stale data
			// Actually fetch
			origState, err = getCurrentState()
			if err != nil {
				return err
			}
			origStateIsCurrent = true
			// Retry
			continue
		}

		ret, ttl, err := s.updateState(origState, tryUpdate)
		if err != nil {
			// If our data is already up to date, return the error
			if origStateIsCurrent {
				return err
			}

			// It's possible we were working with stale data
			// Remember the revision of the potentially stale data and the resulting update error
			cachedRev := origState.rev
			cachedUpdateErr := err

			// Actually fetch
			origState, err = getCurrentState()
			if err != nil {
				return err
			}
			origStateIsCurrent = true

			// it turns out our cached data was not stale, return the error
			if cachedRev == origState.rev {
				return cachedUpdateErr
			}

			// Retry
			continue
		}

		span.AddEvent("About to Encode")
		data, err := runtime.Encode(s.codec, ret)
		if err != nil {
			span.AddEvent("Encode failed", attribute.Int("len", len(data)), attribute.String("err", err.Error()))
			return err
		}
		span.AddEvent("Encode succeeded", attribute.Int("len", len(data)))
		if !origState.stale && bytes.Equal(data, origState.data) {
			// if we skipped the original Get in this loop, we must refresh from
			// etcd in order to be sure the data in the store is equivalent to
			// our desired serialization
			if !origStateIsCurrent {
				origState, err = getCurrentState()
				if err != nil {
					return err
				}
				origStateIsCurrent = true
				if !bytes.Equal(data, origState.data) {
					// original data changed, restart loop
					continue
				}
			}
			// recheck that the data from etcd is not stale before short-circuiting a write
			if !origState.stale {
				err = decode(s.codec, s.versioner, origState.data, destination, origState.rev)
				if err != nil {
					recordDecodeError(s.groupResourceString, preparedKey)
					return err
				}
				return nil
			}
		}

		newData, err := s.transformer.TransformToStorage(ctx, data, transformContext)
		if err != nil {
			span.AddEvent("TransformToStorage failed", attribute.String("err", err.Error()))
			return storage.NewInternalError(err.Error())
		}
		span.AddEvent("TransformToStorage succeeded")

		opts, err := s.ttlOpts(ctx, int64(ttl))
		if err != nil {
			return err
		}
		span.AddEvent("Transaction prepared")

		startTime := time.Now()
		txnResp, err := s.client.KV.Txn(ctx).If(
			clientv3.Compare(clientv3.ModRevision(preparedKey), "=", origState.rev),
		).Then(
			clientv3.OpPut(preparedKey, string(newData), opts...),
		).Else(
			clientv3.OpGet(preparedKey),
		).Commit()
		metrics.RecordEtcdRequestLatency("update", s.groupResourceString, startTime)
		if err != nil {
			span.AddEvent("Txn call failed", attribute.String("err", err.Error()))
			return err
		}
		span.AddEvent("Txn call completed")
		span.AddEvent("Transaction committed")
		if !txnResp.Succeeded {
			getResp := (*clientv3.GetResponse)(txnResp.Responses[0].GetResponseRange())
			klog.V(4).Infof("GuaranteedUpdate of %s failed because of a conflict, going to retry", preparedKey)
			origState, err = s.getState(ctx, getResp, preparedKey, v, ignoreNotFound)
			if err != nil {
				return err
			}
			span.AddEvent("Retry value restored")
			origStateIsCurrent = true
			continue
		}
		putResp := txnResp.Responses[0].GetResponsePut()

		err = decode(s.codec, s.versioner, data, destination, putResp.Header.Revision)
		if err != nil {
			span.AddEvent("decode failed", attribute.Int("len", len(data)), attribute.String("err", err.Error()))
			recordDecodeError(s.groupResourceString, preparedKey)
			return err
		}
		span.AddEvent("decode succeeded", attribute.Int("len", len(data)))
		return nil
	}
}

func getNewItemFunc(listObj runtime.Object, v reflect.Value) func() runtime.Object {
	// For unstructured lists with a target group/version, preserve the group/version in the instantiated list items
	if unstructuredList, isUnstructured := listObj.(*unstructured.UnstructuredList); isUnstructured {
		if apiVersion := unstructuredList.GetAPIVersion(); len(apiVersion) > 0 {
			return func() runtime.Object {
				return &unstructured.Unstructured{Object: map[string]interface{}{"apiVersion": apiVersion}}
			}
		}
	}

	// Otherwise just instantiate an empty item
	elem := v.Type().Elem()
	return func() runtime.Object {
		return reflect.New(elem).Interface().(runtime.Object)
	}
}

func (s *store) Count(key string) (int64, error) {
	preparedKey, err := s.prepareKey(key)
	if err != nil {
		return 0, err
	}

	// We need to make sure the key ended with "/" so that we only get children "directories".
	// e.g. if we have key "/a", "/a/b", "/ab", getting keys with prefix "/a" will return all three,
	// while with prefix "/a/" will return only "/a/b" which is the correct answer.
	if !strings.HasSuffix(preparedKey, "/") {
		preparedKey += "/"
	}

	startTime := time.Now()
	getResp, err := s.client.KV.Get(context.Background(), preparedKey, clientv3.WithRange(clientv3.GetPrefixRangeEnd(preparedKey)), clientv3.WithCountOnly())
	metrics.RecordEtcdRequestLatency("listWithCount", preparedKey, startTime)
	if err != nil {
		return 0, err
	}
	return getResp.Count, nil
}

// GetList implements storage.Interface.
func (s *store) GetList(ctx context.Context, key string, opts storage.ListOptions, listObj runtime.Object) error {
	preparedKey, err := s.prepareKey(key)
	if err != nil {
		return err
	}
	recursive := opts.Recursive
	resourceVersion := opts.ResourceVersion
	match := opts.ResourceVersionMatch
	pred := opts.Predicate
	ctx, span := tracing.Start(ctx, fmt.Sprintf("List(recursive=%v) etcd3", recursive),
		attribute.String("audit-id", audit.GetAuditIDTruncated(ctx)),
		attribute.String("key", key),
		attribute.String("resourceVersion", resourceVersion),
		attribute.String("resourceVersionMatch", string(match)),
		attribute.Int("limit", int(pred.Limit)),
		attribute.String("continue", pred.Continue))
	defer span.End(500 * time.Millisecond)
	listPtr, err := meta.GetItemsPtr(listObj)
	if err != nil {
		return err
	}
	v, err := conversion.EnforcePtr(listPtr)
	if err != nil || v.Kind() != reflect.Slice {
		return fmt.Errorf("need ptr to slice: %v", err)
	}

	// For recursive lists, we need to make sure the key ended with "/" so that we only
	// get children "directories". e.g. if we have key "/a", "/a/b", "/ab", getting keys
	// with prefix "/a" will return all three, while with prefix "/a/" will return only
	// "/a/b" which is the correct answer.
	if recursive && !strings.HasSuffix(preparedKey, "/") {
		preparedKey += "/"
	}
	keyPrefix := preparedKey

	// set the appropriate clientv3 options to filter the returned data set
	var limitOption *clientv3.OpOption
	limit := pred.Limit
	var paging bool
	options := make([]clientv3.OpOption, 0, 4)
	if s.pagingEnabled && pred.Limit > 0 {
		paging = true
		options = append(options, clientv3.WithLimit(limit))
		limitOption = &options[len(options)-1]
	}

	newItemFunc := getNewItemFunc(listObj, v)

	var fromRV *uint64
	if len(resourceVersion) > 0 {
		parsedRV, err := s.versioner.ParseResourceVersion(resourceVersion)
		if err != nil {
			return apierrors.NewBadRequest(fmt.Sprintf("invalid resource version: %v", err))
		}
		fromRV = &parsedRV
	}

	var returnedRV, continueRV, withRev int64
	var continueKey string
	switch {
	case recursive && s.pagingEnabled && len(pred.Continue) > 0:
		continueKey, continueRV, err = storage.DecodeContinue(pred.Continue, keyPrefix)
		if err != nil {
			return apierrors.NewBadRequest(fmt.Sprintf("invalid continue token: %v", err))
		}

		if len(resourceVersion) > 0 && resourceVersion != "0" {
			return apierrors.NewBadRequest("specifying resource version is not allowed when using continue")
		}

		rangeEnd := clientv3.GetPrefixRangeEnd(keyPrefix)
		options = append(options, clientv3.WithRange(rangeEnd))
		preparedKey = continueKey

		// If continueRV > 0, the LIST request needs a specific resource version.
		// continueRV==0 is invalid.
		// If continueRV < 0, the request is for the latest resource version.
		if continueRV > 0 {
			withRev = continueRV
			returnedRV = continueRV
		}
	case recursive && s.pagingEnabled && pred.Limit > 0:
		if fromRV != nil {
			switch match {
			case metav1.ResourceVersionMatchNotOlderThan:
				// The not older than constraint is checked after we get a response from etcd,
				// and returnedRV is then set to the revision we get from the etcd response.
			case metav1.ResourceVersionMatchExact:
				returnedRV = int64(*fromRV)
				withRev = returnedRV
			case "": // legacy case
				if *fromRV > 0 {
					returnedRV = int64(*fromRV)
					withRev = returnedRV
				}
			default:
				return fmt.Errorf("unknown ResourceVersionMatch value: %v", match)
			}
		}

		rangeEnd := clientv3.GetPrefixRangeEnd(keyPrefix)
		options = append(options, clientv3.WithRange(rangeEnd))
	default:
		if fromRV != nil {
			switch match {
			case metav1.ResourceVersionMatchNotOlderThan:
				// The not older than constraint is checked after we get a response from etcd,
				// and returnedRV is then set to the revision we get from the etcd response.
			case metav1.ResourceVersionMatchExact:
				returnedRV = int64(*fromRV)
				withRev = returnedRV
			case "": // legacy case
			default:
				return fmt.Errorf("unknown ResourceVersionMatch value: %v", match)
			}
		}

		if recursive {
			options = append(options, clientv3.WithPrefix())
		}
	}
	if withRev != 0 {
		options = append(options, clientv3.WithRev(withRev))
	}

	// loop until we have filled the requested limit from etcd or there are no more results
	var lastKey []byte
	var hasMore bool
	var getResp *clientv3.GetResponse
	var numFetched int
	var numEvald atomic.Int32
	// Because these metrics are for understanding the costs of handling LIST requests,
	// get them recorded even in error cases.
	defer func() {
		numReturn := v.Len()
		metrics.RecordStorageListMetrics(s.groupResourceString, numFetched, int(numEvald.Load()), numReturn)
	}()
	for {
		startTime := time.Now()
		getResp, err = s.client.KV.Get(ctx, preparedKey, options...)
		if recursive {
			metrics.RecordEtcdRequestLatency("list", s.groupResourceString, startTime)
		} else {
			metrics.RecordEtcdRequestLatency("get", s.groupResourceString, startTime)
		}
		if err != nil {
			return interpretListError(err, len(pred.Continue) > 0, continueKey, keyPrefix)
		}
		kvCount := len(getResp.Kvs)
		numFetched += kvCount
		if err = s.validateMinimumResourceVersion(resourceVersion, uint64(getResp.Header.Revision)); err != nil {
			return err
		}
		hasMore = getResp.More

		if kvCount == 0 && getResp.More {
			return fmt.Errorf("no results were found, but etcd indicated there were more values remaining")
		}

		// avoid small allocations for the result slice, since this can be called in many
		// different contexts and we don't know how significantly the result will be filtered
		if pred.Empty() {
			growSlice(v, kvCount)
		} else {
			growSlice(v, 2048, kvCount)
		}

		// needs to be large enough to give value in terms of parallelizing work
		// but also needs to be small enough to not waste too much work on early exit
		chunkSize := 1 << 5
		if limit > 0 && int64(chunkSize) > limit {
			chunkSize = int(limit)
		}

		var wg sync.WaitGroup
		workChunks := make([]decodeWork, chunkSize)
		// take items from the response until the bucket is full, filtering as we go
	splitChunk:
		for i, chunk := range splitChunks(getResp.Kvs, chunkSize) {
			i := i
			for j, kv := range chunk {
				j, kv := j, kv

				wg.Add(1)
				go func() {
					defer utilruntime.HandleCrash()

					work := decodeWork{key: kv.Key}

					defer func() {
						// this is safe even without a lock because each go routine only modifies its own index
						workChunks[j] = work

						// free kv early. Long lists can take O(seconds) to decode.
						// this is safe even without a lock because each go routine only modifies its own index
						getResp.Kvs[i*chunkSize+j] = nil
						chunk[j] = nil

						wg.Done()
					}()

					data, _, err := s.transformer.TransformFromStorage(ctx, kv.Value, authenticatedDataString(kv.Key))
					if err != nil {
						work.err = storage.NewInternalErrorf("unable to transform key %q: %v", kv.Key, err)
						return
					}

					work.obj, work.err = shouldAppendListItem(data, uint64(kv.ModRevision), pred, s.codec, s.versioner, newItemFunc)
					if work.err != nil {
						recordDecodeError(s.groupResourceString, string(kv.Key))
						return
					}
					numEvald.Add(1)
				}()
			}
			wg.Wait()

			for _, workItem := range workChunks[:len(chunk)] {
				if paging && int64(v.Len()) >= pred.Limit {
					hasMore = true
					break splitChunk
				}
				lastKey = workItem.key

				if err := workItem.err; err != nil {
					return err
				}

				if workItem.obj == nil {
					continue // predicate filtered the item out
				}

				v.Set(reflect.Append(v, reflect.ValueOf(workItem.obj).Elem()))
			}
		}

		// indicate to the client which resource version was returned
		if returnedRV == 0 {
			returnedRV = getResp.Header.Revision
		}

		// no more results remain or we didn't request paging
		if !hasMore || !paging {
			break
		}
		// we're paging but we have filled our bucket
		if int64(v.Len()) >= pred.Limit {
			break
		}

		if limit < maxLimit {
			// We got incomplete result due to field/label selector dropping the object.
			// Double page size to reduce total number of calls to etcd.
			limit *= 2
			if limit > maxLimit {
				limit = maxLimit
			}
			*limitOption = clientv3.WithLimit(limit)
		}
		preparedKey = string(lastKey) + "\x00"
		if withRev == 0 {
			withRev = returnedRV
			options = append(options, clientv3.WithRev(withRev))
		}
	}
	if v.IsNil() {
		// Ensure that we never return a nil Items pointer in the result for consistency.
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
	}

	// instruct the client to begin querying from immediately after the last key we returned
	// we never return a key that the client wouldn't be allowed to see
	if hasMore {
		// we want to start immediately after the last key
		next, err := storage.EncodeContinue(string(lastKey)+"\x00", keyPrefix, returnedRV)
		if err != nil {
			return err
		}
		var remainingItemCount *int64
		// getResp.Count counts in objects that do not match the pred.
		// Instead of returning inaccurate count for non-empty selectors, we return nil.
		// Only set remainingItemCount if the predicate is empty.
		if utilfeature.DefaultFeatureGate.Enabled(features.RemainingItemCount) {
			if pred.Empty() {
				c := int64(getResp.Count - pred.Limit)
				remainingItemCount = &c
			}
		}
		return s.versioner.UpdateList(listObj, uint64(returnedRV), next, remainingItemCount)
	}

	// no continuation
	return s.versioner.UpdateList(listObj, uint64(returnedRV), "", nil)
}

func splitChunks(kvs []*mvccpb.KeyValue, chunkSize int) [][]*mvccpb.KeyValue {
	kvCount := len(kvs)

	if kvCount == 0 {
		return nil
	}

	if kvCount <= chunkSize {
		return [][]*mvccpb.KeyValue{kvs}
	}

	chunks := kvCount / chunkSize
	if kvCount%chunkSize != 0 {
		chunks++
	}

	// TODO remove this allocation by making an iterator struct
	out := make([][]*mvccpb.KeyValue, chunks)
	start, end := 0, chunkSize
	for i := 0; i < chunks; i++ {
		out[i] = kvs[start:end]
		start += chunkSize
		end += chunkSize
		if end > kvCount {
			end = kvCount
		}
	}
	return out
}

type decodeWork struct {
	key []byte
	err error
	obj runtime.Object // nil if filtered out by the predicate
}

// growSlice takes a slice value and grows its capacity up
// to the maximum of the passed sizes or maxCapacity, whichever
// is smaller. Above maxCapacity decisions about allocation are left
// to the Go runtime on append. This allows a caller to make an
// educated guess about the potential size of the total list while
// still avoiding overly aggressive initial allocation. If sizes
// is empty maxCapacity will be used as the size to grow.
func growSlice(v reflect.Value, maxCapacity int, sizes ...int) {
	cap := v.Cap()
	max := cap
	for _, size := range sizes {
		if size > max {
			max = size
		}
	}
	if len(sizes) == 0 || max > maxCapacity {
		max = maxCapacity
	}
	if max <= cap {
		return
	}
	if v.Len() > 0 {
		extra := reflect.MakeSlice(v.Type(), v.Len(), max)
		reflect.Copy(extra, v)
		v.Set(extra)
	} else {
		extra := reflect.MakeSlice(v.Type(), 0, max)
		v.Set(extra)
	}
}

// Watch implements storage.Interface.Watch.
func (s *store) Watch(ctx context.Context, key string, opts storage.ListOptions) (watch.Interface, error) {
	preparedKey, err := s.prepareKey(key)
	if err != nil {
		return nil, err
	}
	rev, err := s.versioner.ParseResourceVersion(opts.ResourceVersion)
	if err != nil {
		return nil, err
	}
	return s.watcher.Watch(ctx, preparedKey, int64(rev), opts.Recursive, opts.ProgressNotify, s.transformer, opts.Predicate)
}

func (s *store) getState(ctx context.Context, getResp *clientv3.GetResponse, key string, v reflect.Value, ignoreNotFound bool) (*objState, error) {
	state := &objState{
		meta: &storage.ResponseMeta{},
	}

	if u, ok := v.Addr().Interface().(runtime.Unstructured); ok {
		state.obj = u.NewEmptyInstance()
	} else {
		state.obj = reflect.New(v.Type()).Interface().(runtime.Object)
	}

	if len(getResp.Kvs) == 0 {
		if !ignoreNotFound {
			return nil, storage.NewKeyNotFoundError(key, 0)
		}
		if err := runtime.SetZeroValue(state.obj); err != nil {
			return nil, err
		}
	} else {
		data, stale, err := s.transformer.TransformFromStorage(ctx, getResp.Kvs[0].Value, authenticatedDataString(key))
		if err != nil {
			return nil, storage.NewInternalError(err.Error())
		}
		state.rev = getResp.Kvs[0].ModRevision
		state.meta.ResourceVersion = uint64(state.rev)
		state.data = data
		state.stale = stale
		if err := decode(s.codec, s.versioner, state.data, state.obj, state.rev); err != nil {
			recordDecodeError(s.groupResourceString, key)
			return nil, err
		}
	}
	return state, nil
}

func (s *store) getStateFromObject(obj runtime.Object) (*objState, error) {
	state := &objState{
		obj:  obj,
		meta: &storage.ResponseMeta{},
	}

	rv, err := s.versioner.ObjectResourceVersion(obj)
	if err != nil {
		return nil, fmt.Errorf("couldn't get resource version: %v", err)
	}
	state.rev = int64(rv)
	state.meta.ResourceVersion = uint64(state.rev)

	// Compute the serialized form - for that we need to temporarily clean
	// its resource version field (those are not stored in etcd).
	if err := s.versioner.PrepareObjectForStorage(obj); err != nil {
		return nil, fmt.Errorf("PrepareObjectForStorage failed: %v", err)
	}
	state.data, err = runtime.Encode(s.codec, obj)
	if err != nil {
		return nil, err
	}
	if err := s.versioner.UpdateObject(state.obj, uint64(rv)); err != nil {
		klog.Errorf("failed to update object version: %v", err)
	}
	return state, nil
}

func (s *store) updateState(st *objState, userUpdate storage.UpdateFunc) (runtime.Object, uint64, error) {
	ret, ttlPtr, err := userUpdate(st.obj, *st.meta)
	if err != nil {
		return nil, 0, err
	}

	if err := s.versioner.PrepareObjectForStorage(ret); err != nil {
		return nil, 0, fmt.Errorf("PrepareObjectForStorage failed: %v", err)
	}
	var ttl uint64
	if ttlPtr != nil {
		ttl = *ttlPtr
	}
	return ret, ttl, nil
}

// ttlOpts returns client options based on given ttl.
// ttl: if ttl is non-zero, it will attach the key to a lease with ttl of roughly the same length
func (s *store) ttlOpts(ctx context.Context, ttl int64) ([]clientv3.OpOption, error) {
	if ttl == 0 {
		return nil, nil
	}
	id, err := s.leaseManager.GetLease(ctx, ttl)
	if err != nil {
		return nil, err
	}
	return []clientv3.OpOption{clientv3.WithLease(id)}, nil
}

// validateMinimumResourceVersion returns a 'too large resource' version error when the provided minimumResourceVersion is
// greater than the most recent actualRevision available from storage.
func (s *store) validateMinimumResourceVersion(minimumResourceVersion string, actualRevision uint64) error {
	if minimumResourceVersion == "" {
		return nil
	}
	minimumRV, err := s.versioner.ParseResourceVersion(minimumResourceVersion)
	if err != nil {
		return apierrors.NewBadRequest(fmt.Sprintf("invalid resource version: %v", err))
	}
	// Enforce the storage.Interface guarantee that the resource version of the returned data
	// "will be at least 'resourceVersion'".
	if minimumRV > actualRevision {
		return storage.NewTooLargeResourceVersionError(minimumRV, actualRevision, 0)
	}
	return nil
}

func (s *store) prepareKey(key string) (string, error) {
	if key == ".." ||
		strings.HasPrefix(key, "../") ||
		strings.HasSuffix(key, "/..") ||
		strings.Contains(key, "/../") {
		return "", fmt.Errorf("invalid key: %q", key)
	}
	if key == "." ||
		strings.HasPrefix(key, "./") ||
		strings.HasSuffix(key, "/.") ||
		strings.Contains(key, "/./") {
		return "", fmt.Errorf("invalid key: %q", key)
	}
	if key == "" || key == "/" {
		return "", fmt.Errorf("empty key: %q", key)
	}
	// We ensured that pathPrefix ends in '/' in construction, so skip any leading '/' in the key now.
	startIndex := 0
	if key[0] == '/' {
		startIndex = 1
	}
	return s.pathPrefix + key[startIndex:], nil
}

// decode decodes value of bytes into object. It will also set the object resource version to rev.
// On success, objPtr would be set to the object.
func decode(codec runtime.Codec, versioner storage.Versioner, value []byte, objPtr runtime.Object, rev int64) error {
	if _, err := conversion.EnforcePtr(objPtr); err != nil {
		return fmt.Errorf("unable to convert output object to pointer: %v", err)
	}
	_, _, err := codec.Decode(value, nil, objPtr)
	if err != nil {
		return err
	}
	// being unable to set the version does not prevent the object from being extracted
	if err := versioner.UpdateObject(objPtr, uint64(rev)); err != nil {
		klog.Errorf("failed to update object version: %v", err)
	}
	return nil
}

// shouldAppendListItem decodes and appends the object (if it passes filter) to v, which must be a slice.  // TODO fix
func shouldAppendListItem(data []byte, rev uint64, pred storage.SelectionPredicate, codec runtime.Codec, versioner storage.Versioner, newItemFunc func() runtime.Object) (runtime.Object, error) {
	obj, _, err := codec.Decode(data, nil, newItemFunc())
	if err != nil {
		return nil, err
	}
	// being unable to set the version does not prevent the object from being extracted
	if err := versioner.UpdateObject(obj, rev); err != nil {
		klog.Errorf("failed to update object version: %v", err)
	}
	if matched, err := pred.Matches(obj); err == nil && matched {
		return obj, nil
	}
	return nil, nil
}

// recordDecodeError record decode error split by object type.
func recordDecodeError(resource string, key string) {
	metrics.RecordDecodeError(resource)
	klog.V(4).Infof("Decoding %s \"%s\" failed", resource, key)
}

func notFound(key string) clientv3.Cmp {
	return clientv3.Compare(clientv3.ModRevision(key), "=", 0)
}

// getTypeName returns type name of an object for reporting purposes.
func getTypeName(obj interface{}) string {
	return reflect.TypeOf(obj).String()
}
