package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Config struct {
	Rpc struct {
		GameServiceAddr     string `yaml:"game_service_addr"`
		CpCenterServiceAddr string `yaml:"cp_center_service_addr"`
	} `yaml:"rpc"`
}

func Init(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&Config); err != nil {
		return err
	}
	return nil
}
