package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nandesh-dev/subtle/internal/routine"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/ent"
)

func initilize() error {
	configDirectoryPath := os.Getenv("CONFIG_DIRECTORY")
	if configDirectoryPath == "" {
		return fmt.Errorf("CONFIG_DIRECTORY environment variable not present")
	}

	if err := config.Init(configDirectoryPath); err != nil {
		return fmt.Errorf("Failed to initilize config: %v", err)
	}

	return nil
}

func main() {
	if err := initilize(); err != nil {
		log.Fatal(err)
	}

	db, err := ent.Open(config.Config().Database.Path)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Schema.Create(context.Background()); err != nil {
		log.Fatal(err)
	}

	routine.Start(db)
}
