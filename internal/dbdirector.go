package main

import (
	"fmt"
	"log/slog"

	//"github.com/vikerian/mongo"
	"github.com/vikerian/mongo"
	// "github.com/vikerian/redis"
)

// types definitions
type DBConn struct {
	//Siri  SiriDB
	//Redis RedisDB
	Mongo mongo.MongoDB
	log   *slog.Logger
}

func NewDBConnections(lg *slog.Logger, cfg *DBAppConfig) (dbs *DBConn, err error) {
	// construct DSN for DBConns and connect to DBConns -> save connections to global dbs var
	//redisDSN := fmt.Sprintf("redis://%s:%s@%s:%d/%d", config.Redis.Username, config.Redis.Password, config.Redis.Host, config.Redis.Port, config.Redis.DBIndex)
	//siriDSN := fmt.Sprintf("siridb://%s:%s@%s:%d/%s", config.SiriDB.Username, config.SiriDB.Password, config.SiriDB.Host, config.SiriDB.Port, config.SiriDB.DBName)
	mongoDSN := mongo.MongoDBCreateDSN(cfg.Mongo.Username, cfg.Mongo.Password, cfg.Mongo.Host, cfg.Mongo.Port, cfg.Mongo.DBName)

	lg.Info("Connecting to mongodb backend database...")
	mg, err := mongo.NewMongoConnection(mongoDSN)
	if err != nil {
		errstr := fmt.Sprintf("Error on connection to mongodb: %v", err)
		lg.Error(errstr)
		return nil, err
	}
	lg.Info("Connection to MongoDB established...")

	return &DBConn{
		log:   lg,
		Mongo: mg,
	}, nil
}

// Create -> create record in collection(table) with key (column) and value
func (dbc *DBConn) Create(table string, key string, value interface{}) (ok bool, err error) {

	return
}

// Read -> Read record from collection(table) with key(column/index), returns value,nil or nil/error
func (dbc *DBConn) Read(table string, key string) (rvalue interface{}, err error) {

	return
}

// Update -> Trying to update record in collection with key
func (dbc *DBConn) Update(table string, key string, newval interface{}) (ok bool, err error) {

	return
}

// Delete -> delete record from collection with key
func (dbc *DBConn) Delete(table string, key string) (ok bool, err error) {

	return
}
