package main

import (
	"inodexia/database"
	"sync"
)

var waitgroup = sync.WaitGroup{}

//Entry Point
func main() {
	//waitgroup.Add(1)
	//go webserver.HTTPServer()

	data := make(chan database.LogPacket)
	go database.ReadFromWAL("/Users/Calvincs/Duplicati/inodexia/testing.dat", data)

	for x := range data {
		println(x.TimeAtIndex, x.IndexPath, x.DataType)
	}

	//waitgroup.Wait()
}
