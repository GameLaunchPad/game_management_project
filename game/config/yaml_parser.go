//go:build !nolint
// +build !nolint

package config

import (
	"gopkg.in/yaml.v3" //nolint:all
)

// parseYAML parses YAML data into the provided interface
//
//nolint:all
func parseYAML(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v) //nolint:all
}
