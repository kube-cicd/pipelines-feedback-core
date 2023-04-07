package store

import (
	"fmt"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Operator struct {
	Store
}

// CountHowManyTimesKubernetesResourceReceived returns count and increases the counter for given resource
func (o *Operator) CountHowManyTimesKubernetesResourceReceived(retrieved *contract.PipelineInfo) int {
	ident := retrieved.GetStoreId() + "/counter"
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
	if setErr := o.Set(ident, fmt.Sprintf("%v", counter)); setErr != nil {
		logrus.Error("cannot save to store", setErr)
	}
	return counter
}

func (o *Operator) WasEventAlreadySent(retrieved contract.PipelineInfo, eventType string) bool {
	ident := retrieved.GetStoreId() + "/" + eventType
	existing, err := o.Get(ident)
	if err != nil && err.Error() == ErrNotFound {
		return false
	}
	return existing == "true"
}

func (o *Operator) RecordEventFiring(retrieved contract.PipelineInfo, eventType string) error {
	ident := retrieved.GetStoreId() + "/" + eventType
	if err := o.Set(ident, "true"); err != nil {
		return errors.Wrap(err, "cannot store information, that event was already fired - RecordEventFiring()")
	}
	return nil
}
