package contract

import (
	"github.com/pkg/errors"
	"net/url"
	"strings"
	"time"
)

// PipelineInfo is a point-in-time Pipeline status including the SCM information
type PipelineInfo struct {
	ctx SCMContext

	name         string
	instanceName string
	namespace    string
	dateStarted  time.Time
	status       Status
	stages       []PipelineStage
	url          string
}

func (pi PipelineInfo) GetSCMContext() SCMContext {
	return pi.ctx
}

// GetStatus is calculating the pipeline status basing on the results of all children stages
func (pi PipelineInfo) GetStatus() Status {
	pending := 0
	succeeded := 0
	running := 0
	allStages := len(pi.stages)

	for _, stage := range pi.stages {
		if stage.Status == Errored {
			return Errored
		}
		if stage.Status == Failed {
			return Failed
		}
		if stage.Status == Pending {
			pending += 1
		}
		if stage.Status == Succeeded {
			succeeded += 1
		}
		if stage.Status == Running {
			running += 1
		}
	}
	if pending > 0 {
		return Pending
	}
	if allStages == succeeded {
		return Succeeded
	}
	if running > 0 {
		return Running
	}
	return Errored
}

func (pi PipelineInfo) GetFullName() string {
	return pi.namespace + "/" + pi.instanceName
}

func (pi PipelineInfo) GetUrl() string {
	return pi.url
}

func (pi PipelineInfo) GetName() string {
	return pi.namespace + "/" + pi.name
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

	OrganizationName string
	RepositoryName   string
}

func NewSCMContext(repoHttpsUrl string) (SCMContext, error) {
	scm := SCMContext{}

	// not matching, not containing the annotation
	if len(repoHttpsUrl) == 0 {
		return scm, nil
	}

	u, err := url.Parse(repoHttpsUrl)
	if err != nil {
		return scm, errors.Wrap(err, "not a valid url")
	}

	nameSplit := strings.Split(u.Path, "/")

	if len(nameSplit) < 3 {
		return scm, errors.New("repository url does not contain valid organization and repository names")
	}

	scm.OrganizationName = nameSplit[1]
	scm.RepositoryName = nameSplit[2]

	return scm, nil
}

func (c SCMContext) IsValid() bool {
	return (c.Commit != "" || c.PrId != "") && c.RepoHttpsUrl != ""
}

func (c SCMContext) GetNameWithOrg() string {
	return c.OrganizationName + "/" + c.RepositoryName
}
