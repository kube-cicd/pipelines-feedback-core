package fake

import (
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
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
