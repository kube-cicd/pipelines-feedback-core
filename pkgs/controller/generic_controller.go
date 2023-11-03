package controller

import (
	"context"
	"fmt"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract/wiring"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/feedback"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/provider"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/store"
	"github.com/pkg/errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"
)

type GenericController struct {
	ObjectType client.Object

	// e.g. Kubernetes batch/v1 Job, Argo Workflow or Tekton Pipeline
	PipelineInfoProvider provider.Provider

	// e.g. a Gitlab, Gitea, Bitbucket, MS Teams, etc.
	FeedbackReceiver feedback.Receiver

	// simple key-value store
	Store store.Operator

	recorder record.EventRecorder

	kubeConfig *rest.Config

	logger *logging.InternalLogger
}

func (gc *GenericController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logging.CreateK8sContextualLogger(ctx, gc.logger, req)

	//
	// Fetch the object from PipelineInfoProvider
	//
	retrieved, retrieveErr := gc.PipelineInfoProvider.ReceivePipelineInfo(ctx, req.Name, req.Namespace)
	if retrieveErr != nil {
		// log: not matched
		if retrieveErr.Error() == provider.ErrNotMatched {
			logger.Debugf("resource not matched. The provider declined to retrieve it")
			return ctrl.Result{}, nil
		}

		logger.Errorf("cannot retrieve resource: %v", retrieveErr.Error())
		return ctrl.Result{Requeue: true}, errors.Wrap(retrieveErr, "cannot receive Pipeline status")
	}

	// record how many times an object was reconciled
	// this info is later used to send update or not
	eventNum := gc.Store.CountHowManyTimesKubernetesResourceReceived(&retrieved)
	logger.Debugf("count(%s) = %v", req.Name, eventNum)
	retrieved.SetRetrievalCount(eventNum)

	//
	// Notify the Feedback Receiver
	//
	if err := gc.updateProgress(ctx, retrieved, logger); err != nil {
		logger.Errorf("cannot update feedback retriever: %s", err.Error())
		return ctrl.Result{RequeueAfter: time.Second * 5}, nil
	}

	return ctrl.Result{}, nil
}

// updateProgress decides when to trigger notification events to the RECEIVER
func (gc *GenericController) updateProgress(ctx context.Context, retrieved contract.PipelineInfo, logger *logging.InternalLogger) error {
	for _, stage := range retrieved.GetStages() {
		logger.Debugf("[%s] %s: %s", retrieved.GetId(), stage.Name, stage.Status.AsHumanReadableDescription())
	}
	logger.Debugf("[%s] status: %s", retrieved.GetId(), retrieved.GetStatus().AsHumanReadableDescription())

	// Always update progress
	logger.Debugf("GenericController -> UpdateProgress(%s)", retrieved.GetId())
	if upErr := gc.FeedbackReceiver.UpdateProgress(ctx, retrieved); upErr != nil {
		return upErr
	}

	// Single-time events
	if retrieved.IsJustCreated() && !retrieved.GetStatus().IsFinished() && !gc.Store.WasEventAlreadySent(retrieved, "created") {
		logger.Debugf("GenericController -> WhenCreated(%s)", retrieved.GetId())
		if createErr := gc.FeedbackReceiver.WhenCreated(ctx, retrieved); createErr != nil {
			return createErr
		}
		if recErr := gc.Store.RecordEventFiring(retrieved, "created"); recErr != nil {
			logger.Warning("cannot record event 'created'")
		}
	}
	if retrieved.GetStatus().IsRunning() && !gc.Store.WasEventAlreadySent(retrieved, "started") {
		logger.Debugf("GenericController -> WhenStarted(%s)", retrieved.GetId())
		if startErr := gc.FeedbackReceiver.WhenStarted(ctx, retrieved); startErr != nil {
			return startErr
		}
		if recErr := gc.Store.RecordEventFiring(retrieved, "started"); recErr != nil {
			logger.Warning("cannot record event 'started'")
		}
	}
	if retrieved.GetStatus().IsFinished() && !gc.Store.WasEventAlreadySent(retrieved, "finished") {
		logger.Debugf("GenericController -> WhenFinished(%s)", retrieved.GetId())
		if finishErr := gc.FeedbackReceiver.WhenFinished(ctx, retrieved); finishErr != nil {
			return finishErr
		}
		if recErr := gc.Store.RecordEventFiring(retrieved, "finished"); recErr != nil {
			logger.Warning("cannot record event 'finished'")
		}
	}
	return nil
}

func (gc *GenericController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(gc.ObjectType).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(e event.DeleteEvent) bool {
				return false
			},
		}).
		Complete(gc)
}

// InjectDependencies is wiring dependencies to all services
func (gc *GenericController) InjectDependencies(recorder record.EventRecorder, kubeConfig *rest.Config,
	logger *logging.InternalLogger, configProvider config.ConfigurationProvider, cfgSchema *config.SchemaValidator) error {

	gc.recorder = recorder
	gc.kubeConfig = kubeConfig
	gc.logger = logger
	sc := wiring.ServiceContext{
		Recorder:     &recorder,
		KubeConfig:   kubeConfig,
		Config:       configProvider,
		Log:          logger.ForkWithFields(context.TODO(), map[string]interface{}{}),
		Store:        &gc.Store,
		ConfigSchema: cfgSchema,
	}
	nErr := func(name string, err error) error {
		return errors.Wrap(err, fmt.Sprintf("cannot inject dependencies to %s", name))
	}
	if _, ok := gc.PipelineInfoProvider.(wiring.WithInitialization); ok {
		if err := gc.PipelineInfoProvider.(wiring.WithInitialization).InitializeWithContext(&sc); err != nil {
			return nErr("PipelineInfoProvider", err)
		}
	}
	if _, ok := gc.FeedbackReceiver.(wiring.WithInitialization); ok {
		if err := gc.FeedbackReceiver.(wiring.WithInitialization).InitializeWithContext(&sc); err != nil {
			return nErr("FeedbackReceiver", err)
		}
	}
	return nil
}
