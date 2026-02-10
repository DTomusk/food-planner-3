package logging

import (
	"log/slog"
	"os"
)

func New(env string) *slog.Logger {
	var handler slog.Handler

	// Default to warn level in production, debug elsewhere
	opts := &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// In non-prod, set log level to debug
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
