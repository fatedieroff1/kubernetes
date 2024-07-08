/*
Copyright The Kubernetes Authors.

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

package v1alpha3

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-codegen.sh

// AUTO-GENERATED FUNCTIONS START HERE. DO NOT EDIT.
var map_AllocationResult = map[string]string{
	"":             "AllocationResult contains attributes of an allocated resource.",
	"devices":      "Devices is the result of allocating devices.",
	"nodeSelector": "NodeSelector defines where the allocated resources are available. If unset, they are available everywhere.",
	"controller":   "Controller is the name of the DRA driver which handled the allocation. That driver is also responsible for deallocating the claim. It is empty when the claim can be deallocated without involving a driver.\n\nA driver may allocate devices provided by other drivers, so this driver name here can be different from the driver names listed for the results.\n\nThis is an alpha field and requires enabling the DRAControlPlaneController feature gate.",
}

func (AllocationResult) SwaggerDoc() map[string]string {
	return map_AllocationResult
}

var map_CELDeviceSelector = map[string]string{
	"":           "CELDeviceSelector contains a CEL expression for selecting a device.",
	"expression": "Expression is a CEL expression which evaluates a single device. It must evaluate to true when the device under consideration satisfies the desired criteria, and false when it does not. Any other result is an error and causes allocation of devices to abort.\n\nThe expression's input is an object named \"device\", which carries the following properties:\n - driver (string): the name of the driver which defines this device.\n - attributes (map[string]object): the device's attributes, grouped by prefix\n   (e.g. device.attributes[\"dra.example.com\"] evaluates to an object with all\n   of the attributes which were prefixed by \"dra.example.com\".\n - capacity (map[string]object): the device's capacities, grouped by prefix.\n\nExample: Consider a device with driver=\"dra.example.com\", which exposes two attributes named \"model\" and \"ext.example.com/family\" and which exposes one capacity named \"modules\". This input to this expression would have the following fields:\n\n    device.driver\n    device.attributes[\"dra.example.com\"].model\n    device.attributes[\"ext.example.com\"].family\n    device.capacity[\"dra.example.com\"].modules\n\nThe device.driver field can be used to check for a specific driver, either as a high-level precondition (i.e. you only want to consider devices from this driver) or as part of a multi-clause expression that is meant to consider devices from different drivers.\n\nThe value type of each attribute is defined by the device definition, and users who write these expressions must consult the documentation for their specific drivers. The value type of each capacity is Quantity.\n\nIf an unknown prefix is used as a lookup in either device.attributes or device.capacity, an empty map will be returned. Any reference to an unknown field will cause an evaluation error and allocation to abort.\n\nA robust expression should check for the existence of attributes before referencing them.\n\nFor ease of use, the cel.bind() function is enabled, and can be used to simplify expressions that access multiple attributes with the same domain. For example:\n\n    cel.bind(dra, device.attributes[\"dra.example.com\"], dra.someBool && dra.anotherBool)",
}

func (CELDeviceSelector) SwaggerDoc() map[string]string {
	return map_CELDeviceSelector
}

var map_ClassConfiguration = map[string]string{
	"": "ClassConfiguration is used in DeviceClass.",
}

func (ClassConfiguration) SwaggerDoc() map[string]string {
	return map_ClassConfiguration
}

var map_Device = map[string]string{
	"":           "Device represents one individual hardware instance that can be selected based on its attributes.",
	"name":       "Name is unique identifier among all devices managed by the driver in the pool. It must be a DNS label.",
	"attributes": "Attributes defines the set of attributes for this device. The name of each attribute must be unique in that set.\n\nThe maximum number of attributes and capacities combined is 32.",
	"capacity":   "Capacity defines the set of capacities for this device. The name of each capacity must be unique in that set.\n\nThe maximum number of attributes and capacities combined is 32.",
}

func (Device) SwaggerDoc() map[string]string {
	return map_Device
}

var map_DeviceAllocationConfiguration = map[string]string{
	"":         "DeviceAllocationConfiguration gets embedded in an AllocationResult.",
	"source":   "Source records whether the configuration comes from a class and thus is not something that a normal user would have been able to set or from a claim.",
	"requests": "Requests lists the names of requests where the configuration applies. If empty, its applies to all requests.",
}

func (DeviceAllocationConfiguration) SwaggerDoc() map[string]string {
	return map_DeviceAllocationConfiguration
}

var map_DeviceAllocationResult = map[string]string{
	"":        "DeviceAllocationResult is the result of allocating devices.",
	"results": "Results lists all allocated devices.",
	"config":  "This field is a combination of all the claim and class configuration parameters. Drivers can distinguish between those based on a flag.\n\nThis includes configuration parameters for drivers which have no allocated devices in the result because it is up to the drivers which configuration parameters they support. They can silently ignore unknown configuration parameters.",
}

func (DeviceAllocationResult) SwaggerDoc() map[string]string {
	return map_DeviceAllocationResult
}

var map_DeviceAttribute = map[string]string{
	"":        "DeviceAttribute must have exactly one field set.",
	"int":     "IntValue is a number.",
	"bool":    "BoolValue is a true/false value.",
	"string":  "StringValue is a string. Must not be longer than 64 characters.",
	"version": "VersionValue is a semantic version according to semver.org spec 2.0.0. Must not be longer than 64 characters.",
}

func (DeviceAttribute) SwaggerDoc() map[string]string {
	return map_DeviceAttribute
}

var map_DeviceCapacity = map[string]string{
	"":         "DeviceCapacity must have exactly one field set.",
	"quantity": "Quantity determines the size of the capacity.",
}

func (DeviceCapacity) SwaggerDoc() map[string]string {
	return map_DeviceCapacity
}

var map_DeviceClaim = map[string]string{
	"":            "DeviceClaim defines how to request devices with a ResourceClaim.",
	"requests":    "Requests represent individual requests for distinct devices which must all be satisfied. If empty, nothing needs to be allocated.",
	"constraints": "These constraints must be satisfied by the set of devices that get allocated for the claim.",
	"config":      "This field holds configuration for multiple potential drivers which could satisfy requests in this claim. It is ignored while allocating the claim.",
}

func (DeviceClaim) SwaggerDoc() map[string]string {
	return map_DeviceClaim
}

var map_DeviceClaimConfiguration = map[string]string{
	"":         "DeviceClaimConfiguration is used for configuration parameters in DeviceClaim.",
	"requests": "Requests lists the names of requests where the configuration applies. If empty, it applies to all requests.",
}

func (DeviceClaimConfiguration) SwaggerDoc() map[string]string {
	return map_DeviceClaimConfiguration
}

var map_DeviceClass = map[string]string{
	"":         "DeviceClass is a vendor or admin-provided resource that contains device configuration and selectors. It can be referenced in the device requests of a claim to apply these presets. Cluster scoped.\n\nThis is an alpha type and requires enabling the DynamicResourceAllocation feature gate.",
	"metadata": "Standard object metadata",
	"spec":     "Spec defines what can be allocated and how to configure it.\n\nThis is mutable. Consumers have to be prepared for classes changing at any time, either because they get updated or replaced. Claim allocations are done once based on whatever was set in classes at the time of allocation.\n\nChanging the spec automatically increments the metadata.generation number.",
}

func (DeviceClass) SwaggerDoc() map[string]string {
	return map_DeviceClass
}

var map_DeviceClassList = map[string]string{
	"":         "DeviceClassList is a collection of classes.",
	"metadata": "Standard list metadata",
	"items":    "Items is the list of resource classes.",
}

func (DeviceClassList) SwaggerDoc() map[string]string {
	return map_DeviceClassList
}

var map_DeviceClassSpec = map[string]string{
	"selectors":     "Each selector must be satisfied by a device which is claimed via this class.",
	"config":        "Config defines configuration parameters that apply to each device that is claimed via this class. Some classses may potentially be satisfied by multiple drivers, so each instance of a vendor configuration applies to exactly one driver.\n\nThey are passed to the driver, but are not considered while allocating the claim.",
	"suitableNodes": "Only nodes matching the selector will be considered by the scheduler when trying to find a Node that fits a Pod when that Pod uses a claim that has not been allocated yet *and* that claim gets allocated through a control plane controller. It is ignored when the claim does not use a control plane controller for allocation.\n\nSetting this field is optional. If unset, all Nodes are candidates.\n\nThis is an alpha field and requires enabling the DRAControlPlaneController feature gate.",
}

func (DeviceClassSpec) SwaggerDoc() map[string]string {
	return map_DeviceClassSpec
}

var map_DeviceConfiguration = map[string]string{
	"":       "DeviceConfiguration must have exactly one field set. It gets embedded inline in some other structs which have other fields, so field names must not conflict with those.",
	"opaque": "Opaque provides driver-specific configuration parameters.",
}

func (DeviceConfiguration) SwaggerDoc() map[string]string {
	return map_DeviceConfiguration
}

var map_DeviceConstraint = map[string]string{
	"":               "DeviceConstraint must have exactly one field set besides Requests.",
	"requests":       "Requests is a list of the one or more requests in this claim which must co-satisfy this constraint. If a request is fulfilled by multiple devices, then all of the devices must satisfy the constraint. If this is not specified, this constraint applies to all requests in this claim.",
	"matchAttribute": "MatchAttribute requires that all devices in question have this attribute and that its type and value are the same across those devices.\n\nFor example, if you specified \"dra.example.com/numa\" (a hypothetical example!), then only devices in the same NUMA node will be chosen. A device which does not have that attribute will not be chosen. All devices should use a value of the same type for this attribute because that is part of its specification, but if one device doesn't, then it also will not be chosen.\n\nMust include the domain qualifier.",
}

func (DeviceConstraint) SwaggerDoc() map[string]string {
	return map_DeviceConstraint
}

var map_DeviceRequest = map[string]string{
	"":     "DeviceRequest is a request for devices required for a claim. This is typically a request for a single resource like a device, but can also ask for several identical devices.\n\nA DeviceClassName is currently required. Clients must check that it is indeed set. It's absence indicates that something changed in a way that is not supported by the client yet, in which case it must refuse to handle the request.",
	"name": "Name can be used to reference this request in a pod.spec.containers[].resources.claims entry and in a constraint of the claim.\n\nMust be a DNS label.",
}

func (DeviceRequest) SwaggerDoc() map[string]string {
	return map_DeviceRequest
}

var map_DeviceRequestAllocationResult = map[string]string{
	"":        "DeviceRequestAllocationResult contains the allocation result for one request.",
	"request": "Request is the name of the request in the claim which caused this device to be allocated. Multiple devices may have been allocated per request.",
	"driver":  "Driver specifies the name of the DRA driver whose kubelet plugin should be invoked to process the allocation once the claim is needed on a node.\n\nMust be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver.",
	"pool":    "This name together with the driver name and the device name field identify which device was allocated (`<driver name>/<pool name>/<device name>`).\n\nMust not be longer than 253 characters and may contain one or more DNS sub-domains separated by slashes.",
	"device":  "Device references one device instance via its name in the driver's resource pool. It must be a DNS label.",
}

func (DeviceRequestAllocationResult) SwaggerDoc() map[string]string {
	return map_DeviceRequestAllocationResult
}

var map_DeviceRequestDetails = map[string]string{
	"":                "DeviceRequestDetails is embedded in DeviceRequest and defines one request.",
	"deviceClassName": "DeviceClassName references a specific DeviceClass, which can define additional configuration and selectors to be inherited by this request.\n\nA class is required. Which classes are available depends on the cluster.\n\nAdministrators may use this to restrict which devices may get requested by only installing classes with selectors for permitted devices. If users are free to request anything without restrictions, then administrators can create an empty DeviceClass for users to reference.",
	"selectors":       "Selectors define criteria which must be satisfied by a specific device in order for that device to be considered for this request. All selectors must be satisfied for a device to be considered.",
	"countMode":       "CountMode and its related fields define how many devices are needed to satisfy this request. Supported values are:\n\n- Exact: This request is for a specific number of devices.\n  This is the default. The exact number is provided in the\n  count field.\n\n- All: This request is for all of the matching devices in a pool.\n  Allocation will fail if some devices are already allocated,\n  unless adminAccess is requested.\n\nIf countMode is not specified, the default countMode is Exact. If countMode is Exact and count is not specified, the default count is one. Any other requests must specify this field.\n\nMore modes may get added in the future. Clients must refuse to handle requests with unknown modes.",
	"count":           "Count is used only when the count mode is \"Exact\". Must be greater than zero. If CountMode is Exact and this field is not specified, the default is one.",
	"adminAccess":     "AdminAccess indicates that this is a claim for administrative access to the device(s). Claims with AdminAccess are expected to be used for monitoring or other management services for a device.  They ignore all ordinary claims to the device with respect to access modes and any resource allocations.",
}

func (DeviceRequestDetails) SwaggerDoc() map[string]string {
	return map_DeviceRequestDetails
}

var map_DeviceSelector = map[string]string{
	"":    "DeviceSelector must have exactly one field set.",
	"cel": "CEL contains a CEL expression for selecting a device.",
}

func (DeviceSelector) SwaggerDoc() map[string]string {
	return map_DeviceSelector
}

var map_OpaqueDeviceConfiguration = map[string]string{
	"":           "OpaqueDeviceConfiguration contains configuration parameters for a driver in a format defined by the driver vendor.",
	"driver":     "Driver is used to determine which kubelet plugin needs to be passed these configuration parameters.\n\nAn admission policy provided by the driver developer could use this to decide whether it needs to validate them.\n\nMust be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver.",
	"parameters": "Parameters can contain arbitrary data. It is the responsibility of the driver developer to handle validation and versioning. Typically this includes self-identification and a version (\"kind\" + \"apiVersion\" for Kubernetes types), with conversion between different versions.",
}

func (OpaqueDeviceConfiguration) SwaggerDoc() map[string]string {
	return map_OpaqueDeviceConfiguration
}

var map_PodSchedulingContext = map[string]string{
	"":         "PodSchedulingContext objects hold information that is needed to schedule a Pod with ResourceClaims that use \"WaitForFirstConsumer\" allocation mode.\n\nThis is an alpha type and requires enabling the DRAControlPlaneController feature gate.",
	"metadata": "Standard object metadata",
	"spec":     "Spec describes where resources for the Pod are needed.",
	"status":   "Status describes where resources for the Pod can be allocated.",
}

func (PodSchedulingContext) SwaggerDoc() map[string]string {
	return map_PodSchedulingContext
}

var map_PodSchedulingContextList = map[string]string{
	"":         "PodSchedulingContextList is a collection of Pod scheduling objects.",
	"metadata": "Standard list metadata",
	"items":    "Items is the list of PodSchedulingContext objects.",
}

func (PodSchedulingContextList) SwaggerDoc() map[string]string {
	return map_PodSchedulingContextList
}

var map_PodSchedulingContextSpec = map[string]string{
	"":               "PodSchedulingContextSpec describes where resources for the Pod are needed.",
	"selectedNode":   "SelectedNode is the node for which allocation of ResourceClaims that are referenced by the Pod and that use \"WaitForFirstConsumer\" allocation is to be attempted.",
	"potentialNodes": "PotentialNodes lists nodes where the Pod might be able to run.\n\nThe size of this field is limited to 128. This is large enough for many clusters. Larger clusters may need more attempts to find a node that suits all pending resources. This may get increased in the future, but not reduced.",
}

func (PodSchedulingContextSpec) SwaggerDoc() map[string]string {
	return map_PodSchedulingContextSpec
}

var map_PodSchedulingContextStatus = map[string]string{
	"":               "PodSchedulingContextStatus describes where resources for the Pod can be allocated.",
	"resourceClaims": "ResourceClaims describes resource availability for each pod.spec.resourceClaim entry where the corresponding ResourceClaim uses \"WaitForFirstConsumer\" allocation mode.",
}

func (PodSchedulingContextStatus) SwaggerDoc() map[string]string {
	return map_PodSchedulingContextStatus
}

var map_ResourceClaim = map[string]string{
	"":         "ResourceClaim describes a request for access to resources in the cluster, for use by workloads. For example, if a workload needs an accelerator device with specific properties, this is how that request is expressed. The status stanza tracks whether this claim has been satisfied and what specific resources have been allocated.\n\nThis is an alpha type and requires enabling the DynamicResourceAllocation feature gate.",
	"metadata": "Standard object metadata",
	"spec":     "Spec describes what is being requested and how to configure it. The spec is immutable.",
	"status":   "Status describes whether the claim is ready to use and what has been allocated.",
}

func (ResourceClaim) SwaggerDoc() map[string]string {
	return map_ResourceClaim
}

var map_ResourceClaimConsumerReference = map[string]string{
	"":         "ResourceClaimConsumerReference contains enough information to let you locate the consumer of a ResourceClaim. The user must be a resource in the same namespace as the ResourceClaim.",
	"apiGroup": "APIGroup is the group for the resource being referenced. It is empty for the core API. This matches the group in the APIVersion that is used when creating the resources.",
	"resource": "Resource is the type of resource being referenced, for example \"pods\".",
	"name":     "Name is the name of resource being referenced.",
	"uid":      "UID identifies exactly one incarnation of the resource.",
}

func (ResourceClaimConsumerReference) SwaggerDoc() map[string]string {
	return map_ResourceClaimConsumerReference
}

var map_ResourceClaimList = map[string]string{
	"":         "ResourceClaimList is a collection of claims.",
	"metadata": "Standard list metadata",
	"items":    "Items is the list of resource claims.",
}

func (ResourceClaimList) SwaggerDoc() map[string]string {
	return map_ResourceClaimList
}

var map_ResourceClaimSchedulingStatus = map[string]string{
	"":                "ResourceClaimSchedulingStatus contains information about one particular ResourceClaim with \"WaitForFirstConsumer\" allocation mode.",
	"name":            "Name matches the pod.spec.resourceClaims[*].Name field.",
	"unsuitableNodes": "UnsuitableNodes lists nodes that the ResourceClaim cannot be allocated for.\n\nThe size of this field is limited to 128, the same as for PodSchedulingSpec.PotentialNodes. This may get increased in the future, but not reduced.",
}

func (ResourceClaimSchedulingStatus) SwaggerDoc() map[string]string {
	return map_ResourceClaimSchedulingStatus
}

var map_ResourceClaimSpec = map[string]string{
	"":           "ResourceClaimSpec defines what is being requested in a ResourceClaim and how to configure it.",
	"devices":    "Devices defines how to request devices.",
	"controller": "Controller is the name of the DRA driver that is meant to handle allocation of this claim. If empty, allocation is handled by the scheduler while scheduling a pod.\n\nMust be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver.\n\nThis is an alpha field and requires enabling the DRAControlPlaneController feature gate.",
}

func (ResourceClaimSpec) SwaggerDoc() map[string]string {
	return map_ResourceClaimSpec
}

var map_ResourceClaimStatus = map[string]string{
	"":                      "ResourceClaimStatus tracks whether the resource has been allocated and what the result of that was.",
	"allocation":            "Allocation is set once the claim has been allocated successfully.",
	"reservedFor":           "ReservedFor indicates which entities are currently allowed to use the claim. A Pod which references a ResourceClaim which is not reserved for that Pod will not be started. A claim that is in use or might be in use because it has been reserved must not get deallocated.\n\nIn a cluster with multiple scheduler instances, two pods might get scheduled concurrently by different schedulers. When they reference the same ResourceClaim which already has reached its maximum number of consumers, only one pod can be scheduled.\n\nBoth schedulers try to add their pod to the claim.status.reservedFor field, but only the update that reaches the API server first gets stored. The other one fails with an error and the scheduler which issued it knows that it must put the pod back into the queue, waiting for the ResourceClaim to become usable again.\n\nThere can be at most 32 such reservations. This may get increased in the future, but not reduced.",
	"deallocationRequested": "Indicates that a claim is to be deallocated. While this is set, no new consumers may be added to ReservedFor.\n\nThis is only used if the claim needs to be deallocated by a DRA driver. That driver then must deallocate this claim and reset the field together with clearing the Allocation field.\n\nThis is an alpha field and requires enabling the DRAControlPlaneController feature gate.",
}

func (ResourceClaimStatus) SwaggerDoc() map[string]string {
	return map_ResourceClaimStatus
}

var map_ResourceClaimTemplate = map[string]string{
	"":         "ResourceClaimTemplate is used to produce ResourceClaim objects.\n\nThis is an alpha type and requires enabling the DynamicResourceAllocation feature gate.",
	"metadata": "Standard object metadata",
	"spec":     "Describes the ResourceClaim that is to be generated.\n\nThis field is immutable. A ResourceClaim will get created by the control plane for a Pod when needed and then not get updated anymore.",
}

func (ResourceClaimTemplate) SwaggerDoc() map[string]string {
	return map_ResourceClaimTemplate
}

var map_ResourceClaimTemplateList = map[string]string{
	"":         "ResourceClaimTemplateList is a collection of claim templates.",
	"metadata": "Standard list metadata",
	"items":    "Items is the list of resource claim templates.",
}

func (ResourceClaimTemplateList) SwaggerDoc() map[string]string {
	return map_ResourceClaimTemplateList
}

var map_ResourceClaimTemplateSpec = map[string]string{
	"":         "ResourceClaimTemplateSpec contains the metadata and fields for a ResourceClaim.",
	"metadata": "ObjectMeta may contain labels and annotations that will be copied into the PVC when creating it. No other fields are allowed and will be rejected during validation.",
	"spec":     "Spec for the ResourceClaim. The entire content is copied unchanged into the ResourceClaim that gets created from this template. The same fields as in a ResourceClaim are also valid here.",
}

func (ResourceClaimTemplateSpec) SwaggerDoc() map[string]string {
	return map_ResourceClaimTemplateSpec
}

var map_ResourcePool = map[string]string{
	"":                   "ResourcePool describes the pool that ResourceSlices belong to.",
	"name":               "Name is used to identify the pool. For node-local devices, this is often the node name, but this is not required.\n\nIt must not be longer than 253 characters and must consist of one or more DNS sub-domains separated by slashes.",
	"generation":         "Generation must get increased in all ResourceSlices of a pool whenever resource definitions change. A consumer must only use definitions from ResourceSlices with the highest generation number and ignore all others.\n\nThis mechanism enables consumers to detect pools that are in an incomplete state. It is not sufficient to determine whether the pool is up-to-date. Because that race is unavoidable, drivers may start publishing a pool with zero as initial generation number when there are no slices for that pool with a higher number.",
	"resourceSliceCount": "ResourceSliceCount is the total number of ResourceSlices in the pool at this generation number. Must be higher than zero.\n\nConsumers can use this to check whether they have seen all ResourceSlices belonging to the same pool.",
}

func (ResourcePool) SwaggerDoc() map[string]string {
	return map_ResourcePool
}

var map_ResourceSlice = map[string]string{
	"":         "ResourceSlice represents one or more resources in a pool of similar resources, managed by a given driver. A pool may span more than one ResourceSlice, and exactly how many ResourceSlices comprise a pool is determined by the driver.\n\nAt the moment, the only supported resources are devices with attributes and capacities. Each device in a given pool, regardless of how many ResourceSlices, must have a unique name. The ResourceSlice in which a device gets published may change over time. The unique identifier for a device is the tuple <driver name>, <pool name>, <device name>.\n\nWhenever a driver needs to update a pool, it increments the pool.Spec.Pool.Generation number and updates all ResourceSlices with that new number and new resource definitions. A consumer must only use ResourceSlices with the highest generation number and ignore all others.\n\nWhen allocating all resources in a pool matching certain criteria or when looking for the best solution among several different alternatives, a consumer should check the number of ResourceSlices in a pool (included in each ResourceSlice) to determine whether its view of a pool is complete and if not, should wait until the driver has completed updating the pool.\n\nFor resources that are not local to a node, the node name is not set. Instead, the driver may use a node selector to specify where the devices are available.\n\nThis is an alpha type and requires enabling the DynamicResourceAllocation feature gate.",
	"metadata": "Standard object metadata",
	"spec":     "Contains the information published by the driver.\n\nChanging the spec automatically increments the metadata.generation number.",
}

func (ResourceSlice) SwaggerDoc() map[string]string {
	return map_ResourceSlice
}

var map_ResourceSliceList = map[string]string{
	"":         "ResourceSliceList is a collection of ResourceSlices.",
	"listMeta": "Standard list metadata",
	"items":    "Items is the list of resource ResourceSlices.",
}

func (ResourceSliceList) SwaggerDoc() map[string]string {
	return map_ResourceSliceList
}

var map_ResourceSliceSpec = map[string]string{
	"":             "ResourceSliceSpec contains the information published by the driver in one ResourceSlice.",
	"driver":       "Driver identifies the DRA driver providing the capacity information. A field selector can be used to list only ResourceSlice objects with a certain driver name.\n\nMust be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver.",
	"pool":         "Pool describes the pool that this ResourceSlice belongs to.",
	"nodeName":     "NodeName identifies the node which provides the resources in this pool. A field selector can be used to list only ResourceSlice objects belonging to a certain node.\n\nThis field can be used to limit access from nodes to ResourceSlices with the same node name. It also indicates to autoscalers that adding new nodes of the same type as some old node might also make new resources available.\n\nExactly one of NodeName, NodeSelector and AllNodes must be set.",
	"nodeSelector": "NodeSelector defines which nodes have access to the resources in the pool.\n\nExactly one of NodeName, NodeSelector and AllNodes must be set.",
	"allNodes":     "AllNodes indicates that all nodes have access to the resources in the pool.\n\nExactly one of NodeName, NodeSelector and AllNodes must be set.",
	"devices":      "Devices lists some or all of the devices in this pool.\n\nMust not have more than 128 entries.",
}

func (ResourceSliceSpec) SwaggerDoc() map[string]string {
	return map_ResourceSliceSpec
}

// AUTO-GENERATED FUNCTIONS END HERE
