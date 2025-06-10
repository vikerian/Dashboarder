package main

//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"

import (
	"context"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	DSN      string
	Clh      *mongo.Client
	Cursor   *mongo.Collection
	Ctx      context.Context
	CancFunc context.CancelFunc
}

// NewMongoClient -> instance of mongodb client connected to mongodb
func NewMongClient(dsn string) (*MongoDB, error) {
	m := new(MongoDB)
	// first parse url params
	u, err := url.Parse(dsn) // u = net.URL struct
	if err != nil {
		return nil, err
	}
	// create context and store it with its cancel function
	ctx, canc := context.WithDeadline(contex.Background, 5*time.Second)
    :x
	m.CancFunc = canc

	// create client, set it up and verify connection
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dsn).SetServerAPIOptions(serverAPI)
	clh, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}
	// verify connection
	err = clh.Ping(m.Ctx, nil)
	if err != nil {
		return nil, err
	}

}
