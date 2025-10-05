package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

var GlobalConfig *Config

type Config struct {
	MySQL struct {
		DSN string `yaml:"dsn"`
	} `yaml:"mysql"`
}

func Init(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	GlobalConfig = &cfg
	return nil
}
