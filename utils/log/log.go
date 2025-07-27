package log

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func InitAsJson() {
	addSource, err := strconv.ParseBool(os.Getenv("LOG_ADD_SOURCE"))
	if err != nil {
		addSource = false // default to false if parsing fails
	}

	var level slog.Level
	levelStr := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(levelStr) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource, // includes function, file and line number in log entries (default: false)
	})

	slog.SetDefault(slog.New(handler))
}
