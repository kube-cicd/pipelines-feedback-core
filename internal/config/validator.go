package config

import "github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"

type Validator interface {
	ValidateRequestedEntry(group string, key string) error
	ValidateConfig(data v1alpha1.Data) error
}

// NullValidator is an empty implementation used in unit tests
type NullValidator struct{}

func (nv *NullValidator) ValidateRequestedEntry(group string, key string) error {
	return nil
}

func (nv *NullValidator) ValidateConfig(data v1alpha1.Data) error {
	return nil
}
