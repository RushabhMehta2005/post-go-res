package main

import (
	"github.com/RushabhMehta2005/post-go-res/memstore"
	"github.com/RushabhMehta2005/post-go-res/server"
	"github.com/RushabhMehta2005/post-go-res/wal"
)

// TODO: Write client library code to interact with our server as well
// TODO: Add checkpointing to WAL
// TODO: Implement new improved open-addressed hash table

func main() {
	// In memory map to store the actual key-value pair data
	var kvstore = store.NewHashMap()

	// Instance of WAL
	const WAL_FILE_PATH = "./wal_files/wal_file"
	var walHandler, _ = wal.NewWAL(WAL_FILE_PATH)

	server := server.NewServer(kvstore, walHandler)

	server.Start()

}
