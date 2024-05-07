package data_collector

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	helmClient "github.com/mittwald/go-helm-client"
	"io"
	corev1 "k8s.io/api/core/v1"
	crdClient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
	metricsClient "k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type DataCollector struct {
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

func NewDataCollector(namespaces ...string) (*DataCollector, error) {

	tmpDir, err := os.MkdirTemp("", "nic-diag")
	if err != nil {
		return nil, fmt.Errorf("unable to create temp directory: %s", err)
	}

	logFile, err := os.OpenFile(filepath.Join(tmpDir, "nic-supportpkg.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create log file: %s", err)
	}

	// Find config
	kubeConfig := os.Getenv("KUBECONFIG")
	if kubeConfig == "" {
		kubeConfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to k8s using file %s: %s", kubeConfig, err)
	}

	dc := DataCollector{
		BaseDir:          tmpDir,
		Namespaces:       namespaces,
		LogFile:          logFile,
		Logger:           log.New(logFile, "", log.LstdFlags|log.LUTC|log.Lmicroseconds|log.Lshortfile),
		K8sHelmClientSet: make(map[string]helmClient.Client),
	}

	//Initialize clients
	dc.K8sRestConfig = config
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

func (c *DataCollector) WrapUp() (string, error) {

	unixTime := time.Now().Unix()
	unixTimeString := strconv.FormatInt(unixTime, 10)
	tarballName := fmt.Sprintf("nic-supportpkg-%s.tar.gz", unixTimeString)
	tarballRootDirName := fmt.Sprintf("nic-supportpkg-%s", unixTimeString)

	c.LogFile.Close()

	file, err := os.Create(tarballName)
	if err != nil {
		return "", err
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

		relativePath, err := filepath.Rel(c.BaseDir, path)
		if err != nil {
			return err
		}
		header.Name = tarballRootDirName + "/" + relativePath

		if err = tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		file, err = os.Open(path)
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
	_ = os.RemoveAll(c.BaseDir)
	return tarballName, nil
}

func (c *DataCollector) PodExecutor(namespace string, pod string, command []string, ctx context.Context) ([]byte, error) {
	req := c.K8sCoreClientSet.CoreV1().RESTClient().Post().
		Namespace(namespace).
		Resource("pods").
		Name(pod).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: command,
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(c.K8sRestConfig, "POST", req.URL())
	if err != nil {
		return nil, err
	}
	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if stdout.Len() > 0 || stderr.Len() > 0 {
		response := append(stdout.Bytes(), stderr.Bytes()...)
		return response, nil
	} else {
		return nil, err
	}
}

func (c *DataCollector) AllNamespacesExist() bool {
	var allExist bool = true
	for _, namespace := range c.Namespaces {
		_, err := c.K8sCoreClientSet.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
		if err != nil {
			c.Logger.Printf("\t%s: %v\n", namespace, err)
			fmt.Printf("\t%s: %v\n", namespace, err)
			allExist = false
		}
	}

	return allExist
}
