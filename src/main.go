package main

// import libs
import (
	"fmt"
	"log/slog"
	"os"
	// "dashboarder/internal/db"
)

// mock constants

// TestDbNAME -> mock database name (example db provided by atlas)
const TestDbNAME string = "sample_mflix"

// TestDbCOLLECTION -> mock database collection provided by atlas
const TestDbCOLLECTION string = "movies"

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
	err := conf.LoadConfiguration()
	if err != nil {
		logger.Error(fmt.Sprintf("Error on configuration load: %+v", err))
		logger.Error("We need these env variables setted:")
		logger.Error("MQTT_URL, MQTT_AUTH, VALKEY_URL, VALKEY_AUTH, SIRIDB_URL, SIRIDB_AUTH")
		logger.Error("AUTH format is token or user:pass, for now in base64 encoding...")
		panic(err)
	}
	// debug needed for data:
	//logger.Info(fmt.Sprintf("conf data: %+v", conf))

	// tryout mongo

	// get mongo url from env
	mongoDB, err := NewMongo(conf.MongoURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Error connecting to mongodb: %v", err))
		panic(err)
	}
	// debug info
	//logger.Info(fmt.Sprintf("MongoCLI data: %v", mongoDB))
	// Set mongo db and collection
	mongoDB.SetDBCollection(TestDbNAME, TestDbCOLLECTION)
	// now show setting for mongodb

	logger.Info(fmt.Sprintf("MongoCLI data: %v", mongoDB))
}
