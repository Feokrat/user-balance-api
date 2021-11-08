package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Feokrat/user-balance-api/internal/config"
)

const configFile = "configs/config"

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	cfg, err := config.Init(configFile, logger)
	if err != nil {
		logger.Fatalf("failed to load application configuration: %s", err)
	}

	fmt.Println(cfg)
}
