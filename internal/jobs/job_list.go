package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nginxinc/kubectl-kic-supportpkg/internal/data_collector"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"path"
)

func JobList() []Job {
	jobList := []Job{
		{
			Name:   "pod-list",
			Global: false,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte {
				jobResults := make(map[string][]byte)
				for _, namespace := range dc.Namespaces {
					result, _ := dc.K8sCoreClientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
					jsonResult, _ := json.MarshalIndent(result, "", "  ")
					jobResults[path.Join(dc.BaseDir, namespace, "pods.json")] = jsonResult
				}

				return jobResults
			},
		},
		{
			Name:   "collect-pods-logs",
			Global: false,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte {
				results := make(map[string][]byte)
				for _, namespace := range dc.Namespaces {
					pods, _ := dc.K8sCoreClientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
					for _, pod := range pods.Items {
						for _, container := range pod.Spec.Containers {
							logFileName := path.Join(dc.BaseDir, namespace, "logs", fmt.Sprintf("%s__%s.txt", pod.Name, container.Name))
							bufferedLogs := dc.K8sCoreClientSet.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{Container: container.Name})
							podLogs, err := bufferedLogs.Stream(context.TODO())
							if err != nil {
								log.Fatal("error in opening stream")
							}
							buf := new(bytes.Buffer)
							_, _ = io.Copy(buf, podLogs)
							podLogs.Close()
							results[logFileName] = buf.Bytes()
						}
					}
				}
				return results
			},
		},
		{
			Name:   "configmap-list",
			Global: false,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte {
				jobResults := make(map[string][]byte)
				for _, namespace := range dc.Namespaces {
					result, _ := dc.K8sCoreClientSet.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
					jsonResult, _ := json.MarshalIndent(result, "", "  ")
					jobResults[path.Join(dc.BaseDir, namespace, "configmaps.json")] = jsonResult
				}

				return jobResults
			},
		},
		{
			Name:   "service-list",
			Global: false,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte {
				jobResults := make(map[string][]byte)
				for _, namespace := range dc.Namespaces {
					result, _ := dc.K8sCoreClientSet.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
					jsonResult, _ := json.MarshalIndent(result, "", "  ")
					jobResults[path.Join(dc.BaseDir, namespace, "services.json")] = jsonResult
				}

				return jobResults
			},
		},
		{
			Name:   "deployment-list",
			Global: false,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte {
				jobResults := make(map[string][]byte)
				for _, namespace := range dc.Namespaces {
					result, _ := dc.K8sCoreClientSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
					jsonResult, _ := json.MarshalIndent(result, "", "  ")
					jobResults[path.Join(dc.BaseDir, namespace, "deployments.json")] = jsonResult
				}

				return jobResults
			},
		},
		{
			Name:   "statefulset-list",
			Global: false,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte {
				jobResults := make(map[string][]byte)
				for _, namespace := range dc.Namespaces {
					result, _ := dc.K8sCoreClientSet.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
					jsonResult, _ := json.MarshalIndent(result, "", "  ")
					jobResults[path.Join(dc.BaseDir, namespace, "statefulsets.json")] = jsonResult
				}

				return jobResults
			},
		},
		{
			Name:   "replicaset-list",
			Global: false,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte {
				jobResults := make(map[string][]byte)
				for _, namespace := range dc.Namespaces {
					result, _ := dc.K8sCoreClientSet.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{})
					jsonResult, _ := json.MarshalIndent(result, "", "  ")
					jobResults[path.Join(dc.BaseDir, namespace, "replicasets.json")] = jsonResult
				}

				return jobResults
			},
		},
		{
			Name:   "lease-list",
			Global: false,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte {
				jobResults := make(map[string][]byte)
				for _, namespace := range dc.Namespaces {
					result, _ := dc.K8sCoreClientSet.CoordinationV1().Leases(namespace).List(ctx, metav1.ListOptions{})
					jsonResult, _ := json.MarshalIndent(result, "", "  ")
					jobResults[path.Join(dc.BaseDir, namespace, "leases.json")] = jsonResult
				}

				return jobResults
			},
		},
		{
			Name:   "k8s-version",
			Global: true,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte {
				jobResults := make(map[string][]byte)
				result, _ := dc.K8sCoreClientSet.ServerVersion()
				jsonResult, _ := json.MarshalIndent(result, "", "  ")
				jobResults[path.Join(dc.BaseDir, "k8s", "version.json")] = jsonResult
				return jobResults
			},
		},
		{
			Name:   "crd-info",
			Global: true,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte {
				jobResults := make(map[string][]byte)
				result, _ := dc.K8sCrdClientSet.ApiextensionsV1().CustomResourceDefinitions().List(ctx, metav1.ListOptions{})
				jsonResult, _ := json.MarshalIndent(result, "", "  ")
				jobResults[path.Join(dc.BaseDir, "k8s", "crd.json")] = jsonResult
				return jobResults
			},
		},
		{
			Name:   "nodes-info",
			Global: true,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte {
				jobResults := make(map[string][]byte)
				result, _ := dc.K8sCoreClientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
				jsonResult, _ := json.MarshalIndent(result, "", "  ")
				jobResults[path.Join(dc.BaseDir, "k8s", "nodes.json")] = jsonResult
				return jobResults
			},
		},
		{
			Name:   "metrics-information",
			Global: true,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte {
				jobResults := make(map[string][]byte)
				nodeMetrics, _ := dc.K8sMetricsClientSet.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
				jsonNodeMetrics, _ := json.MarshalIndent(nodeMetrics, "", "  ")
				jobResults[path.Join(dc.BaseDir, "metrics", "node-resource-list.json")] = jsonNodeMetrics
				for _, namespace := range dc.Namespaces {
					podMetrics, _ := dc.K8sMetricsClientSet.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
					jsonPodMetrics, _ := json.MarshalIndent(podMetrics, "", "  ")
					jobResults[path.Join(dc.BaseDir, "metrics", namespace, "pod-resource-list.json")] = jsonPodMetrics
				}
				return jobResults
			},
		},
		{
			Name:   "helm-information",
			Global: true,
			Execute: func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte {
				jobResults := make(map[string][]byte)
				settings := dc.K8sHelmClientSet.GetSettings()
				//release, _ := dc.K8sHelmClientSet.GetRelease("nginx-ingress-0")
				//fmt.Printf(release.Name)
				jsonSettings, _ := json.MarshalIndent(settings, "", "  ")
				jobResults[path.Join(dc.BaseDir, "helm", "settings.json")] = jsonSettings
				releases, err := dc.K8sHelmClientSet.ListDeployedReleases()
				if err != nil {
					fmt.Printf("Error: %s", err)
				}
				jsonReleases, _ := json.MarshalIndent(releases, "", "  ")
				jobResults[path.Join(dc.BaseDir, "helm", "releases.json")] = jsonReleases
				return jobResults
			},
		},
	}
	return jobList
}
