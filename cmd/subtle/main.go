package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/nandesh-dev/subtle/internal/jobs"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/ent"
	"github.com/nandesh-dev/subtle/pkgs/logging"
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

	logFile, err := os.OpenFile(c.Logging.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("cannot open log file", err)
	}

	defer logFile.Close()

	logger := logging.New(logging.Options{
		ConsoleLevel:  c.Logging.ConsoleLevel,
		FileLevel:     c.Logging.FileLevel,
		ConsoleWriter: os.Stdout,
		FileWriter:    logFile,
	})

	db, err := ent.Open(c.Database.Path)
	if err != nil {
		log.Fatal("Cannot open database", err)
	}

	if err := db.Schema.Create(context.Background()); err != nil {
		log.Fatal("Cannot migrate database", err)
	}

	jobs.Init(logger, config, db)
}
