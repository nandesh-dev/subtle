package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nandesh-dev/subtle/internal/routine"
	"github.com/nandesh-dev/subtle/internal/server"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/logger"
)

func initilize() error {
	configDirectoryPath := os.Getenv("CONFIG_DIRECTORY")
	if configDirectoryPath == "" {
		return fmt.Errorf("CONFIG_DIRECTORY environment variable not present")
	}

	if err := config.Init(configDirectoryPath); err != nil {
		return fmt.Errorf("Failed to initilize config: %v", err)
	}

	if err := database.Init(); err != nil {
		return fmt.Errorf("Failed to initilize database: %v", err)
	}

	logger.Init()
	return nil
}

func main() {
	if err := initilize(); err != nil {
		log.Fatal(err)
	}

	go routine.Start()

	server.New().Listen()
}
