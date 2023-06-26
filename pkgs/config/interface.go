package config

import (
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
)

type ConfigurationCollector interface {
	contract.Pluggable
	CollectInitially() ([]*v1alpha1.PFConfig, error)
	CollectOnRequest(info contract.PipelineInfo) ([]*v1alpha1.PFConfig, error)
	SetLogger(logger *logging.InternalLogger)
}
