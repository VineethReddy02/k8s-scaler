/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	kube_client "github.com/VineethReddy02/k8s-scaler/pkg/kube-client"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "To list namespaces, deployments, pods, daemonsets.",
	Long: `List will list the number of pods, deployments, daemonsets, jobs, cronjobs, statefulsets, replicationcontrollers, replicasets per namespace
across different namespaces`,
	Run: func(cmd *cobra.Command, args []string) {
		config, _ := rootCmd.PersistentFlags().GetString("kubeconfig")
		kubeClient := kube_client.NewKubeClient()
		clientInfo := kubeClient.GetKubeClient(config)
		kubeClient.Client = clientInfo
		kubeClient.ListResources()
	},
	Example: `
# You can provide path to the KUBECONFIG using --kubeconfig flag
# If not provided k8s-scaler reads the KUBECONFIG environment variable
# If KUBECONFIG env is not set tries find InClusterConfig using k8s client-go

# To list namespaces, deployments, pods, daemonsets, jobs, cronjobs, replicationcontrollers, replicasets
./k8s-scaler list
`,
}

func init() {
	rootCmd.AddCommand(listCmd)
}
