package main

import (
	"fmt"
	"github.com/nginxinc/kubectl-kic-supportpkg/internal/data_collector"
	"github.com/nginxinc/kubectl-kic-supportpkg/internal/jobs"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {

	home := homedir.HomeDir()
	kubeConfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)

	namespaces := []string{"nginx-ingress-0", "observability"}

	collector := data_collector.NewDataCollector(config, namespaces...)

	if err != nil {
		panic(err.Error())
	}

	for _, job := range jobs.JobList() {
		fmt.Printf("Running %s and collecting the output in %s...\n", job.Name, collector.BaseDir)
		job.Collect(collector)
	}

	//collector.WrapUp()

}
