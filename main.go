package main

import (
	"github.com/RushabhMehta2005/post-go-res/memstore"
	"github.com/RushabhMehta2005/post-go-res/server"
	"github.com/RushabhMehta2005/post-go-res/wal"
)

// TODO: Write client library code to interact with our db, multiple programming languages
// TODO: Implement new improved open-addressed hash table, maybe?
// TODO: Add cmd line flags to set port, wal file location, hashmap initial size, etc params
// TODO: Add compaction to WAL
// TODO: Improve responses to client

func main() {
	// In memory map to store the actual key-value pair data
	var kvstore = store.NewHashMap()

	// Instance of WAL
	const WAL_FILE_PATH = "./wal_files/wal_file"
	var walHandler, _ = wal.NewWAL(WAL_FILE_PATH)

	server := server.NewServer(kvstore, walHandler)
	server.Start()
}
