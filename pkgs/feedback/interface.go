package feedback

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

// Receiver is a side that receives the update about the Pipeline status. It may be a Gitea, Gitlab, GitHub or other
type Receiver interface {
	Update(status contract.PipelineInfo) error
	InjectDependencies(recorder record.EventRecorder, kubeconfig *rest.Config) error
}
