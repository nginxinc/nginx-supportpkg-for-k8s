/**

Copyright 2024 F5, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

**/

package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nginxinc/nginx-k8s-supportpkg/pkg/crds"
	"github.com/nginxinc/nginx-k8s-supportpkg/pkg/data_collector"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path"
	"strings"
	"time"
)

func NICJobList() []Job {
	jobList := []Job{
		{
			Name:    "pod-list",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					result, err := dc.K8sCoreClientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve pod list for namespace %s: %v\n", namespace, err)
					} else {
						jsonResult, _ := json.MarshalIndent(result, "", "  ")
						jobResult.Files[path.Join(dc.BaseDir, "resources", namespace, "pods.json")] = jsonResult
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "collect-pods-logs",
			Timeout: time.Second * 30,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					pods, err := dc.K8sCoreClientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve pod list for namespace %s: %v\n", namespace, err)
					}
					for _, pod := range pods.Items {
						for _, container := range pod.Spec.Containers {
							logFileName := path.Join(dc.BaseDir, "logs", namespace, fmt.Sprintf("%s__%s.txt", pod.Name, container.Name))
							bufferedLogs := dc.K8sCoreClientSet.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{Container: container.Name})
							podLogs, err := bufferedLogs.Stream(context.TODO())
							if err != nil {
								dc.Logger.Printf("\tCould not get logs for pod %s/%s: %v\n", namespace, pod.Name, err)
							} else {
								buf := new(bytes.Buffer)
								_, err := io.Copy(buf, podLogs)
								if err != nil {
									dc.Logger.Printf("\tCould not copy log buffer for pod %s/%s: %v\n", namespace, pod.Name, err)
								} else {
									jobResult.Files[logFileName] = buf.Bytes()
								}
								podLogs.Close()
							}
						}
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "events-list",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					result, err := dc.K8sCoreClientSet.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve events list for namespace %s: %v\n", namespace, err)
					} else {
						jsonResult, _ := json.MarshalIndent(result, "", "  ")
						jobResult.Files[path.Join(dc.BaseDir, "resources", namespace, "events.json")] = jsonResult
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "configmap-list",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					result, err := dc.K8sCoreClientSet.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve configmap list for namespace %s: %v\n", namespace, err)
					} else {
						jsonResult, _ := json.MarshalIndent(result, "", "  ")
						jobResult.Files[path.Join(dc.BaseDir, "resources", namespace, "configmaps.json")] = jsonResult
					}
				}

				ch <- jobResult
			},
		},
		{
			Name:    "service-list",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					result, err := dc.K8sCoreClientSet.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve services list for namespace %s: %v\n", namespace, err)
					} else {
						jsonResult, _ := json.MarshalIndent(result, "", "  ")
						jobResult.Files[path.Join(dc.BaseDir, "resources", namespace, "services.json")] = jsonResult
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "deployment-list",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					result, err := dc.K8sCoreClientSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve deployments list for namespace %s: %v\n", namespace, err)
					} else {
						jsonResult, _ := json.MarshalIndent(result, "", "  ")
						jobResult.Files[path.Join(dc.BaseDir, "resources", namespace, "deployments.json")] = jsonResult
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "statefulset-list",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					result, err := dc.K8sCoreClientSet.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve statefulsets list for namespace %s: %v\n", namespace, err)
					} else {
						jsonResult, _ := json.MarshalIndent(result, "", "  ")
						jobResult.Files[path.Join(dc.BaseDir, "resources", namespace, "statefulsets.json")] = jsonResult
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "daemonsets-list",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					result, err := dc.K8sCoreClientSet.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve daemonsets list for namespace %s: %v\n", namespace, err)
					} else {
						jsonResult, _ := json.MarshalIndent(result, "", "  ")
						jobResult.Files[path.Join(dc.BaseDir, "resources", namespace, "daemonsets.json")] = jsonResult
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "replicaset-list",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					result, err := dc.K8sCoreClientSet.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve replicasets list for namespace %s: %v\n", namespace, err)
					} else {
						jsonResult, _ := json.MarshalIndent(result, "", "  ")
						jobResult.Files[path.Join(dc.BaseDir, "resources", namespace, "replicasets.json")] = jsonResult
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "lease-list",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					result, err := dc.K8sCoreClientSet.CoordinationV1().Leases(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve leases list for namespace %s: %v\n", namespace, err)
					} else {
						jsonResult, _ := json.MarshalIndent(result, "", "  ")
						jobResult.Files[path.Join(dc.BaseDir, "resources", namespace, "leases.json")] = jsonResult
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "roles-list",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					result, err := dc.K8sCoreClientSet.RbacV1().Roles(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve roles list for namespace %s: %v\n", namespace, err)
					} else {
						jsonResult, _ := json.MarshalIndent(result, "", "  ")
						jobResult.Files[path.Join(dc.BaseDir, "k8s", "rbac", namespace, "roles.json")] = jsonResult
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "serviceaccounts-list",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					result, err := dc.K8sCoreClientSet.CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve serviceaccounts list for namespace %s: %v\n", namespace, err)
					} else {
						jsonResult, _ := json.MarshalIndent(result, "", "  ")
						jobResult.Files[path.Join(dc.BaseDir, "k8s", "rbac", namespace, "serviceaccounts.json")] = jsonResult
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "rolebindings-list",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					result, err := dc.K8sCoreClientSet.RbacV1().RoleBindings(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve role bindings list for namespace %s: %v\n", namespace, err)
					} else {
						jsonResult, _ := json.MarshalIndent(result, "", "  ")
						jobResult.Files[path.Join(dc.BaseDir, "k8s", "rbac", namespace, "rolebindings.json")] = jsonResult
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "k8s-version",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				result, err := dc.K8sCoreClientSet.ServerVersion()
				if err != nil {
					dc.Logger.Printf("\tCould not retrieve server version: %v\n", err)
				} else {
					jsonResult, _ := json.MarshalIndent(result, "", "  ")
					jobResult.Files[path.Join(dc.BaseDir, "k8s", "version.json")] = jsonResult
				}
				ch <- jobResult
			},
		},
		{
			Name:    "crd-info",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				result, err := dc.K8sCrdClientSet.ApiextensionsV1().CustomResourceDefinitions().List(ctx, metav1.ListOptions{})
				if err != nil {
					dc.Logger.Printf("\tCould not retrieve crd data: %v\n", err)
				} else {
					jsonResult, _ := json.MarshalIndent(result, "", "  ")
					jobResult.Files[path.Join(dc.BaseDir, "k8s", "crd.json")] = jsonResult
				}
				ch <- jobResult
			},
		},
		{
			Name:    "clusterroles-info",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				result, err := dc.K8sCoreClientSet.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
				if err != nil {
					dc.Logger.Printf("\tCould not retrieve clusterroles data: %v\n", err)
				} else {
					jsonResult, _ := json.MarshalIndent(result, "", "  ")
					jobResult.Files[path.Join(dc.BaseDir, "k8s", "rbac", "clusterroles.json")] = jsonResult
				}
				ch <- jobResult
			},
		},
		{
			Name:    "clusterroles-bindings-info",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				result, err := dc.K8sCoreClientSet.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
				if err != nil {
					dc.Logger.Printf("\tCould not retrieve clusterroles binding data: %v\n", err)
				} else {
					jsonResult, _ := json.MarshalIndent(result, "", "  ")
					jobResult.Files[path.Join(dc.BaseDir, "k8s", "rbac", "clusterrolesbindings.json")] = jsonResult
				}
				ch <- jobResult
			},
		},
		{
			Name:    "nodes-info",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				result, err := dc.K8sCoreClientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
				if err != nil {
					dc.Logger.Printf("\tCould not retrieve nodes information: %v\n", err)
				} else {
					jsonResult, _ := json.MarshalIndent(result, "", "  ")
					jobResult.Files[path.Join(dc.BaseDir, "k8s", "nodes.json")] = jsonResult
				}
				ch <- jobResult
			},
		},
		{
			Name:    "metrics-info",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				nodeMetrics, err := dc.K8sMetricsClientSet.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
				if err != nil {
					dc.Logger.Printf("\tCould not retrieve nodes metrics: %v\n", err)
				} else {
					jsonNodeMetrics, _ := json.MarshalIndent(nodeMetrics, "", "  ")
					jobResult.Files[path.Join(dc.BaseDir, "metrics", "node-resource-list.json")] = jsonNodeMetrics
				}
				for _, namespace := range dc.Namespaces {
					podMetrics, _ := dc.K8sMetricsClientSet.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve pods metrics for namespace %s: %v\n", namespace, err)
					} else {
						jsonPodMetrics, _ := json.MarshalIndent(podMetrics, "", "  ")
						jobResult.Files[path.Join(dc.BaseDir, "metrics", namespace, "pod-resource-list.json")] = jsonPodMetrics
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "helm-info",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				settings := dc.K8sHelmClientSet[dc.Namespaces[0]].GetSettings()
				jsonSettings, err := json.MarshalIndent(settings, "", "  ")
				if err != nil {
					dc.Logger.Printf("\tCould not retrieve helm information: %v\n", err)
				} else {
					jobResult.Files[path.Join(dc.BaseDir, "helm", "settings.json")] = jsonSettings
				}
				ch <- jobResult
			},
		},
		{
			Name:    "helm-deployments",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					releases, err := dc.K8sHelmClientSet[namespace].ListDeployedReleases()
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve helm deployments for namespace %s: %v\n", namespace, err)
					} else {
						for _, release := range releases {
							jsonRelease, _ := json.MarshalIndent(release, "", "  ")
							jobResult.Files[path.Join(dc.BaseDir, "helm", namespace, release.Name+"_release.json")] = jsonRelease
							jobResult.Files[path.Join(dc.BaseDir, "helm", namespace, release.Name+"_manifest.txt")] = []byte(release.Manifest)
						}
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "exec-nginx-t",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				//command := []string{"/bin/sh", "-c", "nginx -T"}
				command := []string{"/usr/sbin/nginx", "-T"}
				for _, namespace := range dc.Namespaces {
					pods, err := dc.K8sCoreClientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
					if err != nil {
						dc.Logger.Printf("\tCould not retrieve pod list for namespace %s: %v\n", namespace, err)
					} else {
						for _, pod := range pods.Items {
							if strings.Contains(pod.Name, "ingress") {
								res, err := dc.PodExecutor(namespace, pod.Name, command, ctx)
								if err != nil {
									dc.Logger.Printf("\tCommand execution %s failed for pod %s in namespace %s: %v\n", command, pod.Name, namespace, err)
								} else {
									jobResult.Files[path.Join(dc.BaseDir, "exec", namespace, pod.Name+"__nginx-t.txt")] = res
								}
							}
						}
					}
				}
				ch <- jobResult
			},
		},
		{
			Name:    "crd-objects",
			Timeout: time.Second * 10,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
				jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
				for _, namespace := range dc.Namespaces {
					for _, crd := range crds.GetCRDList() {
						result, err := dc.QueryCRD(crd, namespace, ctx)
						if err != nil {
							dc.Logger.Printf("\tCRD %s.%s/%s could not be collected in namespace %s: %v\n", crd.Resource, crd.Group, crd.Version, namespace, err)
						} else {
							var jsonResult bytes.Buffer
							_ = json.Indent(&jsonResult, result, "", "  ")
							jobResult.Files[path.Join(dc.BaseDir, "crds", namespace, crd.Resource+".json")] = jsonResult.Bytes()
						}
					}
				}
				ch <- jobResult
			},
		},
	}
	return jobList
}
