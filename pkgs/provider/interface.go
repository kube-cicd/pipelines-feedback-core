package provider

import (
	"context"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
)

const (
	ErrNotMatched = "not-matched"
)

type Provider interface {
	ReceivePipelineInfo(ctx context.Context, name string, namespace string) (contract.PipelineInfo, error)
}
