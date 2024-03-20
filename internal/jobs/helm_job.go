package jobs

import (
	"context"
	"encoding/json"
	helmClient "github.com/mittwald/go-helm-client"
	"os"
	"path"
	"time"
)

type HelmJob struct {
	Name             string
	OutputFile       string
	RetrieveFunction func(helmClient.Client, context.Context) []byte
}

func (j HelmJob) Collect(baseFolder string, cs helmClient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := j.RetrieveFunction(cs, ctx)

	fullPathFile := path.Join(baseFolder, j.OutputFile)
	os.MkdirAll(path.Dir(fullPathFile), os.ModePerm)

	file, _ := os.Create(fullPathFile)
	defer file.Close()
	_, _ = file.Write(result)
}

func HelmJobList() []HelmJob {
	jobList := []HelmJob{
		{
			Name:       "helm-settings",
			OutputFile: "/helm/settings.json",
			RetrieveFunction: func(c helmClient.Client, ctx context.Context) []byte {
				settings := c.GetSettings()
				jsonSettings, _ := json.MarshalIndent(settings, "", "  ")
				return jsonSettings
			},
		},
		{
			Name:       "helm-releases",
			OutputFile: "/helm/releases.json",
			RetrieveFunction: func(c helmClient.Client, ctx context.Context) []byte {
				releases, _ := c.ListDeployedReleases()
				jsonSettings, _ := json.MarshalIndent(releases, "", "  ")
				return jsonSettings
			},
		},
	}
	return jobList

}
