package main

import (
	"fmt"
	"log/slog"
	"os"
)

func init() {
		lf, err := os.OpenFile("application.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0640)
	if err != nil {
		panic(err)
	}
	// multiwriter
	mw := io.MultiWriter(os.Stdout, lf)
	// our logger setup
	applog = slog.New(slog.NewJSONHandler(mw, nil))
}

func main() {

}
