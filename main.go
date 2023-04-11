package main

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/cli"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/implementation/batchjob"
	"os"
)

func main() {
	app := cli.PipelinesFeedbackApp{
		JobController: batchjob.CreateJobController(),
	}
	cmd := cli.NewRootCommand(&app)
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
