package provider

import (
	"context"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

const (
	ErrNotMatched = "not-matched"
)

type Provider interface {
	ReceivePipelineInfo(ctx context.Context, name string, namespace string) (contract.PipelineInfo, error)
	InjectDependencies(recorder record.EventRecorder, kubeconfig *rest.Config) error
}
