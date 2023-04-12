package config

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

type LocalFileConfigurationProvider struct {
}

func (lf *LocalFileConfigurationProvider) InjectDependencies(recorder record.EventRecorder, kubeconfig *rest.Config) error {
	return nil
}

func (lf *LocalFileConfigurationProvider) Provide(info contract.PipelineInfo) Configuration {
	return Configuration{}
}

func (lf *LocalFileConfigurationProvider) CanHandle(adapterName string) bool {
	return adapterName == lf.GetImplementationName()
}

func (lf *LocalFileConfigurationProvider) GetImplementationName() string {
	return "local"
}
