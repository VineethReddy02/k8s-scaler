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
	v1 "k8s.io/api/core/v1"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// chaosCmd represents the chaos command
var chaosCmd = &cobra.Command{
	Use:   "chaos",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		count, _ := cmd.Flags().GetInt32("scale")
		replicas, _ := cmd.Flags().GetInt32("replicas")
		containers, _ := cmd.Flags().GetInt32("containers")
		config, _ := rootCmd.PersistentFlags().GetString("kubeconfig")
		namespace, _ := cmd.Flags().GetString("namespace")
		namespaces, _ := cmd.Flags().GetString("exclude-namespaces")
		excludeNamespaces := strings.Split(namespaces, ",")
		nodeselector, _ := cmd.Flags().GetString("node-selector")
		toleration, _ := cmd.Flags().GetString("toleration")
		timeInterval, _ := cmd.Flags().GetInt32("time")
		if nodeselector != "" {
			ns := strings.Split(nodeselector, "=")
			kube_client.Config.NodeSelector = map[string]string{ns[0]:ns[1]}
		} else {
			kube_client.Config.NodeSelector = nil
		}
		if toleration != "" {
			ts := strings.Split(toleration, "=")
			toleration := &v1.Toleration{
				Key:               ts[0],
				Operator:          "Equal",
				Value:             ts[1],
				Effect:            "NoSchedule",
			}
			kube_client.Config.Tolerations = append(kube_client.Config.Tolerations, *toleration)
		}else {
			kube_client.Config.Tolerations = nil
		}

		kubeClient := kube_client.NewKubeClient()
		clientInfo := kubeClient.GetKubeClient(config)
		kubeClient.Client = clientInfo
		if resourceType == "deployments" ||  resourceType == "d" {
			kubeClient.CreateChaosWithDeployments(count, replicas, containers, namespace, excludeNamespaces, timeInterval)
		} else {
			log.Fatal("k8s-scaler only creating chaos with deployments")
		}
	},
}

func init() {
	rootCmd.AddCommand(chaosCmd)

	chaosCmd.Flags().Int32P("scale", "s", 1, "number of instances.")
	chaosCmd.Flags().Int32P("replicas", "r", 1, "number of replicas per instance.")
	chaosCmd.Flags().Int32P("containers", "c", 1, "number of containers per pod.")
	chaosCmd.Flags().StringP("namespace", "n", "", "specify the namespace")
	chaosCmd.Flags().String("node-selector", "", "specify the node selector as map key=value")
	chaosCmd.Flags().StringP("toleration", "t", "", "specify the toleration for a specific node as map key=value")
	chaosCmd.Flags().StringP("exclude-namespaces", "e", "", "specify namespaces that needs to be excluded during creation.")
	chaosCmd.Flags().Int32P("time", "", 300, "time interval between creation & deletion, default to 5 minutes.")
}
