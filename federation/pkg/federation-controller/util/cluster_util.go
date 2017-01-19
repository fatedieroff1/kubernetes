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

package util

import (
	"fmt"
	"net"
	"os"
	"time"

	"encoding/json"
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/wait"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	federation_v1beta1 "k8s.io/kubernetes/federation/apis/federation/v1beta1"
	fedclientset "k8s.io/kubernetes/federation/client/clientset_generated/federation_clientset"
	"k8s.io/kubernetes/pkg/api"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

const (
	KubeAPIQPS              = 20.0
	KubeAPIBurst            = 30
	KubeconfigSecretDataKey = "kubeconfig"
	getSecretTimeout        = 1 * time.Minute
)

func BuildClusterConfig(c *federation_v1beta1.Cluster) (*restclient.Config, error) {
	var serverAddress string
	var clusterConfig *restclient.Config
	hostIP, err := utilnet.ChooseHostInterface()
	if err != nil {
		return nil, err
	}

	for _, item := range c.Spec.ServerAddressByClientCIDRs {
		_, cidrnet, err := net.ParseCIDR(item.ClientCIDR)
		if err != nil {
			return nil, err
		}
		myaddr := net.ParseIP(hostIP.String())
		if cidrnet.Contains(myaddr) == true {
			serverAddress = item.ServerAddress
			break
		}
	}
	if serverAddress != "" {
		if c.Spec.SecretRef == nil {
			glog.Infof("didn't find secretRef for cluster %s. Trying insecure access", c.Name)
			clusterConfig, err = clientcmd.BuildConfigFromFlags(serverAddress, "")
		} else {
			kubeconfigGetter := KubeconfigGetterForCluster(c)
			clusterConfig, err = clientcmd.BuildConfigFromKubeconfigGetter(serverAddress, kubeconfigGetter)
		}
		if err != nil {
			return nil, err
		}
		clusterConfig.QPS = KubeAPIQPS
		clusterConfig.Burst = KubeAPIBurst
	}
	return clusterConfig, nil
}

// This is to inject a different kubeconfigGetter in tests.
// We don't use the standard one which calls NewInCluster in tests to avoid having to setup service accounts and mount files with secret tokens.
var KubeconfigGetterForCluster = func(c *federation_v1beta1.Cluster) clientcmd.KubeconfigGetter {
	return func() (*clientcmdapi.Config, error) {
		secretRefName := ""
		if c.Spec.SecretRef != nil {
			secretRefName = c.Spec.SecretRef.Name
		} else {
			glog.Infof("didn't find secretRef for cluster %s. Trying insecure access", c.Name)
		}
		return KubeconfigGetterForSecret(secretRefName)()
	}
}

// KubeconfigGetterForSecret is used to get the kubeconfig from the given secret.
var KubeconfigGetterForSecret = func(secretName string) clientcmd.KubeconfigGetter {
	return func() (*clientcmdapi.Config, error) {
		var data []byte
		if secretName != "" {
			// Get the namespace this is running in from the env variable.
			namespace := os.Getenv("POD_NAMESPACE")
			if namespace == "" {
				return nil, fmt.Errorf("unexpected: POD_NAMESPACE env var returned empty string")
			}
			// Get a client to talk to the k8s apiserver, to fetch secrets from it.
			cc, err := restclient.InClusterConfig()
			if err != nil {
				return nil, fmt.Errorf("error in creating in-cluster client: %s", err)
			}
			client, err := clientset.NewForConfig(cc)
			if err != nil {
				return nil, fmt.Errorf("error in creating in-cluster client: %s", err)
			}
			data = []byte{}
			var secret *api.Secret
			err = wait.PollImmediate(1*time.Second, getSecretTimeout, func() (bool, error) {
				secret, err = client.Core().Secrets(namespace).Get(secretName, metav1.GetOptions{})
				if err == nil {
					return true, nil
				}
				glog.Warningf("error in fetching secret: %s", err)
				return false, nil
			})
			if err != nil {
				return nil, fmt.Errorf("timed out waiting for secret: %s", err)
			}
			if secret == nil {
				return nil, fmt.Errorf("unexpected: received null secret %s", secretName)
			}
			ok := false
			data, ok = secret.Data[KubeconfigSecretDataKey]
			if !ok {
				return nil, fmt.Errorf("secret does not have data with key: %s", KubeconfigSecretDataKey)
			}
		}
		return clientcmd.Load(data)
	}
}

// Returns Clientset for the given cluster.
func GetClientsetForCluster(cluster *federation_v1beta1.Cluster) (*fedclientset.Clientset, error) {
	clusterConfig, err := BuildClusterConfig(cluster)
	if err != nil && clusterConfig != nil {
		clientset := fedclientset.NewForConfigOrDie(restclient.AddUserAgent(clusterConfig, userAgentName))
		return clientset, nil
	}
	return nil, err
}

// Returns a bool to indicate if the object should be forwarded to a cluster
func SendToCluster(cluster *federation_v1beta1.Cluster, objMeta metav1.ObjectMeta) bool {
	forwardObject := true
	// Check if there is an annotation for ClusterSelector
	if val, ok := objMeta.Annotations[federation_v1beta1.FederationClusterSelector]; ok {
		// Check if the Annotation contains valid data
		selector := labels.NewSelector()
		requirements := make([]federation_v1beta1.ClusterSelectorRequirement, 0)
		if err := json.Unmarshal([]byte(val), &requirements); err == nil {
			for _, requirement := range requirements {
				r, err := labels.NewRequirement(requirement.Key, ConvertOperator(requirement.Operator), requirement.Values)
				if err != nil {
					glog.V(8).Infof("Unable to convert ClusterSelector Annotation to Requirement: %+v, %s", requirement, err.Error())
				}
				selector = selector.Add(*r)
			}
			if !selector.Matches(labels.Set(cluster.Labels)) {
				glog.V(8).Infof("Object %s Selector: %+v does not match Cluster %s labels: %+v", objMeta.Name, selector.String(), cluster.ObjectMeta.Name, cluster.Labels)
				forwardObject = false
			}
		} else {
			glog.V(8).Infof("Unable to parse ClusterSelector Annotation: %s", err.Error())
		}
	}
	return forwardObject
}

func ConvertOperator(source string) selection.Operator {
	var op selection.Operator
	switch source {
	case "!", "DoesNotExist":
		op = selection.DoesNotExist
	case "=":
		op = selection.Equals
	case "==":
		op = selection.DoubleEquals
	case "in", "In":
		op = selection.In
	case "!=":
		op = selection.NotEquals
	case "notin", "NotIn":
		op = selection.NotIn
	case "exists", "Exists":
		op = selection.Exists
	case "gt", "Gt":
		op = selection.GreaterThan
	case "lt", "Lt":
		op = selection.LessThan
	}
	return op
}
