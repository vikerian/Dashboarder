package main

import (
	"dashboarder/config"
	"fmt"
	"log/slog"
	"os"
)

var Log *slog.Logger

func init() {
	//slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.SetLogLoggerLevel(slog.LevelDebug)
	Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func main() {
	Log.Info("Dashboarder starting up...")
	slog.Debug("Getting configuration...")
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	dbgMsg := fmt.Sprintf("Config loaded: %+v", cfg)
	slog.Debug(dbgMsg)
	slog.Debug("Configuration loaded, connecting to DBs...")

}
