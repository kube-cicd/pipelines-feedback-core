package controller

import (
	"context"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/config"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/feedback"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/provider"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

	// can read configuration from various sources
	ConfigProvider config.ConfigurationProvider

	recorder record.EventRecorder

	kubeconfig *rest.Config
}

func (gc *GenericController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := createLogger(ctx, req)

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

		logger.Errorf("cannot retrieve resource")
		return ctrl.Result{Requeue: true}, errors.Wrap(retrieveErr, "cannot receive Pipeline status")
	}

	//
	// Notify the Feedback Receiver
	//
	if err := gc.FeedbackReceiver.Update(retrieved); err != nil {
		logger.Errorf("cannot update feedback retriever")
		return ctrl.Result{RequeueAfter: time.Second * 30}, errors.Wrap(err, "cannot update feedback receiver")
	}

	return ctrl.Result{}, nil
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

func (gc *GenericController) InjectDependencies(recorder record.EventRecorder, kubeconfig *rest.Config) error {
	gc.recorder = recorder
	gc.kubeconfig = kubeconfig

	if err := gc.ConfigProvider.InjectDependencies(recorder, kubeconfig); err != nil {
		return errors.Wrap(err, "cannot inject dependencies to ConfigProvider")
	}
	if err := gc.PipelineInfoProvider.InjectDependencies(recorder, kubeconfig); err != nil {
		return errors.Wrap(err, "cannot inject dependencies to ConfigProvider")
	}
	if err := gc.FeedbackReceiver.InjectDependencies(recorder, kubeconfig); err != nil {
		return errors.Wrap(err, "cannot inject dependencies to ConfigProvider")
	}
	return nil
}

func createLogger(ctx context.Context, req ctrl.Request) *logrus.Entry {
	id, _ := uuid.NewUUID()
	return logrus.WithContext(ctx).WithFields(map[string]interface{}{
		"request": id,
		"name":    req.NamespacedName,
	})
}
