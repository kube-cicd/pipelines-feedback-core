package jxscm

import (
	"context"
	"fmt"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/config"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract/wiring"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/templating"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

const defaultProgressComment = `
:rocket: The Pipeline {{ .pipeline.GetStatus.AsHumanReadableDescription }} {{ if .pipeline.GetStatus.IsNotStarted }}Not started{{ else if .pipeline.GetStatus.IsRunning }}:hourglass_flowing_sand:{{ else if .pipeline.GetStatus.IsErroredOrFailed }}:x:{{ else if .pipeline.GetStatus.IsSucceeded }}:white_check_mark:{{ end }}
--------------------------------------

| Stage | Status |
|-------|--------|
{{- range $stage := .pipeline.GetStages }}
| {{ $stage.Name }} |  {{ if $stage.Status.IsNotStarted }}Pending{{ else if $stage.Status.IsRunning }}:hourglass_flowing_sand:{{ else if $stage.Status.IsErroredOrFailed }}:x:{{ else if $stage.Status.IsSucceeded }}:white_check_mark:{{ end }}  |
{{- end }}
`

const defaultFinishedComment = `
The Pipeline finished with status '{{ .pipeline.GetStatus }}' {{ if .pipeline.GetStatus.IsErroredOrFailed }}:x:{{ else if .pipeline.GetStatus.IsSucceeded }}:white_check_mark:{{ end }}
--------------------
`

const markingBodyPart = `

<details>
    <summary>Build details</summary>
    commentId: {{ .commentId }}
</details>
`

type Receiver struct {
	sc *wiring.ServiceContext
}

// updatePRStatusComment is keeping the PR comment up-to-date with the detailed status of the Pipeline. The comment
//
//	will be created, and then edited multiple times
func (jx *Receiver) updatePRStatusComment(ctx context.Context, cfg config.Data, pipeline contract.PipelineInfo) error {
	// if we are not in context of a PR, then it makes no sense to proceed
	if pipeline.GetSCMContext().PrId == "" {
		return nil
	}

	client, clientErr := jx.createClient(ctx, cfg, pipeline)
	if clientErr != nil {
		return errors.Wrap(clientErr, "cannot create/update PR status comment, SCM client error")
	}

	prId, _ := strconv.Atoi(pipeline.GetSCMContext().PrId)
	markingPart := "(pfc-id=" + pipeline.GetId() + "/updatePRStatusComment)" // we identify a comment by this marking

	// 1. Find existing comment in the cache
	commentId := jx.sc.Store.GetStatusPRCommentId(pipeline)

	// 2. If not, then search through the comments in the PR
	if commentId == "" {
		commentId = jx.findCommentIdByMarking(ctx, markingPart, pipeline, prId, client)
	}

	// 2.1. Check cache - skip if last status is the same as current (then we do not need to edit anything)
	if jx.sc.Store.GetLastRecordedPipelineStatus(pipeline) == string(pipeline.GetStatus()) {
		jx.sc.Log.Debugf("Skipping update, status already wrote to SCM for '%s'", pipeline.GetId())
		return nil
	}

	content, tplErr := templating.TemplateProgressComment(
		createTemplate(cfg.GetOrDefault("progress-comment", defaultProgressComment)),
		pipeline,
		markingPart,
	)
	if tplErr != nil {
		return errors.Wrap(tplErr, "cannot create a comment from template")
	}

	// 3. Create new comment
	if commentId == "" {
		_, _, createErr := client.PullRequests.CreateComment(ctx, pipeline.GetSCMContext().GetNameWithOrg(), prId, &scm.CommentInput{
			Body: content,
		})
		if createErr != nil {
			return errors.Wrap(createErr, "cannot create a comment on a Pull Request")
		}
		jx.sc.Store.RecordInfoAboutLastComment(pipeline, commentId)
	} else {
		// 4. Update existing comment
		commentIdInt, _ := strconv.Atoi(commentId)
		_, _, editErr := client.PullRequests.EditComment(ctx, pipeline.GetSCMContext().GetNameWithOrg(), prId, commentIdInt, &scm.CommentInput{
			Body: content,
		})
		if editErr != nil {
			return errors.Wrap(editErr, "cannot edit existing comment on a Pull Request")
		}
		jx.sc.Store.RecordInfoAboutLastComment(pipeline, commentId)
	}
	return nil
}

func (jx *Receiver) InitializeWithContext(sc *wiring.ServiceContext) error {
	sc.Log.Info("Initializing JX SCM Receiver")
	jx.sc = sc
	return nil
}

func (jx *Receiver) WhenCreated(ctx context.Context, pipeline contract.PipelineInfo) error {
	return nil
}

func (jx *Receiver) WhenStarted(ctx context.Context, pipeline contract.PipelineInfo) error {
	return nil
}

// WhenFinished is creating a final comment on the PR to make sure user is notified about the final status
func (jx *Receiver) WhenFinished(ctx context.Context, pipeline contract.PipelineInfo) error {
	// Skip if not in a PR context
	if pipeline.GetSCMContext().PrId == "" {
		return nil
	}

	cfg := jx.sc.Config.FetchContextual("jxscm", pipeline.GetNamespace(), pipeline)
	client, clientErr := jx.createClient(ctx, cfg, pipeline)
	if clientErr != nil {
		return errors.Wrap(clientErr, "cannot create/update PR status comment, SCM client error")
	}

	prId, _ := strconv.Atoi(pipeline.GetSCMContext().PrId)
	markingPart := "(pfc-id=" + pipeline.GetId() + "/WhenFinished)"

	// Do not send the same comment twice
	if jx.sc.Store.WasSummaryCommentCreated(pipeline) {
		jx.sc.Log.Debugf("Skipping update, status already written to SCM for '%s'", pipeline.GetId())
		return nil

	} else {
		// Fallback - in case there was no cache
		if jx.findCommentIdByMarking(ctx, markingPart, pipeline, prId, client) != "" {
			jx.sc.Log.Debugf("Skipping update, status already written to SCM for '%s'", pipeline.GetId())
			return nil
		}
	}

	// Template a comment body
	content, tplErr := templating.TemplateSummaryComment(
		createTemplate(cfg.GetOrDefault("finished-comment", defaultFinishedComment)),
		pipeline,
		markingPart,
	)
	if tplErr != nil {
		return errors.Wrap(tplErr, "cannot create a comment from template")
	}

	// Send comment to SCM
	_, _, createErr := client.PullRequests.CreateComment(ctx, pipeline.GetSCMContext().GetNameWithOrg(), prId, &scm.CommentInput{
		Body: content,
	})
	if createErr != nil {
		return errors.Wrap(createErr, "cannot create a comment on a Pull Request")
	}
	jx.sc.Store.RecordSummaryCommentCreated(pipeline)
	return nil
}

// UpdateProgress is keeping commit & PR up-to-date with the progress by creating & updating statuses
func (jx *Receiver) UpdateProgress(ctx context.Context, pipeline contract.PipelineInfo) error {
	cfg := jx.sc.Config.FetchContextual("jxscm", pipeline.GetNamespace(), pipeline)
	client, clientErr := jx.createClient(ctx, cfg, pipeline)
	if clientErr != nil {
		return errors.Wrap(clientErr, "cannot create/update PR status comment, SCM client error")
	}

	scmCtx := pipeline.GetSCMContext()
	ourStatus := pipeline.GetStatus()
	overallStatus := jx.translateStatus(ourStatus)

	var commitStatusErr error = nil
	var prCommentStatusErr error = nil

	if commentStatusErr := jx.updatePRStatusComment(ctx, cfg, pipeline); commentStatusErr != nil {
		prCommentStatusErr = errors.Wrap(commentStatusErr, "cannot create/update status comment in Pull Request")
	}

	// Update commit status
	if client.Repositories != nil {
		_, _, commitStatusErr = client.Repositories.CreateStatus(ctx, pipeline.GetSCMContext().GetNameWithOrg(),
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

	if commitStatusErr != nil {
		return commitStatusErr
	}
	if prCommentStatusErr != nil {
		return prCommentStatusErr
	}
	return nil
}

func (jx *Receiver) translateStatus(status contract.Status) scm.State {
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

func (jx *Receiver) CanHandle(name string) bool {
	return name == jx.GetImplementationName()
}

func (jx *Receiver) GetImplementationName() string {
	return "jxscm"
}

func createTemplate(bodyTemplate string) string {
	return bodyTemplate + markingBodyPart
}

// findCommentIdByMarking finds a SCM comment id in a pull request by looking for a text (markingPart) in a comment body. Returns first occurrence
func (jx *Receiver) findCommentIdByMarking(ctx context.Context, markingPart string, pipeline contract.PipelineInfo, prId int, client *scm.Client) string {
	comments, _, _ := client.PullRequests.ListComments(ctx, pipeline.GetSCMContext().GetNameWithOrg(), prId, &scm.ListOptions{Size: 100})
	for _, comment := range comments {
		if strings.Contains(comment.Body, markingPart) {
			return fmt.Sprintf("%v", comment.ID)
		}
	}
	return ""
}

func (jx *Receiver) createClient(ctx context.Context, data config.Data, pipeline contract.PipelineInfo) (*scm.Client, error) {
	// will first try to fetch GIT token from "jxscm.token" (plaintext in configuration)
	// fallbacks to looking for a `kind: Secret` specified by name in "jxscm.token-secret-name", and there it will look for a key specified by "jxscm.token-secret-key"
	gitToken, err := jx.sc.Config.FetchFromFieldOrSecret(ctx, &data, pipeline.GetNamespace(), "token", "token-secret-key", "token-secret-name")
	if err != nil {
		return nil, errors.Wrap(err, "cannot create a JX SCM client - cannot fetch a GIT token neither from 'jxscm.token' as plaintext neither from a `kind: Secret` referenced in 'jxscm.token-secret-name'")
	}
	if gitToken == "" {
		return nil, errors.New("cannot create a JX SCM client - cannot fetch a GIT token neither from 'jxscm.token' as plaintext neither from a `kind: Secret` referenced in 'jxscm.token-secret-name'")
	}

	// constructs a client
	return newClientFromConfig(data, gitToken)
}
