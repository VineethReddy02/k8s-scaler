package kube_client

import (
	"k8s.io/api/batch/v1beta1"
	"sync"

	"go.uber.org/zap"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ctx *KubeClient) CreateCronJobs(count, containers int32, namespace string, excludeNamespaces []string) {
	var created int
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	for counter := 0; counter < int(count); counter++ {
		created++
		if namespace == "" {
			namespace = ctx.generateNamespace(excludeNamespaces)
		}

		cronjobsClient := ctx.Client.BatchV1beta1().CronJobs(namespace)
		go func() {
			defer syncer.Done()
			cronjob := generateCronJobSpec(containers)
			_, err := cronjobsClient.Create(cronjob)
			if err != nil {
				created--
				ctx.Logger.Error("Failed to create cron job ", zap.String("name", cronjob.Name), zap.Error(err))
			}
		}()
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully created", zap.Int("CronJobs", created))
}

func generateCronJobSpec(containers int32) *v1beta1.CronJob {
	name := generateName()
	labels := generateLabels(jobName, name)
	cronjob := &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: v1beta1.CronJobSpec{
			Schedule: "*/30 * * * *",
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: v1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers:    generateContainersForJobs(containers, name),
							RestartPolicy: "Never",
						},
					},
				},
			},
		},
	}

	return cronjob
}

func (ctx *KubeClient) DeleteCronJobs(count int32, providedNamespace string, excludeNamespaces []string) {
	var counter int
	var calcNamespace []string
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	cronjobClient := ctx.Client.BatchV1beta1().CronJobs(providedNamespace)
	cronjobList, err := cronjobClient.List(metav1.ListOptions{})
	if err != nil {
		ctx.Logger.Error("Unable to list cron jobs", zap.Error(err))
	}

	if providedNamespace != "" {
		if len(cronjobList.Items) < int(count) {
			panic("The number of cron jobs in the provided namespaces are lesser than the the provided scale to delete")
		}
		for _, d := range cronjobList.Items {
			counter++
			// This is to make sure goroutines get the actual name.
			name := d.Name
			go func() {
				defer syncer.Done()
				err = ctx.Client.BatchV1beta1().CronJobs(providedNamespace).Delete(name, &metav1.DeleteOptions{})
				if err != nil {
					counter--
					ctx.Logger.Error("Failed to delete cron jobs ", zap.String("name", name), zap.Error(err))
				}
			}()
			if counter == int(count) {
				break
			}
		}
	} else {
		calcNamespace = ctx.namespacesForDeletion(excludeNamespaces)
		for _, namespace := range calcNamespace {
			cronjobClient = ctx.Client.BatchV1beta1().CronJobs(namespace)
			cronjobList, err = cronjobClient.List(metav1.ListOptions{})
			if cronjobList != nil && len(cronjobList.Items) > 0 {
				for _, d := range cronjobList.Items {
					// This is to make sure goroutines get the actual name.
					name := d.Name
					counter++
					go func() {
						defer syncer.Done()
						err = cronjobClient.Delete(name, &metav1.DeleteOptions{})
						if err != nil {
							counter--
							ctx.Logger.Error("Failed to delete cron jobs ", zap.String("name", name), zap.Error(err))
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
			ctx.Logger.Error("Unable to delete the cron jobs. As provided scale of cron jobs doesn't exits across the namespaces.", zap.Int("Deleted", counter))
		}
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully deleted ", zap.Int("CronJobs", counter))
}
