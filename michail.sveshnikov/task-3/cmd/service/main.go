package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/svesh3000/task-3/internal/centralbank"
	"github.com/svesh3000/task-3/internal/config"
)

func main() {
	configPath := flag.String("config", "", "path to config file (required)")
	flag.Parse()
	if *configPath == "" {
		panic("ERROR: Flag --config is required!")
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(fmt.Sprintf("ERROR: Failed to load config: %v", err))
	}

	dir := filepath.Dir(cfg.OutputFile)
	os.MkdirAll(dir, os.ModePerm)

	// Шаг 2: Чтение и парсинг XML файла с валютами от ЦБ РФ
	currencies, err := centralbank.ParseXMLFile(cfg.InputFile)
	if err != nil {
		panic(fmt.Sprintf("ERROR: Failed to parse XML file: %v", err))
	}

	fmt.Println(currencies[0].CharCode)
}
