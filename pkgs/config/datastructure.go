package config

type Data map[string]string

func (d Data) HasKey(keyName string) bool {
	_, hasKey := d[keyName]
	return hasKey
}

// Get retrieves a configuration value. For non-existing key it returns an empty string
func (d Data) Get(keyName string) string {
	return d.GetOrDefault(keyName, "")
}

// GetOrDefault retrieves a configuration value. For non-existing keys it returns a default value defined in `defaultVal` parameter
func (d Data) GetOrDefault(keyName string, defaultVal string) string {
	val, hasKey := d[keyName]
	if hasKey {
		return val
	}
	return defaultVal
}
