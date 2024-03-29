package controller

import (
	"context"
	configinternal "github.com/kube-cicd/pipelines-feedback-core/internal/config"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	v1alpha1client "github.com/kube-cicd/pipelines-feedback-core/pkgs/client/clientset/versioned"
	pipelinesfeedbackv1alpha1 "github.com/kube-cicd/pipelines-feedback-core/pkgs/client/clientset/versioned/typed/pipelinesfeedback.keskad.pl/v1alpha1"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/store"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

// ConfigurationController is reconciling CRD that provides configuration
type ConfigurationController struct {
	Provider config.ConfigurationProviderInterface
	docs     configinternal.IndexedDocumentStore
	client   pipelinesfeedbackv1alpha1.PipelinesfeedbackV1alpha1Interface
	logger   *logging.InternalLogger
}

func (cc *ConfigurationController) Initialize(kubeConfig *rest.Config, collector config.ConfigurationCollector,
	logger *logging.InternalLogger, kvStore store.Operator, schema *config.SchemaValidator) error {

	cc.logger = logger
	client, err := v1alpha1client.NewForConfig(kubeConfig)
	if err != nil {
		return errors.Wrap(err, "cannot initialize BatchV1JobProvider")
	}

	// storage
	cc.docs = configinternal.CreateIndexedDocumentStore(schema)

	// API interface for components
	coreV1Client, err := corev1.NewForConfig(kubeConfig)
	if err != nil {
		return errors.Wrap(err, "cannot construct ConfigurationProvider, "+
			"Kubernetes Core API v1 construction error")
	}
	cc.Provider, err = config.NewConfigurationProvider(cc.docs, logger, coreV1Client, kvStore, schema)
	if err != nil {
		return errors.Wrap(err, "cannot initialize ConfigurationController")
	}

	cc.client = client.PipelinesfeedbackV1alpha1()
	return nil
}

func (cc *ConfigurationController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	cfg, err := cc.client.PFConfigs(req.Namespace).Get(ctx, req.Name, v1.GetOptions{})
	if err != nil {
		// not found anymore?
		cc.docs.Delete(req.Namespace, req.Name)
		return ctrl.Result{}, err
	}

	if pErr := cc.docs.Push(cfg); pErr != nil {
		cc.logger.Errorf("cannot load configuration '%s' coming from Kubernetes: %s", req.NamespacedName, pErr.Error())
		return ctrl.Result{RequeueAfter: 300}, pErr
	}
	cc.logger.Infof("loaded configuration '%s' from Kubernetes", req.NamespacedName)

	return ctrl.Result{}, nil
}

func (cc *ConfigurationController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PFConfig{}).
		Complete(cc)
}

func (cc *ConfigurationController) CollectInitially(collector config.ConfigurationCollector) error {
	docs, err := collector.CollectInitially()
	if err != nil {
		return errors.Wrap(err, "cannot initially read configuration")
	}
	for _, doc := range docs {
		if pErr := cc.docs.Push(doc); pErr != nil {
			return pErr
		}
	}
	return nil
}
