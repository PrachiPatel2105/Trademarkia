package logger

import (
	"log/slog"
	"os"
)

// New creates a new structured logger
func New() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}
