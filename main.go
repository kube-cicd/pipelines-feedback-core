package main

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/app"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/cli"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/controller"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/implementation/batchjob"
	"os"
)

func main() {
	pfcApp := app.PipelinesFeedbackApp{
		JobController:    batchjob.CreateJobController(),
		ConfigController: &controller.ConfigurationController{},
	}
	cmd := cli.NewRootCommand(&pfcApp)
	args := os.Args
	if args != nil {
		args = args[1:]
		cmd.SetArgs(args)
	}
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
