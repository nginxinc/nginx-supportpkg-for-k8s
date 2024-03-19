package jobs

import (
	"context"
	"encoding/json"
	crdClient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	metricsClient "k8s.io/metrics/pkg/client/clientset/versioned"
)

func K8sJobList() []Job {
	jobList := []Job{
		{
			Name:       "pod-list",
			OutputFile: "/list/pods.json",
			RetrieveFunction: func(c *kubernetes.Clientset, ctx context.Context) []byte {
				pods, _ := c.CoreV1().Pods("").List(ctx, v1.ListOptions{})
				jsonPods, _ := json.MarshalIndent(pods, "", "  ")
				return jsonPods
			},
		},
		{
			Name:       "configmap-list",
			OutputFile: "/list/configmaps.json",
			RetrieveFunction: func(c *kubernetes.Clientset, ctx context.Context) []byte {
				pods, _ := c.CoreV1().ConfigMaps("").List(ctx, v1.ListOptions{})
				jsonPods, _ := json.MarshalIndent(pods, "", "  ")
				return jsonPods
			},
		},
		{
			Name:       "service-list",
			OutputFile: "/list/services.json",
			RetrieveFunction: func(c *kubernetes.Clientset, ctx context.Context) []byte {
				pods, _ := c.CoreV1().Services("").List(ctx, v1.ListOptions{})
				jsonPods, _ := json.MarshalIndent(pods, "", "  ")
				return jsonPods
			},
		},
		{
			Name:       "deployment-list",
			OutputFile: "/list/deployments.json",
			RetrieveFunction: func(c *kubernetes.Clientset, ctx context.Context) []byte {
				leases, _ := c.AppsV1().Deployments("").List(ctx, v1.ListOptions{})
				jsonLeases, _ := json.MarshalIndent(leases, "", "  ")
				return jsonLeases
			},
		},
		{
			Name:       "statefulset-list",
			OutputFile: "/list/StatefulSets.json",
			RetrieveFunction: func(c *kubernetes.Clientset, ctx context.Context) []byte {
				leases, _ := c.AppsV1().StatefulSets("").List(ctx, v1.ListOptions{})
				jsonLeases, _ := json.MarshalIndent(leases, "", "  ")
				return jsonLeases
			},
		},
		{
			Name:       "server-version",
			OutputFile: "/k8s/server_version.json",
			RetrieveFunction: func(c *kubernetes.Clientset, ctx context.Context) []byte {
				serverVersion, _ := c.ServerVersion()
				jsonServerVersion, _ := json.MarshalIndent(serverVersion, "", "  ")
				return jsonServerVersion
			},
		},
		{
			Name:       "lease-list",
			OutputFile: "/list/leases.json",
			RetrieveFunction: func(c *kubernetes.Clientset, ctx context.Context) []byte {
				leases, _ := c.CoordinationV1().Leases("").List(ctx, v1.ListOptions{})
				jsonLeases, _ := json.MarshalIndent(leases, "", "  ")
				return jsonLeases
			},
		},
	}
	return jobList
}

func K8sCustomJobList() []CustomJob {
	jobList := []CustomJob{
		{
			Name:       "crd-list",
			OutputFile: "/list/crd.json",
			RetrieveFunction: func(c *crdClient.Clientset, ctx context.Context) []byte {
				crds, _ := c.ApiextensionsV1().CustomResourceDefinitions().List(ctx, v1.ListOptions{})
				jsonCrds, _ := json.MarshalIndent(crds, "", "  ")
				return jsonCrds
			},
		},
	}
	return jobList
}

func K8sMetricsJobList() []MetricsJob {
	jobList := []MetricsJob{
		{
			Name:       "pod-resource-list",
			OutputFile: "/list/pod_resource_list.json",
			RetrieveFunction: func(c *metricsClient.Clientset, ctx context.Context) []byte {
				podMetrics, _ := c.MetricsV1beta1().PodMetricses("").List(ctx, v1.ListOptions{})
				jsonPodMetrics, _ := json.MarshalIndent(podMetrics, "", "  ")
				return jsonPodMetrics
			},
		},
		{
			Name:       "node-resource-list",
			OutputFile: "/list/node_resource_list.json",
			RetrieveFunction: func(c *metricsClient.Clientset, ctx context.Context) []byte {
				podMetrics, _ := c.MetricsV1beta1().NodeMetricses().List(ctx, v1.ListOptions{})
				jsonPodMetrics, _ := json.MarshalIndent(podMetrics, "", "  ")
				return jsonPodMetrics
			},
		},
	}
	return jobList
}
