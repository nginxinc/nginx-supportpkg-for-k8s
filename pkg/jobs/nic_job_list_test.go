package jobs

import (
	"context"
	"github.com/nginxinc/nginx-k8s-supportpkg/pkg/data_collector"
	"reflect"
	"testing"
	"time"
)

func TestNICJobList(t *testing.T) {
	tests := []struct {
		name string
		want []Job
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			want: []Job{
				{
					Name:    "pod-list",
					Timeout: time.Second * 10,
					Execute: func(dc *data_collector.DataCollector, ctx context.Context, ch chan JobResult) {
						jobResult := JobResult{Files: make(map[string][]byte), Error: nil}
						ch <- jobResult
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NICJobList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NICJobList() = %v, want %v", got, tt.want)
			}
		})
	}
}
