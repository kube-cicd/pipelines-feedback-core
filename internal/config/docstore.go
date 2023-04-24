package config

import "github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"

func CreateDocumentStore() DocumentStore {
	return DocumentStore{
		namespaces: make(map[string]NamespacedDocuments, 0),
		global:     make(map[string]*v1alpha1.PFConfig, 0),
	}
}

// DocumentStore is a storage for configuration files structured as CRD
//
//	Every PFConfig document has meta attributes like Job selector, namespace
//	so a ConfigurationService can serve a contextual documentation - in context of a Job
//	or in global context.
type DocumentStore struct {
	namespaces map[string]NamespacedDocuments
	global     map[string]*v1alpha1.PFConfig
}

type NamespacedDocuments map[string]*v1alpha1.PFConfig

// Push is adding or overwriting a document in DocumentStore
func (ds *DocumentStore) Push(cfg *v1alpha1.PFConfig) {
	// global
	if cfg.IsGlobalConfiguration() {
		ds.global[cfg.Name] = cfg
		return
	}

	// namespaced
	nsName := cfg.Namespace
	if _, exists := ds.namespaces[nsName]; !exists {
		ds.namespaces[nsName] = NamespacedDocuments{}
	}

	ns := ds.namespaces[nsName]
	ns[cfg.Name] = cfg
	ds.namespaces[nsName] = ns
}

// Delete is deleting an element from DocumentStore
func (ds *DocumentStore) Delete(namespace string, name string) {
	// namespaced
	if namespace != "" {
		if _, exists := ds.namespaces[namespace]; exists {
			ns := ds.namespaces[namespace]
			if _, docExists := ns[name]; docExists {
				delete(ds.namespaces[namespace], name)
				return
			}
		}
	}
	// global
	if _, globalExists := ds.global[name]; globalExists {
		delete(ds.global, name)
		return
	}
}
