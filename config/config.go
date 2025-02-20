package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* Import libs */

/* Types */
type Config struct {
	ServerCfg struct {
		ListenIP   string
		ListenPort uint
		MaxWorker  uint
	}

	SiriDB struct {
		Database      string
		Host          string
		Port          uint
		Username      string
		Password      string
		AdminUsername string
		AdminPassword string
	}

	MongoDB struct {
		Url         string
		Options     *options.ClientOptions
		Collections []*mongo.Collection
		Client      *mongo.Client
	}

	OpenWeather struct {
		Longitude float64
		Latitude  float64
		Token     string
	}
	CTX    context.Context
	Cancel context.CancelFunc
}

/* Global vars */

/* functions and methods */

// New -> create clean new instance of configuration
func New() *Config {
	return &Config{}
}

// GetConfig -> read configuration, returns setuped instance of configuration struct
// for now mock setting for devel, next read from environment vars
func GetConfig() (*Config, error) {
	cfg := new(Config)
	cfg.CTX, cfg.Cancel = context.WithTimeout(context.Background(), 5*time.Second)
	cfg.ServerCfg.ListenIP = "0.0.0.0"
	cfg.ServerCfg.ListenPort = 3500
	cfg.SiriDB.Host = "localhost"
	cfg.SiriDB.Port = 9000
	cfg.SiriDB.Username = "iris"
	cfg.SiriDB.Password = "siri"
	cfg.SiriDB.AdminUsername = "sa"
	cfg.SiriDB.AdminPassword = "siri"
	cfg.SiriDB.Database = "devel"
	cfg.MongoDB.Url = "mongodb://localhost:27017"
	cfg.MongoDB.Options = options.Client().ApplyURI(cfg.MongoDB.Url)
	client, err := mongo.Connect(cfg.CTX, cfg.MongoDB.Options)
	if err != nil {
		return nil, err
	}
	cfg.MongoDB.Client = client
	return cfg, nil
}
