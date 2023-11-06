package contract

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"net/url"
	"strings"
	"time"
)

// PipelineInfo is a point-in-time Pipeline status including the SCM information
type PipelineInfo struct {
	ctx JobContext

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

func (pi PipelineInfo) GetSCMContext() JobContext {
	return pi.ctx
}

// GetStatus is calculating the pipeline status basing on the results of all children stages
func (pi PipelineInfo) GetStatus() Status {
	pending := 0
	succeeded := 0
	running := 0
	cancelled := 0
	allStages := len(pi.stages)

	for _, stage := range pi.stages {
		if stage.Status == PipelineErrored {
			return PipelineErrored
		}
		if stage.Status == PipelineFailed {
			return PipelineFailed
		}
		if stage.Status == PipelinePending {
			pending += 1
		}
		if stage.Status == PipelineSucceeded {
			succeeded += 1
		}
		if stage.Status == PipelineRunning {
			running += 1
		}
		if stage.Status == PipelineCancelled {
			cancelled += 1
		}
	}
	if cancelled > 0 {
		return PipelineCancelled
	}
	if running > 0 {
		return PipelineRunning
	}
	if pending > 0 {
		return PipelinePending
	}
	if allStages == succeeded {
		return PipelineSucceeded
	}
	return PipelineErrored
}

func (pi PipelineInfo) ToHash() string {
	sum := fmt.Sprintf("summary=%s\n", pi.GetStatus().AsHumanReadableDescription())
	for _, stage := range pi.stages {
		sum += fmt.Sprintf("stage=%s,status=%s\n", stage.Name, stage.Status)
	}
	hasher := sha256.New()
	hasher.Write([]byte(sum))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetFullName returns a namespace, object name and its instance name (often uid or generated name)
func (pi PipelineInfo) GetFullName() string {
	return pi.namespace + "/" + pi.name + "/" + pi.instanceName
}

// GetInstanceName is returning execution name, a short name
func (pi PipelineInfo) GetInstanceName() string {
	return pi.instanceName
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

// PipelineInfoWithLogsCollector should return whole Pipeline logs on demand. Implement is as lazy-fetch function
func PipelineInfoWithLogsCollector(collector func() string) func(pipelineInfo *PipelineInfo) {
	return func(pipelineInfo *PipelineInfo) {
		pipelineInfo.logs = collector
	}
}

// PipelineInfoWithUrl is setting optionally a URL pointing to a Pipeline visualization
func PipelineInfoWithUrl(url string) func(pipelineInfo *PipelineInfo) {
	return func(pipelineInfo *PipelineInfo) {
		pipelineInfo.url = url
	}
}

func NewPipelineInfo(scm JobContext, namespace string, name string, instanceName string, dateStarted time.Time,
	stages []PipelineStage, labels labels.Labels, annotations labels.Labels, options ...func(info *PipelineInfo)) *PipelineInfo {
	pi := PipelineInfo{
		ctx:          scm,
		name:         name,
		instanceName: instanceName,
		namespace:    namespace,
		dateStarted:  dateStarted,
		stages:       stages,
		labels:       labels,
		annotations:  annotations,
		url:          "",
		logs: func() string {
			return ""
		},
		_logs: "",
	}
	for _, option := range options {
		option(&pi)
	}
	return &pi
}

type Status string

const (
	PipelineRunning   Status = "running"
	PipelineFailed    Status = "failed"
	PipelinePending   Status = "pending"
	PipelineErrored   Status = "errored"
	PipelineSucceeded Status = "succeeded"
	PipelineCancelled Status = "cancelled"
)

func (s Status) IsFinished() bool {
	return s == PipelineFailed || s == PipelineErrored || s == PipelineSucceeded
}

func (s Status) IsRunning() bool {
	return s == PipelineRunning
}

func (s Status) IsErroredOrFailed() bool {
	return s == PipelineFailed || s == PipelineErrored
}

func (s Status) IsSucceeded() bool {
	return s == PipelineSucceeded
}

func (s Status) IsCancelled() bool {
	return s == PipelineCancelled
}

func (s Status) IsNotStarted() bool {
	return s != PipelineRunning && s != PipelineFailed && s != PipelineErrored && s != PipelineSucceeded && s != PipelineCancelled
}

func (s Status) AsHumanReadableDescription() string {
	if s == PipelineRunning || s == PipelinePending {
		return "is " + string(s)
	} else if s == PipelineFailed || s == PipelineErrored || s == PipelineSucceeded || s == PipelineCancelled {
		return string(s)
	}
	return "is in unknown state"
}

// PipelineStage represents a status of a particular Pipeline stage (naming: in Jenkins "Stage", in Tekton "Task")
type PipelineStage struct {
	Name   string
	Status Status
}

type JobContext struct {
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

	// When a job does not have a SCM context, but has a technical-job annotation
	TechnicalJob string
}

func NewSCMContext(repoHttpsUrl string) (JobContext, error) {
	scm := JobContext{}
	scm.RepoHttpsUrl = repoHttpsUrl
	scm.OrganizationName = ""
	scm.RepositoryName = ""

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
	scm.RepositoryName = strings.TrimSuffix(nameSplit[2], ".git")

	return scm, nil
}

func (c JobContext) IsTechnicalJob() bool {
	return c.TechnicalJob != ""
}

func (c JobContext) IsValid() bool {
	if c.TechnicalJob != "" {
		return true
	}
	return (c.Commit != "" || c.PrId != "") && c.RepoHttpsUrl != ""
}

func (c JobContext) GetNameWithOrg() string {
	return c.OrganizationName + "/" + c.RepositoryName
}
