package kube_client

import (
	"sync"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ctx *KubeClient) CreateReplicaSet(count, replicas, containers int32, namespace string, excludeNamespaces []string) {
	var created int
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	for counter := 0; counter < int(count); counter++ {
		created++
		if namespace == "" {
			namespace = ctx.generateNamespace(excludeNamespaces)
		}

		rsClient := ctx.Client.AppsV1().ReplicaSets(namespace)
		go func() {
			defer syncer.Done()
			rs := generateReplicaSetSpec(containers, replicas)
			_, err := rsClient.Create(rs)
			if err != nil {
				created--
				ctx.Logger.Error("Failed to create replicaset", zap.String("name", rs.Name), zap.Error(err))
			}
		}()
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully created", zap.Int("ReplicaSets", created))
}

func generateReplicaSetSpec(containers, replicas int32) *appsv1.ReplicaSet {
	name := generateName()
	labels := generateLabels(replicasetName, name)
	rc := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.ReplicaSetSpec{
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

	return rc
}

func (ctx *KubeClient) DeleteReplicaSets(count int32, providedNamespace string, excludeNamespaces []string) {
	var counter int
	var calcNamespace []string
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	rsClient := ctx.Client.AppsV1().ReplicaSets(providedNamespace)
	rsList, err := rsClient.List(metav1.ListOptions{})
	if err != nil {
		ctx.Logger.Error("Unable to list replicaset", zap.Error(err))
	}

	if providedNamespace != "" {
		if len(rsList.Items) < int(count) {
			panic("The number of replicasets in the provided namespaces are lesser than the the provided scale to delete")
		}
		for _, d := range rsList.Items {
			counter++
			// This is to make sure goroutines get the actual name.
			name := d.Name
			go func() {
				defer syncer.Done()
				err = ctx.Client.AppsV1().ReplicaSets(providedNamespace).Delete(name, &metav1.DeleteOptions{})
				if err != nil {
					counter--
					ctx.Logger.Error("Failed to delete replicasets ", zap.String("name", name), zap.Error(err))
				}
			}()
			if counter == int(count) {
				break
			}
		}
	} else {
		calcNamespace = ctx.namespacesForDeletion(excludeNamespaces)
		for _, namespace := range calcNamespace {
			rsClient = ctx.Client.AppsV1().ReplicaSets(namespace)
			rsList, err = rsClient.List(metav1.ListOptions{})
			if rsList != nil && len(rsList.Items) > 0 {
				for _, d := range rsList.Items {
					// This is to make sure goroutines get the actual name.
					name := d.Name
					counter++
					go func() {
						defer syncer.Done()
						err = rsClient.Delete(name, &metav1.DeleteOptions{})
						if err != nil {
							counter--
							ctx.Logger.Error("Failed to delete replicasets ", zap.String("name", name), zap.Error(err))
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
			ctx.Logger.Error("Unable to delete the replicasets. As provided scale of replicasets doesn't exits across the namespaces.", zap.Int("Deleted", counter))
		}
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully deleted ", zap.Int("ReplicaSets", counter))
}
