package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	//"go.mongodb.org/mongo-driver/v2/bson"
	//"go.mongodb.org/mongo-driver/v2/mongo"
	//"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// types

// MongoCLI -> for mongo client instance based on https://www.mongodb.com/docs/drivers/go/current/fundamentals/connections/connection-guide/#std-label-golang-connection-guide
type MongoCLI struct {
	Url            string // url in format mongodb://user:pass@fqdn:port/?connection_options
	Ctx            context.Context
	Cancel         context.CancelFunc
	Clh            *mongo.Client
	DatabaseName   string
	Database       *mongo.Database
	CollectionName string
	Collection     *mongo.Collection
}

// NewMongoCLI -> constructor
func NewMongoCLI(url string) (mc *MongoCLI, err error) {
	mc = new(MongoCLI)
	mc.Url = url
	mc.Ctx, mc.Cancel = context.WithTimeout(context.Background(), 5*time.Second)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	OPTS := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
	mc.Clh, err = mongo.Connect(OPTS)
	if err != nil {
		return nil, err
	}
	err = mc.Clh.Ping(mc.Ctx, readpref.Primary())
	return
}

// CloseMongo -> destructor
func (mc *MongoCLI) CloseMongo() {
	mc.Clh.Disconnect(mc.Ctx)
}

// SetDbAndCollection -> set database name and collection
func (mc *MongoCLI) SetDbAndCollection(dbNAME, dbCOLL string) {
	mc.DatabaseName = dbNAME
	mc.Database = mc.Clh.Database(dbNAME)
	mc.CollectionName = dbCOLL
	mc.Collection = mc.Database.Collection(dbCOLL)
	return
}

// GetDB -> get database name
func (mc *MongoCLI) GetDB() string {
	return mc.DatabaseName
}

// GetCollection -> returns collection name
func (mc *MongoCLI) GetCollection() string {
	return mc.CollectionName
}

// GetDatabases -> return list of databases as string array
func (mc *MongoCLI) GetDatabases() (dblist []string, err error) {
	dblist, err = mc.Clh.ListDatabaseNames(mc.Ctx, bson.M{})
	return
}

// CreateDoc -> create document with specified key into collection
func (mc *MongoCLI) CreateDoc(key string, doc interface{}) (docid interface{}, err error) {
	res, err := mc.Collection.InsertOne(mc.Ctx, doc)
	if err != nil {
		return
	}
	docid = res.InsertedID
	return
}

// GetAllDocs -> get all docs from collection
func (mc *MongoCLI) GetAllDocs() (docs []interface{}, err error) {
	cursor, err := mc.Collection.Find(mc.Ctx, bson.M{})
	if err != nil {
		return
	}
	err = cursor.All(mc.Ctx, &docs)
	return
}

// GetDocKV -> get document by key:value filter
func (mc *MongoCLI) GetDocKV(key string, val string) (doc []interface{}, err error) {
	filter, err := mc.Collection.Find(mc.Ctx, bson.M{key: val})
	if err != nil {
		return
	}
	err = filter.All(mc.Ctx, &doc)
	return
}

// FindIdsByKV -> get only slice of doc ids, filter by key:value
func (mc *MongoCLI) FindIdsByKV(key string, val string) (docids []string, err error) {
	docs, err := mc.GetDocKV(key, val)
	if err != nil {
		return
	}
	for _, document := range docs {
		fmt.Printf("document: %+v", document.(bson.M)["_id"])
		docids = append(docids, (document.(bson.M)["_id"]))
	}
	return
}
