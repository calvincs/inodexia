package webserver

import (
	"fmt"
	"strings"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// WebServer Setup
func WebServer() {
	router := routing.New()
	router.Get("/index/*", func(c *routing.Context) error {
		fmt.Fprintf(c, "PATH %q", c.Path())
		index := strings.Replace(string(c.Path()), "/index/", "", -1)
		println("string: ", index)
		rawheaders := string(c.Request.Header.Peek("red"))
		if rawheaders != "" {
			println(rawheaders)
		} else {
			println("No headers")
		}

		return nil
	})

	panic(fasthttp.ListenAndServe(":8080", router.HandleRequest))
}
