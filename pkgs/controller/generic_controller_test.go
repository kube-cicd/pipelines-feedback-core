package controller_test

import (
	"context"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/controller"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/fake"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/store"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"testing"
	"time"
)

func TestGenericController_ReconcileWithRetries(t *testing.T) {
	pipeline := contract.NewPipelineInfo(
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
	)

	receiver := &fake.Receiver{}

	gc := controller.GenericController{
		PipelineInfoProvider:        &fake.Provider{Pipeline: *pipeline, Error: nil},
		FeedbackReceiver:            receiver,
		ObjectType:                  &v1.Job{},
		Store:                       store.Operator{Store: store.NewMemory()},
		DelayAfterErrorNum:          5,
		RequeueDelaySecs:            35,
		StopProcessingAfterErrorNum: 30,
	}
	_ = gc.InjectDependencies(
		&fake.Recorder{},
		&rest.Config{},
		logging.CreateLogger(false),
		&fake.ConfigurationProvider{
			Contextual: config.Data{},
			Global:     config.Data{},
		},
		&fake.NullValidator{},
	)

	// Mock: make UpdateProgress() always return error
	receiver.UpdateProgressReturns = errors.New("blah blah blah")

	for i := 1; i <= 5; i++ {
		result, err := gc.Reconcile(context.TODO(), controllerruntime.Request{
			NamespacedName: types.NamespacedName{
				Name:      "book",
				Namespace: "bookchin",
			},
		})
		assert.NotNil(t, result)
		assert.Nil(t, err)
	}

	// 5th time (DelayAfterErrorNum >= 5)
	result, err := gc.Reconcile(context.TODO(), controllerruntime.Request{
		NamespacedName: types.NamespacedName{
			Name:      "book",
			Namespace: "bookchin",
		},
	})
	assert.Equal(t, time.Second*35, result.RequeueAfter, "After 5th retry the reconciliation should be delayed by 35 seconds")
	assert.Nil(t, err)
	logrus.Println(result, err)

	// Go to the limit StopProcessingAfterErrorNum == 30
	for i := 1; i <= 24; i++ {
		_, _ = gc.Reconcile(context.TODO(), controllerruntime.Request{
			NamespacedName: types.NamespacedName{
				Name:      "book",
				Namespace: "bookchin",
			},
		})
	}

	// StopProcessingAfterErrorNum + 1: So now, the reconciliation should be cancelled
	result, err = gc.Reconcile(context.TODO(), controllerruntime.Request{
		NamespacedName: types.NamespacedName{
			Name:      "book",
			Namespace: "bookchin",
		},
	})
	assert.False(t, result.Requeue, "Reconciliation should be cancelled, so no requeue should be returned")
	assert.Equal(t, time.Duration(0), result.RequeueAfter, "Reconciliation should be cancelled, so no requeue should be returned")
	assert.Nil(t, err)
}
