package main

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func main() {
	configPath := flag.String("config", "", "path to config file (required)")
	flag.Parse()
	if *configPath == "" {
		panic("ERROR: Flag --config is required!")
	}

	config, err := loadConfig(*configPath)
	if err != nil {
		panic(fmt.Sprintf("ERROR: Failed to load config: %v", err))
	}

	fmt.Println(config.InputFile)
	fmt.Println(config.OutputFile)
}

func loadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("cannot parse YAML config: %w", err)
	}

	return &config, nil
}
