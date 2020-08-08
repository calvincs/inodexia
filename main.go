package main

import (
	"inodexia/webserver"
	"sync"
)

var waitgroup = sync.WaitGroup{}

//Entry Point
func main() {
	waitgroup.Add(1)
	go webserver.HTTPServer()
	waitgroup.Wait()
}
