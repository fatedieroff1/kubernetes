/*
Copyright 2017 The Kubernetes Authors.

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

package azure

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// VMSet defines functions all vmsets (including scale set and availability
// set) should be implemented.
type VMSet interface {
	// GetInstanceIDByNodeName gets the cloud provider ID by node name.
	// It must return ("", cloudprovider.ErrInstanceNotFound) if the instance does
	// not exist or is no longer running.
	GetInstanceIDByNodeName(name string) (string, error)
	// GetInstanceTypeByNodeName gets the instance type by node name.
	GetInstanceTypeByNodeName(name string) (string, error)
	// GetIPByNodeName gets machine private IP and public IP by node name.
	GetIPByNodeName(name string) (string, string, error)
	// GetPrimaryInterface gets machine primary network interface by node name.
	GetPrimaryInterface(nodeName string) (network.Interface, error)
	// GetNodeNameByProviderID gets the node name by provider ID.
	GetNodeNameByProviderID(providerID string) (types.NodeName, error)

	// GetZoneByNodeName gets cloudprovider.Zone by node name.
	GetZoneByNodeName(name string) (cloudprovider.Zone, error)

	// GetPrimaryVMSetName returns the VM set name depending on the configured vmType.
	// It returns config.PrimaryScaleSetName for vmss and config.PrimaryAvailabilitySetName for standard vmType.
	GetPrimaryVMSetName() string
	// GetVMSetNames selects all possible availability sets or scale sets
	// (depending vmType configured) for service load balancer, if the service has
	// no loadbalancer mode annotation returns the primary VMSet. If service annotation
	// for loadbalancer exists then return the eligible VMSet.
	GetVMSetNames(service *v1.Service, nodes []*v1.Node) (availabilitySetNames *[]string, err error)
	// EnsureHostsInPool ensures the given Node's primary IP configurations are
	// participating in the specified LoadBalancer Backend Pool.
	EnsureHostsInPool(service *v1.Service, nodes []*v1.Node, backendPoolID string, vmSetName string, isInternal bool) error
	// EnsureBackendPoolDeleted ensures the loadBalancer backendAddressPools deleted from the specified vmSet.
	EnsureBackendPoolDeleted(service *v1.Service, poolID, vmSetName string, backendAddressPools *[]network.BackendAddressPool) error

	// AttachDisk attaches a vhd to vm. The vhd must exist, can be identified by diskName, diskURI, and lun.
	AttachDisk(isManagedDisk bool, diskName, diskURI string, nodeName types.NodeName, lun int32, cachingMode compute.CachingTypes) error
	// DetachDiskByName detaches a vhd from host. The vhd can be identified by diskName or diskURI.
	DetachDiskByName(diskName, diskURI string, nodeName types.NodeName) error
	// GetDataDisks gets a list of data disks attached to the node.
	GetDataDisks(nodeName types.NodeName) ([]compute.DataDisk, error)

	// GetPowerStatusByNodeName returns the power state of the specified node.
	GetPowerStatusByNodeName(name string) (string, error)
}
