package main

// github.com/SiriDB/go-siridb-connector
import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/SiriDB/go-siridb-connector"
	// "github.com/kos-v/dsnparser"
	"github.com/kos-v/dsnparser"
)

const SIRIDB_BUFFER_SIZE = 8192

// SiriCon - structure for siriDB connection interface
type SiriCon struct {
	Conn          *siridb.Connection
	Options       map[string]interface{}
	Database      string
	Table         string
	AdminUsername string
	AdminPassword string
}

// SiriDB - interface to our Siri Database
type SiriDB interface {
	//NewConnection(string) (SiriCon, error)
	Create(string, interface{}) error
	Read(string) (interface{}, error)
	Update(string, interface{}) (bool, error)
	Delete(string) (bool, error)
	Close() error
}

// NewConnection - constructor for our siridb connection
func NewSiriDBConnection(dsnstr string) (SiriCon, error) {
	// dsn: siridb://user:pass@host:port/db -> now parse it!
	dsn := dsnparser.Parse(dsnstr)
	host := dsn.GetHost()
	port := dsn.GetPort()
	iport, err := strconv.Atoi(port)
	if err != nil {
		errstr := fmt.Sprintf("Error on parsing DSN string: %v", err)
		return SiriCon{}, errors.New(errstr)
	}
	username := dsn.GetUser()
	password := dsn.GetPassword()
	dbname := dsn.GetPath()
	// parsed values used here
	var sdb SiriCon //:= new(SiriCon)
	sdb.Conn = siridb.NewConnection(host, uint16(iport))
	sdb.Options = make(map[string]interface{})
	sdb.Options["dbname"] = dbname
	sdb.Options["time_precision"] = "ms"
	sdb.Options["buffer_size"] = SIRIDB_BUFFER_SIZE
	sdb.Database = dbname
	sdb.AdminUsername = os.Getenv("SIRI_ADMIN_USER")
	sdb.AdminPassword = os.Getenv("SIRI_ADMIN_PASSWORD")
	if sdb.AdminUsername == "" { // if not known via env, try default
		sdb.AdminUsername = "sa"
	}
	if sdb.AdminPassword == "" { // if not known via env, try default
		sdb.AdminPassword = "siri"
	}
	if err := sdb.Conn.Connect(username, password, dbname); err != nil {
		var errstr string
		errstr = fmt.Sprintf("Connection to siridb on host %s port %s database %s failed: %v", host, port, dbname, err)
		if strings.Contains(err.Error(), "unknown database") { // try to create db if unknown db
			// using the default service account 'sa' and password 'siri'
			if _, err := sdb.Conn.Manage(sdb.AdminUsername, sdb.AdminPassword, siridb.AdminNewDatabase, sdb.Options); err != nil {
				errstr = fmt.Sprintf("Connection to siridb on host %s port %s database %s failed: %v", host, port, dbname, err)
				return SiriCon{}, errors.New(errstr)
			} else {
				return sdb, nil
			}
		}
		return SiriCon{}, errors.New(errstr) // other error then unknown db
	}

	return sdb, nil
}

// Create - create value function
func (sc SiriCon) Create(tkey string, val interface{}) error {

	return nil
}

// Read - read value by key
func (sc SiriCon) Read(key string) (interface{}, error) {

	return nil, nil
}

func (sc SiriCon) Update(key string, value interface{}) (bool, error) {

	return true, nil
}

func (sc SiriCon) Delete(key string) (bool, error) {

	return true, nil
}

func (sc SiriCon) Close() error {
	sc.Conn.Close()
	return nil
}
