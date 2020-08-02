package main

import (
	"fmt"
	"sync"

	"github.com/calvincs/inodexia/webserver"
)

var waitgroup = sync.WaitGroup{}

//Entry Point
func main() {
	fmt.Println("hello")
	webserver.HTTPServer()
	waitgroup.Add(1)
	waitgroup.Wait()
}
