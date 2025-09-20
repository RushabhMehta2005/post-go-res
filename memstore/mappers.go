package store

import "hash/crc32"

// MapToN is a function type that takes a string key
// and returns an integer (which will be normalized to [0, size) by Mapper).
type MapToN func(string) int

// Mapper encapsulates key-to-shard index mapping logic.
// It applies the provided mapping function and normalizes
// its output to a valid shard index.
type Mapper struct {
	mapping MapToN
	size    int
}

// NewMapper creates a new Mapper using the given mapping function and size.
func NewMapper(mapping MapToN, size int) *Mapper {
	return &Mapper{
		mapping: mapping,
		size:    size,
	}
}

// GetMapping returns the shard index for the given key.
// The returned value is always in [0, size).
func (m *Mapper) GetMapping(key string) int {
	return (m.mapping(key)%m.size + m.size) % m.size
}

// SimpleSumMap is a trivial MapToN implementation that
// sums the Unicode code points of the key's characters.
func SimpleSumMap(key string) int {
	sum := 0
	for _, char := range key {
		sum += int(char)
	}
	return sum
}

// the classic DJB2 string hash algorithm.
func DJB2Hash(key string) int {
	var h uint32 = 5381
	for i := 0; i < len(key); i++ {
		h = ((h << 5) + h) + uint32(key[i]) // h * 33 + c
	}
	return int(h)
}

// the IEEE CRC32 checksum of the key as an int.
func CRC32Hash(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key)))
}
