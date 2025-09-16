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
	// Flush the log to disk to ensure durability, this is disk IO and it is expensive
	if err := w.logFile.Sync(); err != nil {
		panic(err)
	}
}

func (w *WAL) ReBuild(store store.InMemStore) {
	scanner := bufio.NewScanner(w.logFile)

	// Iterate through the lines
	for scanner.Scan() {
		line := scanner.Text()
		// operation does not matter for now
		_, key, value, err := parseLine(line)
		if err != nil {
			log.Fatal(err)
		}
		store.Set(key, value)
	}	
}

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
