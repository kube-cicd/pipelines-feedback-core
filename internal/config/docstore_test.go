package config_test

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/internal/config"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIndexedDocumentStoreFlow_Namespaced(t *testing.T) {
	store := config.CreateIndexedDocumentStore(&config.NullValidator{})

	// No document stored yet
	assert.Empty(t, store.GetForNamespace("social"))

	// Namespace: social
	pfc := v1alpha1.NewPFConfig()
	pfc.Name = "conquest-of-bread"
	pfc.Namespace = "social"
	pfc.Data = v1alpha1.Data{
		"book": "bread",
	}
	assert.Nil(t, store.Push(&pfc))

	// Namespace: anti-social
	pfcAS := v1alpha1.NewPFConfig()
	pfcAS.Name = "conquest-of-money"
	pfcAS.Namespace = "anti-social"
	pfcAS.Data = v1alpha1.Data{
		"book": "dollar",
	}
	assert.Nil(t, store.Push(&pfcAS))

	// After Push() it should be returned by GetForNamespace()
	assert.NotEmpty(t, store.GetForNamespace("social"))

	// Namespace: social should have only one entry
	assert.Equal(t, 1, len(store.GetForNamespace("social")))

	// Namespace: anti-social also should have only one entry
	assert.Equal(t, 1, len(store.GetForNamespace("anti-social")))

	// Deletion of entry in "anti-social" namespace should result in emptying the namespace
	store.Delete("anti-social", "conquest-of-money")
	assert.Empty(t, store.GetForNamespace("anti-social"))
	assert.NotEmpty(t, store.GetForNamespace("social")) // "social" should still not be empty
}

func TestIndexedDocumentStoreFlow_ClusterScope(t *testing.T) {
	store := config.CreateIndexedDocumentStore(&config.NullValidator{})

	// Namespace: social
	pfc := v1alpha1.NewPFConfig()
	pfc.Name = "conquest-of-bread"
	pfc.Namespace = "social"
	pfc.Data = v1alpha1.Data{
		"book": "bread",
	}
	assert.Nil(t, store.Push(&pfc))

	// No namespace: global
	pfcG := v1alpha1.NewPFConfig()
	pfcG.Name = "graeber"
	pfcG.Data = v1alpha1.Data{
		"book": "debt",
	}
	assert.Nil(t, store.Push(&pfcG))

	assert.Equal(t, 2, len(store.GetForNamespace("social")))  // there we have 1 namespaced + 1 global
	assert.Equal(t, 1, len(store.GetForNamespace("default"))) // there we didn't push any PFConfig
}
