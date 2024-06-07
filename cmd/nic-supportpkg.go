/**

Copyright 2024 F5, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

**/

package cmd

import (
	"fmt"
	"github.com/nginxinc/nginx-k8s-supportpkg/pkg/data_collector"
	"github.com/nginxinc/nginx-k8s-supportpkg/pkg/jobs"
	"github.com/nginxinc/nginx-k8s-supportpkg/pkg/version"
	"github.com/spf13/cobra"
	"os"
)

func Execute() {

	var namespaces []string

	var rootCmd = &cobra.Command{
		Use:   "nic-supportpkg",
		Short: "nic-supportpkg - a tool to create Ingress Controller diagnostics package",
		Long:  `nic-supportpkg - a tool to create Ingress Controller diagnostics package`,
		Run: func(cmd *cobra.Command, args []string) {

			collector, err := data_collector.NewDataCollector(namespaces...)
			if err != nil {
				fmt.Println(fmt.Errorf("unable to start data collector: %s", err))
				os.Exit(1)
			}

			collector.Logger.Printf("Starting kubectl-nic-suportpkg - version: %s - build: %s", version.Version, version.Build)

			if collector.AllNamespacesExist() == true {
				for _, job := range jobs.JobList() {
					fmt.Printf("Running job %s...", job.Name)
					err = job.Collect(collector)
					if err != nil {
						fmt.Printf(" Error: %s\n", err)
					} else {
						fmt.Print(" OK\n")
					}
				}

				tarFile, err := collector.WrapUp()
				if err != nil {
					fmt.Println(fmt.Errorf("error when wrapping up: %s", err))
					os.Exit(1)
				} else {
					fmt.Printf("Supportpkg successfully generated: %s\n", tarFile)
				}
			} else {
				fmt.Println(" Error: Some namespaces do not exist")
			}
		},
	}

	rootCmd.Flags().StringSliceVarP(&namespaces, "namespace", "n", []string{}, "list of namespaces to collect information from")
	if err := rootCmd.MarkFlagRequired("namespace"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.SetUsageTemplate("Usage: \n nic supportpkg [-n|--namespace] ns1 [-n|--namespace] ns2 ...\n nic supportpkg [-n|--namespace] ns1,ns2 ...\n")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
