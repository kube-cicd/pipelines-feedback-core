package feedback

import (
	"context"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract/wiring"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/pkg/errors"
	"strconv"
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

func (jx *JXSCMReceiver) WhenCreated(ctx context.Context, pipeline contract.PipelineInfo) error {
	if pipeline.GetSCMContext().PrId == "" {
		jx.sc.Log.Debug("Missing PR id, skipping")
		return nil
	}
	prId, convErr := strconv.Atoi(pipeline.GetSCMContext().PrId)
	if convErr != nil {
		jx.sc.Log.Warnf("PR id has wrong format, should be an integer, got '%v'", pipeline.GetSCMContext().PrId)
		return nil
	}

	// :white_check_mark: / :x: / :hourglass_flowing_sand:

	pr, _, _ := jx.client.PullRequests.ListComments()

	// Create PR comment
	_, _, err := jx.client.PullRequests.CreateComment(ctx, pipeline.GetSCMContext().GetNameWithOrg(), prId, &scm.CommentInput{
		Body: "The Pipeline " + pipeline.GetStatus().AsHumanReadableDescription(),
	})
	if err != nil {
		return errors.Wrap(err, "cannot update Pull Request status")
	}
	return nil
}

func (jx *JXSCMReceiver) WhenStarted(ctx context.Context, pipeline contract.PipelineInfo) error {
	return nil
}

func (jx *JXSCMReceiver) WhenFinished(ctx context.Context, pipeline contract.PipelineInfo) error {
	return nil
}

func (jx *JXSCMReceiver) UpdateProgress(ctx context.Context, pipeline contract.PipelineInfo) error {
	jx.sc.Log.Infoln("JXSCMReceiver", pipeline)

	scmCtx := pipeline.GetSCMContext()
	ourStatus := pipeline.GetStatus()
	overallStatus := jx.translateStatus(ourStatus)

	var commitStatusErr error = nil
	var prStatusErr error = nil

	if jx.client == nil {
		return errors.New("jx.client is nil")
	}

	if jx.client.Repositories != nil {
		// todo: Create status for multiple stages?
		// todo: Choose between agregated status vs per-stage status?

		_, _, commitStatusErr = jx.client.Repositories.CreateStatus(ctx, pipeline.GetSCMContext().GetNameWithOrg(),
			scmCtx.Commit, &scm.StatusInput{
				State:  overallStatus,
				Label:  "pipelines-feedback",
				Desc:   ourStatus.AsHumanReadableDescription(),
				Target: "", // todo: URL
			},
		)
	} else {
		jx.sc.Log.Warning("jx.client.Repositories is nil. No support for commit status update for this SCM provider in jx go-scm?")
	}

	// Pull/Merge Requests
	if jx.client.PullRequests != nil {
		// todo: Update PR status
	} else {
		jx.sc.Log.Warning("jx.client.PullRequests is nil. No support for manipulating Pull/Merge Requests for this SCM provider in jx go-scm? Cannot operate on Pull/Merge Requests")
	}

	if commitStatusErr != nil {
		return commitStatusErr
	}
	if prStatusErr != nil {
		return prStatusErr
	}
	return nil
}

func (jx *JXSCMReceiver) translateStatus(status contract.Status) scm.State {
	switch status {
	case contract.Running:
		return scm.StateRunning
	case contract.Pending:
		return scm.StatePending
	case contract.Succeeded:
		return scm.StateSuccess
	case contract.Errored:
		return scm.StateError
	case contract.Failed:
		return scm.StateFailure
	default:
		return scm.StateUnknown
	}
}
