package jobs

import (
	"context"
	"errors"
	"fmt"
	"github.com/nginxinc/kubectl-kic-supportpkg/pkg/data_collector"
	"os"
	"path"
	"time"
)

type Job struct {
	Name    string
	Global  bool
	Timeout time.Duration
	Execute func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult)
}

type JobResult struct {
	Files map[string][]byte
	Error error
}

func (j Job) Collect(dc *data_collector.DataCollector) error {
	ch := make(chan JobResult, 1)

	ctx, cancel := context.WithTimeout(context.Background(), j.Timeout)
	defer cancel()

	go j.Execute(dc, ctx, ch)

	select {
	case <-ctx.Done():
		return errors.New(fmt.Sprintf("Context cancelled: %v", ctx.Err()))

	case jobResults := <-ch:
		if jobResults.Error != nil {
			return jobResults.Error
		}

		for fileName, fileValue := range jobResults.Files {
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
}