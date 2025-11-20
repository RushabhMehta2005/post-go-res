package main

import (
	"flag"
	"log"

	"github.com/RushabhMehta2005/post-go-res/memstore"
	"github.com/RushabhMehta2005/post-go-res/server"
	"github.com/RushabhMehta2005/post-go-res/wal"
)

// TODO: Write client library code to interact with our db, multiple programming languages
// TODO: Implement new improved open-addressed hash table, maybe?
// TODO: Add compaction to WAL
// TODO: Persistence should be optional, maybe off by default?


// main is the program entry point.
//
// This program starts a simple TCP key-value server backed by an in-memory
// store and a write-ahead log (WAL) for durability. The server supports a
// tiny text protocol with SET/GET/DEL commands.
//
// Command-line flags:
//   -port   : TCP port to listen on (default 4242)
//   -size   : initial capacity hint for non-sharded HashMap (default 64)
//   -wal    : path to the WAL file (default "./wal_files/wal_file")
//   -shards : number of shards to use for the in-memory store (default 8)
//
// Notes:
// - If shards == 1 a single HashMap is used; otherwise a ShardedMap is created.
// - The WAL is replayed at server startup to rebuild the in-memory state.
func main() {
	// Parse command-line flags
	port := flag.Int("port", 4242, "Port for the database server")
	initialSize := flag.Int("size", 64, "Initial number of key-value pairs ")
	walPath := flag.String("wal", "./wal_files/wal_file", "Path to the WAL file")
	numShards := flag.Int("shards", 8, "The number of in-memory maps")
	persistenceMode := flag.Bool("persist", true, "Enable persistent database mode")

	flag.Parse()

	// Validate flags
	if *port <= 0 || *port > 65535 {
		log.Fatalf("Invalid port: %d. Must be between 1 and 65535.", *port)
	}

	if *initialSize <= 0 {
		log.Fatalf("Initial size must be positive, got %d", *initialSize)
	}

	if *numShards <= 0 {
		log.Fatalf("Invalid shard count: %d. Must be positive.", *numShards)
	}

	// Choose an in-memory store implementation based on shard count.
	var kvstore store.InMemStore
	if *numShards == 1 {
		kvstore = store.NewHashMap(*initialSize)
	} else {
		kvstore = store.NewShardedMap(*numShards, store.DJB2Hash)
	}

	// Instance of WAL
	var walHandler wal.WAL
	if *persistenceMode {
		walHandler, _ = wal.NewFileWAL(*walPath)
	} else {
		walHandler, _ = wal.NewNoOpWAL()
	}

	// Construct and start the server (blocking call)
	server := server.NewServer(kvstore, walHandler, *port)
	server.Start()
}
