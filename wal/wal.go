package wal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"unicode"

	"github.com/RushabhMehta2005/post-go-res/memstore"
)

// TODO: Improve error handling here, decide who should handle the error and where

// WAL implements a simple write-ahead log stored in a file.
//
// Each mutation (SET/DEL) should be appended to the WAL before being
// applied to the in-memory store. On startup the WAL can be replayed
// to rebuild the in-memory store by calling ReBuild.
//
// WAL provides basic concurrency protection for appending log entries.
type WAL struct {
	logFilePath string
	logFile     *os.File
	mu          sync.Mutex
}

// NewWAL opens (or creates) the WAL file at logFilePath and returns a WAL.
// This function will only be called once by the main server thread
func NewWAL(logFilePath string) (*WAL, error) {
	w := new(WAL)
	w.logFilePath = logFilePath
	var err error
	w.logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	return w, err
}

// Log appends the provided LogEntry to the WAL and flushes it to stable storage.
//
// Log is safe for concurrent callers: it serializes append+sync operations
// with an internal mutex to ensure entries are written and persisted in order.
func (w *WAL) Log(entry LogEntry) {
	// Convert the log to a slice of bytes
	logEntryBytes := entry.ToBytes()
	w.logFile.Write(logEntryBytes)

	// Make appending to log file thread safe
	w.mu.Lock()
	defer w.mu.Unlock()

	// Flush the log to disk to ensure durability, this is disk IO and it is expensive
	if err := w.logFile.Sync(); err != nil {
		panic(err)
	}
}

// ReBuild replays the WAL and applies each logged mutation to the provided store.
func (w *WAL) ReBuild(store store.InMemStore) {
	scanner := bufio.NewScanner(w.logFile)

	// Iterate through the lines
	for scanner.Scan() {
		line := scanner.Text()
		operation, key, value, err := parseLine(line)
		if err != nil {
			log.Fatal(err)
		}

		switch operation {
		case "SET":
			store.Set(key, value)
		case "DEL":
			store.Delete(key)
		}
	}
}

// parseLine parses a single WAL line into operation, key and value.
//
// The WAL format expected by parseLine is a sequence of three length-prefixed fields:
//   <opLen><opBytes><keyLen><keyBytes><valueLen><valueBytes>
// where each length is ASCII digits (e.g. "3SET5hello4data") and the parser reads
// the number, then consumes exactly that many bytes for the field.
func parseLine(line string) (operation, key, value string, err error) {
	readLen := func(s string, i *int) (int, error) {
		start := *i
		for *i < len(s) && unicode.IsDigit(rune(s[*i])) {
			*i++
		}
		if start == *i {
			return 0, fmt.Errorf("expected number at position %d", start)
		}
		n, err := strconv.Atoi(s[start:*i])
		if err != nil {
			return 0, err
		}
		return n, nil
	}

	i := 0
	// Parse operation
	opLen, err := readLen(line, &i)
	if err != nil {
		return "", "", "", err
	}
	if i+opLen > len(line) {
		return "", "", "", fmt.Errorf("operation length exceeds input")
	}
	operation = line[i : i+opLen]
	i += opLen

	// Parse key
	keyLen, err := readLen(line, &i)
	if err != nil {
		return "", "", "", err
	}
	if i+keyLen > len(line) {
		return "", "", "", fmt.Errorf("key length exceeds input")
	}
	key = line[i : i+keyLen]
	i += keyLen

	// Parse value
	valueLen, err := readLen(line, &i)
	if err != nil {
		return "", "", "", err
	}
	if i+valueLen > len(line) {
		return "", "", "", fmt.Errorf("value length exceeds input")
	}
	value = line[i : i+valueLen]
	i += valueLen

	return operation, key, value, nil
}
