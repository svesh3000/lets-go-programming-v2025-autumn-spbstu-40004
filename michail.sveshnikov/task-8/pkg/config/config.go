package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (*Config, error) {
	var cfg Config

	err := yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}

	return &cfg, nil
}
