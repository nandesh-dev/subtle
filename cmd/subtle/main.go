package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/nandesh-dev/subtle/internal/routine"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/ent"
)

func main() {
	configDirectoryPath := os.Getenv("CONFIG_DIRECTORY")
	if configDirectoryPath == "" {
		log.Fatal("CONFIG_DIRECTORY environment variable not present")
	}

	config, err := config.Open(filepath.Join(configDirectoryPath, "config.yaml"))
	if err != nil {
		log.Fatal("Cannot open config", err)
	}

	c, err := config.Read()
	if err != nil {
		log.Fatal("Cannot read config", err)
	}

	db, err := ent.Open(c.Database.Path)
	if err != nil {
		log.Fatal("Cannot open database", err)
	}

	if err := db.Schema.Create(context.Background()); err != nil {
		log.Fatal("Cannot migrate database", err)
	}

	routine.Start(config, db)
}
