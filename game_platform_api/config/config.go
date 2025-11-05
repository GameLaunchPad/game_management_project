package config

import (
	"encoding/json"
	"os"
	"strings"
)

var Config struct {
	Rpc struct {
		GameServiceAddr     string `yaml:"game_service_addr" json:"game_service_addr"`
		CpCenterServiceAddr string `yaml:"cp_center_service_addr" json:"cp_center_service_addr"`
	} `yaml:"rpc" json:"rpc"`
}

func Init(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Auto-detect file format based on extension or content
	if strings.HasSuffix(path, ".json") || strings.HasPrefix(strings.TrimSpace(string(data)), "{") {
		// JSON format
		return json.Unmarshal(data, &Config)
	}

	// YAML format - use lazy loading
	return unmarshalYAML(data, &Config)
}

// unmarshalYAML is a separate function to isolate yaml dependency
// The actual implementation is in yaml_parser.go
func unmarshalYAML(data []byte, v interface{}) error {
	// This calls parseYAML from yaml_parser.go
	return parseYAML(data, v)
}
