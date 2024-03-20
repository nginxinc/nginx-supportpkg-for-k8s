package main

import (
	"fmt"
	helmClient "github.com/mittwald/go-helm-client"
	"github.com/nginxinc/kubectl-kic-supportpkg/internal/jobs"
	crdClient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsClient "k8s.io/metrics/pkg/client/clientset/versioned"
	"os"
	"path/filepath"
)

func main() {

	tmpDir, err := os.MkdirTemp("", "kic-diag")
	//defer os.RemoveAll(tmpDir)

	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// Create Kubernetes clientSet
	clientSet, _ := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	for _, job := range jobs.K8sJobList() {
		fmt.Printf("Running %s and collecting the output in %s/%s...\n", job.Name, tmpDir, job.OutputFile)
		job.Collect(tmpDir, clientSet)
	}

	// Create a new clientset for CRDs
	crdClientset, err := crdClient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for _, job := range jobs.K8sCustomJobList() {
		fmt.Printf("Running %s and collecting the output in %s/%s...\n", job.Name, tmpDir, job.OutputFile)
		job.CustomCollect(tmpDir, crdClientset)
	}

	metricsClientset, err := metricsClient.NewForConfig(config)
	for _, job := range jobs.K8sMetricsJobList() {
		fmt.Printf("Running %s and collecting the output in %s/%s...\n", job.Name, tmpDir, job.OutputFile)
		job.MetricsCollect(tmpDir, metricsClientset)

		// Create a Helm client
		opts := &helmClient.RestConfClientOptions{
			Options:    &helmClient.Options{},
			RestConfig: config,
		}
		hc, _ := helmClient.NewClientFromRestConf(opts)
		for _, job := range jobs.HelmJobList() {
			fmt.Printf("Running %s and collecting the output in %s/%s...\n", job.Name, tmpDir, job.OutputFile)
			job.Collect(tmpDir, hc)

		}
	}
}
