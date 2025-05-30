package main

import (
	"log/slog"
	"os"
)

var sl *slog.Logger

func initLog() *slog.Logger {
	lg := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	sl = lg
	return lg
}
