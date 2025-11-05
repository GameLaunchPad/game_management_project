package config

import (
	"bufio"
	"encoding/json"
	"fmt"
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

	// YAML format - use simple built-in parser (no external dependencies)
	return parseSimpleYAML(data)
}

// parseSimpleYAML parses a simple YAML format without external dependencies
// Supports basic key-value pairs, sufficient for config files
func parseSimpleYAML(data []byte) error {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	var currentSection string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for section headers (e.g., "rpc:")
		if strings.HasSuffix(line, ":") && !strings.Contains(line, " ") {
			currentSection = strings.TrimSuffix(line, ":")
			continue
		}

		// Parse key-value pairs
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				// Remove quotes if present
				value = strings.Trim(value, `"'`)

				// Set values based on section and key
				if currentSection == "rpc" {
					if key == "game_service_addr" {
						Config.Rpc.GameServiceAddr = value
					} else if key == "cp_center_service_addr" {
						Config.Rpc.CpCenterServiceAddr = value
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading YAML: %w", err)
	}

	return nil
}
