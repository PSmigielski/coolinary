// Package log sets up JSON logging on stdout (see docs/service-conventions.md).
package log

import (
	"log/slog"
	"os"
	"strings"
)

func New(service, level string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "msg"
			}
			return a
		},
	})
	logger := slog.New(h).With("service", service)
	slog.SetDefault(logger)
	return logger
}
