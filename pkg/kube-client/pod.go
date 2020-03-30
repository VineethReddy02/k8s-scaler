package kube_client

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ctx *KubeClient) CreatePods(count int32, namespace string, excludeNamespaces []string) {
	var created int
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	for counter := 0; counter < int(count); counter++ {
		created++
		if namespace == "" {
			namespace = ctx.generateNamespace(excludeNamespaces)
		}
		podClient := ctx.Client.CoreV1().Pods(namespace)
		go func() {
			defer syncer.Done()
			pod := generatePodSpec(int32(counter))
			_, err := podClient.Create(pod)
			if err != nil {
				created--
				ctx.Logger.Error("Failed to create pod ", zap.String("name", pod.Name), zap.Error(err))
			}
		}()
	}
	syncer.Wait()
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
	var syncer sync.WaitGroup
	syncer.Add(int(count))

	podClient := ctx.Client.CoreV1().Pods(providedNamespace)
	podList, err := podClient.List(metav1.ListOptions{})
	if err != nil {
		ctx.Logger.Error("Unable to list pods", zap.Error(err))
	}

	if providedNamespace != "" {
		if len(podList.Items) < int(count) {
			panic("The number of pods in the provided namespaces are lesser than the the provided scale to delete")
		}
		for _, d := range podList.Items {
			counter++
			// This is to make sure goroutines get the actual name.
			name := d.Name
			go func() {
				defer syncer.Done()
				err = podClient.Delete(name, &metav1.DeleteOptions{})
				if err != nil {
					counter--
					ctx.Logger.Error("Failed to delete pods ", zap.String("name", name), zap.Error(err))
				}
			}()
			if counter == int(count) {
				break
			}
		}
	} else {
		calcNamespace = ctx.namespacesForDeletion(excludeNamespaces)
		for _, namespace := range calcNamespace {
			podClient = ctx.Client.CoreV1().Pods(namespace)
			podList, err = podClient.List(metav1.ListOptions{})

			if podList != nil && len(podList.Items) > 0 {
				for _, d := range podList.Items {
					counter++
					name := d.Name
					go func() {
						defer syncer.Done()
						err = podClient.Delete(name, &metav1.DeleteOptions{})
						if err != nil {
							counter--
							ctx.Logger.Error("Failed to delete pods ", zap.String("name", name), zap.Error(err))
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
			ctx.Logger.Error("Unable to delete the pods. As provided scale of pods doesn't exits across the namespaces.", zap.Int("Deleted", counter))
		}
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully deleted ", zap.Int("pods", counter))
}
