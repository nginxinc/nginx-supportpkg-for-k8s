package data_collector

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	helmClient "github.com/mittwald/go-helm-client"
	"io"
	crdClient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsClient "k8s.io/metrics/pkg/client/clientset/versioned"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type DataCollector struct {
	BaseDir             string
	Namespaces          []string
	K8sCoreClientSet    *kubernetes.Clientset
	K8sCrdClientSet     *crdClient.Clientset
	K8sMetricsClientSet *metricsClient.Clientset
	K8sHelmClientSet    map[string]helmClient.Client
}

func NewDataCollector(namespaces ...string) (*DataCollector, error) {

	tmpDir, err := os.MkdirTemp("", "kic-diag")

	if err != nil {
		return nil, fmt.Errorf("unable to create temp directory: %s", err)
	}

	// Find config
	home := homedir.HomeDir()
	kubeConfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)

	if err != nil {
		return nil, fmt.Errorf("unable to create k8s config: %s", err)
	}

	dc := DataCollector{
		BaseDir:          tmpDir,
		Namespaces:       namespaces,
		K8sHelmClientSet: make(map[string]helmClient.Client),
	}

	//Initialize clients
	dc.K8sCoreClientSet, _ = kubernetes.NewForConfig(config)
	dc.K8sCrdClientSet, _ = crdClient.NewForConfig(config)
	dc.K8sMetricsClientSet, _ = metricsClient.NewForConfig(config)
	for _, namespace := range dc.Namespaces {
		dc.K8sHelmClientSet[namespace], _ = helmClient.NewClientFromRestConf(&helmClient.RestConfClientOptions{
			Options:    &helmClient.Options{Namespace: namespace},
			RestConfig: config,
		})
	}

	return &dc, nil
}

func (c *DataCollector) WrapUp() error {

	// Create the tarball file
	//wd, _ := os.Getwd()
	unixTime := time.Now().Unix()
	unixTimeString := strconv.FormatInt(unixTime, 10)
	tarballName := fmt.Sprintf("kic-supportpkg-%s.tar.gz", unixTimeString)

	file, err := os.Create(tarballName)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	filepath.Walk(c.BaseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = path
		//TODO: correct this to relative path

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tw, file)
		if err != nil {
			return err
		}

		return nil
	})

	fmt.Println("Tarball created successfully:", tarballName)
	os.RemoveAll(c.BaseDir)
	return nil
}
