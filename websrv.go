package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
)

const DefaultWSConfig string = "webserver.json"

type WebServer struct {
	ListenIP   string            `json:"listen_ip"`
	ListenPort uint32            `json:"listen_port`
	Routes     map[string]string `json:"route_map,omitempty"` // map route:pageName
	Ctx        context.Context   `json:"omitempty"`
	Pages      []*WebPage        `json:"omitempty"`
}

type WebPage struct {
	Id         uuid.UUID `bson:"id"`
	Name       string    `bson:"name"`
	Content    []byte    `bson:"content"`
	Accessible bool      `bson:"accessible"`
}

// NewWebServer - Web server instance construtor
// takes params for web server run
// returns web server instance (not runnning, just filled struct)
func NewWebServer(listenIP string, listenPort uint32, routes map[string]string) *WebServer {
	ws := new(WebServer)
	err := ws.LoadConfig("webserver.json")
	if err != nil {
		errstr := fmt.Sprintf("Fatal error, can't read configuration file for webserver: %v", err)
		slog.Error(errstr)
		panic(err)
	}
	ctx := context.Background()
	ws.Ctx = ctx
	return ws
}

func (ws *WebServer) LoadConfig(cfgFileName string) error {
	if cfgFileName == "" {
		cfgFileName = DefaultWSConfig
	}
	// read configuration - open file - with check
	cf, err := os.Open(cfgFileName)
	if err != nil {
		errstr := fmt.Sprintf("Error on opening configuration file %s: %v", cfgFileName, err)
		return errors.New(errstr)
	}
	defer cf.Close()
	// read configuration - read whole file into config
	cfg := ws
	cfgData, err := io.ReadAll(cf)
	if err != nil {
		errstr := fmt.Sprintf("Error on reading configuration data: %v", err)
		return errors.New(errstr)
	}
	err = json.Unmarshal(cfgData, cfg)
	if err != nil {
		errstr := fmt.Sprintf("Error on decoding configuration: %v", err)
		return errors.New(errstr)
	}
	return err
}

func (ws *WebServer) RunServer() error {
	listenSTR := fmt.Sprintf("%s:%d", ws.ListenIP, ws.ListenPort)
	err := http.ListenAndServe(listenSTR, nil)
	return err
}
