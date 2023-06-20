package store_test

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/fields"
	"testing"
	"time"
)
import "github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/store"
import internalstore "github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/internal/store"

func createBreadBookPipeline() *contract.PipelineInfo {
	scm, _ := contract.NewSCMContext("https://gitlab.com/aaa/bbb.git")
	return contract.NewPipelineInfo(
		scm,
		"default",
		"hello-kropotkin",
		"the-conquest-of-bread",
		time.Now(),
		contract.Running,
		[]contract.PipelineStage{},
		"https://dashboard.tekton.local/pipeline-some/pipeline",
		fields.Set{},
		fields.Set{},
		func() string {
			return "Baked!"
		},
	)
}

func TestOperator_CountHowManyTimesKubernetesResourceReceived(t *testing.T) {
	o := store.Operator{Store: internalstore.NewMemory()}
	scm, _ := contract.NewSCMContext("https://gitlab.com/aaa/bbb.git")

	firstPipeline := contract.NewPipelineInfo(
		scm,
		"default",
		"hello-kropotkin",
		"the-conquest-of-bread",
		time.Now(),
		contract.Running,
		[]contract.PipelineStage{},
		"https://dashboard.tekton.local/pipeline-some/pipeline",
		fields.Set{},
		fields.Set{},
		func() string {
			return "Baked!"
		},
	)
	for _, _ = range []int{1, 2, 3} {
		o.CountHowManyTimesKubernetesResourceReceived(firstPipeline)
	}
	assert.Equal(t, 4, o.CountHowManyTimesKubernetesResourceReceived(firstPipeline))
	assert.Equal(t, 5, o.CountHowManyTimesKubernetesResourceReceived(firstPipeline))

	// second try is to check that one Pipeline does not impact other Pipeline
	secondPipeline := contract.NewPipelineInfo(
		scm,
		"default",
		"hello-francisco-ferrer",
		"the-barcelona-school",
		time.Now(),
		contract.Running,
		[]contract.PipelineStage{},
		"https://dashboard.tekton.local/pipeline-some/pipeline",
		fields.Set{},
		fields.Set{},
		func() string {
			return "Created!"
		},
	)
	assert.Equal(t, 1, o.CountHowManyTimesKubernetesResourceReceived(secondPipeline))
	assert.Equal(t, 2, o.CountHowManyTimesKubernetesResourceReceived(secondPipeline))
}

func TestOperator_RecordEventFiring_WasEventAlreadySent(t *testing.T) {
	o := store.Operator{Store: internalstore.NewMemory()}

	assert.False(t, o.WasEventAlreadySent(*createBreadBookPipeline(), "start"))

	_ = o.RecordEventFiring(*createBreadBookPipeline(), "start")
	assert.True(t, o.WasEventAlreadySent(*createBreadBookPipeline(), "start"))
}
