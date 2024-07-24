package contract_test

import (
	"testing"
	"time"

	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/fake"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/labels"
)

func TestPipelineInfo_GetLogs(t *testing.T) {
	cfg := config.NewData("", map[string]string{}, &fake.NullValidator{}, logging.NewInternalLogger())
	pi := contract.NewPipelineInfo(
		contract.JobContext{
			Commit:           "123",
			Reference:        "test",
			RepoHttpsUrl:     "",
			PrId:             "",
			OrganizationName: "",
			RepositoryName:   "",
			TechnicalJob:     "",
		},
		"test-ns",
		"bread-pipeline",
		"a-slice-123",
		time.Now(),
		[]contract.PipelineStage{
			{Name: "clone", Status: contract.PipelineSucceeded},
		},
		labels.Set{},
		labels.Set{},
		&cfg,
		contract.PipelineInfoWithLogsCollector(func() string {
			return "test123"
		}),
	)

	assert.Equal(t, "test123", pi.GetLogs(), "Expecting logs returned, as by default logs are enabled - 'global.logs-enable == false' not defined")

	// ------
	// Step 2
	// ------
	// and now we disable logs by overwriting configuration by pointer
	cfg = config.NewData("", map[string]string{
		"logs-enabled": "false",
	}, &fake.NullValidator{}, logging.NewInternalLogger())

	assert.Equal(t, "", pi.GetLogs(), "Expecting empty logs - when 'global.logs-enable == false'")
}

func TestNewSCMContext_SupportsNestedDirectoriesStructureInGITUrl(t *testing.T) {
	ctx, err := contract.NewSCMContext("https://git.example.org/books/anarchism/bakunin.git")

	assert.Equal(t, "books/anarchism", ctx.OrganizationName)
	assert.Equal(t, "bakunin", ctx.RepositoryName)
	assert.Nil(t, err)
}

func TestNewSCMContext_WithSimpleGITUrl(t *testing.T) {
	ctx, err := contract.NewSCMContext("https://git.example.org/books/bakunin.git")

	assert.Equal(t, "books", ctx.OrganizationName)
	assert.Equal(t, "bakunin", ctx.RepositoryName)
	assert.Nil(t, err)
}
