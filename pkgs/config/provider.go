package config

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
)

type ConfigurationProvider interface {
	contract.Pluggable
	Provide(info contract.PipelineInfo) Configuration
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
