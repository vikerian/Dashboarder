package main

import (
	"io"
	"log/slog"
	"os"
)

var applog *slog.Logger

func init() {
	lf, err := os.OpenFile("application.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0640)
	if err != nil {
		panic(err)
	}
	// multiwriter
	mw := io.MultiWriter(os.Stdout, lf)
	// our logger setup
	applog = slog.New(slog.NewJSONHandler(mw, nil))
}

func main() {
	// first info about start / end
	applog.Info("Dashboarder starting up...")
	defer applog.Info("Dashboarder ending...")
}
