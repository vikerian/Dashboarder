package main

import "net/url"

//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"

// DatabaseType represents different types of db
type DatabaseType int

const (
	MongoDB DatabaseType = iota
	SiriDB
	ValkeyDB
)

// DatabaseConnection is an interface for different database connection
type DatabaseConnection interface {
	Connect(string, string, string) error
	Disconnect() error
	Store(data interface{}) (interface{}, error)
	Get(string) (DataItem, error)
}

type DatabaseConnectionInfo struct {
	DbTYPE int
	DSN    string // dsn = database connection string in url format: mongodb://user@pass@fqdn:port/database
	URI    url.URL
	Client interface{}
}
