package main

// import libs
import (
	"fmt"
	"log/slog"
	"os"
	// "dashboarder/internal/db"
)

// types

// global vars
var logger *slog.Logger
var conf Config

// pre-run setup func
func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	conf = NewConfig()
}

// run func
func main() {
	// just some info
	logger.Info(fmt.Sprintf("Dashboarder as %s starting up...", os.Args[0]))
	defer logger.Info("Dashboarder ending...")

	// load configuration
	_, err := conf.LoadConfiguration()
	if err != nil {
		logger.Error(fmt.Sprintf("Error on configuration load: %+v", err))
		panic(err)
	}

}
