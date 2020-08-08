package webserver

import (
	"fmt"
	"inodexia/database"
	"log"
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
	Content-Type: application/x-ndjson
	Content-Length: <byte_size>

	<json_encoded_log>
	<json_encoded_log>
	<json_encoded_log>
*/
func WriteHandler(ctx *routing.Context) error {

	//Get the Index Path of the write
	indexHead := strings.Replace(string(ctx.Path()), "/write/", "", -1)
	//log.Print(index)
	//database.WriteToWAL(index, ctx.PostBody())

	//Get the Header for the content type, how are we going to handle this
	rawHeaders := string(ctx.Request.Header.Peek("Content-Type"))
	// switch rawHeaders {
	// case "application/x-ndjson":
	// 	//new line delimited json entries
	// 	log.Print(rawHeaders, " New Line JSON")
	// 	//WriteToWAL(rawheaders)
	// case "application/json":
	// 	log.Print(rawHeaders, " JSON")
	// default:
	// 	log.Print("No Valid entry found, sending error")
	// 	ctx.Error("invalid headers detected", 503)
	// }

	//Send the data onward to the Ingestion Engine for indexing
	database.IngestionEngine(database.LogPacket{
		TimeAtIndex: time.Now().Unix(),
		IndexHead:   indexHead,
		IndexPath:   indexHead,
		DataBlob:    ctx.PostBody(),
		DataType:    rawHeaders})

	return nil
}

// GetInformation GET
/*
	Future things for here
*/
func GetInformation(ctx *routing.Context) error {
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
