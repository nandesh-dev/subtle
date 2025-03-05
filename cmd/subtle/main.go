package main

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nandesh-dev/subtle/internal/api"
	"github.com/nandesh-dev/subtle/internal/jobs"
	"github.com/nandesh-dev/subtle/pkgs/configuration"
	"github.com/nandesh-dev/subtle/pkgs/ent"
	"github.com/nandesh-dev/subtle/pkgs/env"
	"github.com/nandesh-dev/subtle/pkgs/logging"
)

func main() {
	logFile, err := os.OpenFile(env.LogFilepath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("cannot open log file", err)
	}

	defer logFile.Close()

	logger := logging.New(logging.Options{
		ConsoleLevel:  env.ConsoleLogLevel(),
		FileLevel:     env.FileLogLevel(),
		ConsoleWriter: os.Stdout,
		FileWriter:    logFile,
	})

	logger.Info("reading config file")
	configFile, err := configuration.Open(env.ConfigFilepath())
	if err != nil {
		logger.Error("cannot open config file", "err", err, "path", env.ConfigFilepath())
		return
	}

	logger.Info("opening database file")
	db, err := ent.Open(env.DatabaseFilepath())
	if err != nil {
		logger.Error("cannot open database file", "err", err, "database_filepath", env.DatabaseFilepath())
		return
	}

	logger.Info("migrating database")
	if err := db.Schema.Create(context.Background()); err != nil {
		logger.Error("cannot migrate database", "err", err)
		return
	}

	logger.Info("setting up jobs in database")
	if err := jobs.SetupDatabase(db); err != nil {
		logger.Error("error setting up jobs in database", "err", err)
		return
	}

	logger.Info("running jobs")
	go jobs.StartJobRunTicker(context.Background(), logger, configFile, db)

	logger.Info("creating api server")
	apiServer := api.NewAPIServer(context.Background(), configFile, db, api.APIServerOptions{
		EnableGRPCReflection: env.EnableGRPCReflection(),
	})

	if err := apiServer.ListenAndServe(env.WebServerAddress()); err != nil {
		logger.Error("cannot start the api server", "err", err)
	}
}
