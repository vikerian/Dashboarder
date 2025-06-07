package main

import (
	"context"
	"fmt"
	"time"

	//	"go.mongodb.org/mongo-driver/bson"
	//	"go.mongodb.org/mongo-driver/bson/primitive"
	//	"go.mongodb.org/mongo-driver/mongo"
	//	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseType represents different types of db
type DatabaseType int

const (
	MongoDB DatabaseType = iota
	SiriDB
	ValkeyDB
)

// DataItem represents a generic data item to be stored
type DataItem struct {
	Key interface{}
	Data interface{}
}

// MongoDOC represents a document to store to MongoDB
type MongoDODC struct {
	Name string `bson:"_name"`

// DatabaseConnection is an interface for different database connection
type DatabaseConnection interface {
	Connect(string,string,string) error
	Disconnect() error
	Store(data interface{}) (interface{}, error)
	Get(string) (DataItem, error)
}

// MongoDBConnection implements DatabaseConnection for MongoDB
type MongoDBConnection struct {
	URI        string
	Database   string
	Collection string
	Ctx        *context.Context
	Cnc        *context.CancelFunc
	Clh        *mongo.Client
	Cursor	   *mongo.Collection
}

// NewDataItem -> create dataitem instance, if specified key, it will be used in Store, otherwise key will be filled in Store method
func NewDataItem(key interface{}, data interface{}) (*DataItem) {
	return &DataItem{
		Key: key,
		Data: data,
	}
}

// Connect - create connection to mongodb and verify it
func (m *MongoDBConnection) Connect(uri string, dbname string, collection string) (err error) {
	m.URI = uri
	m.Database = dbname
	m.Collection = collection
	fmt.Println("Connecting  to MongoDB:", m.URI)
	ctx, canc := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	m.Ctx = &ctx
	m.Cnc = &canc
	// create client with url - connection setup
	m.clh,err = mongo.Connect(m.Ctx, options.Client().ApplyURI(m.URI))
	if err != nil {
		err = errors.New(fmt.Sprintf("Error on connection to MongoDB: %v", err))
		return 
	}
	// verify connection as is
	err = m.Clh.Ping(m.Ctx, nil)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error on verifying connection to MongoDB: %v", err))
	}
	//create cursor
	col, err := m.Clh.Database(m.Database).Collection(m.Collection)
	m.Cursor = col
	return 
}

// Disconnect -> disconnect actual mongo connection
func (m *MongoDBConnection) Disconnect() (err error) {
	fmt.Println("Disconnecting from MongoDB")
	err = m.Clh.Disconnect(m.Ctx)
	return 
}

// Store data into mongodb, returns key (_id) of documen /  error
func (m *MongoDBConnection) Store(data DataItem) (key interface{}, err error) {

	return 
}

func (m *MongoDBConnection) Get(key interface{}) (data DataItem, err error) {


	return nil
}
