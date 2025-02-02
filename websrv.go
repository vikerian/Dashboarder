package main

import (
	"context"

	"github.com/google/uuid"
)

type webserver struct {
	listenIP   string
	listenPort uint32
	routes     map[string]string // map route:function
	ctx        context.Context
	pages      []webPage
}

type webPage struct {
	id         uuid.UUID `bson:"id"`
	name       string    `bson:"name"`
	content    []byte    `bson:"content"`
	accessible bool      `bson:"accessible"`
}
