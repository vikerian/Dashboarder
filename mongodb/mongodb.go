package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Url              string
	Client           *mongo.Client
	Database         *mongo.Database
	Collections      []mongo.Collection
	ActiveCollection mongo.Collection
}

/* New(string) -> Constructor */
func New(ctx context.Context, uri string) (mg *MongoDB, err error) {
	mg = new(MongoDB)
	mg.Url = uri

	mg.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(mg.Url))
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
func (mg *MongoDB) GetAllDocumentsByCollection(ctx context.Context, colname string) (docs documents, err error) {
	if colname != "" {
		mg.ActiveCollection = *mg.Database.Collection(colname)
	}
	cursor, err := mg.ActiveCollection.Find(ctx, bson.D{})
	if err != nil {
		return
	}

	var singleResult interface{}
	for cursor.Next(ctx) {
		if err = cursor.Decode(&singleResult); err != nil {
			return
		}
		sr := singleResult.(interface{})
		docs = append(docs, sr)
	}
	return
}
