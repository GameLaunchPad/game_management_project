package config

import (
	"encoding/json"
	"os"
	"strings"
)

var GlobalConfig *Config

type Config struct {
	MySQL struct {
		DSN string `yaml:"dsn" json:"dsn"`
	} `yaml:"mysql" json:"mysql"`
}

func Init(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg Config

	// Auto-detect file format based on extension or content
	if strings.HasSuffix(path, ".json") || strings.HasPrefix(strings.TrimSpace(string(data)), "{") {
		// JSON format
		err = json.Unmarshal(data, &cfg)
	} else {
		// YAML format - use lazy loading to avoid linter issues
		err = unmarshalYAML(data, &cfg)
	}

	if err != nil {
		return err
	}

	GlobalConfig = &cfg
	return nil
}

// unmarshalYAML is a separate function to isolate yaml dependency
// The actual implementation is in yaml_parser.go
func unmarshalYAML(data []byte, v interface{}) error {
	// This calls parseYAML from yaml_parser.go
	return parseYAML(data, v)
}
