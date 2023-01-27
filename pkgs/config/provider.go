package config

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

type ConfigurationProvider interface {
	Provide(info contract.PipelineInfo) Configuration
	InjectDependencies(recorder record.EventRecorder, kubeconfig *rest.Config) error
}

type Configuration struct {
	// how to connect to e.g. gitlab (jx scm)
	keys map[string]string
}

func (c *Configuration) GetOrDefault(keyName string, defaultValue string) string {
	if val, exists := c.keys[keyName]; exists {
		return val
	}
	return defaultValue
}
