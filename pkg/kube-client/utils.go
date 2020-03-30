package kube_client

import (
	"math/rand"
	"time"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	daemonsetName  = "daemonset"
	deploymentName = "deployment"
	podName        = "pod"
	namespaceName  = "namespace"
	letterBytes    = "abcdefghijklmnopqrstuvwxyz"
	stringSize     = 9
)

var images = []string{"nginx:latest", "vineeth0297/languages:1.0", "tomcat:latest", "httpd:latest"}
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func randStringBytes() string {
	b := make([]byte, stringSize)
	for i := range b {
		b[i] = letterBytes[seededRand.Intn(len(letterBytes))]
	}
	return string(b)
}

func generateName() string {
	return randStringBytes()
}

func generateImage() string {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(images)
	imageName := images[n]
	return imageName
}

func (ctx *KubeClient) generateNamespace(excludeNamespaces []string) string {
	// This loop makes sure no creation/deletion on resources from kube-system namespace.
	var namespaces []string
	for {
		if len(excludeNamespaces) == 0 {
			namespacesList, err := ctx.Client.CoreV1().Namespaces().List(metav1.ListOptions{})
			if err != nil {
				ctx.Logger.Error("Unable to list namespaces from the cluster", zap.Error(err))
			}
			for _, namespace := range namespacesList.Items {
				namespaces = append(namespaces, namespace.Name)
			}
		} else {
			namespaces = ctx.namespacesForDeletion(excludeNamespaces)
		}
		n := seededRand.Int() % len(namespaces)
		namespace := namespaces[n]
		if namespace != "kube-system" && namespace != "kube-node-lease" && namespace != "kube-public" {
			return namespace
		}
	}
}

func (ctx *KubeClient) namespacesForDeletion(excludeNamespaces []string) []string {
	namespaces, err := ctx.Client.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		ctx.Logger.Error("Unable to list namespaces from the cluster", zap.Error(err))
	}
	var result []string
	namespacesMap := make(map[string]int)

	for _, k := range namespaces.Items {
		namespacesMap[k.Name] = 1
	}

	for _, k := range excludeNamespaces {
		v := namespacesMap[k]
		v++
		namespacesMap[k] = v
	}

	for key, value := range namespacesMap {
		//This check avoids messing up with kubernetes in-built namespaces
		if key != "kube-system" && key != "kube-node-lease" && key != "kube-public" && value != 2 {
			result = append(result, key)
		}
	}
	return result
}

func generateLabels(resourceType, name string) map[string]string {
	return map[string]string{resourceType: name}
}
