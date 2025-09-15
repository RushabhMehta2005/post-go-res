package wal

import (
	"os"
	"strconv"
	"sync"
)

// This file contains the internal implementation
// of the Write-Ahead Logging API

type WAL struct {
	logFilePath string
	logFile *os.File
	mu sync.Mutex
}

// This function will only be called once by the main server thread
func NewWAL(logFilePath string) (*WAL, error) {
	w := new(WAL)
	w.logFilePath = logFilePath
	var err error
	w.logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	return w, err
}

func (w *WAL) Log(operation, key, value string) {
	// Make appending to log file thread safe
	w.mu.Lock()
	defer w.mu.Unlock()
	logEntry := strconv.Itoa(len(operation)) + operation + strconv.Itoa(len(key)) + key + strconv.Itoa(len(value)) + value + "\n"
	logEntryBytes := []byte(logEntry)
	w.logFile.Write(logEntryBytes)
	// Flush the log to disk to ensure durability
	if err := w.logFile.Sync(); err != nil {
		panic(err)
	}
}
