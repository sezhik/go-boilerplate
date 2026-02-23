package configure

import (
	"log/slog"
	"os"
)

func MustLog() *slog.Logger {
	lvl := os.Getenv("LOG_LEVEL")

	var level slog.Level
	switch lvl {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelError
	}

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(l)

	return l
}
