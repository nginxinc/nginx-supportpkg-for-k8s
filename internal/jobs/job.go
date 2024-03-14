package jobs

import (
	"context"
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
