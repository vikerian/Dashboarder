package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	//"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9"
)

// RedisCon - structure for redis client operations
type RedisCon struct {
	Client    *redis.Client
	CTX       context.Context
	CTXCANCEL context.CancelFunc
	DbIndex   int
	DSN       string
	Username  string
	Password  string
	Host      string
	Port      string
}

type RedisDB interface {
	//NewRedisConnection(string) (*RedisCon, error)
	Create(string, interface{}) error
	Read(string) (interface{}, error)
	Update(string, interface{}) (bool, error)
	Delete(string) (bool, error)
	Close() error
}

// NewRedisConnection - constructor for redis connections
func NewRedisConnection(dsn string) (*RedisCon, error) {
	rdc := new(RedisCon)

	// construct connection string
	//	connstr := fmt.Sprintf("redis://%s:%s@%s:%s/%d", rdc.Username, rdc.Password, rdc.Host, rdc.Port, rdc.DbIndex)
	rdc.DSN = dsn

	rdc.Password = ""
	opt, err := redis.ParseURL(rdc.DSN)
	if err != nil {
		errstr := fmt.Sprintf("Error on construction/parsing of redis connection dsn: %v", err)
		return nil, errors.New(errstr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	rdc.CTX = ctx
	rdc.CTXCANCEL = cancel
	rdc.Client = redis.NewClient(opt)

	// verify connection
	_, err = rdc.Client.Ping(rdc.CTX).Result()
	if err != nil {
		errstr := fmt.Sprintf("Error on verify redis connection: %v", err)
		return nil, errors.New(errstr)
	}

	return rdc, nil
}

// Create - just key:val, without persistence time (infinite)
func (rdc *RedisCon) Create(key string, val interface{}) error {
	err := rdc.Client.Set(rdc.CTX, key, val, 0)
	if err != nil {
		errstr := fmt.Sprintf("Error on setting key:val : %v", err.Err())
		return errors.New(errstr)
	}
	return nil
}

// Read for a key
func (rdc *RedisCon) Read(key string) (val interface{}, err error) {
	val, err = rdc.Client.Get(rdc.CTX, key).Result()
	if err != nil {
		errstr := fmt.Sprintf("Error on getting value with key %s: %v", key, err)
		return nil, errors.New(errstr)
	}
	return val, nil
}

// Update - update value if exists
func (rdc *RedisCon) Update(key string, val interface{}) (bool, error) {
	/*oldval, err := rdc.Read(key)
	if val == oldval {
		return false, errors.New("Values are the same!")
	}*/
	return true, nil
}

// Delete - delete key/val
func (rdc *RedisCon) Delete(key string) (bool, error) {

	return true, nil
}

func (rdc *RedisCon) Close() error {
	rdc.Client.Close()
	return nil
}
