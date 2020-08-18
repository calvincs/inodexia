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

	//data := make(chan database.LogPacket)
	//go database.ReadFromWAL("/Users/Calvincs/Duplicati/inodexia/testing.dat", data)

	// for x := range data {
	// 	println("*****************************************************", x.TimeAtIndex, x.IndexPath, x.DataType)
	// 	println(string(x.DataBlob))
	// }

	//Latest Hour
	// tenMinBucket := (time.Now().Unix() / 600) * 600
	// println(tenMinBucket)

}
