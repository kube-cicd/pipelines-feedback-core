package batchjob

import (
	"context"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract/wiring"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/k8s"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/provider"
	"github.com/pkg/errors"
	v1model "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

type BatchV1JobProvider struct {
	client *v1.BatchV1Client
}

func (bjp *BatchV1JobProvider) InitializeWithContext(sc *wiring.ServiceContext) error {
	client, err := v1.NewForConfig(sc.KubeConfig)
	if err != nil {
		return errors.Wrap(err, "cannot initialize BatchV1JobProvider")
	}
	bjp.client = client
	return nil
}

// ReceivePipelineInfo is tracking batch/v1, kind: Job type objects
func (bjp *BatchV1JobProvider) ReceivePipelineInfo(ctx context.Context, name string, namespace string) (contract.PipelineInfo, error) {
	// find an object
	job, err := bjp.client.Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return contract.PipelineInfo{}, errors.Wrap(err, "cannot fetch batch/v1 Job")
	}

	// validate
	if !k8s.HasUsableAnnotations(job.ObjectMeta) {
		return contract.PipelineInfo{}, errors.New(provider.ErrNotMatched)
	}

	// translate its status
	scm := k8s.CreateSCMContextFromKubernetesAnnotations(job.ObjectMeta)
	jobStatus := translateJobStatus(job)
	pi := contract.NewPipelineInfo(
		scm,
		name,
		job.Status.StartTime.Time,
		jobStatus,
		[]contract.PipelineStage{
			{Name: "Job", Status: jobStatus},
		},
	)

	return *pi, nil
}

// translateJobStatus translates status from batch/v1 Job format to contract.Status
func translateJobStatus(job *v1model.Job) contract.Status {
	if job.Status.Failed > 0 {
		return contract.Failed
	}
	if job.Status.Active > 0 {
		return contract.Running
	}
	if job.Status.Succeeded > 0 {
		return contract.Succeeded
	}
	return contract.Pending
}
