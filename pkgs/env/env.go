package env

import (
	"log"
	"log/slog"
	"os"
)

const (
	configFilepathKey   = "SUBTLE_CONFIG_FILEPATH"
	databaseFilepathKey = "SUBTLE_DATABASE_FILEPATH"
	logFilepathKey      = "SUBTLE_LOG_FILEPATH"
	fileLogLevelKey     = "SUBTLE_FILE_LOG_LEVEL"
	consoleLogLevelKey  = "SUBTLE_CONSOLE_LOG_LEVEL"
)

func ConfigFilepath() string {
	path := os.Getenv(configFilepathKey)

	if path == "" {
		log.Fatalf("config filepath not set using \"%s\"", configFilepathKey)
	}

	return path
}

func DatabaseFilepath() string {
	path := os.Getenv(databaseFilepathKey)

	if path == "" {
		log.Fatalf("database filepath not set using \"%s\"", databaseFilepathKey)
	}

	return path
}

func LogFilepath() string {
	path := os.Getenv(logFilepathKey)

	if path == "" {
		log.Fatalf("log filepath not set using \"%s\"", logFilepathKey)
	}

	return path
}

func FileLogLevel() slog.Level {
	rawLogLevel := os.Getenv(fileLogLevelKey)

	if rawLogLevel == "" {
		log.Fatalf("file log level not set using \"%s\"", fileLogLevelKey)
	}

	switch rawLogLevel {
	case "INFO":
		return slog.LevelInfo
  case "WARN":
    return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	case "DEBUG":
		return slog.LevelDebug
	}

	log.Fatalf("invalid file log level; %s", rawLogLevel)
	return slog.LevelInfo
}

func ConsoleLogLevel() slog.Level {
	rawLogLevel := os.Getenv(consoleLogLevelKey)

	if rawLogLevel == "" {
		log.Fatalf("console log level not set using \"%s\"", consoleLogLevelKey)
	}

	switch rawLogLevel {
	case "INFO":
		return slog.LevelInfo
  case "WARN":
    return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	case "DEBUG":
		return slog.LevelDebug
	}

	log.Fatalf("invalid console log level; %s", rawLogLevel)
	return slog.LevelInfo
}
