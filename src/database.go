package main

import "fmt"

// DatabaseType represents different types of db
type DatabaseType int

const (
	MongoDB DatabaseType = iota
	SiriDB
	ValkeyDB
)

// DataItem represents a generic data item to be stored
type DataItem struct {
	Data interface{}
}

// DatabaseConnection is an interface for different database connection
type DatabaseConnection interface {
	Connect() error
	Disconnect() error
	Store(data interface{}) (interface{}, error)
	Get(key string) (interface{}, error)
}

// MongoDBConnection implements DatabaseConnection for MongoDB
type MongoDBConnection struct {
	URI        string
	Database   string
	Collection string
}

func (m *MongoDBConnection) Connect() error {
	fmt.Println("Connecting  to MongoDB:", m.URI)

	return nil
}

func (m *MongoDBConnection) Disconnect() error {
	fmt.Println("Disconnecting from MongoDB")

	return nil
}

func (m *MongoDBConnection) Store(data interface{}) error {
	fmt.Println("Storing complex JSON data in MongoDB: ", m.Database)

	return nil
}

func (m *MongoDBConnection) Get(key interface{}) error {

	return nil
}
