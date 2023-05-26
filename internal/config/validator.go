package config

import "github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"

type Validator interface {
	ValidateConfig(data v1alpha1.Data) error
}
