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
	"strings"

	kube_client "github.com/VineethReddy02/k8s-scaler/pkg/kube-client"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "To create deployments/daemonsets/pods/namespaces",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		count, _ := cmd.Flags().GetInt32("scale")
		replicas, _ := cmd.Flags().GetInt32("replicas")
		config, _ := rootCmd.PersistentFlags().GetString("kubeconfig")
		namespace, _ := cmd.Flags().GetString("namespace")
		namespaces, _ := cmd.Flags().GetString("exclude-namespaces")
		excludeNamespaces := strings.Split(namespaces, ",")
		kubeClient := kube_client.NewKubeClient()
		clientInfo := kubeClient.GetKubeClient(config)
		kubeClient.Client = clientInfo
		if resourceType == "deployments" {
			kubeClient.CreateDeployments(count, replicas, namespace, excludeNamespaces)
		} else if resourceType == "pods" {
			kubeClient.CreatePods(count, namespace, excludeNamespaces)
		} else if resourceType == "daemonsets" {
			kubeClient.CreateDaemonsets(count, namespace, excludeNamespaces)
		} else if resourceType == "namespaces" {
			kubeClient.CreateNamespaces(count)
		} else {
			panic("Invalid resource with create cmd")
		}
	},
	Example: `
# You can provide path to the KUBECONFIG using --kubeconfig flag
# If not provided k8s-scaler reads the KUBECONFIG environment variable
# If KUBECONFIG env is not set tries find InClusterConfig using k8s client-go

# To create deployments randomly across different namespaces
./k8s-scaler create deployments --scale 10 --replicas 3

# To create deployments in a specific namespace
./k8s-scaler create deployments --scale 10 --replicas 3 --namespace k8s-scaler

# To create deployments and exclude some specific namespaces for deployment creation
./k8s-scaler create deployments --scale 10 --replicas 3 --exclude-namespaces namespace01,namespace02

# To create deployments randomly across different namespaces and load provided KUBECONFIG
./k8s-scaler create deployments --scale 10 --replicas 3 --kubeconfig /home/vineeth/gke.yaml

Note: The above provided examples are also applicable for pods.

# To create namespaces
./k8s-scaler create namespaces --scale 10

# To create daemonsets across multiple namespaces.
./k8s-scaler create daemonsets --scale 10

# To create daemonsets across multiple namespaces but exclude couple of namespaces.
./k8s-scaler create daemonsets --scale 10 --exclude-namespaces namespace01,namespace02

`,
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().Int32P("scale", "s", 1, "number of instances.")
	createCmd.Flags().Int32P("replicas", "r", 1, "number of replicas per instance.")
	createCmd.Flags().StringP("namespace", "n", "", "specify the namespace")
	createCmd.Flags().StringP("exclude-namespaces", "e", "", "specify namespaces that needs to be excluded during creation.")
}
