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
	Long: `To create deployments, daemonsets, statefulsets, jobs, cronjobs, pods and you also configure
number of instances using --scale per resource also, replicas if resource can manage replicas and number 
of containers are also configurable. Resources can be created in the desired namespaces and desired namespaces
can also be excluded if creation is performed on a random namespace without specifying the namespace.
`,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		count, _ := cmd.Flags().GetInt32("scale")
		replicas, _ := cmd.Flags().GetInt32("replicas")
		containers, _ := cmd.Flags().GetInt32("containers")
		config, _ := rootCmd.PersistentFlags().GetString("kubeconfig")
		namespace, _ := cmd.Flags().GetString("namespace")
		namespaces, _ := cmd.Flags().GetString("exclude-namespaces")
		excludeNamespaces := strings.Split(namespaces, ",")
		kubeClient := kube_client.NewKubeClient()
		clientInfo := kubeClient.GetKubeClient(config)
		kubeClient.Client = clientInfo
		if resourceType == "deployments" {
			kubeClient.CreateDeployments(count, replicas, containers, namespace, excludeNamespaces)
		} else if resourceType == "pods" {
			kubeClient.CreatePods(count, containers, namespace, excludeNamespaces)
		} else if resourceType == "daemonsets" {
			kubeClient.CreateDaemonsets(count, containers, namespace, excludeNamespaces)
		} else if resourceType == "namespaces" {
			kubeClient.CreateNamespaces(count)
		} else if resourceType == "statefulsets" {
			kubeClient.CreateStatefulsets(count, containers, replicas, namespace, excludeNamespaces)
		} else if resourceType == "jobs" {
			kubeClient.CreateJobs(count, containers, namespace, excludeNamespaces)
		} else if resourceType == "cronjobs" {
			kubeClient.CreateCronJobs(count, containers, namespace, excludeNamespaces)
		} else {
			panic("Invalid resource with create cmd")
		}
	},
	Example: `
# You can provide path to the KUBECONFIG using --kubeconfig flag
# If not provided k8s-scaler reads the KUBECONFIG environment variable
# If KUBECONFIG env is not set tries find InClusterConfig using k8s client-go

# To create deployments in a random namespace
./k8s-scaler create deployments --scale 10 --replicas 3 --containers 15 --containers 15

# To create deployments in a specific namespace
./k8s-scaler create deployments --scale 10 --replicas 3 --containers 15 --namespace k8s-scaler

# To create deployments and exclude some specific namespaces for deployment creation
./k8s-scaler create deployments --scale 10 --replicas 3 --containers 15 --exclude-namespaces namespace01,namespace02

# To create deployments in a random namespace and load provided KUBECONFIG
./k8s-scaler create deployments --scale 10 --replicas 3 --containers 15 --kubeconfig /home/vineeth/gke.yaml

Note: The above provided examples are also applicable for pods.

# To create namespaces
./k8s-scaler create namespaces --scale 10

# To create daemonsets in a random namespace.
./k8s-scaler create daemonsets --scale 10 --containers 15

# To create daemonsets in a random namespace but exclude couple of namespaces.
./k8s-scaler create daemonsets --scale 10 --containers 15 --exclude-namespaces namespace01,namespace02

# To create statefulsets in a random namespace.
./k8s-scaler create statefulsets --scale 10 --containers 15

# To create statefulsets in a random namespace but exclude couple of namespaces.
./k8s-scaler create statefulsets --scale 10 --containers 15 --exclude-namespaces namespace01,namespace02

Note: All the jobs created are by default configured to sleep for 1 minute and move to completed state.

# To create jobs in a random namespace namespaces.
./k8s-scaler create jobs --scale 10 --containers 15

# To create jobs in a random namespace but exclude couple of namespaces.
./k8s-scaler create jobs --scale 10 --containers 15 --exclude-namespaces namespace01,namespace02

Note: All the cron jobs created are by default configured to sleep for 1 minute and to run for every 30 minutes.

# To create cronjobs in a random namespace namespaces.
./k8s-scaler create cronjobs --scale 10 --containers 15

# To create cronjobs in a random namespace but exclude couple of namespaces.
./k8s-scaler create cronjobs --scale 10 --containers 15 --exclude-namespaces namespace01,namespace02

`,
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().Int32P("scale", "s", 1, "number of instances.")
	createCmd.Flags().Int32P("replicas", "r", 1, "number of replicas per instance.")
	createCmd.Flags().Int32P("containers", "c", 1, "number of containers per pod.")
	createCmd.Flags().StringP("namespace", "n", "", "specify the namespace")
	createCmd.Flags().StringP("exclude-namespaces", "e", "", "specify namespaces that needs to be excluded during creation.")
}
