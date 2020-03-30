package kube_client

import (
	"fmt"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ctx *KubeClient) CreatePods(count int32, namespace string, excludeNamespaces []string) {
	created := 1
	for counter := 0; counter < int(count); counter++ {
		created++
		if namespace == "" {
			namespace = ctx.generateNamespace(excludeNamespaces)
		}
		if namespace == "" {
			panic("Couldn't find the namespace for resource creation.")
		}
		podClient := ctx.Client.CoreV1().Pods(namespace)
		pod := generatePodSpec(int32(counter))
		_, err := podClient.Create(pod)
		if err != nil {
			created--
			ctx.Logger.Error("Failed to create pod ", zap.String("name", pod.Name), zap.Error(err))
		}
	}
	ctx.Logger.Info("Successfully created", zap.Int("Pods", created))
}

func generatePodSpec(counter int32) *corev1.Pod {
	name := generateName()
	image := generateImage()
	labels := generateLabels(podName, name)
	containerName := name + fmt.Sprint(counter)
	pod := &corev1.Pod{
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
	}
	return pod
}

func (ctx *KubeClient) DeletePods(count int32, providedNamespace string, excludeNamespaces []string) {
	var counter int
	var calcNamespace []string
	if providedNamespace == "" {
		calcNamespace = ctx.namespacesForDeletion(excludeNamespaces)
	}

	podClient := ctx.Client.CoreV1().Pods(providedNamespace)
	podList, err := podClient.List(metav1.ListOptions{})
	if err != nil {
		ctx.Logger.Error("Unable to list pods", zap.Error(err))
	}

	for _, namespace := range calcNamespace {
		if providedNamespace != "" {
			if len(podList.Items) < int(count) {
				panic("The number of pods in the provided namespaces are lesser than the the provided scale to delete")
			}
		} else {
			podClient = ctx.Client.CoreV1().Pods(namespace)
			podList, err = podClient.List(metav1.ListOptions{})
		}

		if podList != nil && len(podList.Items) > 0 {
			for b, _ := range podList.Items {
				err = podClient.Delete(podList.Items[b].Name, &metav1.DeleteOptions{})
				counter++
				if err != nil {
					counter--
					ctx.Logger.Error("Failed to delete pods ", zap.String("name", podList.Items[b].Name), zap.Error(err))
				}
				if counter == int(count) {
					break
				}
			}
		} else {
			counter--
		}
	}

	ctx.Logger.Info("Successfully deleted ", zap.Int("pods", counter))
}
