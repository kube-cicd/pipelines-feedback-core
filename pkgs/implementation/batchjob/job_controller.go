package batchjob

import (
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/controller"
	v1 "k8s.io/api/batch/v1"
)

func CreateJobController() *controller.GenericController {
	return &controller.GenericController{
		PipelineInfoProvider: &BatchV1JobProvider{},
		ObjectType:           &v1.Job{},
	}
}
