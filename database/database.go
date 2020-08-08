package database

import (
	"bytes"
	"fmt"
	"os"

	xdr "github.com/davecgh/go-xdr/xdr2"
	"github.com/klauspost/compress/snappy"
)

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
	//fmt.Println("encoded data:", string(encodedData))

	filename := "testing.dat"

	//Write data to file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()

	check(err)

	//Lets compress the text
	encoded := snappy.Encode(nil, encodedData)

	f.Write(encoded)

	f.WriteString("0\r0\r")
	f.Sync()

	//Decode Testing
	// var h LogPacket
	// bytesRead, err := xdr.Unmarshal(bytes.NewReader(encodedData), &h)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println("bytes read:", bytesRead)
	// fmt.Println("size:", len(snappy.Encode(nil, encodedData)))
	// fmt.Printf("h: %+v", h)

	//OLD STUFF
	// Current Time in MS
	//nanos := time.Now().Unix()

	//Lets compress the data
	//compressedData := snappy.Encode(nil, data)

	//indexHead := strings.Split(index, "/")[0]
	//We only want the base of the index, all data will be written here
	//println("WAL: ", nanos, " ", indexHead, " ", string(compressedData))
	// var msg strings.Builder
	// msg.WriteString(strconv.FormatInt(nanos, 10))
	// msg.WriteString(indexHead)
	// msg.WriteString(string(compressedData))
	// msg.WriteString("\n")

	// println(msg.String())
	//:= string(nanos , ":" . indexHead + ":" + string(data)

	// file, error := os.OpenFile(indexHead, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	// check(error)
	// defer file.Close()

	// file.Write(msg

}
