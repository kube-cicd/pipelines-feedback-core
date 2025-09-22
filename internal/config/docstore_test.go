package config_test

import (
	"testing"

	"github.com/kube-cicd/pipelines-feedback-core/internal/config"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/fake"
	"github.com/stretchr/testify/assert"
)

func TestIndexedDocumentStoreFlow_Namespaced(t *testing.T) {
	store := config.CreateIndexedDocumentStore(&fake.NullValidator{})

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
	store := config.CreateIndexedDocumentStore(&fake.NullValidator{})

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

func TestGetForNamespace_SortsByPriorityWeightAscending(t *testing.T) {
	store := config.CreateIndexedDocumentStore(&fake.NullValidator{})

	cfgLow := v1alpha1.NewPFConfig()
	cfgLow.Name = "low"
	cfgLow.Namespace = "test"
	cfgLow.Spec.PriorityWeight = 10

	cfgHigh := v1alpha1.NewPFConfig()
	cfgHigh.Name = "high"
	cfgHigh.Namespace = "test"
	cfgHigh.Spec.PriorityWeight = 100

	cfgMid := v1alpha1.NewPFConfig()
	cfgMid.Name = "mid"
	cfgMid.Namespace = "test"
	cfgMid.Spec.PriorityWeight = 50

	store.Push(&cfgHigh)
	store.Push(&cfgLow)
	store.Push(&cfgMid)

	result := store.GetForNamespace("test")
	assert.Equal(t, 3, len(result))
	assert.Equal(t, "low", result[0].Name)
	assert.Equal(t, "mid", result[1].Name)
	assert.Equal(t, "high", result[2].Name)
}

func TestGetForNamespace_DefaultPriorityWeightIsZero(t *testing.T) {
	store := config.CreateIndexedDocumentStore(&fake.NullValidator{})

	cfgZero := v1alpha1.NewPFConfig()
	cfgZero.Name = "zero"
	cfgZero.Namespace = "test"
	// PriorityWeight is not specified, should default to 0

	cfgTen := v1alpha1.NewPFConfig()
	cfgTen.Name = "ten"
	cfgTen.Namespace = "test"
	cfgTen.Spec.PriorityWeight = 10

	store.Push(&cfgTen)
	store.Push(&cfgZero)

	result := store.GetForNamespace("test")
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "zero", result[0].Name)
	assert.Equal(t, "ten", result[1].Name)
}
