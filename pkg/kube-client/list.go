package kube_client

import (
	"fmt"
	"strconv"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ctx *KubeClient) ListResources() {
	namespaces, err := ctx.Client.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic("Unable to list the namespaces")
	}
	template := "%-18s%-16s%-16s%-16s%-12s%-12s%-12s\n"
	fmt.Printf(template, "NAMESPACE", "DAEMONSETS", "DEPLOYMENTS", "STATEFULSETS", "PODS", "JOBS", "CRONJOBS")
	var syncer sync.WaitGroup
	for _, namespace := range namespaces.Items {
		syncer.Add(6)
		var numberOfPods, numberOfDeployments, numberOfDaemonsets, numberOfStatefulsets, numberOfJobs, numberOfCronJobs int
		go func() {
			defer syncer.Done()
			pods, err := ctx.Client.CoreV1().Pods(namespace.Name).List(metav1.ListOptions{})
			if err != nil {
				panic("Unable to list the pods")
			}
			numberOfPods = len(pods.Items)
		}()

		go func() {
			defer syncer.Done()
			deployments, err := ctx.Client.AppsV1().Deployments(namespace.Name).List(metav1.ListOptions{})
			if err != nil {
				panic("Unable to list the deployments")
			}
			numberOfDeployments = len(deployments.Items)
		}()

		go func() {
			defer syncer.Done()
			daemonsets, err := ctx.Client.AppsV1().DaemonSets(namespace.Name).List(metav1.ListOptions{})
			if err != nil {
				panic("Unable to list the daemomsets")
			}
			numberOfDaemonsets = len(daemonsets.Items)
		}()

		go func() {
			defer syncer.Done()
			statefulsets, err := ctx.Client.AppsV1().StatefulSets(namespace.Name).List(metav1.ListOptions{})
			if err != nil {
				panic("Unable to list the statefulsets")
			}
			numberOfStatefulsets = len(statefulsets.Items)
		}()

		go func() {
			defer syncer.Done()
			jobs, err := ctx.Client.BatchV1().Jobs(namespace.Name).List(metav1.ListOptions{})
			if err != nil {
				panic("Unable to list the jobs")
			}
			numberOfJobs = len(jobs.Items)
		}()

		go func() {
			defer syncer.Done()
			cronjobs, err := ctx.Client.BatchV1beta1().CronJobs(namespace.Name).List(metav1.ListOptions{})
			if err != nil {
				panic("Unable to list the cron jobs")
			}
			numberOfCronJobs = len(cronjobs.Items)
		}()

		syncer.Wait()
		fmt.Printf(template, namespace.Name, strconv.Itoa(numberOfDeployments), strconv.Itoa(numberOfDaemonsets), strconv.Itoa(numberOfStatefulsets), strconv.Itoa(numberOfPods), strconv.Itoa(numberOfJobs), strconv.Itoa(numberOfCronJobs))
	}
}
