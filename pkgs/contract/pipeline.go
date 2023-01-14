package contract

import "time"

// PipelineInfo is a point-in-time Pipeline status including the SCM information
type PipelineInfo struct {
	ctx SCMContext

	name        string
	dateStarted time.Time
	status      Status
	stages      []PipelineStage
}

func NewPipelineInfo(scm SCMContext, name string, dateStarted time.Time, status Status, stages []PipelineStage) *PipelineInfo {
	return &PipelineInfo{
		ctx:         scm,
		name:        name,
		dateStarted: dateStarted,
		status:      status,
		stages:      stages,
	}
}

type Status string

const (
	Running   Status = "running"
	Failed    Status = "failed"
	Pending   Status = "pending"
	Errored   Status = "errored"
	Succeeded Status = "succeeded"
)

// PipelineStage represents a status of a particular Pipeline stage (naming: in Jenkins "Stage", in Tekton "Task")
type PipelineStage struct {
	Name   string
	Status Status
}

type SCMContext struct {
	// Commit represents long commit hash
	Commit string

	// Reference represents a full GIT reference e.g. refs/heads/v1.6.1 or refs/heads/release-1.3.1.2
	Reference string

	// RepoHttpsUrl describes a full URL to the repository in HTTPS format
	RepoHttpsUrl string

	// PrId is a pull/merge request id
	PrId string
}

func (c SCMContext) IsValid() bool {
	return (c.Commit != "" || c.PrId != "") && c.RepoHttpsUrl != ""
}
