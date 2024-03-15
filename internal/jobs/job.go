package jobs

import (
	"context"
	crdClient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"os"
	"path"
	"time"
)

type Job struct {
	Name             string
	OutputFile       string
	RetrieveFunction func(*kubernetes.Clientset, context.Context) []byte
}

type CustomJob struct {
	Name             string
	OutputFile       string
	RetrieveFunction func(*crdClient.Clientset, context.Context) []byte
}

func (j Job) Collect(baseFolder string, cs *kubernetes.Clientset) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := j.RetrieveFunction(cs, ctx)

	fullPathFile := path.Join(baseFolder, j.OutputFile)
	os.MkdirAll(path.Dir(fullPathFile), os.ModePerm)

	file, _ := os.Create(fullPathFile)
	defer file.Close()
	_, _ = file.Write(result)
}

func (j CustomJob) CustomCollect(baseFolder string, cs *crdClient.Clientset) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := j.RetrieveFunction(cs, ctx)

	fullPathFile := path.Join(baseFolder, j.OutputFile)
	os.MkdirAll(path.Dir(fullPathFile), os.ModePerm)

	file, _ := os.Create(fullPathFile)
	defer file.Close()
	_, _ = file.Write(result)
}
