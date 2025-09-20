package store

import "sync"

// ShardedMap is a concurrent, sharded key-value store.
// It splits keys into multiple independent shards to reduce
// lock contention under heavy concurrent access.
type ShardedMap struct {
	shardLocks []sync.RWMutex // one RWMutex per shard
	shardData  []map[string]string // one map per shard
	mapper     Mapper // maps keys to shard indices
}

// NewShardedMap creates a new ShardedMap with numShards shards.
func NewShardedMap(numShards int, mapping MapToN) *ShardedMap {
	sm := &ShardedMap{
		shardLocks: make([]sync.RWMutex, numShards),
		shardData:  make([]map[string]string, numShards),
		mapper:     *NewMapper(mapping, numShards),
	}

	for i := range numShards {
		sm.shardData[i] = make(map[string]string)
	}
	return sm
}

// Get retrieves the value for the given key.
func (sm *ShardedMap) Get(key string) (string, bool) {
	shardIndex := sm.mapper.GetMapping(key)
	sm.shardLocks[shardIndex].RLock()
	defer sm.shardLocks[shardIndex].RUnlock()
	value, found := sm.shardData[shardIndex][key]
	return value, found
}

// Set sets the value for the given key.
func (sm *ShardedMap) Set(key, value string) {
	shardIndex := sm.mapper.GetMapping(key)
	sm.shardLocks[shardIndex].Lock()
	defer sm.shardLocks[shardIndex].Unlock()
	sm.shardData[shardIndex][key] = value
}

// Delete removes the given key from the map.
func (sm *ShardedMap) Delete(key string) {
	shardIndex := sm.mapper.GetMapping(key)
	sm.shardLocks[shardIndex].Lock()
	defer sm.shardLocks[shardIndex].Unlock()
	delete(sm.shardData[shardIndex], key)
}
