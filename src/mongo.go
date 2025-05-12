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
	URL        string // url in format mongodb://user:pass@fqdn:port/?connection_options
	CTX        context.Context
	CANC       context.CancelFunc
	CLH        *mongo.Client
	OPTS       *options.ClientOptions
	DB         *mongo.Database
	COL        *mongo.Collection
	DBNAME     string
	COLLECTION string
}

// NewMongo -> constructor
func NewMongo(url string) (mc *MongoCLI, err error) {
	mc = new(MongoCLI)
	mc.URL = url
	mc.CTX, mc.CANC = context.WithTimeout(context.Background(), 5*time.Second)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mc.OPTS = options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
	mc.CLH, err = mongo.Connect(mc.OPTS)
	if err != nil {
		return nil, err
	}
	err = mc.CLH.Ping(mc.CTX, readpref.Primary())
	return
}

// CloseMongo -> destructor
func (mc *MongoCLI) CloseMongo() {
	mc.CLH.Disconnect(mc.CTX)
}

// SetDBCollection -> set database name and collection
func (mc *MongoCLI) SetDBCollection(dbNAME, dbCOLL string) {
	mc.DBNAME = dbNAME
	mc.COLLECTION = dbCOLL
	mc.COL = mc.CLH.Database(dbNAME).Collection(dbCOLL)
	return
}

// GetDB -> get database name
func (mc *MongoCLI) GetDB() string {
	return mc.DBNAME
}

// GetCollection -> returns collection name
func (mc *MongoCLI) GetCollection() string {
	return mc.COLLECTION
}

// GetDatabases -> return list of databases as string array
func (mc *MongoCLI) GetDatabases() (dblist []string, err error) {
	dblist, err = mc.CLH.ListDatabaseNames(mc.CTX, bson.M{})
	return
}

// CreateDoc -> create document with specified key into collection
func (mc *MongoCLI) CreateDoc(key string, doc interface{}) (docid interface{}, err error) {
	res, err := mc.COL.InsertOne(mc.CTX, doc)
	if err != nil {
		return
	}
	docid = res.InsertedID
	return
}

// GetAllDocs -> get all docs from collection
func (mc *MongoCLI) GetAllDocs() (docs []interface{}, err error) {
	cursor, err := mc.COL.Find(mc.CTX, bson.M{})
	if err != nil {
		return
	}
	err = cursor.All(mc.CTX, &docs)
	return
}

// GetDocKV -> get document by key:value filter
func (mc *MongoCLI) GetDocKV(key string, val string) (doc []interface{}, err error) {
	filter, err := mc.COL.Find(mc.CTX, bson.M{key: val})
	if err != nil {
		return
	}
	err = filter.All(mc.CTX, &doc)
	return
}

// GetDocIDKV -> get only doc id , filter by key:value
func (mc *MongoCLI) GetDocIDKV(key string, val string) (docids []string, err error) {
	docs, err := mc.GetDocKV(key, val)
	if err != nil {
		return
	}
	for _, document := range docs {
		fmt.Printf("document: %+v", document.(bson.M)["_id"])
		//docids = append(docids, document.(bson.D{primitive.E{Key: "_id", Value: oId}).String())
	}
	return
}
