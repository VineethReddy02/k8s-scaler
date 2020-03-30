package kube_client

import (
	"fmt"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ctx *KubeClient) CreateDeployments(count, replicas int32, namespace string, excludeNamespaces []string) {
	var created int
	for counter := 0; counter < int(count); counter++ {
		created++
		if namespace == "" {
			namespace = ctx.generateNamespace(excludeNamespaces)
		}
		if namespace == "" {
			panic("Couldn't find the namespace for resource creation.")
		}
		deploymentsClient := ctx.Client.AppsV1().Deployments(namespace)
		deployment := generateDeploymentSpec(int32(counter), replicas)
		_, err := deploymentsClient.Create(deployment)
		if err != nil {
			created--
			ctx.Logger.Error("Failed to create deployment ", zap.String("name", deployment.Name), zap.Error(err))
		}
	}
	ctx.Logger.Info("Successfully created", zap.Int("Deployments", created))
}

func generateDeploymentSpec(counter, replicas int32) *appsv1.Deployment {
	name := generateName()
	image := generateImage()
	labels := generateLabels(deploymentName, name)
	containerName := name + fmt.Sprint(counter)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: image,
						Name:  containerName,
						Args:  []string{},
					}},
				},
			},
		},
	}

	return deployment
}

func (ctx *KubeClient) DeleteDeployments(count int32, providedNamespace string, excludeNamespaces []string) {
	var counter int
	var calcNamespace []string
	if providedNamespace == "" {
		calcNamespace = ctx.namespacesForDeletion(excludeNamespaces)
	}

	deploymentClient := ctx.Client.AppsV1().Deployments(providedNamespace)
	deploymentList, err := deploymentClient.List(metav1.ListOptions{})
	if err != nil {
		ctx.Logger.Error("Unable to list deployments", zap.Error(err))
	}

	for _, namespace := range calcNamespace {
		if providedNamespace != "" {
			if len(deploymentList.Items) < int(count) {
				panic("The number of deployments in the provided namespaces are lesser than the the provided scale to delete")
			}
		} else {
			deploymentClient = ctx.Client.AppsV1().Deployments(namespace)
			deploymentList, err = deploymentClient.List(metav1.ListOptions{})
		}

		if deploymentList != nil && len(deploymentList.Items) > 0 {
			for b, _ := range deploymentList.Items {
				err = deploymentClient.Delete(deploymentList.Items[b].Name, &metav1.DeleteOptions{})
				counter++
				if err != nil {
					counter--
					ctx.Logger.Error("Failed to delete deployments ", zap.String("name", deploymentList.Items[b].Name), zap.Error(err))
				}
				if counter == int(count) {
					break
				}
			}
		} else {
			counter--
		}
	}

	ctx.Logger.Info("Successfully deleted ", zap.Int("deployments", counter))
}
