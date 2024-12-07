package logging

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
)

func NewRoutineLogger(name string) *slog.Logger {
	return slog.New(tint.NewHandler(os.Stdout, &tint.Options{AddSource: true, Level: slog.LevelDebug.Level()}))
}

func NewManagerLogger(name string) *slog.Logger {
	return slog.New(tint.NewHandler(os.Stdout, &tint.Options{AddSource: true, Level: slog.LevelDebug.Level()}))
}

const DatabaseTransactionCommitError = "cannot commit transaction"
const DatabaseTransactionRollbackError = "cannot rollback transaction"
const DatabaseTransactionCreateError = "cannot create transaction"
