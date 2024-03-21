package data_collector

import (
	"fmt"
	helmClient "github.com/mittwald/go-helm-client"
	crdClient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metricsClient "k8s.io/metrics/pkg/client/clientset/versioned"
	"os"
)

type DataCollector struct {
	BaseDir             string
	Namespaces []string
	K8sCoreClientSet    *kubernetes.Clientset
	K8sCrdClientSet     *crdClient.Clientset
	K8sMetricsClientSet *metricsClient.Clientset
	K8sHelmClientSet    helmClient.Client
}

func NewDataCollector(config *rest.Config, namespaces... string) *DataCollector {

	tmpDir, err := os.MkdirTemp("", "kic-diag")

	if err != nil {
		panic(fmt.Sprintf("%s: Unable to create temporary directory.\n", err))
	}

	dc := DataCollector{
		BaseDir: tmpDir,
		Namespaces : namespaces,
	}

	//Initialize clients
	dc.K8sCoreClientSet, _ = kubernetes.NewForConfig(config)
	dc.K8sCrdClientSet, _ = crdClient.NewForConfig(config)
	dc.K8sMetricsClientSet, _ = metricsClient.NewForConfig(config)
	dc.K8sHelmClientSet, _ = helmClient.NewClientFromRestConf(&helmClient.RestConfClientOptions{
		Options:    &helmClient.Options{},
		RestConfig: config,
	})

	return &dc
}

func (c *DataCollector) WrapUp() {
	os.RemoveAll(c.BaseDir)
}
