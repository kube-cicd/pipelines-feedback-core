package config

import "github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"

type DocumentStore struct {
	namespaces map[string]NamespacedDocuments
}

type NamespacedDocuments map[string][]v1alpha1.PFConfig
