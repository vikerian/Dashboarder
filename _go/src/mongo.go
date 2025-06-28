package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	DSN              string // dsn format: scheme://user:password@host:port/database
	DatabaseNAME     string
	CollectionNAME   string
	Client           *mongo.Client
	Collection       *mongo.Collection
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
	// create url for connection
	uri := fmt.Sprintf("%s://%s:%s/", m.URL.Scheme, m.URL.Host, m.URL.Port())
	opts := options.Client().ApplyURI(uri)
	if m.URL.User != nil {

	}
	return err
}

func (m *mongoDB) Disconnect() {
	m.Client.Disconnect(m.ClientCTX)
}

func (m *mongoDB) GetDoc(key interface{}) (doc interface{}, err error) {

	return
}

func (m *mongoDB) SaveDoc(doc interface{}) (key interface{}, err error) {
	return
}
