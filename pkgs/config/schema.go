package config

import (
	"fmt"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	"github.com/pkg/errors"
	"strings"
)

type SchemaValidator struct {
	schema map[string]Schema
}

func (sm *SchemaValidator) Add(schema Schema) {
	if len(sm.schema) == 0 {
		sm.schema = make(map[string]Schema)
	}
	sm.schema[schema.Name] = schema
}

func (sm *SchemaValidator) ValidateRequestedEntry(group string, key string) error {
	if _, exists := sm.schema[group]; !exists {
		return errors.New(fmt.Sprintf("component '%s' has no registered its schema", group))
	}
	grp := sm.schema[group]
	for _, fieldName := range grp.AllowedFields {
		if key == fieldName {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("field '%s' is not a valid field for component '%s'", key, group))
}

func (sm *SchemaValidator) ValidateConfig(data v1alpha1.Data) error {
	for fullKey, _ := range data {
		group, key := parseEntryIntoGroupAndKey(fullKey)
		if err := sm.ValidateRequestedEntry(group, key); err != nil {
			return err
		}
	}
	return nil
}

func parseEntryIntoGroupAndKey(entry string) (string, string) {
	parts := strings.Split(entry, ".")
	key := parts[len(parts)-1:][0]

	// remove key from parts
	parts = parts[:len(parts)-1]

	group := strings.Join(parts, ".")
	if group == "" {
		group = "global"
	}

	return group, key
}

type Schema struct {
	Name          string
	AllowedFields []string
}
