package main

import (
	"dashboarder/config"
	"fmt"
	"log/slog"
	"os"
	//"github.com/k0kubun/pp"
)

var log *slog.Logger

// initialization
func init() {
	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func main() {

	// some welcomes and so
	welcomeMsg := fmt.Sprintf("Dashboarder started up from binary %s...", os.Args[0])
	log.Info(welcomeMsg)
	defer log.Info("Dashboarder quitting...")

	// config get
	conf, err := config.GetConfig()
	if err != nil {
		errstr := fmt.Sprintf("Error on reading configuration: %v", err)
		log.Error(errstr)
		panic(err)
	}

	// for now print config
	confstr := fmt.Sprintf("Configuration: \n %v+ \n", conf)
	log.Info(confstr)

}
