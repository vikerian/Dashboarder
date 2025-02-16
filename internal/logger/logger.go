package logger

import (
	"log/slog"
	"os"
)

var instance *slog.Logger

// GetLogger returns a singleton instance of logger
func GetLogger() *slog.Logger {
	if instance == nil {
		instance = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return instance
}
