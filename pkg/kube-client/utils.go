package kube_client

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"math/rand"
	"time"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	daemonsetName             = "daemonset"
	deploymentName            = "deployment"
	statefulsetName           = "statefulset"
	replicationcontrollerName = "replicationcontroller"
	replicasetName            = "replicaset"
	jobName                   = "job"
	podName                   = "pod"
	namespaceName             = "namespace"
	letterBytes               = "abcdefghijklmnopqrstuvwxyz"
	stringSize                = 12
)

type GlobalConfig struct {
	NodeSelector map[string]string
	Tolerations  []corev1.Toleration
}

var Config GlobalConfig

var images = []string{"nginx:latest", "vineeth0297/languages:1.0", "tomcat:latest", "httpd:latest", "postgres:latest", "cassandra:latest"}
var seededRand = rand.New(
	rand.NewSource(time.Now().Unix()))
var random int

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
	random++
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
		n := seededRand.Intn(len(namespaces))
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

func generateContainers(count int32, name string) (containers []corev1.Container) {
	for i := 0; i < int(count); i++ {
		containerName := name + fmt.Sprint(i)
		image := generateImage()
		container := corev1.Container{
			Name:  containerName,
			Image: image,
			Args:  []string{"sleep", "50000"},
		}
		containers = append(containers, container)
	}
	return containers
}

func generateContainersForJobs(count int32, name string) (containers []corev1.Container) {
	for i := 0; i < int(count); i++ {
		containerName := name + fmt.Sprint(i)
		image := generateImage()
		container := corev1.Container{
			Name:  containerName,
			Image: image,
			Args:  []string{"sleep", "1m"},
		}
		containers = append(containers, container)
	}
	return containers
}

func generateLabels(resourceType, name string) map[string]string {
	return map[string]string{resourceType: name}
}
