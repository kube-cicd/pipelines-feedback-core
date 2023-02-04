package feedback

import (
	"context"
	"fmt"
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

func (jx *JXSCMReceiver) Update(ctx context.Context, pipeline contract.PipelineInfo) error {
	jx.sc.Log.Infoln("JXSCMReceiver", pipeline)

	scmCtx := pipeline.GetSCMContext()
	overallStatus := jx.translateStatus(pipeline.GetStatus())

	if jx.client == nil {
		return errors.New("jx.client is nil")
	}

	if jx.client.Commits == nil {
		return errors.New("jx.client.Commits is nil")
	}

	_, response, err := jx.client.Commits.UpdateCommitStatus(ctx, scmCtx.GetNameWithOrg(), scmCtx.Commit, &scm.CommitStatusUpdateOptions{
		ID:          scmCtx.PrId,
		Sha:         pipeline.GetSCMContext().Commit,
		State:       overallStatus,
		Ref:         scmCtx.Reference,
		Name:        pipeline.GetFullName(),
		TargetURL:   pipeline.GetUrl(),
		Description: "",
	})

	if err != nil {
		return errors.Wrap(err, "cannot update commit status")
	}
	if response.Status > 299 {
		return errors.New(fmt.Sprintf("cannot update commit status: got HTTP %v error", response.Status))
	}

	return nil
}

func (jx *JXSCMReceiver) translateStatus(status contract.Status) string {
	switch status {
	case contract.Running:
		return "running"
	case contract.Pending:
		return "pending"
	case contract.Succeeded:
		return "success"
	case contract.Errored:
		return "error"
	case contract.Failed:
		return "failed"
	default:
		return "pending"
	}
}
