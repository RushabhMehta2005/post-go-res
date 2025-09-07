package store

import "sync"

// The 'InMemStore' struct implements a thread safe map
type InMemStore struct {
	rwMutex sync.RWMutex
	data    map[string]string
}

// "Constructor" function to return new instance of InMemStore
// Currently the size is fixed to 32
// TODO: Make this configurable through cmd line flag
func NewInMemStore() *InMemStore {
	return &InMemStore{
		data: make(map[string]string, 32),
	}
}

// Gets value from map given the key.
// If key does not exist in map, false is returned alongside the empty string
func (s *InMemStore) Get(key string) (string, bool) {
	s.rwMutex.RLock()
	val, found := s.data[key]
	s.rwMutex.RUnlock()
	return val, found
}

// Set key to value in map
func (s *InMemStore) Set(key, value string) {
	s.rwMutex.Lock()
	s.data[key] = value
	s.rwMutex.Unlock()
}
