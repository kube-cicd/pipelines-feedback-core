package config

import (
	"github.com/sirupsen/logrus"
)

func NewData(componentName string, kv map[string]string, validator *SchemaValidator) Data {
	return Data{
		component: componentName,
		kv:        kv,
		validator: validator,
	}
}

type Data struct {
	component string
	kv        map[string]string
	validator *SchemaValidator
}

func (d *Data) HasKey(keyName string) bool {
	_, hasKey := d.kv[keyName]
	return hasKey
}

// Get retrieves a configuration value. For non-existing key it returns an empty string
func (d *Data) Get(keyName string) string {
	return d.GetOrDefault(keyName, "")
}

// GetOrDefault retrieves a configuration value. For non-existing keys it returns a default value defined in `defaultVal` parameter
func (d *Data) GetOrDefault(keyName string, defaultVal string) string {
	// validate: if the code is not using an unknown (undocumented) configuration option
	if err := d.validator.ValidateRequestedEntry(d.component, keyName); err != nil {
		// that's a development stage error, should not occur on production build
		logrus.Fatalf("code contains undocumented configuration option: %s.%s, please register it in schema so the users will be aware of it", d.component, keyName)
		return ""
	}
	// fetch key, if not exists, then return default
	val, hasKey := d.kv[keyName]
	if hasKey {
		return val
	}
	return defaultVal
}
