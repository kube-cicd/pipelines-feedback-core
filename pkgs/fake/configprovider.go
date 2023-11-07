package fake

import (
	"context"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
)

type ConfigurationProvider struct {
	Contextual config.Data
	Global     config.Data
}

func (cp *ConfigurationProvider) FetchContextual(component string, namespace string, pipeline contract.PipelineInfo) config.Data {
	return cp.Contextual
}

func (cp *ConfigurationProvider) FetchGlobal(component string) config.Data {
	return cp.Global
}

func (cp *ConfigurationProvider) FetchSecretKey(ctx context.Context, name string, namespace string, key string, cache bool) (string, error) {
	return "", nil
}

func (cp *ConfigurationProvider) FetchFromFieldOrSecret(ctx context.Context, data *config.Data, namespace string, fieldKey string, referenceKey string, referenceSecretNameKey string) (string, error) {
	return "", nil
}
