package config

import (
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
	"github.com/pkg/errors"
)

func CreateMultipleCollector(collectors []ConfigurationCollector, logger *logging.InternalLogger) *MultipleCollector {
	return &MultipleCollector{collectors: collectors, logger: logger}
}

// MultipleCollector is an "adapter of adapters" pattern that lets you use multiple configuration sources at a single time
type MultipleCollector struct {
	collectors []ConfigurationCollector
	logger     *logging.InternalLogger
}

func (mc *MultipleCollector) SetLogger(logger *logging.InternalLogger) {
	mc.logger = logger
}

func (mc *MultipleCollector) CanHandle(adapterName string) bool {
	return true
}

func (mc *MultipleCollector) GetImplementationName() string {
	return "multiple"
}

func (mc *MultipleCollector) CollectInitially() ([]*v1alpha1.PFConfig, error) {
	all := make([]*v1alpha1.PFConfig, 0)
	for _, collector := range mc.collectors {
		collected, err := collector.CollectInitially()
		if err != nil {
			return all, errors.Wrapf(err, "one of configuration collectors - '%v' failed", collector.GetImplementationName())
		}
		all = append(all, collected...)
	}
	return all, nil
}

func (mc *MultipleCollector) CollectOnRequest(info contract.PipelineInfo) ([]*v1alpha1.PFConfig, error) {
	all := make([]*v1alpha1.PFConfig, 0)
	for _, collector := range mc.collectors {
		collected, err := collector.CollectOnRequest(info)
		if err != nil {
			return all, errors.Wrapf(err, "one of configuration collectors - '%v' failed", collector.GetImplementationName())
		}
		all = append(all, collected...)
	}
	return all, nil
}
