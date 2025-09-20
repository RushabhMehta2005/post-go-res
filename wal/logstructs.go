package wal

import "strconv"

// LogEntry is the interface implemented by WAL entries.
// Implementations must provide a ToBytes method that returns the
// length-prefixed textual representation appended to the WAL file.
type LogEntry interface {
	ToBytes() []byte
}


// SetEntry represents a SET operation in the WAL.
// The serialized form is:
//   <len("SET")>"SET"<len(key)><key><len(value)><value>\n
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

// DelEntry represents a DEL operation in the WAL.
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
