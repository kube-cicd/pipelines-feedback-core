package config

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/logging"
)

func NewData(componentName string, kv map[string]string, validator Validator, logger logging.FatalLogger) Data {
	return Data{
		logger:    logger,
		component: componentName,
		kv:        kv,
		validator: validator,
	}
}

type Data struct {
	logger    logging.FatalLogger
	component string
	kv        map[string]string
	validator Validator
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
		d.logger.Fatalf("code contains undocumented configuration option: %s.%s, please register it in schema so the users will be aware of it", d.component, keyName)
		return ""
	}
	// fetch key, if not exists, then return default
	if d.HasKey(keyName) {
		val, _ := d.kv[keyName]
		return val
	}
	return defaultVal
}
