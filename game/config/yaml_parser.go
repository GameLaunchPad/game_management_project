//go:build !nolint
// +build !nolint

package config

import (
	"gopkg.in/yaml.v3"
)

// parseYAML parses YAML data into the provided interface
func parseYAML(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
