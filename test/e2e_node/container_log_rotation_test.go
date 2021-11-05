/*
Copyright 2018 The Kubernetes Authors.

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

package e2enode

import (
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	kubelogs "k8s.io/kubernetes/pkg/kubelet/logs"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
	"k8s.io/kubernetes/test/e2e/framework"
	e2eskipper "k8s.io/kubernetes/test/e2e/framework/skipper"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	rotationPollInterval        = 5 * time.Second
	rotationEventuallyTimeout   = 3 * time.Minute
	rotationConsistentlyTimeout = 2 * time.Minute
)

var _ = SIGDescribe("ContainerLogRotation [Slow] [Serial] [Disruptive]", func() {
	f := framework.NewDefaultFramework("container-log-rotation-test")
	ginkgo.Context("when a container generates a lot of log", func() {
		var maxLogFiles int

		ginkgo.BeforeEach(func() {
			if framework.TestContext.ContainerRuntime != kubetypes.RemoteContainerRuntime {
				e2eskipper.Skipf("Skipping ContainerLogRotation test since the container runtime is not remote")
			}

			kubeletConfig, err := getCurrentKubeletConfig()
			if err != nil {
				e2eskipper.Skipf("Skipping as kubelet config is unavailable")
			}

			maxLogFiles = int(kubeletConfig.ContainerLogMaxFiles)
			if maxLogFiles == 0 || kubeletConfig.ContainerLogMaxSize == "" {
				e2eskipper.Skipf("Skipping ContainerLogRotation test as rotation is not configured")
			}
		})

		ginkgo.It("should be rotated and limited to a fixed amount of files", func() {
			ginkgo.By("create log container")
			pod := &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-container-log-rotation",
				},
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
					Containers: []v1.Container{
						{
							Name:  "log-container",
							Image: busyboxImage,
							Command: []string{
								"sh",
								"-c",
								// ~480Kb/s. Log rotation period is 10 seconds.
								"while true; do echo the quick brown fox jumped over the lonely wolf; sleep 0.0001; done;",
							},
						},
					},
				},
			}
			pod = f.PodClient().CreateSync(pod)
			ginkgo.By("get container log path")
			framework.ExpectEqual(len(pod.Status.ContainerStatuses), 1)
			id := kubecontainer.ParseContainerID(pod.Status.ContainerStatuses[0].ContainerID).ID
			r, _, err := getCRIClient()
			framework.ExpectNoError(err)
			status, err := r.ContainerStatus(id)
			framework.ExpectNoError(err)
			logPath := status.GetLogPath()
			ginkgo.By("wait for container log being rotated to max file limit")
			gomega.Eventually(func() (int, error) {
				logs, err := kubelogs.GetAllLogs(logPath)
				if err != nil {
					return 0, err
				}
				return len(logs), nil
			}, rotationEventuallyTimeout, rotationPollInterval).Should(gomega.Equal(maxLogFiles), "should eventually rotate to max file limit")
			ginkgo.By("make sure container log number won't exceed max file limit")
			gomega.Consistently(func() (int, error) {
				logs, err := kubelogs.GetAllLogs(logPath)
				if err != nil {
					return 0, err
				}
				return len(logs), nil
			}, rotationConsistentlyTimeout, rotationPollInterval).Should(gomega.BeNumerically("<=", maxLogFiles), "should never exceed max file limit")
		})
	})
})
