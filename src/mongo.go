package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	//"go.mongodb.org/mongo-driver/v2/bson"
	//"go.mongodb.org/mongo-driver/v2/mongo"
	//"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// types

// MongoCLI -> for mongo client instance based on https://www.mongodb.com/docs/drivers/go/current/fundamentals/connections/connection-guide/#std-label-golang-connection-guide
type MongoCLI struct {
	URL  string // url in format mongodb://user:pass@fqdn:port/?connection_options
	CTX  context.Context
	CANC context.CancelFunc
	CLH  *mongo.Client
	OPTS *options.ClientOptions
}

// NewMongo -> constructor
func NewMongo(url string) (mc *MongoCLI, err error) {
	mc = new(MongoCLI)
	mc.URL = url
	mc.CTX, mc.CANC = context.WithTimeout(context.Background(), 5*time.Second)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mc.OPTS = options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
	mc.CLH, err = mongo.Connect(mc.OPTS)
	return
}
