package batchjob

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/controller"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/store"
	v1 "k8s.io/api/batch/v1"
)

func CreateJobController() *controller.GenericController {
	return &controller.GenericController{
		PipelineInfoProvider: &BatchV1JobProvider{},
		ObjectType:           &v1.Job{},
		Store:                store.Operator{Store: store.NewMemory()},
		// todo: schema provider
	}
}
