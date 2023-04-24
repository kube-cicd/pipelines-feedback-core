package config

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/internal/config"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/logging"
)

type ConfigurationProvider struct {
	DocStore config.DocumentStore
	Logger   logging.Logger
}

func (cp *ConfigurationProvider) FetchContextual(namespace string, pipeline contract.PipelineInfo) map[string]string {
	cp.Logger.Debugf("fetchContextual(%s, %s)", namespace, pipeline.GetFullName())
	endMap := make(map[string]string)
	for _, doc := range cp.DocStore.GetForNamespace(namespace) {
		cp.Logger.Debugf("fetchContextual => config '%s' available for this namespace, checking if matches", doc.Name)
		if doc.IsForPipeline(pipeline) {
			cp.Logger.Debugf("fetchContextual(%s, %s) => config '%s'", namespace, pipeline.GetFullName(), doc.Name)
			endMap = mergeMaps(endMap, doc.Data)
		}
	}
	return endMap
}

func mergeMaps(m1 map[string]string, m2 map[string]string) map[string]string {
	merged := make(map[string]string)
	for k, v := range m1 {
		merged[k] = v
	}
	for key, value := range m2 {
		merged[key] = value
	}
	return merged
}
