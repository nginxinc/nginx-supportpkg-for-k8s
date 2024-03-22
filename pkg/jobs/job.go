package jobs

import (
	"context"
	"github.com/nginxinc/kubectl-kic-supportpkg/pkg/data_collector"
	"os"
	"path"
	"time"
)

type Job struct {
	Name    string
	Global  bool
	Execute func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte
	//TODO: execute function must return an error
}

func (j Job) Collect(dc *data_collector.DataCollector) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobResults := j.Execute(dc, ctx)

	for fileName, fileValue := range jobResults {
		err := os.MkdirAll(path.Dir(fileName), os.ModePerm)
		if err != nil {
			return err
		}
		file, _ := os.Create(fileName)
		_, err = file.Write(fileValue)
		if err != nil {
			return err
		}
		_ = file.Close()
	}
	return nil
}
