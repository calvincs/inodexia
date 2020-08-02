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

	waitgroup.Add(1)
	webserver.HTTPServer()
	waitgroup.Wait()
}
