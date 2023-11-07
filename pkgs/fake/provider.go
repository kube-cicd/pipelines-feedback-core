package fake

import (
	"context"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
)

type Provider struct {
	Pipeline contract.PipelineInfo
	Error    error
}

func (fp *Provider) ReceivePipelineInfo(ctx context.Context, name string, namespace string, log *logging.InternalLogger) (contract.PipelineInfo, error) {
	return fp.Pipeline, fp.Error
}
