package store

import "sync"

// The 'HashMap' struct implements a thread safe in memory store using
// the map data structure in Go
type HashMap struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewHashMap(initialSize int) *HashMap {
	return &HashMap{
		data: make(map[string]string, initialSize),
	}
}

// Gets value from map given the key.
// If key does not exist in map, false is returned alongside the empty string
func (s *HashMap) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, found := s.data[key]
	return val, found
}

// Set key to value in map
func (s *HashMap) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Delete key from map
func (s *HashMap) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}
