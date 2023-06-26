package config_test

import (
	intConfig "github.com/kube-cicd/pipelines-feedback-core/internal/config"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestData_GetOrDefault(t *testing.T) {
	cfg := config.NewData("jxscm", map[string]string{
		"hello": "bread",
	}, &intConfig.NullValidator{}, &logging.InternalLogger{})

	assert.Equal(t, "my-default", cfg.GetOrDefault("this-key-does-not-exist", "my-default"))
	assert.Equal(t, "bread", cfg.GetOrDefault("hello", "something"))
}

func TestData_GetOrDefault_WithValidationFailure(t *testing.T) {
	logger := FakeLogger{}
	cfg := config.NewData("jxscm", map[string]string{
		"hello": "bread",
	}, &config.SchemaValidator{}, &logger)

	cfg.GetOrDefault("hello", "")
	assert.True(t, logger.Called)
}

type FakeLogger struct {
	Called bool
}

func (l *FakeLogger) Fatalf(format string, args ...interface{}) {
	l.Called = true
}
