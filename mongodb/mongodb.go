package mongodb

import (

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"

	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Clh    *mongo.Client
	Ctx    context.Context
	Cancel context.CancelFunc
}

// Constructor for our connection
func New(connstr string) (*MongoDB, error) {
	mc = new(MongoDB)
	cli, err := mongo.NewClient(options.Client().ApplyURI(connstr))
	if err != nil {
		return nil, err
	}
	mc.Clh = cli

	return mc, nil
}

// Connect
func (mdb *MongoDB) Connect() (ok bool, err error) {

	return
}
