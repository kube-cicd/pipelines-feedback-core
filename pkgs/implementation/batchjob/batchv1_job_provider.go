package batchjob

import (
	"context"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract/wiring"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/k8s"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/provider"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/store"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1model "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	"time"
)

type BatchV1JobProvider struct {
	client *v1.BatchV1Client
	store  *store.Operator
	logger *logrus.Entry
}

func (bjp *BatchV1JobProvider) InitializeWithContext(sc *wiring.ServiceContext) error {
	client, err := v1.NewForConfig(sc.KubeConfig)
	if err != nil {
		return errors.Wrap(err, "cannot initialize BatchV1JobProvider")
	}
	bjp.client = client
	bjp.store = sc.Store
	bjp.logger = sc.Log
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
	if ok, err := k8s.HasUsableAnnotations(job.ObjectMeta); !ok {
		if err != nil {
			return contract.PipelineInfo{}, err
		}
		return contract.PipelineInfo{}, errors.New(provider.ErrNotMatched)
	}

	// translate its status
	scm, _ := k8s.CreateSCMContextFromKubernetesAnnotations(job.ObjectMeta)
	jobStatus := translateJobStatus(job)

	// start time
	var startTime time.Time
	if job.Status.StartTime != nil {
		startTime = job.Status.StartTime.Time
	}

	pi := contract.NewPipelineInfo(
		scm,
		job.Namespace,
		job.Name,
		string(job.UID),
		startTime,
		jobStatus,
		[]contract.PipelineStage{
			{Name: "Job", Status: jobStatus},
		},
	)

	eventNum := bjp.store.CountHowManyTimesKubernetesResourceReceived(pi)
	bjp.logger.Debugf("count(%s) = %v", job.Name, eventNum)
	pi.SetRetrievalCount(eventNum)

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
