package kube_client

import (
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ctx *KubeClient) CreateNamespaces(count int32) {
	namespaceClient := ctx.Client.CoreV1().Namespaces()
	created := 0
	for counter := 0; counter < int(count); counter++ {
		created++
		namespace := generateNamespaceSpec(int32(counter))
		_, err := namespaceClient.Create(namespace)
		if err != nil {
			created--
			ctx.Logger.Error("Failed to create namespace ", zap.String("name", namespace.Name), zap.Error(err))
		}
	}
	ctx.Logger.Info("Successfully created", zap.Int("Namespaces", created))
}

func generateNamespaceSpec(counter int32) *corev1.Namespace {
	name := generateName()
	labels := generateLabels(namespaceName, name)
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}
	return namespace
}

func (ctx *KubeClient) DeleteNamespaces(count int32, excludeNamspaces []string) {
	namespaceClient := ctx.Client.CoreV1().Namespaces()
	namespaces := ctx.namespacesForDeletion(excludeNamspaces)
	var deleted int
	if len(namespaces) < int(count) {
		panic("The provided scale of namespaces doesn't exist to delete.")
	}
	for counter := 0; counter < int(count); counter++ {
		deleted++
		err := namespaceClient.Delete(namespaces[counter], &metav1.DeleteOptions{})
		if err != nil {
			deleted--
			ctx.Logger.Error("Failed to delete namespace ", zap.String("name", namespaces[counter]), zap.Error(err))
		}
	}
	ctx.Logger.Info("Successfully created", zap.Int("Namespaces", deleted))
}
