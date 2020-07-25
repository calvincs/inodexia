package main

import (
	"fmt"
	"sync"
)

var waitgroup = sync.WaitGroup{}

//Entry Point
func main() {
	fmt.Println("hello")

	waitgroup.Add(1)

	waitgroup.Wait()
}
