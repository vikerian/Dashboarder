package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Url         string
	Options     *options.ClientOptions
	Collections []*mongo.Collection
	Client      *mongo.Client
}

func New() *MongoDB {
	mg := new(MongoDB)
	return mg
}
