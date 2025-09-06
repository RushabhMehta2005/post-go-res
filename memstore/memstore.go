package store

import "sync"

// The 'Store' struct implements a thread safe map
type InMemStore struct {
	rwMutex sync.RWMutex
	data map[string]string
}

// Set key to value in map
func (s *InMemStore) Set(key, value string) {
	s.rwMutex.Lock()
	s.data[key] = value
	s.rwMutex.Unlock()
}

// Gets value from map given the key.
// If key does not exist in map, false is returned alongside the empty string
func (s *InMemStore) Get(key string) (string, bool) {
	s.rwMutex.RLock()
	val, found := s.data[key]
	s.rwMutex.RUnlock()
	return val, found
}
