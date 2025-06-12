package main

import (
	"net/url"

	"go.mongodb.org/mongo-driver/mongo"
)

//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"

type mongoDB struct {
	URL    url.URL
	DSN    string
	Client mongo.Client
	Cursor mongo.Collection
}

type MongoClient interface {
	Connect(dsn string) MongoClient
	Disconnect()
}
