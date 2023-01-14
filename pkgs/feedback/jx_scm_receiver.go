package feedback

import "github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"

type JXSCMReceiver struct {
}

func (JXSCMReceiver) Update(status contract.PipelineInfo) error {
	return nil
}
