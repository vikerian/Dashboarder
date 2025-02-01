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
	Siri  SiriDB
	Redis RedisDB
	Mongo mongo.MongoDB
	log   *slog.Logger
}

func NewDBConnections(lg *slog.Logger, cfg *DBAppConfig) (dbs *DBConn, err error) {
	//construct DSN for DBConns and connect to DBConns -> save connections to global dbs var
	redisDSN := fmt.Sprintf("redis://%s:%s@%s:%d/%d", cfg.Redis.Username, cfg.Redis.Password, cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DBIndex)
	siriDSN := fmt.Sprintf("siridb://%s:%s@%s:%d/%s", cfg.SiriDB.Username, cfg.SiriDB.Password, cfg.SiriDB.Host, cfg.SiriDB.Port, cfg.SiriDB.DBName)
	mongoDSN := mongo.MongoDBCreateDSN(cfg.Mongo.Username, cfg.Mongo.Password, cfg.Mongo.Host, cfg.Mongo.Port, cfg.Mongo.DBName)

	lg.Info("Connecting to mongodb backend database...")
	mg, err := NewMongoConnection(mongoDSN, lg)
	if err != nil {
		errstr := fmt.Sprintf("Error on connection to mongodb: %v", err)
		lg.Error(errstr)
		return nil, err
	}
	lg.Info("Connection to MongoDB established...")

	lg.Info("Connecting to siriDB database...")
	sdb, err := NewSiriDBConnection(siriDSN)
	if err != nil {
		errstr := fmt.Sprintf("Error on connecting to siridb: %v", err)
		lg.Error(errstr)
		return nil, err
	}
	lg.Info("SiriDB connection established...")

	lg.Info("Connecting to redis database...")
	rdb, err := NewRedisConnection(redisDSN)
	if err != nil {
		errstr := fmt.Sprintf("Error on connecting to redis : %v", err)
		lg.Error(errstr)
		return nil, err
	}

	lg.Info("RedisDB connection established...")

	return &DBConn{
		log:   lg,
		Mongo: mg,
		Redis: rdb,
		Siri:  sdb,
	}, nil
}

// Create -> create record in collection(table) with key (column) and value
func (dbc *DBConn) Create(table string, key string, value interface{}) (ok bool, err error) {
	ok, err = dbc.Create(table, key, value)
	return
}

// Read -> Read record from collection(table) with key(column/index), returns value,nil or nil/error
func (dbc *DBConn) Read(table string, key string) (rvalue interface{}, err error) {
	rvalue, err = dbc.Read(table, key)
	return
}

// Update -> Trying to update record in collection with key
func (dbc *DBConn) Update(table string, key string, newval interface{}) (ok bool, err error) {

	return
}

// Delete -> delete record from collection with key
func (dbc *DBConn) Delete(table string, key string) (ok bool, err error) {
	ok, err = dbc.Delete(table, key)
	return
}
