package jxscm_test

import (
	"testing"

	"github.com/kube-cicd/pipelines-feedback-core/internal/feedback/jxscm"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
	"github.com/stretchr/testify/assert"
)

func TestNewClientFromConfig_InvalidGitServerURL(t *testing.T) {
	data := config.NewData("jx-scm", map[string]string{
		"git-kind":   "github",
		"git-server": "ht!tp://bad_url",
		"git-token":  "dummy-token",
	}, &config.FakeValidator{}, nil)

	_, err := jxscm.NewClientFromConfig(data, "")
	assert.NotNil(t, err)
	assert.Equal(t, "invalid git-server URL: \"ht!tp://bad_url\". valid values are empty or a correctly formatted URL address", err.Error())
}

func TestNewClientFromConfig_ValidGitServerURL(t *testing.T) {
	data := config.NewData("jx-scm", map[string]string{
		"git-kind":   "github",
		"git-server": "github.com",
		"git-token":  "dummy-token",
	}, &config.FakeValidator{}, nil)

	client, err := jxscm.NewClientFromConfig(data, "")

	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestNewClientFromConfig_EmptyGitServerURL(t *testing.T) {
	data := config.NewData("jx-scm", map[string]string{
		"git-kind":   "github",
		"git-server": "",
		"git-token":  "dummy-token",
	}, &config.FakeValidator{}, nil)

	client, err := jxscm.NewClientFromConfig(data, "")

	assert.Nil(t, err)
	assert.NotNil(t, client)
}
