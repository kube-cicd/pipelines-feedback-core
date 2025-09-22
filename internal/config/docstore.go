package config

import (
	"sort"

	"github.com/kube-cicd/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	"github.com/pkg/errors"
)

func CreateIndexedDocumentStore(cfgSchema Validator) IndexedDocumentStore {
	return IndexedDocumentStore{
		namespaces: make(map[string]NamespacedDocuments, 0),
		global:     make(map[string]*v1alpha1.PFConfig, 0),
		cfgSchema:  cfgSchema,
	}
}

// IndexedDocumentStore is a storage for configuration files structured as CRD
//
//	Every PFConfig document has meta attributes like Job selector, namespace
//	so a ConfigurationService can serve a contextual documentation - in context of a Job
//	or in global context.
type IndexedDocumentStore struct {
	namespaces map[string]NamespacedDocuments
	global     map[string]*v1alpha1.PFConfig
	cfgSchema  Validator
}

type NamespacedDocuments map[string]*v1alpha1.PFConfig

func (ds *IndexedDocumentStore) GetForNamespace(namespace string) []*v1alpha1.PFConfig {
	// first provide global configuration
	docsForNs := make([]*v1alpha1.PFConfig, 0)
	for _, doc := range ds.global {
		docsForNs = append(docsForNs, doc)
	}
	// then provide namespaced configuration
	if docs, exists := ds.namespaces[namespace]; exists {
		for _, doc := range docs {
			docsForNs = append(docsForNs, doc)
		}
	}
	// sort by PriorityWeight ascending
	sort.SliceStable(docsForNs, func(i, j int) bool {
		return docsForNs[i].Spec.PriorityWeight < docsForNs[j].Spec.PriorityWeight
	})
	return docsForNs
}

// Push is adding or overwriting a document in IndexedDocumentStore
func (ds *IndexedDocumentStore) Push(cfg *v1alpha1.PFConfig) error {
	if vErr := ds.cfgSchema.ValidateConfig(cfg.Data); vErr != nil {
		return errors.Wrap(vErr, "cannot push PFConfig to IndexedDocumentStore")
	}

	// global
	if cfg.IsGlobalConfiguration() {
		ds.global[cfg.Name] = cfg
		return nil
	}

	// namespaced
	nsName := cfg.Namespace
	if _, exists := ds.namespaces[nsName]; !exists {
		ds.namespaces[nsName] = NamespacedDocuments{}
	}

	ns := ds.namespaces[nsName]
	ns[cfg.Name] = cfg
	ds.namespaces[nsName] = ns
	return nil
}

// Delete is deleting an element from IndexedDocumentStore
func (ds *IndexedDocumentStore) Delete(namespace string, name string) {
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
