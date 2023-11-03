package feedback

import (
	"context"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
)

// Receiver is a side that receives the update about the Pipeline status. It may be a Gitea, Gitlab, GitHub or other
type Receiver interface {
	contract.Pluggable

	// UpdateProgress is called each time a status is changed
	UpdateProgress(ctx context.Context, status contract.PipelineInfo, log *logging.InternalLogger) error

	// WhenCreated is an event, when a Pipeline was created and is in Pending or already in Running state
	WhenCreated(ctx context.Context, status contract.PipelineInfo, log *logging.InternalLogger) error

	// WhenStarted is an event, when a Pipeline is started
	WhenStarted(ctx context.Context, status contract.PipelineInfo, log *logging.InternalLogger) error

	// WhenFinished is an event, when a Pipeline is finished - Failed, Errored, Aborted or Succeeded
	WhenFinished(ctx context.Context, status contract.PipelineInfo, log *logging.InternalLogger) error
}
