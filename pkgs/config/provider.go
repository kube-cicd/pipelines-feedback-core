package config

import "github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"

type ConfigurationProvider interface {
	Provide(info contract.PipelineInfo) Configuration
}

type Configuration struct {
	// how to connect to e.g. gitlab (jx scm)
}
