package store

import "sync"

// HashMap is a simple, thread-safe in-memory key-value store.
// HashMap is best suited for smaller maps or low-contention workloads.
// For high-concurrency scenarios, consider using ShardedMap.
type HashMap struct {
	mu   sync.RWMutex // read/write mutex to guard access to data
	data map[string]string // underlying key-value store
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
