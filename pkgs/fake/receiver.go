package fake

import (
	"context"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
)

type Receiver struct {
	UpdateProgressReturns error
}

// UpdateProgress is called each time a status is changed
func (r *Receiver) UpdateProgress(ctx context.Context, status contract.PipelineInfo, log *logging.InternalLogger) error {
	return r.UpdateProgressReturns
}

// WhenCreated is an event, when a Pipeline was created and is in Pending or already in Running state
func (r *Receiver) WhenCreated(ctx context.Context, status contract.PipelineInfo, log *logging.InternalLogger) error {
	return nil
}

// WhenStarted is an event, when a Pipeline is started
func (r *Receiver) WhenStarted(ctx context.Context, status contract.PipelineInfo, log *logging.InternalLogger) error {
	return nil
}

// WhenFinished is an event, when a Pipeline is finished - Failed, Errored, Aborted or Succeeded
func (r *Receiver) WhenFinished(ctx context.Context, status contract.PipelineInfo, log *logging.InternalLogger) error {
	return nil
}

func (r *Receiver) CanHandle(adapterName string) bool {
	return true
}

func (r *Receiver) GetImplementationName() string {
	return "fake"
}
