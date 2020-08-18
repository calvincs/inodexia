package webserver

import (
	"fmt"
	"inodexia/database"
	"log"
	"regexp"
	"strings"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// Handle errors, panic for now
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// WriteHandler POST
/*
	POST <uri> HTTP/1.1
	Host: <uri.host>
	Content-Type: application/x-ndjson || application/json
	Content-Length: <byte_size>
*/
func WriteHandler(ctx *routing.Context) error {

	//Get the Path values
	indexPath := strings.Replace(string(ctx.Path()), "/write/", "", -1)
	if indexPath == "" {
		indexPath = "default"
	}
	tmp := strings.SplitAfter(indexPath, "/")
	indexHead := strings.TrimRight(tmp[0], "/")
	indexPath = strings.TrimRight(indexPath, "/")

	//Ensure pathing is proper
	isValidPathChar := regexp.MustCompile(`^[A-Za-z0-9\/\.\-\_]+$`).MatchString
	for _, pathchar := range []string{indexPath} {
		if !isValidPathChar(pathchar) {
			ctx.Error("invalid path detected", 400)
			return nil
		}
	}

	//Get the Header, validate type, push to Ingestion Enging
	rawHeaders := string(ctx.Request.Header.Peek("Content-Type"))
	if rawHeaders == "application/x-ndjson" || rawHeaders == "application/json" {
		//Send the data onward to the Ingestion Engine for indexing
		database.IngestionEngine(database.LogPacket{
			TimeAtIndex: time.Now().Unix(),
			IndexHead:   indexHead,
			IndexPath:   indexPath,
			DataBlob:    ctx.PostBody(),
			DataType:    rawHeaders})
	} else {
		ctx.Error("invalid headers detected", 415)
	}

	return nil
}

// GetInformation GET
/*
	Future things for here
*/
func GetInformation(ctx *routing.Context) error {
	// var path string = "/Users/Calvincs/Duplicati/inodexia/testing.dat"
	// data := database.ReadFromWAL(path)

	fmt.Fprintf(ctx, "stuff %s", ctx.Path())
	return nil
}

// HTTPServer Setup
func HTTPServer() {
	router := routing.New()
	router.Post("/write/*", WriteHandler)
	router.Get("/info/*", GetInformation)

	error := fasthttp.ListenAndServe(":8080", router.HandleRequest)
	if error != nil {
		log.Fatal(error)
	}
}
