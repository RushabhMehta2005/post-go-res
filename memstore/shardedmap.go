package store

import "sync"

type ShardedMap struct {
	shardLocks []sync.RWMutex
	shardData  []map[string]string
	mapper     Mapper
}

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

func (sm *ShardedMap) Get(key string) (string, bool) {
	shardIndex := sm.mapper.GetMapping(key)
	sm.shardLocks[shardIndex].RLock()
	defer sm.shardLocks[shardIndex].RUnlock()
	value, found := sm.shardData[shardIndex][key]
	return value, found
}

func (sm *ShardedMap) Set(key, value string) {
	shardIndex := sm.mapper.GetMapping(key)
	sm.shardLocks[shardIndex].Lock()
	defer sm.shardLocks[shardIndex].Unlock()
	sm.shardData[shardIndex][key] = value
}

func (sm *ShardedMap) Delete(key string) {
	shardIndex := sm.mapper.GetMapping(key)
	sm.shardLocks[shardIndex].Lock()
	defer sm.shardLocks[shardIndex].Unlock()
	delete(sm.shardData[shardIndex], key)
}
