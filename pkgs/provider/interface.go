package provider

import (
	"context"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
)

const (
	ErrNotMatched = "not-matched"
)

type Provider interface {
	ReceivePipelineInfo(ctx context.Context, name string, namespace string, log *logging.InternalLogger) (contract.PipelineInfo, error)
}
