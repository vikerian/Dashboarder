package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
)

// ConfigurationFILE - what configuration file we use
const ConfigurationFILE string = "config.json"

// AppConfig - application configuration data
type DBAppConfig struct {
	Redis struct {
		Username string `json:"redis_username"` // default none
		Password string `json:"redis_password"` // default none
		Host     string `json:"redis_host"`     // default localhost
		Port     uint16 `json:"redis_port"`     // default 6379
		DBIndex  uint16 `json:"redis_database_index"`
	} `json:"redis"`
	Mongo struct {
		Username   string `json:"mongo_username"`
		Password   string `json:"mongo_password"`
		Host       string `json:"mongo_host"`          // default localhost
		Port       string `json:"mongo_port"`          // default 27017
		DBName     string `json:"mongo_db_name"`       // default "appdb"
		Collection string `json:"mongo_db_collection"` // default dashboard_collection
	} `json:"mongo"`
	SiriDB struct {
		Host     string `json:"siri_host"`
		Port     uint16 `json:"siri_port"`
		Username string `json:"siri_username"` // default "iris"
		Password string `json:"siri_password"` // default siri
		DBName   string `json:"siri_db"`       // default dbtest
	} `json:"siridb"`
}

type OpenWeatherConfig struct {
	Id                uuid.UUID `json:"uuid"`
	LocationName      string    `json:"location_name"`
	LocationLongitude float64   `json:"location_longitude"`
	LocationLatitude  float64   `json:"location_latitude"`
}

/*
type dataSource struct {
	ID     uint32    `json:"datasource_id"`
	UUID   uuid.UUID `json:"datasource_unique_id,omitempty"`
	Name   string    `json:"datasource_name"`
	URL    string    `json:"datasource_url"`
	Setted bool      `json:"datasource_already_setted"`
}
*/

// ReadConfig - read configuration file, returns AppConfig structure with values or error
func ReadDBConfig(confFile string) (*DBAppConfig, error) {
	var ConfFilePath string
	// first which conf file we use
	if confFile == "" {
		ConfFilePath = ConfigurationFILE
	} else {
		ConfFilePath = confFile
	}

	// read configuration - open file - with check
	cf, err := os.Open(ConfFilePath)
	if err != nil {
		errstr := fmt.Sprintf("Error on opening configuration file %s: %v", ConfFilePath, err)
		return nil, errors.New(errstr)
	}
	defer cf.Close()
	// read configuration - read whole file into config
	cfg := new(DBAppConfig)
	cfgData, err := io.ReadAll(cf)

	if err != nil {
		errstr := fmt.Sprintf("Error on reading configuration data: %v", err)
		return nil, errors.New(errstr)
	}
	err = json.Unmarshal(cfgData, cfg)
	if err != nil {
		errstr := fmt.Sprintf("Error on decoding configuration: %v", err)
		return nil, errors.New(errstr)
	}
	// else
	// now uuid data sources
	/*for pos := range cfg.DataSources {
		uid, err := uuid.NewRandom()
		if err != nil {
			errstr := fmt.Sprintf("Error on uuid generation for data sources: %v", err)
			return nil, errors.New(errstr)
		}
		cfg.DataSources[pos].UUID = uid
		cfg.DataSources[pos].Setted = true
	}*/
	return cfg, nil
}
