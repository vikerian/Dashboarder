package mongodb

// Import necessary libraries (those commented out are just hints, our vim-go sometimes put them out :D )
import (

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"github.com/google/uuid"
)

// MongoDB -> structure for our Mongo Api Client to work (settings for mongo communications)
type MongoDB struct {
	URL            string
	Clh            *mongo.Client
	Ctx            context.Context
	Cancel         context.CancelFunc
	DatabaseName   string
	CollectionName string
	Cursor         *mongo.Collection
}

// Document -> structure to hold document and metadata to it
type Document struct {
	ID      string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string    `json:"name" bson:"name"`
	Content []byte    `json:"content" bson:"content"`
	Tags    []string  `json:"tags,omitempty" bson:"tags,omitempty"`
	Created time.Time `json:"created" bson:"created"`
}

// New - Constructor for our connection
// take params:
// connstr -> string ->  connection url like mongodb://localhost:27017
// return instance of MongoDB, error
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
// take params:
// dbName -> string -> database name 
// colname -> string -> collection name
func (mdb *MongoDB) SetDBCollection(dbName string, colname string) {
	mdb.DatabaseName = dbName
	mdb.CollectionName = colname
	mdb.Cursor = mdb.Clh.Database(dbName).Collection(colname)
}

// InsertDoc -> insert document to database
// take params:
// id -> string -> prefered id (if we have prefference), or give "", func will generate one
// name -> string -> document name or if "", func will generate document-Year-Month-Day_Minute-Second
// docContent -> []byte -> raw document
// tags -> []string -> array of tags for document
// returns "",err if error, 
// otherwise string ObjectID("[id of our stored doc]")
func (mdb *MongoDB) InsertDoc(id string, name string,docContent []byte,tags []string) (saveID string, err error) {
	// actual time 
	atime := time.Now()
	// set id if none
	if id == "" {
		id = uuid.New().String()
	}
	// set name if none
	if name == "" {
		name = fmt.Sprintf("document-%d-%d-%d_%d-%d",atime.Year(),atime.Month(),atime.Day(),atime.Minute(),atime.Second())
	}
	var doc = Document{
		ID: id,
		Name: name,
		Content: docContent,
		Tags: tags,
		Created: atime,
	}				
		
	// save document to db
	result, err := mdb.Cursor.InsertOne(mdb.Ctx, doc)
	if err != nil {
		return "", err
	}
	// get saved object id 
	saveID = fmt.Sprintf("%v",result.InsertedID)
	return saveID, nil
}
