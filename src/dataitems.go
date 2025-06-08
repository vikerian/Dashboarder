package main

//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"

// DataItem represents a generic data item to be stored
type DataItem struct {
	Key  interface{}
	Data interface{}
}

// NewDataItem -> create dataitem instance, if specified key, it will be used in Store, otherwise key will be filled in Store method
func NewDataItem(key interface{}, data interface{}) *DataItem {
	return &DataItem{
		Key:  key,
		Data: data,
	}
}

// MongoDOC represents a document to store to MongoDB
type MongoDOC struct {
	Name     string   `bson:"_name"`
	Key      string   `bson:"_id,omitempty"`
	Document DataItem `bson:"data_item"`
}

// NewMongoDOC -> create mongo document instance
func NewMongoDOC(name string, key string, doc interface{}) *MongoDOC {
	return &MongoDOC{
		Name:     name,
		Key:      key,
		Document: doc.(DataItem),
	}
}
