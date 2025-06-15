package main

import (
	"context"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	//"golang.org/x/net/context"
)

const (
	APPNAME string = "Dashboarder/v0.01"
	COLNAME string = "dshbrcol"
	DBNAME  string = "DashboarderDB"
	TIMEOUT        = 5 * time.Second
)

//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"

type mongoDB struct {
	URL              url.URL
	DSN              string
	DatabaseNAME     string
	CollectionNAME   string
	Client           *mongo.Client
	Cursor           *mongo.Collection
	ClientCTX        context.Context
	ClientCancelFUNC context.CancelFunc
}

type MongoClient interface {
	Connect(dsn string) error
	Disconnect()
	SaveDoc(doc interface{}) (interface{}, error)
	GetDoc(key interface{}) (interface{}, error)
}

// NewMongoClient -> returns interface instance
func NewMongoClient() MongoClient {
	m := &mongoDB{
		DatabaseNAME:   DBNAME,
		CollectionNAME: COLNAME,
	}
	return m
}

// Connect - connect to db and verify connection
func (m *mongoDB) Connect(dsn string) error {
	// parse connection url
	url, err := url.Parse(dsn)
	if err != nil {
		return err
	}
	// setup our instance
	m.URL = *url
	m.DSN = dsn
	ctx, canc := context.WithTimeout(context.Background(), TIMEOUT)
	m.ClientCTX = ctx
	m.ClientCancelFUNC = canc

	// connect

	return err
}

func (m *mongoDB) Disconnect() {
	m.Client.Disconnect(m.ClientCTX)
}
