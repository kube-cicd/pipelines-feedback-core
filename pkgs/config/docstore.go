package config

import "github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"

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
