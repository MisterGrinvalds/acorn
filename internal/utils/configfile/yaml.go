package configfile

import (
	"gopkg.in/yaml.v3"
)

// YAMLWriter implements the Writer interface for YAML format.
// Used for k9s config, Kubernetes manifests, and similar YAML config files.
type YAMLWriter struct{}

func init() {
	Register(&YAMLWriter{})
}

// Format returns the format identifier.
func (w *YAMLWriter) Format() string {
	return "yaml"
}

// Write generates YAML content from values.
// The values map is directly marshaled to YAML.
func (w *YAMLWriter) Write(values map[string]interface{}) ([]byte, error) {
	return yaml.Marshal(values)
}
