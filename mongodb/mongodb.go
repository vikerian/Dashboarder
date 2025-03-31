package mongodb

import (

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mc *mongo.Client

// Constructor for our connection
func New(connstr string) (*mongo.Client, error) {
	cli, err := mongo.NewClient(options.Client().ApplyURI(connstr))
	if err != nil {
		return nil, err
	}
	mc = cli
	return cli, nil
}
