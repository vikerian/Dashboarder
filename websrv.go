package main

import (
	"context"

	"github.com/google/uuid"
)

type WebServer struct {
	ListenIP   string
	ListenPort uint32
	Routes     map[string]string // map route:pageName
	Ctx        context.Context
	Pages      []*WebPage
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
	return &WebServer{
		ListenIP:   listenIP,
		ListenPort: listenPort,
		Routes:     routes,
	}
}
