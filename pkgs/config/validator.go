package config

import "github.com/kube-cicd/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"

type Validator interface {
	ValidateRequestedEntry(group string, key string) error
	ValidateConfig(data v1alpha1.Data) error
}
