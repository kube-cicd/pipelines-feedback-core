package batchjob

import (
	store2 "github.com/kube-cicd/pipelines-feedback-core/internal/store"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/controller"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/store"
	v1 "k8s.io/api/batch/v1"
)

func CreateJobController() *controller.GenericController {
	return &controller.GenericController{
		PipelineInfoProvider: &BatchV1JobProvider{},
		ObjectType:           &v1.Job{},
		Store:                store.Operator{Store: store2.NewMemory()},
		// todo: schema provider
	}
}
