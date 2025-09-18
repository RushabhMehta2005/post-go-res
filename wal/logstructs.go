package wal

import "strconv"

type LogEntry interface {
	ToBytes() []byte
}

type SetEntry struct {
	operation string
	key       *string
	value     *string
}

func NewSetEntry(key, value *string) *SetEntry {
	return &SetEntry{
		operation: "SET",
		key:       key,
		value:     value,
	}
}

func (e *SetEntry) ToBytes() []byte {
	logEntryString := strconv.Itoa(len(e.operation)) + e.operation
	logEntryString += strconv.Itoa(len(*e.key)) + *e.key
	logEntryString += strconv.Itoa(len(*e.value)) + *e.value + "\n"
	logEntryBytes := []byte(logEntryString)
	return logEntryBytes
}

type DelEntry struct {
	operation string
	key       *string
}

func NewDelEntry(key *string) *DelEntry {
	return &DelEntry{
		operation: "DEL",
		key:       key,
	}
}

func (e *DelEntry) ToBytes() []byte {
	logEntryString := strconv.Itoa(len(e.operation)) + e.operation
	logEntryString += strconv.Itoa(len(*e.key)) + *e.key + "\n"
	logEntryBytes := []byte(logEntryString)
	return logEntryBytes
}
