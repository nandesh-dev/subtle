package logging

import (
	"log/slog"
	"os"
)

func NewRoutineLogger(name string) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})).With("type", "routine", "name", name)
}

func NewManagerLogger(name string) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})).With("type", "manager", "name", name)
}
