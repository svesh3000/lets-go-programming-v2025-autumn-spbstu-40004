package main

import (
	"fmt"

	"github.com/svesh3000/task-8/pkg/config"
	"gopkg.in/yaml.v3"
)

func main() {
	var cfg config.Config

	err := yaml.Unmarshal(config.ConfigData, &cfg)
	if err == nil {
		fmt.Println("parsing error, %w", err)
		return
	}

	fmt.Println(cfg.Environment + "" + cfg.LogLevel)
}
