package main

import (
	"dashboarder/config"
	"dashboarder/mongodb"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

var Log *slog.Logger

var mockDocument = map[string]string{
	"name":          "simpletestdocument",
	"page":          "sinister.html",
	"sometext":      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed sit amet pharetra ligula, nec ultricies est. Morbi ac erat nec nunc fringilla scelerisque rutrum vitae nisl. Ut nec orci luctus, vulputate nisi volutpat, tincidunt neque. Sed ac libero vel tellus tincidunt ultrices in nec mauris. Morbi dictum faucibus nisl et ultrices. Sed vitae sapien iaculis, gravida augue eu, cursus tortor. Suspendisse potenti. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Nunc nulla magna, volutpat in nulla id, tempor posuere lectus. In ac nunc vel massa suscipit ultrices. Nam ultricies venenatis neque id pulvinar. Donec molestie scelerisque pharetra. Etiam orci risus, porta ac scelerisque nec, condimentum non mi. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Curabitur porttitor nisl nec lectus elementum consequat. Curabitur egestas, lectus eu blandit dapibus, sem nunc blandit justo, efficitur fermentum erat augue eget turpis.",
	"someothertext": "Duis porta eu odio vitae elementum. Vivamus orci tellus, semper in diam eu, tristique vulputate est. Proin imperdiet auctor nulla, eu rutrum dui maximus eget. Integer augue libero, placerat ac nunc ut, tristique euismod erat. Cras sit amet sapien facilisis est hendrerit molestie. Sed bibendum pellentesque metus, non semper nisi efficitur quis. Maecenas laoreet orci non egestas mollis. Sed et condimentum sem. Donec ac leo interdum, vulputate risus accumsan, vestibulum nibh. Praesent non venenatis quam, sed fermentum sapien. Integer fermentum id ligula et tempor. Donec tincidunt augue eu consectetur feugiat.",
}

func init() {
	//slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.SetLogLoggerLevel(slog.LevelDebug)
	Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func main() {
	Log.Info("Dashboarder starting up...")
	slog.Debug("Getting configuration...")
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	dbgMsg := fmt.Sprintf("Config loaded: %+v", cfg)
	slog.Debug(dbgMsg)
	slog.Debug("Configuration loaded, connecting to DBs...")

	slog.Debug(fmt.Sprintf("Testing mongodb package..."))
	mdb, err := mongodb.New("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	slog.Debug("Setting up database and collection, getting cursor...")
	mdb.SetDBCollection("testdb", "testcollection")
	slog.Debug(fmt.Sprintf("Actual mdb: %+v", mdb))

	slog.Debug(fmt.Sprintf("inserting testing doc into db..."))
	mockDocConverted, err := json.Marshal(mockDocument)
	if err != nil {
		panic(err)
	}
	tags := []string{"simpledocument", "mockdocument", "jsondocument"}
	docID, err := mdb.InsertDoc("", "simpledoc", mockDocConverted, tags)
	if err != nil {
		panic(err)
	}
	slog.Debug(fmt.Sprintf("Seems doc is saved as: %v", docID))
}
