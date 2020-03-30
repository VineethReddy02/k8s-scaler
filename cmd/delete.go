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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "To delete deployments/daemonsets/pods/namespaces",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		count, _ := cmd.Flags().GetInt32("scale")
		config, _ := rootCmd.PersistentFlags().GetString("kubeconfig")
		namespace, _ := cmd.Flags().GetString("namespace")
		namespaces, _ := cmd.Flags().GetString("exclude-namespaces")
		excludeNamespaces := strings.Split(namespaces, ",")
		kubeClient := kube_client.NewKubeClient()
		clientInfo := kubeClient.GetKubeClient(config)
		kubeClient.Client = clientInfo
		if resourceType == "deployments" {
			kubeClient.DeleteDeployments(count, namespace, excludeNamespaces)
		} else if resourceType == "pods" {
			kubeClient.DeletePods(count, namespace, excludeNamespaces)
		} else if resourceType == "daemonsets" {
			kubeClient.DeleteDaemonsets(count, namespace, excludeNamespaces)
		} else if resourceType == "namespaces" {
			kubeClient.DeleteNamespaces(count, excludeNamespaces)
		} else {
			panic("Invalid resource with delete cmd")
		}
	},
	Example: `
# You can provide path to the KUBECONFIG using --kubeconfig flag
# If not provided k8s-scaler reads the KUBECONFIG environment variable
# If KUBECONFIG env is not set tries find InClusterConfig using k8s client-go

# To delete deployments randomly across different namespaces
./k8s-scaler delete deployments --scale 10 --replicas 3

# To delete deployments in a specific namespace
./k8s-scaler delete deployments --scale 10 --replicas 3 --namespace k8s-scaler

# To delete deployments and exclude some specific namespaces for deployment creation
./k8s-scaler delete deployments --scale 10 --replicas 3 --exclude-namespaces namespace01,namespace02

# To delete deployments randomly across different namespaces and load provided KUBECONFIG
./k8s-scaler delete deployments --scale 10 --replicas 3 --kubeconfig /home/vineeth/gke.yaml

Note: The above provided examples are also applicable for pods.

# To delete namespaces
./k8s-scaler delete namespaces --scale 10

# To delete daemonsets across multiple namespaces.
./k8s-scaler delete --scale 5 daemonsets

# To delete daemonsets across multiple namespaces but exclude couple of namespaces.
./k8s-scaler delete daemonsets --scale 5 --exclude-namespaces namespace01,namespace02


`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().Int32P("scale", "s", 1, "number of instances.")
	deleteCmd.Flags().StringP("namespace", "n", "", "specify the namespace")
	deleteCmd.Flags().StringP("exclude-namespaces", "e", "", "specify namespaces that needs to be excluded for deletion")
}
