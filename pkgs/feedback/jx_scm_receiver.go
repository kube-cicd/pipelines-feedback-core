package feedback

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract/wiring"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/pkg/errors"
)

type JXSCMReceiver struct {
	client *scm.Client
	sc     *wiring.ServiceContext
}

func (jx *JXSCMReceiver) InitializeWithContext(sc *wiring.ServiceContext) error {
	sc.Log.Infoln("Initializing JXSCMReceiver")

	client, err := factory.NewClientFromEnvironment()
	jx.client = client
	jx.sc = sc

	if err != nil {
		return errors.Wrap(err, "cannot initialize Jenkins X's go-scm client")
	}
	return nil
}

func (jx *JXSCMReceiver) Update(status contract.PipelineInfo) error {
	jx.sc.Log.Infoln("JXSCMReceiver", status)
	return nil
}
