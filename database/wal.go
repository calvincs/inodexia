package database

// WriteToWAL This take entries and writes them directly to the WAL file for a given index
func WriteToWAL(index string) {

	//We only want the base of the index, all data will be written here
	println("WAL: ", index)
}
