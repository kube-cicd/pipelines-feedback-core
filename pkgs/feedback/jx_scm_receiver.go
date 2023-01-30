package feedback

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract/wiring"
)

type JXSCMReceiver struct {
}

func (JXSCMReceiver) InitializeWithContext(sc *wiring.ServiceContext) error {
	return nil
}

func (JXSCMReceiver) Update(status contract.PipelineInfo) error {
	return nil
}
