package fake

import (
	"context"

	"github.com/kube-cicd/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
)

// NullValidator is an empty implementation used in unit tests
type NullValidator struct{}

func (nv *NullValidator) ValidateRequestedEntry(group string, key string) error {
	return nil
}

func (nv *NullValidator) ValidateConfig(data v1alpha1.Data) error {
	return nil
}

func (nv *NullValidator) Add(schema config.Schema) {
}

type FakeConfigurationProvider struct {
}

func (f *FakeConfigurationProvider) FetchContextual(component string, namespace string, pipeline contract.PipelineInfo) config.Data {
	return config.NewData(component, map[string]string{}, &NullValidator{}, logging.CreateLogger(true))
}

func (f *FakeConfigurationProvider) FetchGlobal(component string) config.Data {
	return config.NewData(component, map[string]string{}, &NullValidator{}, logging.CreateLogger(true))
}

func (f *FakeConfigurationProvider) FetchSecretKey(ctx context.Context, name string, namespace string, key string, cache bool) (string, error) {
	return "", nil
}

func (f *FakeConfigurationProvider) FetchFromFieldOrSecret(ctx context.Context, data *config.Data, namespace string, fieldKey string, referenceKey string, referenceSecretNameKey string) (string, error) {
	return "", nil
}
