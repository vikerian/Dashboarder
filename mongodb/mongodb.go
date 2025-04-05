package mongodb

import (

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	URL            string
	Clh            *mongo.Client
	Ctx            context.Context
	Cancel         context.CancelFunc
	DatabaseName   string
	CollectionName string
	Cursor         *mongo.Collection
}

type Document struct {
	ID      string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string    `json:"name" bson:"name"`
	Content []byte    `json:"content" bson:"content"`
	Tags    []string  `json:"tags,omitempty" bson:"tags,omitempty"`
	Created time.Time `json:"created" bson:"created"`
}

// Constructor for our connection
func New(connstr string) (mdb *MongoDB, err error) {
	mdb = new(MongoDB)
	mdb.URL = connstr
	mdb.Ctx, mdb.Cancel = context.WithTimeout(context.Background(), 5*time.Second)
	clOpts := options.Client().ApplyURI(connstr)
	if err != nil {
		return nil, err
	}
	// create client and connect
	mdb.Clh, err = mongo.Connect(mdb.Ctx, clOpts)
	if err != nil {
		return nil, err
	}
	// check connection
	err = mdb.Clh.Ping(mdb.Ctx, nil)
	if err != nil {
		return nil, err
	}

	return mdb, nil
}

// SetDBbCollection -> set which db and connection we are using right now
func (mdb *MongoDB) SetDBCollection(databaseName string, colname string) {
	mdb.DatabaseName = databaseName
	mdb.CollectionName = colname
	mdb.Cursor = mdb.Clh.Database(databaseName).Collection(colname)
}

// InsertDoc
func (mdb *MongoDB) InsertDoc(doc []byte) (saveID string, err error) {

	return
}
