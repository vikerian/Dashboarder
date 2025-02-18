package main

/* import libs */
import (
	"fmt"
	"log/slog"
	"os"

	"dashboarder/config"

	"github.com/k0kubun/pp/v3"
)

/* global vars (no fuj, ale pro singleton ideal) */
var log *slog.Logger


/* init func */
func init() {
	log = slog.New(slog.NewJSONHandler(os.Stdout,nil))
}


/* global startup function */
func main() {
	/* just some welcome/... */
	welcomeMsg := fmt.Sprintf("Dashboarder asi binary %s is starting up...", os.Args[0])
	log.Info(welcomeMsg)
	defer log.Info("Dashboarder ending")

	/* get configuration */
	conf,err := config.GetConfig()
	if err != nil {
		panic (err)
	}

	log.Info(pp.Print(conf))
}
