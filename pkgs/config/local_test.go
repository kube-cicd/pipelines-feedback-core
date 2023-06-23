package config_test

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/config"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/logging"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLocalCollector_CollectInitially_NotSetUp(t *testing.T) {
	local := config.NewLocalFileConfigurationCollector(logging.NewInternalLogger(), "")
	configs, err := local.CollectInitially()
	assert.Nil(t, err)
	assert.Empty(t, configs)
}

func TestLocalCollector_CollectInitially_WithConfigFound(t *testing.T) {
	pwd, _ := os.Getwd()
	local := config.NewLocalFileConfigurationCollector(logging.NewInternalLogger(), pwd+"/testdata/pipelines-feedback.json")
	configs, err := local.CollectInitially()
	assert.Nil(t, err)
	assert.NotEmpty(t, configs)
}
