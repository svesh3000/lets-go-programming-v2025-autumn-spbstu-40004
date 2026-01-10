package main

import (
	"fmt"

	"github.com/svesh3000/task-8/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("load config error, %w", err)

		return
	}

	fmt.Print(cfg.Environment + " " + cfg.LogLevel)
}
