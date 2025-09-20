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
// TODO: Write README.md

func main() {
	// Deal with command-line flags
	port := flag.Int("port", 4242, "Port for the database server")
	initialSize := flag.Int("size", 64, "Initial number of key-value pairs ")
	walPath := flag.String("wal", "./wal_files/wal_file", "Path to the WAL file")
	numShards := flag.Int("shards", 8, "The number of in-memory maps")

	flag.Parse()

	if *port <= 0 || *port > 65535 {
		log.Fatalf("Invalid port: %d. Must be between 1 and 65535.", *port)
	}

	if *initialSize <= 0 {
		log.Fatalf("Initial size must be positive, got %d", *initialSize)
	}

	if *numShards <= 0 {
		log.Fatalf("Invalid shard count: %d. Must be positive.", *numShards)
	}

	// In-memory structures to store the actual key-value pair data

	var kvstore store.InMemStore
	if *numShards == 1 {
		kvstore = store.NewHashMap(*initialSize)
	} else {
		kvstore = store.NewShardedMap(*numShards, store.DJB2Hash)
	}

	// Instance of WAL
	var walHandler, _ = wal.NewWAL(*walPath)

	server := server.NewServer(kvstore, walHandler, *port)
	server.Start()
}
