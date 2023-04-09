package feedback

import (
	"context"
	"fmt"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract/wiring"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/templating"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

const defaultProgressComment = `
:rocket: The Pipeline {{ .pipeline.GetStatus.AsHumanReadableDescription }} {{ if .pipeline.GetStatus.IsNotStarted }}Not started{{ else if .pipeline.GetStatus.IsRunning }}:hourglass_flowing_sand:{{ else if .pipeline.GetStatus.IsErroredOrFailed }}:x:{{ else if .pipeline.GetStatus.IsSucceeded }}:white_check_mark:{{ end }}

| Stage | Status |
|-------|--------|
{{- range $stage := .pipeline.GetStages }}
| {{ $stage.Name }} |  {{ if $stage.Status.IsNotStarted }}Pending{{ else if $stage.Status.IsRunning }}:hourglass_flowing_sand:{{ else if $stage.Status.IsErroredOrFailed }}:x:{{ else if $stage.Status.IsSucceeded }}:white_check_mark:{{ end }}  |
{{- end }}

<details>
    <summary>Build</summary>
    Build: {{ .buildId }}
</details>
`

type JXSCMReceiver struct {
	client *scm.Client
	sc     *wiring.ServiceContext
}

// updatePRStatusComment is keeping the PR comment up-to-date with the detailed status of the Pipeline. The comment
//
//	will be created, and then edited multiple times
//
// todo: check in cache if the last status == current status - if yes, then do not update twice
func (jx *JXSCMReceiver) updatePRStatusComment(ctx context.Context, pipeline contract.PipelineInfo) error {
	// :white_check_mark: / :x: / :hourglass_flowing_sand:
	idPart := "(pfc-comment-id=" + pipeline.GetId() + ")"

	if pipeline.GetSCMContext().PrId == "" {
		return nil
	}
	prId, _ := strconv.Atoi(pipeline.GetSCMContext().PrId)

	// 1. Find existing comment in the cache
	commentId := jx.sc.Store.GetStatusPRCommentId(pipeline)

	// 2. If not, then search through the comments in the PR
	if commentId == "" {
		comments, _, _ := jx.client.PullRequests.ListComments(ctx, pipeline.GetSCMContext().GetNameWithOrg(), prId, &scm.ListOptions{Size: 100})
		for _, comment := range comments {
			if strings.Contains(comment.Body, idPart) {
				commentId = fmt.Sprintf("%v", comment.ID)
				break
			}
		}
	}

	// 2.1. Check cache - skip if last status is the same as current (then we do not need to edit anything)
	if jx.sc.Store.GetLastRecordedPipelineStatus(pipeline) == string(pipeline.GetStatus()) {
		jx.sc.Log.Debugf("Skipping update, status already wrote to SCM for '%s'", pipeline.GetId())
		return nil
	}

	content, tplErr := templating.TemplateProgressComment(defaultProgressComment, pipeline, idPart)
	if tplErr != nil {
		return errors.Wrap(tplErr, "cannot create a comment from template")
	}

	// 3. Create new comment
	if commentId == "" {
		_, _, createErr := jx.client.PullRequests.CreateComment(ctx, pipeline.GetSCMContext().GetNameWithOrg(), prId, &scm.CommentInput{
			Body: content,
		})
		if createErr != nil {
			return errors.Wrap(createErr, "cannot create a comment on a Pull Request")
		}
		jx.sc.Store.RecordInfoAboutLastComment(pipeline, commentId)
	} else {
		// 4. Update existing comment
		commentIdInt, _ := strconv.Atoi(commentId)
		_, _, editErr := jx.client.PullRequests.EditComment(ctx, pipeline.GetSCMContext().GetNameWithOrg(), prId, commentIdInt, &scm.CommentInput{
			Body: content,
		})
		if editErr != nil {
			return errors.Wrap(editErr, "cannot edit existing comment on a Pull Request")
		}
		jx.sc.Store.RecordInfoAboutLastComment(pipeline, commentId)
	}
	return nil
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

	if commentStatusErr := jx.updatePRStatusComment(ctx, pipeline); commentStatusErr != nil {
		return errors.Wrap(commentStatusErr, "cannot create/update status comment in Pull Request")
	}

	// Update commit status
	if jx.client.Repositories != nil {
		_, _, commitStatusErr = jx.client.Repositories.CreateStatus(ctx, pipeline.GetSCMContext().GetNameWithOrg(),
			scmCtx.Commit, &scm.StatusInput{
				State:  overallStatus,
				Label:  "Pipeline - " + pipeline.GetFullName(),
				Desc:   ourStatus.AsHumanReadableDescription(),
				Target: pipeline.GetUrl(),
			},
		)
	} else {
		jx.sc.Log.Warning("jx.client.Repositories is nil. No support for commit status update for this SCM provider in jx go-scm?")
	}

	// Pull/Merge Requests merge status
	if jx.client.PullRequests != nil {
		// todo: Update PR status. Make it not able to merge due to running Pipeline
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
