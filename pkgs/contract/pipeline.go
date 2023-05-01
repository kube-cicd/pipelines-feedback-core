package contract

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
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
	retrievalNum int
	labels       labels.Labels
	annotations  labels.Labels
	logs         func() string
	_logs        string
}

// GetId is returning execution ID, unique for a single Pipeline execution
func (pi PipelineInfo) GetId() string {
	return pi.namespace + "/" + pi.name + "/" + pi.instanceName
}

// IsJustCreated tells us if the resource was retrieved from the cluster first time
func (pi PipelineInfo) IsJustCreated() bool {
	return pi.retrievalNum < 2
}

// GetStages returns a stage list with statuses for each
func (pi PipelineInfo) GetStages() []PipelineStage {
	return pi.stages
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

// GetFullName returns a namespace, object name and its instance name (often uid or generated name)
func (pi PipelineInfo) GetFullName() string {
	return pi.namespace + "/" + pi.name + "/" + pi.instanceName
}

// GetUrl returns a URL to some dashboard, where the pipeline could be looked up
func (pi PipelineInfo) GetUrl() string {
	return pi.url
}

// GetName returns a full object name, including namespace
func (pi PipelineInfo) GetName() string {
	return pi.namespace + "/" + pi.name
}

// SetRetrievalCount (for internal use only)
func (pi PipelineInfo) SetRetrievalCount(num int) {
	pi.retrievalNum = num
}

// GetLabels returns Kubernetes object .metadata.labels
func (pi PipelineInfo) GetLabels() labels.Labels {
	return pi.labels
}

// GetAnnotations returns Kubernetes object .metadata.annotations
func (pi PipelineInfo) GetAnnotations() labels.Labels {
	return pi.annotations
}

func (pi PipelineInfo) GetNamespace() string {
	return pi.namespace
}

// GetLogs is returning truncated logs. It is a lazy-loaded method, fetches logs on demand. After first fetch logs are kept in the memory
func (pi PipelineInfo) GetLogs() string {
	if pi._logs == "" {
		pi._logs = pi.logs()
	}
	return pi._logs
}

func NewPipelineInfo(scm SCMContext, namespace string, name string, instanceName string, dateStarted time.Time,
	status Status, stages []PipelineStage, url string, labels labels.Labels, annotations labels.Labels, logs func() string) *PipelineInfo {

	return &PipelineInfo{
		ctx:          scm,
		name:         name,
		instanceName: instanceName,
		namespace:    namespace,
		dateStarted:  dateStarted,
		status:       status,
		stages:       stages,
		labels:       labels,
		annotations:  annotations,
		url:          url,
		logs:         logs,
		_logs:        "",
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

func (s Status) IsFinished() bool {
	return s == Failed || s == Errored || s == Succeeded
}

func (s Status) IsRunning() bool {
	return s == Running
}

func (s Status) IsErroredOrFailed() bool {
	return s == Failed || s == Errored
}

func (s Status) IsSucceeded() bool {
	return s == Succeeded
}

func (s Status) IsNotStarted() bool {
	return s != Running && s != Failed && s != Errored && s != Succeeded
}

func (s Status) AsHumanReadableDescription() string {
	if s == Running || s == Pending {
		return "is " + string(s)
	} else if s == Failed || s == Errored || s == Succeeded {
		return string(s)
	}
	return "is in unknown state"
}

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
