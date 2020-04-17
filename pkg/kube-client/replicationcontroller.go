package kube_client

import (
	"sync"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ctx *KubeClient) CreateReplicationControllers(count, replicas, containers int32, namespace string, excludeNamespaces []string) {
	var created int
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	for counter := 0; counter < int(count); counter++ {
		created++
		if namespace == "" {
			namespace = ctx.generateNamespace(excludeNamespaces)
		}

		rcClient := ctx.Client.CoreV1().ReplicationControllers(namespace)
		go func() {
			defer syncer.Done()
			rc := generateReplicationControllerSpec(containers, replicas)
			_, err := rcClient.Create(rc)
			if err != nil {
				created--
				ctx.Logger.Error("Failed to create replication controller ", zap.String("name", rc.Name), zap.Error(err))
			}
		}()
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully created", zap.Int("ReplicationControllers", created))
}

func generateReplicationControllerSpec(containers, replicas int32) *corev1.ReplicationController {
	name := generateName()
	labels := generateLabels(replicationcontrollerName, name)
	rc := &corev1.ReplicationController{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.ReplicationControllerSpec{
			Replicas: &replicas,
			Selector: labels,
			Template: &corev1.PodTemplateSpec{
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

func (ctx *KubeClient) DeleteReplicationControllers(count int32, providedNamespace string, excludeNamespaces []string) {
	var counter int
	var calcNamespace []string
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	rcClient := ctx.Client.CoreV1().ReplicationControllers(providedNamespace)
	rcList, err := rcClient.List(metav1.ListOptions{})
	if err != nil {
		ctx.Logger.Error("Unable to list replication controllers", zap.Error(err))
	}

	if providedNamespace != "" {
		if len(rcList.Items) < int(count) {
			panic("The number of replication controllers in the provided namespaces are lesser than the the provided scale to delete")
		}
		for _, d := range rcList.Items {
			counter++
			// This is to make sure goroutines get the actual name.
			name := d.Name
			go func() {
				defer syncer.Done()
				err = ctx.Client.CoreV1().ReplicationControllers(providedNamespace).Delete(name, &metav1.DeleteOptions{})
				if err != nil {
					counter--
					ctx.Logger.Error("Failed to delete replication controllers ", zap.String("name", name), zap.Error(err))
				}
			}()
			if counter == int(count) {
				break
			}
		}
	} else {
		calcNamespace = ctx.namespacesForDeletion(excludeNamespaces)
		for _, namespace := range calcNamespace {
			rcClient = ctx.Client.CoreV1().ReplicationControllers(namespace)
			rcList, err = rcClient.List(metav1.ListOptions{})
			if rcList != nil && len(rcList.Items) > 0 {
				for _, d := range rcList.Items {
					// This is to make sure goroutines get the actual name.
					name := d.Name
					counter++
					go func() {
						defer syncer.Done()
						err = rcClient.Delete(name, &metav1.DeleteOptions{})
						if err != nil {
							counter--
							ctx.Logger.Error("Failed to delete replication controllers ", zap.String("name", name), zap.Error(err))
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
			ctx.Logger.Error("Unable to delete the replication controllers. As provided scale of replication controllers doesn't exits across the namespaces.", zap.Int("Deleted", counter))
		}
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully deleted ", zap.Int("Replication Controllers", counter))
}
