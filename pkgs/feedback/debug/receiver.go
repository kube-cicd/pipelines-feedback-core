package debug

import (
	"context"

	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract/wiring"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
	"github.com/sirupsen/logrus"
)

type Receiver struct{}

func (d *Receiver) InitializeWithContext(sc *wiring.ServiceContext) error {
	return nil
}

func (d *Receiver) WhenCreated(ctx context.Context, pipeline contract.PipelineInfo, log *logging.InternalLogger) error {
	log.Info("debug.WhenCreated()")

	return nil
}

func (d *Receiver) WhenStarted(ctx context.Context, pipeline contract.PipelineInfo, log *logging.InternalLogger) error {
	log.Info("debug.WhenStarted()")

	return nil
}

func (d *Receiver) UpdateProgress(ctx context.Context, pipeline contract.PipelineInfo, log *logging.InternalLogger) error {
	log.Info("debug.UpdateProgress()")
	logrus.Info(pipeline)

	return nil
}

func (d *Receiver) WhenFinished(ctx context.Context, pipeline contract.PipelineInfo, log *logging.InternalLogger) error {
	log.Info("debug.WhenFinished()")

	return nil
}

func (d *Receiver) CanHandle(name string) bool {
	return true
}

func (d *Receiver) GetImplementationName() string {
	return "debug"
}
