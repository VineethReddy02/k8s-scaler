package kube_client

import "time"

func (ctx *KubeClient) CreateChaosWithDeployments(count, replicas, containers int32, namespace string, excludeNamespaces []string, timeInterval int32) {
	ticker := time.NewTicker(time.Duration(timeInterval) * time.Second)
	for {
		ctx.CreateDeployments(count, replicas, containers, namespace, excludeNamespaces)
		<-ticker.C
		ctx.DeleteDeployments(count, namespace, excludeNamespaces)
		<-ticker.C
	}
}