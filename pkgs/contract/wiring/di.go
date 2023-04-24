package wiring

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/config"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/logging"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/store"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

type WithInitialization interface {
	InitializeWithContext(sc *ServiceContext) error
}

type ServiceContext struct {
	Recorder   *record.EventRecorder
	KubeConfig *rest.Config
	Config     config.ConfigurationProvider
	Log        logging.Logger
	Store      *store.Operator
}
