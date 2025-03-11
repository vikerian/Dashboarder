package main

import (
	"context"
	"dashboarder/config"
	mg "dashboarder/mongo"
//	mqtt "dashboarder/mqtt"
	"fmt"
	"log/slog"
	"os"
	"time"
	//"github.com/k0kubun/pp"
)

var Log *slog.Logger
var ctx context.Context
var CancelFunc context.CancelFunc

// initialization
func init() {
	Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func main() {
	ctx, CancelFunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer CancelFunc()

	// some welcomes and so
	welcomeMsg := fmt.Sprintf("Dashboarder started up from binary %s...", os.Args[0])
	Log.Info(welcomeMsg)
	defer Log.Info("Dashboarder quitting...")

	// config get
	infomsg := "Loading configuration ..."
	Log.Info(infomsg)
	var conf *config.Config
	conf, err := config.GetConfig()
	if err != nil {
		errstr := fmt.Sprintf("Error on reading configuration: %v", err)
		Log.Error(errstr)
		panic(err)
	}

	// for now print config
	confstr := fmt.Sprintf("Configuration: \n %v+ \n", conf)
	Log.Info(confstr)

	// create mongo connection
	infomsg = "Setting up connection to mongo database..."
	Log.Info(infomsg)
	ourmongo, err := mg.New(ctx, conf.MongoDB.Url)
	if err != nil {
		errstr := fmt.Sprintf("Error on connecting to mongo database: %v", err)
		Log.Error(errstr)
		panic(err)
	}

	ourmongo.SetDatabaseAndCollection(conf.MongoDB.DatabaseName, conf.MongoDB.CollectionSTR)

	// now for debug - print mongo setup
	mongoinfostr := fmt.Sprintf("MongoDB client settings: \n %+v \n", ourmongo)
	Log.Info(mongoinfostr)

	// get documents from atlasian mongo db (url MONGODB_URL in environment)
	mongodocs, err := ourmongo.GetAllDocumentsByCollection(ctx, conf.MongoDB.CollectionSTR)
	if err != nil {
		panic(err)
	}
	infomsg = fmt.Sprintf("Documents from %s.%s: \n", conf.MongoDB.DatabaseName, conf.MongoDB.CollectionSTR)
	Log.Info(infomsg)
	infomsg = fmt.Sprintf("%+v", mongodocs)
	Log.Info(infomsg)

	// now try mqtt
}
