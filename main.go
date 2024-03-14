package main

import (
	"fmt"
	"github.com/nginxinc/kubectl-kic-supportpkg/internal/jobs"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
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
}
