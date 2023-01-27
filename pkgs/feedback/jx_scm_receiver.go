package feedback

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

type JXSCMReceiver struct {
}

func (JXSCMReceiver) InjectDependencies(recorder record.EventRecorder, kubeconfig *rest.Config) error {
	return nil
}

func (JXSCMReceiver) Update(status contract.PipelineInfo) error {
	return nil
}
