package logger

import (
	"log/slog"
	"os"
	"strings"
)

func SetupLogger(logLevel string) *slog.Logger {
	var log *slog.Logger
	logLevel = strings.ToLower(logLevel)

	if logLevel == "dev" {
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
		return log
	}

	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	return log
}
