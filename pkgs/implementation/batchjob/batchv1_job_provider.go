package batchjob

import (
	"context"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/config"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract/wiring"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/k8s"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/logging"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/provider"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/store"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/templating"
	"github.com/pkg/errors"
	v1model "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"time"
)

type BatchV1JobProvider struct {
	batchV1Client *v1.BatchV1Client
	coreV1Client  *v1core.CoreV1Client
	store         *store.Operator
	logger        logging.Logger
	confProvider  *config.ConfigurationProvider
}

func (bjp *BatchV1JobProvider) InitializeWithContext(sc *wiring.ServiceContext) error {
	client, err := v1.NewForConfig(sc.KubeConfig)
	if err != nil {
		return errors.Wrap(err, "cannot initialize BatchV1JobProvider")
	}
	bjp.batchV1Client = client
	coreClient, err := v1core.NewForConfig(sc.KubeConfig)
	if err != nil {
		return errors.Wrap(err, "cannot initialize BatchV1JobProvider")
	}
	bjp.coreV1Client = coreClient
	bjp.store = sc.Store
	bjp.logger = sc.Log
	bjp.confProvider = &sc.Config
	return nil
}

// ReceivePipelineInfo is tracking batch/v1, kind: Job type objects
func (bjp *BatchV1JobProvider) ReceivePipelineInfo(ctx context.Context, name string, namespace string) (contract.PipelineInfo, error) {
	globalCfg := bjp.confProvider.FetchGlobal("global")

	// find an object
	job, err := bjp.batchV1Client.Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
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
	scm, _ := k8s.CreateJobContextFromKubernetesAnnotations(job.ObjectMeta)
	jobStatus := translateJobStatus(job)

	// start time
	var startTime time.Time
	if job.Status.StartTime != nil {
		startTime = job.Status.StartTime.Time
	}

	dashboardUrl, dashboardTplErr := templating.TemplateDashboardUrl(globalCfg.Get("dashboard-url"), job, job.TypeMeta)
	if dashboardTplErr != nil {
		bjp.logger.Warningf("Cannot render dashboard template URL '%s': '%s'", dashboardUrl, dashboardTplErr.Error())
	}

	// logs are lazy-fetched on demand
	logs := func() string { return bjp.fetchLogs(ctx, job, globalCfg) }

	// create an universal PipelineInfo object
	pi := contract.NewPipelineInfo(
		scm,
		job.Namespace,
		job.Name,
		string(job.UID),
		startTime,
		jobStatus,
		[]contract.PipelineStage{
			{Name: "job/" + job.Name, Status: jobStatus},
		},
		dashboardUrl,
		labels.Set(job.Labels),
		labels.Set(job.Annotations),
		logs,
	)

	return *pi, nil
}

func (bjp *BatchV1JobProvider) fetchLogs(ctx context.Context, job *v1model.Job, data config.Data) string {
	return k8s.TruncateLogs(
		k8s.FindAndReadLogsFromLastPod(ctx, bjp.coreV1Client.Pods(job.Namespace), labels.Set(job.Spec.Selector.MatchLabels).String()),
		data,
	)
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
