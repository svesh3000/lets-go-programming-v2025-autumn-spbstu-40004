package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrConfigFileNotExist = errors.New("config file does not exist")
	ErrInputFileRequired  = errors.New("input-file field is required in config")
	ErrOutputFileRequired = errors.New("output-file field is required in config")
	ErrInputFileNotExist  = errors.New("input file does not exist")
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(configPath string) (*Config, error) {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", ErrConfigFileNotExist, configPath)
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
		return nil, ErrInputFileRequired
	}

	if config.OutputFile == "" {
		return nil, ErrOutputFileRequired
	}

	_, err = os.Stat(config.InputFile)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", ErrInputFileNotExist, config.InputFile)
	}

	return &config, nil
}
