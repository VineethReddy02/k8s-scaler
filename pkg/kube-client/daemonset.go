package kube_client

import (
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ctx *KubeClient) CreateDaemonsets(count int32, namespace string, excludeNamespaces []string) {
	var created int
	for i := 0; i < int(count); i++ {
		created++
		if namespace == "" {
			namespace = ctx.generateNamespace(excludeNamespaces)
		}
		if namespace == "" {
			panic("Couldn't find the namespace for resource creation.")
		}
		daemonsetClient := ctx.Client.AppsV1().DaemonSets(namespace)
		daemonset := generateDaemonsetSpec()
		_, err := daemonsetClient.Create(daemonset)
		if err != nil {
			created--
			ctx.Logger.Error("Failed to create daemonsets ", zap.String("name", daemonset.Name), zap.Error(err))
		}
	}
	ctx.Logger.Info("Successfully created ", zap.Int("daemonsets", created))
}

func generateDaemonsetSpec() *appsv1.DaemonSet {
	name := generateName()
	image := generateImage()
	labels := generateLabels(daemonsetName, name)
	containerName := name
	pod := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   name,
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  containerName,
							Image: image,
						},
					},
				},
			},
		},
	}

	return pod
}

func (ctx *KubeClient) DeleteDaemonsets(count int32, providedNamespace string, excludeNamespaces []string) {
	var counter int
	var calcNamespace []string
	if providedNamespace == "" {
		calcNamespace = ctx.namespacesForDeletion(excludeNamespaces)
	}

	daemonsetClient := ctx.Client.AppsV1().DaemonSets(providedNamespace)
	daemonsetList, err := daemonsetClient.List(metav1.ListOptions{})
	if err != nil {
		ctx.Logger.Error("Unable to list daemonsets", zap.Error(err))
	}

	for _, namespace := range calcNamespace {
		if providedNamespace != "" {
			if len(daemonsetList.Items) < int(count) {
				panic("The number of daemonsets in the provided namespaces are lesser than the the provided scale to delete")
			}
		} else {
			daemonsetClient = ctx.Client.AppsV1().DaemonSets(namespace)
			daemonsetList, err = daemonsetClient.List(metav1.ListOptions{})
		}

		if daemonsetList != nil && len(daemonsetList.Items) > 0 {
			for b, _ := range daemonsetList.Items {
				err = daemonsetClient.Delete(daemonsetList.Items[b].Name, &metav1.DeleteOptions{})
				counter++
				if err != nil {
					counter--
					ctx.Logger.Error("Failed to delete daemonsets ", zap.String("name", daemonsetList.Items[b].Name), zap.Error(err))
				}
				if counter == int(count) {
					break
				}
			}
		} else {
			counter--
		}
	}

	ctx.Logger.Info("Successfully deleted ", zap.Int("daemonsets", counter))
}
