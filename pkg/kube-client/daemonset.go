package kube_client

import (
	"sync"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ctx *KubeClient) CreateDaemonsets(count, containers int32, namespace string, excludeNamespaces []string) {
	var created int
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	for i := 0; i < int(count); i++ {
		created++
		if namespace == "" {
			namespace = ctx.generateNamespace(excludeNamespaces)
		}

		daemonsetClient := ctx.Client.AppsV1().DaemonSets(namespace)
		go func() {
			defer syncer.Done()
			daemonset := generateDaemonsetSpec(containers)
			_, err := daemonsetClient.Create(daemonset)
			if err != nil {
				created--
				ctx.Logger.Error("Failed to create daemonsets ", zap.String("name", daemonset.Name), zap.Error(err))
			}
		}()
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully created ", zap.Int("daemonsets", created))
}

func generateDaemonsetSpec(containers int32) *appsv1.DaemonSet {
	name := generateName()
	labels := generateLabels(daemonsetName, name)
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
					Containers: generateContainers(containers, name),
				},
			},
		},
	}

	return pod
}

func (ctx *KubeClient) DeleteDaemonsets(count int32, providedNamespace string, excludeNamespaces []string) {
	var counter int
	var calcNamespace []string
	var syncer sync.WaitGroup
	syncer.Add(int(count))

	daemonsetClient := ctx.Client.AppsV1().DaemonSets(providedNamespace)
	daemonsetList, err := daemonsetClient.List(metav1.ListOptions{})
	if err != nil {
		ctx.Logger.Error("Unable to list daemonsets", zap.Error(err))
	}

	if providedNamespace != "" {
		if len(daemonsetList.Items) < int(count) {
			panic("The number of daemonsets in the provided namespaces are lesser than the the provided scale to delete")
		}
		for _, d := range daemonsetList.Items {
			// This is to make sure goroutines get the actual name.
			name := d.Name
			go func() {
				defer syncer.Done()
				err = daemonsetClient.Delete(name, &metav1.DeleteOptions{})
				counter++
				if err != nil {
					counter--
					ctx.Logger.Error("Failed to delete daemonsets ", zap.String("name", name), zap.Error(err))
				}
			}()
			if counter == int(count) {
				break
			}
		}
	} else {
		calcNamespace = ctx.namespacesForDeletion(excludeNamespaces)
		for _, namespace := range calcNamespace {
			daemonsetClient = ctx.Client.AppsV1().DaemonSets(namespace)
			daemonsetList, err = daemonsetClient.List(metav1.ListOptions{})
			if daemonsetList != nil && len(daemonsetList.Items) > 0 {
				for _, d := range daemonsetList.Items {
					counter++
					// This is to make sure goroutines get the actual name.
					name := d.Name
					go func() {
						defer syncer.Done()
						err = daemonsetClient.Delete(name, &metav1.DeleteOptions{})
						if err != nil {
							counter--
							ctx.Logger.Error("Failed to delete daemonsets ", zap.String("name", name), zap.Error(err))
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
			ctx.Logger.Error("Unable to delete the daemonsets. As provided scale of daemonsets doesn't exits across the namespaces.", zap.Int("Deleted", counter))
		}
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully deleted ", zap.Int("daemonsets", counter))
}
