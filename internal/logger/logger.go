package logger

import (
	"log/slog"
	"os"
)

func Init() {
	w := os.Stdout
	options := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(w, options)

	logger := slog.New(handler)

	slog.SetDefault(logger)

	slog.Debug("Logger initialized")
}
