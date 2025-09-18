package store

import "sync"

// The 'HashMap' struct implements a thread safe in memory store using
// the map data structure in Go
type HashMap struct {
	rwMutex sync.RWMutex
	data    map[string]string
}

// "Constructor" function to return new instance of HashMap
// Currently the initial size is fixed to 32
// TODO: Make this configurable through cmd line flag
func NewHashMap() *HashMap {
	return &HashMap{
		data: make(map[string]string, 32),
	}
}

// Gets value from map given the key.
// If key does not exist in map, false is returned alongside the empty string
func (s *HashMap) Get(key string) (string, bool) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	val, found := s.data[key]
	return val, found
}

// Set key to value in map
func (s *HashMap) Set(key, value string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	s.data[key] = value
}

// Delete key from map
func (s *HashMap) Delete(key string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	delete(s.data, key)
}
