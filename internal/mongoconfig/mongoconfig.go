package mongoconfig

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Url        string
	Ctx        context.Context
	Cancel     context.CancelFunc
	Options    *options.ClientOptions
	Client     *mongo.Client
	Collection *mongo.Collection
}

// GetMongoConfig - read mongo configuration from environment (after all we will be in kubernetes!)
func GetMongo() *MongoConfig {
	mgc := new(MongoConfig)
	ctx, canc := context.WithTimeout(context.Background(), 5*time.Second)
	mgc.Ctx = ctx
	mgc.Cancel = canc
	mgc.Url = os.Getenv("MONGO_URL")
	mgc.Options = options.Client().ApplyURI(mgc.Url)
	return mgc
}
