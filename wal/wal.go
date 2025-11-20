package wal

import (
	"github.com/RushabhMehta2005/post-go-res/memstore"
)

// WAL defines the interface for write-ahead logging operations.
type WAL interface {
	Log(entry LogEntry) error
	ReBuild(store store.InMemStore) error
}
