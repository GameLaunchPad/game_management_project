package config

import (
	"bufio"
	"encoding/json"
	"fmt"
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
		// YAML format - use simple built-in parser (no external dependencies)
		err = parseSimpleYAML(data, &cfg)
	}

	if err != nil {
		return err
	}

	GlobalConfig = &cfg
	return nil
}

// parseSimpleYAML parses a simple YAML format without external dependencies
// Supports basic key-value pairs, sufficient for config files
func parseSimpleYAML(data []byte, cfg *Config) error {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	var currentSection string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for section headers (e.g., "mysql:")
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
				if currentSection == "mysql" && key == "dsn" {
					cfg.MySQL.DSN = value
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading YAML: %w", err)
	}

	return nil
}
