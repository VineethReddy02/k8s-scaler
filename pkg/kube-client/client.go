package kube_client

import (
	"os"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {
	Client *kubernetes.Clientset
	Logger *zap.Logger
}

func NewKubeClient() *KubeClient {
	logger, _ := zap.NewProduction()
	return &KubeClient{
		Logger: logger,
	}
}

func (ctx *KubeClient) GetKubeClient(kubeConfig string) *kubernetes.Clientset {
	config := &rest.Config{}
	var err error
	if kubeConfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			ctx.Logger.Error("unable to get provided KUBECONFIG", zap.Error(err))
		}
	} else if kubeConfig == "" {
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		if err != nil {
			ctx.Logger.Error("unable to find KUBECONFIG environment variable")
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic("Unable to find the incluster config")
		}
	}
	ctx.Client, err = kubernetes.NewForConfig(config)
	res, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic("Unable to create NewForConfig")
	}
	return res
}
