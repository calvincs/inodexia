package webserver

import (
	"fmt"
	"log"
	"strings"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

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
	index := strings.Replace(string(ctx.Path()), "/write/", "", -1)
	log.Print(index)

	//Get the Header for the content type, how are we going to handle this
	rawheaders := string(ctx.Request.Header.Peek("Content-Type"))
	switch rawheaders {
	case "application/x-ndjson":
		//new line delimited json entries
		log.Print(rawheaders, " New Line JSON")
		//WriteToWAL(rawheaders)
	case "application/json":
		log.Print(rawheaders, " JSON")
	default:
		log.Print("No Valid entry found, sending error")
		ctx.Error("invalid headers detected", 503)
	}

	//Data output from post
	postbody := string(ctx.PostBody())
	log.Print(postbody)

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
