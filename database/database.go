package database

import (
	"bufio"
	"bytes"
	"fmt"
	_ "fmt" //Non used return value exclusion
	"os"
	"strings"
	"time"

	xdr "github.com/davecgh/go-xdr/xdr2"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/klauspost/compress/snappy"
)

// Globals
var cfgDb ConfigInodexiaDatabase

func init() {
	// Import configuration file
	err := cleanenv.ReadConfig("config.yml", &cfgDb)
	if err != nil {
		panic(err)
	}
}

// ConfigInodexiaDatabase - Configuration settings
type ConfigInodexiaDatabase struct {
	IndexPath string `yaml:"indexPath"`
}

// LogPacket - Raw incoming data packets for processing
type LogPacket struct {
	TimeAtIndex int64  // Time of the recorded ingestion
	IndexHead   string // Index Head
	IndexPath   string // Full Path of data
	DataBlob    []byte // Compressed Snappy Data
	DataType    string // Content-Type
}

// Handle errors, panic for now
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// IngestionEngine - Were we need to feed DataPackets for processing
func IngestionEngine(packet LogPacket) error {
	//First we send the packet off to the WAL to ensure its recorded incase of failure
	WriteToWAL(packet)
	//fmt.Println(packet.DataType)

	return nil
}

// WriteToWAL This take entries and writes them directly to the WAL file for a given index
func WriteToWAL(packet LogPacket) {

	var packetBuffer bytes.Buffer
	bytesWritten, error := xdr.Marshal(&packetBuffer, &packet)
	check(error)

	encodedData := packetBuffer.Bytes()
	fmt.Println("bytes written:", bytesWritten)
	tenMinBucket := (time.Now().Unix() / 600) * 600
	filename := fmt.Sprintf("%v.snap", tenMinBucket)

	//Set the basepath
	basePath := cfgDb.IndexPath + string(os.PathSeparator) + packet.IndexHead

	//Make the indexing directory if it doesnt already exist
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.Mkdir(basePath, os.ModeDir)
	}

	//Construct the fullPath of the WAL file to be written
	fullPath := basePath + string(os.PathSeparator) + filename

	//Write data to file
	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	check(err)

	//Lets compress the text
	compressed := snappy.Encode(nil, encodedData)
	output := append(compressed, []byte("00\n")...)
	f.Write(output)
	//Ensure its written to the disk
	f.Sync()
}

// A fucntion for the scanner.Scan for reading the WAL file
func walSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {

	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := strings.Index(string(data), "00\n"); i >= 0 {
		return i + 3, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return
}

// ReadFromWAL - Get data from the WAL
func ReadFromWAL(filePath string, data chan<- LogPacket) {

	// Read the file
	file, error := os.Open(filePath)
	check(error)
	defer file.Close()

	// Create the Scanner
	scanner := bufio.NewScanner(file)
	scanner.Split(walSplitFunc)

	//Recreate the objects from the WAL file for reading
	for scanner.Scan() {
		text := scanner.Text()
		expanded, err := snappy.Decode(nil, []byte(text))
		check(err)

		//Decode DataType LogPacket
		var h LogPacket
		bytesRead, err := xdr.Unmarshal(bytes.NewReader(expanded), &h)
		_ = bytesRead
		if err != nil {
			fmt.Println(err)
		}
		data <- h
	}
	close(data)
}
