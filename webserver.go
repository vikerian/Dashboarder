package main

import (
	"fmt"
	"net/http"
	"html/template"
)


func static() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
}

