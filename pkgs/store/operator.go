package store

import (
	"fmt"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strconv"
)

const StatusCacheTtl = 86400 * 30
const StatusLongCacheTtl = 86400 * 365 * 10

// SecretCacheTtl is used to not fetch the same `kind: Secret` multiple times during at least one Pipeline lifecycle
const SecretCacheTtl = 120

type Operator struct {
	Store
}

// CountHowManyTimesKubernetesResourceReceived returns count and increases the counter for given resource
func (o *Operator) CountHowManyTimesKubernetesResourceReceived(retrieved *contract.PipelineInfo) int {
	ident := retrieved.GetId() + "/counter"
	existing, err := o.Get(ident)
	counter := 0

	if err == nil {
		c, cErr := strconv.Atoi(existing)
		if cErr != nil {
			c = 0
		}
		counter = c
	}
	counter += 1
	if setErr := o.Set(ident, fmt.Sprintf("%v", counter), StatusCacheTtl); setErr != nil {
		logrus.Error("cannot save to store", setErr)
	}
	return counter
}

func (o *Operator) WasEventAlreadySent(retrieved contract.PipelineInfo, eventType string) bool {
	ident := retrieved.GetId() + "/" + eventType
	existing, err := o.Get(ident)
	if err != nil && err.Error() == ErrNotFound {
		return false
	}
	return existing == "true"
}

func (o *Operator) RecordEventFiring(retrieved contract.PipelineInfo, eventType string) error {
	ident := retrieved.GetId() + "/" + eventType
	if err := o.Set(ident, "true", StatusLongCacheTtl); err != nil {
		return errors.Wrap(err, "cannot store information, that event was already fired - RecordEventFiring()")
	}
	return nil
}

func (o *Operator) GetStatusPRCommentId(pipeline contract.PipelineInfo) string {
	return o.readOrEmpty(pipeline, "PRCommentId")
}

func (o *Operator) GetLastRecordedPipelineStatus(pipeline contract.PipelineInfo) string {
	return o.readOrEmpty(pipeline, "PRLastStatus")
}

func (o *Operator) RecordInfoAboutLastComment(pipeline contract.PipelineInfo, commentId string) {
	_ = o.Set(pipeline.GetId()+"/PRCommentId", commentId, StatusCacheTtl)
	_ = o.Set(pipeline.GetId()+"/PRLastStatus", string(pipeline.GetStatus()), StatusCacheTtl)
}

func (o *Operator) RecordSummaryCommentCreated(pipeline contract.PipelineInfo) {
	_ = o.Set(pipeline.GetId()+"/PRSummaryCreated", "true", StatusLongCacheTtl)
}

func (o *Operator) WasSummaryCommentCreated(pipeline contract.PipelineInfo) bool {
	value, _ := o.Get(pipeline.GetId() + "/PRSummaryCreated")
	return value == "true"
}

func (o *Operator) readOrEmpty(pipeline contract.PipelineInfo, key string) string {
	ident := pipeline.GetId() + "/" + key
	existing, err := o.Get(ident)
	if err != nil && err.Error() == ErrNotFound {
		return ""
	}
	return existing
}

func (o *Operator) GetConfigSecretKey(namespace string, refKey string, refSecretKey string) string {
	ident := namespace + "/" + refKey + "/" + refSecretKey
	existing, _ := o.Get(ident)
	return existing
}

func (o *Operator) PushConfigSecretKey(namespace string, refKey string, refSecretKey string, val string) {
	ident := namespace + "/" + refKey + "/" + refSecretKey
	_ = o.Set(ident, val, SecretCacheTtl)
}
