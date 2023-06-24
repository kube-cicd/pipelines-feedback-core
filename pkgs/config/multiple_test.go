package config_test

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/config"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/logging"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultipleCollector_CollectInitially(t *testing.T) {
	logger := logging.NewInternalLogger()
	collector := config.CreateMultipleCollector(
		[]config.ConfigurationCollector{
			config.NewLocalFileConfigurationCollector(logger, ""),
			&config.FakeCollector{},
		},
		logger,
	)
	cfgs, err := collector.CollectInitially()

	assert.Nil(t, err)
	assert.NotEmpty(t, cfgs)
	assert.Equal(t, "bread", cfgs[0].Name)
}

func TestMultipleCollector_CollectOnRequest(t *testing.T) {
	logger := logging.NewInternalLogger()
	collector := config.CreateMultipleCollector(
		[]config.ConfigurationCollector{
			config.NewLocalFileConfigurationCollector(logger, ""),
			&config.FakeCollector{},
		},
		logger,
	)
	cfgs, err := collector.CollectOnRequest(contract.PipelineInfo{})

	assert.Nil(t, err)
	assert.NotEmpty(t, cfgs)
	assert.Equal(t, "mutual-aid-a-factor-of-revolution", cfgs[0].Name)
}
