package config_test

import (
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSchemaValidator_ValidateRequestedEntry(t *testing.T) {
	validator := config.SchemaValidator{}
	validator.Add(config.Schema{
		Name:          "kropotkin",
		AllowedFields: []string{"bread", "book"},
	})
	assert.Nil(t, validator.ValidateRequestedEntry("kropotkin", "bread"))
	assert.NotNil(t, validator.ValidateRequestedEntry("kropotkin", "this-does-not-exist"))
}

func TestSchemaValidator_ValidateConfig_Scoped(t *testing.T) {
	validator := config.SchemaValidator{}
	validator.Add(config.Schema{
		Name:          "kropotkin",
		AllowedFields: []string{"bread", "book"},
	})
	assert.Nil(t, validator.ValidateConfig(v1alpha1.Data{
		"kropotkin.bread": "The Conquest of bread",
	}))
	assert.Contains(t, validator.ValidateConfig(v1alpha1.Data{
		"kropotkin.blabla": "This key does not exist",
	}).Error(), "field 'blabla' is not a valid field for component 'kropotkin'")
}

func TestSchemaValidator_ValidateConfig_Global(t *testing.T) {
	validator := config.SchemaValidator{}
	validator.Add(config.Schema{
		Name:          "global",
		AllowedFields: []string{"bread", "book"},
	})
	assert.Nil(t, validator.ValidateConfig(v1alpha1.Data{
		"bread": "The Conquest of bread",
	}))
}
