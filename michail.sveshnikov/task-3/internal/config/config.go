package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(configPath string) (*Config, error) {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("cannot parse YAML config: %w", err)
	}

	if config.InputFile == "" {
		return nil, fmt.Errorf("input-file field is required in config")
	}

	if config.OutputFile == "" {
		return nil, fmt.Errorf("output-file field is required in config")
	}

	_, err = os.Stat(config.InputFile)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("input file does not exist: %s", config.InputFile)
	}

	return &config, nil
}
