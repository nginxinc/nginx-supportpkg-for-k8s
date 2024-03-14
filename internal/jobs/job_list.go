package jobs

import (
	"context"
	"encoding/json"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
