package jobs

import (
	"context"
	"github.com/nginxinc/kubectl-kic-supportpkg/internal/data_collector"
	"os"
	"path"
	"time"
)

type Job struct {
	Name    string
	Global  bool
	Execute func(dc *data_collector.DataCollector, ctx context.Context) map[string][]byte
}

func (j Job) Collect(dc *data_collector.DataCollector) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobResults := j.Execute(dc, ctx)

	for fileName, fileValue := range jobResults {
		os.MkdirAll(path.Dir(fileName), os.ModePerm)
		file, _ := os.Create(fileName)
		_, _ = file.Write(fileValue)
		file.Close()
	}
}
