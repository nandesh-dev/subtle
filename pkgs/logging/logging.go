package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type Options struct {
	ConsoleLevel slog.Level
	FileLevel    slog.Level

	ConsoleWriter *os.File
	FileWriter    *os.File
}

func New(opt Options) *slog.Logger {
	return slog.New(
		&handler{
			level: min(opt.ConsoleLevel, opt.FileLevel),
			handlers: []slog.Handler{
				tint.NewHandler(opt.ConsoleWriter, &tint.Options{AddSource: true, Level: opt.ConsoleLevel}),
				slog.NewJSONHandler(opt.FileWriter, &slog.HandlerOptions{AddSource: true, Level: opt.FileLevel}),
			},
		},
	)
}

type handler struct {
	level    slog.Level
	handlers []slog.Handler
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r); err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, 0, len(h.handlers))

	for _, handler := range h.handlers {
		newHandlers = append(newHandlers, handler.WithAttrs(attrs))
	}

	return &handler{
		level:    h.level,
		handlers: newHandlers,
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, 0, len(h.handlers))

	for _, handler := range h.handlers {
		newHandlers = append(newHandlers, handler.WithGroup(name))
	}

	return &handler{
		level:    h.level,
		handlers: newHandlers,
	}
}

const DatabaseTransactionCommitError = "cannot commit transaction"
const DatabaseTransactionRollbackError = "cannot rollback transaction"
const DatabaseTransactionCreateError = "cannot create transaction"
