package feedback

import (
	"context"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
)

// Receiver is a side that receives the update about the Pipeline status. It may be a Gitea, Gitlab, GitHub or other
type Receiver interface {
	Update(ctx context.Context, status contract.PipelineInfo) error
}
