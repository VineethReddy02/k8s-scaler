package kube_client

import (
	"sync"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ctx *KubeClient) CreateDeployments(count, replicas, containers int32, namespace string, excludeNamespaces []string) {
	var created int
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	for counter := 0; counter < int(count); counter++ {
		created++
		if namespace == "" {
			namespace = ctx.generateNamespace(excludeNamespaces)
		}

		deploymentsClient := ctx.Client.AppsV1().Deployments(namespace)
		go func() {
			defer syncer.Done()
			deployment := generateDeploymentSpec(containers, replicas)
			_, err := deploymentsClient.Create(deployment)
			if err != nil {
				created--
				ctx.Logger.Error("Failed to create deployment ", zap.String("name", deployment.Name), zap.Error(err))
			}
		}()
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully created", zap.Int("Deployments", created))
}

func generateDeploymentSpec(containers, replicas int32) *appsv1.Deployment {
	name := generateName()
	labels := generateLabels(deploymentName, name)
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
					Containers: generateContainers(containers, name),
				},
			},
		},
	}

	return deployment
}

func (ctx *KubeClient) DeleteDeployments(count int32, providedNamespace string, excludeNamespaces []string) {
	var counter int
	var calcNamespace []string
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	deploymentClient := ctx.Client.AppsV1().Deployments(providedNamespace)
	deploymentList, err := deploymentClient.List(metav1.ListOptions{})
	if err != nil {
		ctx.Logger.Error("Unable to list deployments", zap.Error(err))
	}

	if providedNamespace != "" {
		if len(deploymentList.Items) < int(count) {
			panic("The number of deployments in the provided namespaces are lesser than the the provided scale to delete")
		}
		for _, d := range deploymentList.Items {
			counter++
			// This is to make sure goroutines get the actual name.
			name := d.Name
			go func() {
				defer syncer.Done()
				err = ctx.Client.AppsV1().Deployments(providedNamespace).Delete(name, &metav1.DeleteOptions{})
				if err != nil {
					counter--
					ctx.Logger.Error("Failed to delete deployments ", zap.String("name", name), zap.Error(err))
				}
			}()
			if counter == int(count) {
				break
			}
		}
	} else {
		calcNamespace = ctx.namespacesForDeletion(excludeNamespaces)
		for _, namespace := range calcNamespace {
			deploymentClient = ctx.Client.AppsV1().Deployments(namespace)
			deploymentList, err = deploymentClient.List(metav1.ListOptions{})
			if deploymentList != nil && len(deploymentList.Items) > 0 {
				for _, d := range deploymentList.Items {
					// This is to make sure goroutines get the actual name.
					name := d.Name
					counter++
					go func() {
						defer syncer.Done()
						err = deploymentClient.Delete(name, &metav1.DeleteOptions{})
						if err != nil {
							counter--
							ctx.Logger.Error("Failed to delete deployments ", zap.String("name", name), zap.Error(err))
						}
					}()
					if counter == int(count) {
						break
					}
				}
			}
			if counter == int(count) {
				break
			}
		}
		if counter != int(count) {
			ctx.Logger.Error("Unable to delete the deployments. As provided scale of deployments doesn't exits across the namespaces.", zap.Int("Deleted", counter))
		}
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully deleted ", zap.Int("deployments", counter))
}
