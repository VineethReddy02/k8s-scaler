package kube_client

import (
	"sync"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ctx *KubeClient) CreateStatefulsets(count, containers, replicas int32, namespace string, excludeNamespaces []string) {
	var created int
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	for counter := 0; counter < int(count); counter++ {
		created++
		if namespace == "" {
			namespace = ctx.generateNamespace(excludeNamespaces)
		}

		statefulsetsClient := ctx.Client.AppsV1().StatefulSets(namespace)
		go func() {
			defer syncer.Done()
			statefulset := generateStatefulSetSpec(containers, replicas)
			_, err := statefulsetsClient.Create(statefulset)
			if err != nil {
				created--
				ctx.Logger.Error("Failed to create statefulset ", zap.String("name", statefulset.Name), zap.Error(err))
			}
		}()
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully created", zap.Int("StatefulSets", created))
}

func generateStatefulSetSpec(containers, replicas int32) *appsv1.StatefulSet {
	name := generateName()
	labels := generateLabels(statefulsetName, name)
	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: appsv1.StatefulSetSpec{
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
					Tolerations: Config.Tolerations,
					NodeSelector:  Config.NodeSelector,
				},
			},
		},
	}

	return statefulset
}

func (ctx *KubeClient) DeleteStatefulSets(count int32, providedNamespace string, excludeNamespaces []string) {
	var counter int
	var calcNamespace []string
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	statefulsetClient := ctx.Client.AppsV1().StatefulSets(providedNamespace)
	statefulsetList, err := statefulsetClient.List(metav1.ListOptions{})
	if err != nil {
		ctx.Logger.Error("Unable to list statefulsets", zap.Error(err))
	}

	if providedNamespace != "" {
		if len(statefulsetList.Items) < int(count) {
			panic("The number of statefulsets in the provided namespaces are lesser than the the provided scale to delete")
		}
		for _, d := range statefulsetList.Items {
			counter++
			// This is to make sure goroutines get the actual name.
			name := d.Name
			go func() {
				defer syncer.Done()
				err = ctx.Client.AppsV1().StatefulSets(providedNamespace).Delete(name, &metav1.DeleteOptions{})
				if err != nil {
					counter--
					ctx.Logger.Error("Failed to delete statefulsets ", zap.String("name", name), zap.Error(err))
				}
			}()
			if counter == int(count) {
				break
			}
		}
	} else {
		calcNamespace = ctx.namespacesForDeletion(excludeNamespaces)
		for _, namespace := range calcNamespace {
			statefulsetClient = ctx.Client.AppsV1().StatefulSets(namespace)
			statefulsetList, err = statefulsetClient.List(metav1.ListOptions{})
			if statefulsetList != nil && len(statefulsetList.Items) > 0 {
				for _, d := range statefulsetList.Items {
					// This is to make sure goroutines get the actual name.
					name := d.Name
					counter++
					go func() {
						defer syncer.Done()
						err = statefulsetClient.Delete(name, &metav1.DeleteOptions{})
						if err != nil {
							counter--
							ctx.Logger.Error("Failed to delete statefulsets ", zap.String("name", name), zap.Error(err))
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
			ctx.Logger.Error("Unable to delete the statefulsets. As provided scale of statefulsets doesn't exits across the namespaces.", zap.Int("Deleted", counter))
		}
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully deleted ", zap.Int("statefulsets", counter))
}
