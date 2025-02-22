package main

import (
	"dashboarder/config"
	"dashboarder/mongo"
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
	infomsg := "Loading configuration ..."
	log.Info(infomsg)
	var conf *config.Config
	conf, err := config.GetConfig()
	if err != nil {
		errstr := fmt.Sprintf("Error on reading configuration: %v", err)
		log.Error(errstr)
		panic(err)
	}

	// for now print config
	confstr := fmt.Sprintf("Configuration: \n %v+ \n", conf)
	log.Info(confstr)

	// create mongo connection
	infomsg = "Setting up connection to mongo database..."
	log.Info(infomsg)
	ourmongo, err := mongo.New(conf.MongoDB.Url)
	if err != nil {
		errstr := fmt.Sprintf("Error on connecting to mongo database: %v", err)
		log.Error(errstr)
		panic(err)
	}

	ourmongo.SetDatabaseAndCollection(conf.MongoDB.DatabaseName, conf.MongoDB.CollectionSTR)

	// now for debug - print mongo setup
	mongoinfostr := fmt.Sprintf("MongoDB client settings: \n %+v \n", ourmongo)
	log.Info(mongoinfostr)
}
