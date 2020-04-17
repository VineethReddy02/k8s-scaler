package kube_client

import (
	"sync"

	"go.uber.org/zap"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ctx *KubeClient) CreateJobs(count, containers int32, namespace string, excludeNamespaces []string) {
	var created int
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	for counter := 0; counter < int(count); counter++ {
		created++
		if namespace == "" {
			namespace = ctx.generateNamespace(excludeNamespaces)
		}

		jobsClient := ctx.Client.BatchV1().Jobs(namespace)
		go func() {
			defer syncer.Done()
			job := generateJobSpec(containers)
			_, err := jobsClient.Create(job)
			if err != nil {
				created--
				ctx.Logger.Error("Failed to create job ", zap.String("name", job.Name), zap.Error(err))
			}
		}()
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully created", zap.Int("Jobs", created))
}

func generateJobSpec(containers int32) *v1.Job {
	name := generateName()
	labels := generateLabels(jobName, name)
	job := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: v1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: generateContainersForJobs(containers, name),
					Tolerations: Config.Tolerations,
					NodeSelector:  Config.NodeSelector,
					RestartPolicy: "Never",
				},
			},
		},
	}

	return job
}

func (ctx *KubeClient) DeleteJobs(count int32, providedNamespace string, excludeNamespaces []string) {
	var counter int
	var calcNamespace []string
	var syncer sync.WaitGroup
	syncer.Add(int(count))
	jobClient := ctx.Client.BatchV1().Jobs(providedNamespace)
	jobList, err := jobClient.List(metav1.ListOptions{})
	if err != nil {
		ctx.Logger.Error("Unable to list jobs", zap.Error(err))
	}

	if providedNamespace != "" {
		if len(jobList.Items) < int(count) {
			panic("The number of jobs in the provided namespaces are lesser than the the provided scale to delete")
		}
		for _, d := range jobList.Items {
			counter++
			// This is to make sure goroutines get the actual name.
			name := d.Name
			go func() {
				defer syncer.Done()
				err = ctx.Client.BatchV1().Jobs(providedNamespace).Delete(name, &metav1.DeleteOptions{})
				if err != nil {
					counter--
					ctx.Logger.Error("Failed to delete jobs ", zap.String("name", name), zap.Error(err))
				}
			}()
			if counter == int(count) {
				break
			}
		}
	} else {
		calcNamespace = ctx.namespacesForDeletion(excludeNamespaces)
		for _, namespace := range calcNamespace {
			jobClient = ctx.Client.BatchV1().Jobs(namespace)
			jobList, err = jobClient.List(metav1.ListOptions{})
			if jobList != nil && len(jobList.Items) > 0 {
				for _, d := range jobList.Items {
					// This is to make sure goroutines get the actual name.
					name := d.Name
					counter++
					go func() {
						defer syncer.Done()
						err = jobClient.Delete(name, &metav1.DeleteOptions{})
						if err != nil {
							counter--
							ctx.Logger.Error("Failed to delete jobs ", zap.String("name", name), zap.Error(err))
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
			ctx.Logger.Error("Unable to delete the jobs. As provided scale of jobs doesn't exits across the namespaces.", zap.Int("Deleted", counter))
		}
	}
	syncer.Wait()
	ctx.Logger.Info("Successfully deleted ", zap.Int("Jobs", counter))
}
