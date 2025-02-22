package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Url              string
	Client           *mongo.Client
	Ctx              context.Context
	CancelFunc       context.CancelFunc
	Database         *mongo.Database
	Collections      []mongo.Collection
	ActiveCollection mongo.Collection
}

/* New(string) -> Constructor */
func New(uri string) (mg *MongoDB, err error) {
	mg = new(MongoDB)
	mg.Url = uri
	mg.Ctx, mg.CancelFunc = context.WithTimeout(context.Background(), 5*time.Second)
	mg.Client, err = mongo.Connect(mg.Ctx, options.Client().ApplyURI(mg.Url))
	if err != nil {
		return
	}
	return
}

/* SetDatabaseAndCollection -> setup database and collection to use now */
func (mg *MongoDB) SetDatabaseAndCollection(databasestr, collectionstr string) {
	mg.Database = mg.Client.Database(databasestr)
	mg.ActiveCollection = *mg.Database.Collection(collectionstr)

}

type documents []interface{}

/* GetAllDocumentsByCollection -> get all docs from collection */
func (mg *MongoDB) GetAllDocumentsByCollection(colname string) (docs documents, err error) {
	if colname != "" {
		mg.ActiveCollection = *mg.Database.Collection(colname)
	}
	cursor, err := mg.ActiveCollection.Find(mg.Ctx, bson.D{})
	if err != nil {
		return
	}
	for cursor.Next(mg.Ctx) {
		var singleResult interface{}
		if err = cursor.Decode(&singleResult); err != nil {
			return
		}
		docs = docs.append(singleResult)
	}
	return
}
