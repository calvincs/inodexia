package main

import (
	"inodexia/webserver"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

// ConfigInodexia - Configuration settings
type ConfigInodexia struct {
	StartWeb bool `yaml:"startWeb"`
}

//Entry Point
func main() {
	// Import configuration file
	var cfg ConfigInodexia
	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		panic(err)
	}

	//Create the wait Group
	var waitgroup = sync.WaitGroup{}

	if cfg.StartWeb {
		waitgroup.Add(1)
		go webserver.HTTPServer()
		waitgroup.Wait()
	}

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
