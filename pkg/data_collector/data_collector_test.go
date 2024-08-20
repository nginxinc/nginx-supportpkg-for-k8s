package data_collector

import (
	helmClient "github.com/mittwald/go-helm-client"
	crdClient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metricsClient "k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"os"
	"testing"
)

func TestDataCollector_AllNamespacesExist(t *testing.T) {
	type fields struct {
		BaseDir             string
		Namespaces          []string
		Logger              *log.Logger
		LogFile             *os.File
		K8sRestConfig       *rest.Config
		K8sCoreClientSet    *kubernetes.Clientset
		K8sCrdClientSet     *crdClient.Clientset
		K8sMetricsClientSet *metricsClient.Clientset
		K8sHelmClientSet    map[string]helmClient.Client
	}

	var (
		logger     *log.Logger  = log.New(os.Stdout, "TEST: ", log.LstdFlags)
		logFile, _              = os.Create("/path/to/logfile") // Make sure to handle the error in real code.
		restConfig *rest.Config = &rest.Config{ /* ... */ }
		//coreClientSet *kubernetes.Clientset = /* ... */
		//crdClientSet *crdClient.Clientset = crdClient.NewForConfig(config)
		//metricsClientSet *metricsClient.Clientset = /* ... */
		//helmClientSets map[string]helmClient.Client = /* ... */
	)

	//crdClientSet *crdClient.Clientset = crdClient.NewForConfig(config)

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Test Case 1",
			fields: fields{
				BaseDir:       "/path/to/base",
				Namespaces:    []string{"default", "kube-system"},
				Logger:        logger,
				LogFile:       logFile,
				K8sRestConfig: restConfig,
				//K8sCoreClientSet: coreClientSet,
				//K8sCrdClientSet:     crdClientSet,
				//K8sMetricsClientSet: metricsClientSet,
				//K8sHelmClientSet:    helmClientSets,
			},
			want: true,
		},
		// You can add more test cases as needed.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DataCollector{
				BaseDir:             tt.fields.BaseDir,
				Namespaces:          tt.fields.Namespaces,
				Logger:              tt.fields.Logger,
				LogFile:             tt.fields.LogFile,
				K8sRestConfig:       tt.fields.K8sRestConfig,
				K8sCoreClientSet:    tt.fields.K8sCoreClientSet,
				K8sCrdClientSet:     tt.fields.K8sCrdClientSet,
				K8sMetricsClientSet: tt.fields.K8sMetricsClientSet,
				K8sHelmClientSet:    tt.fields.K8sHelmClientSet,
			}
			if got := c.AllNamespacesExist(); got != tt.want {
				t.Errorf("AllNamespacesExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
