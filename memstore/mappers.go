package store

import "hash/crc32"

// A class of functions which take a string as input
// and return an integer which is their hash
type MapToN func(string) int

// A structure to encapsulate the
// key -> index mapping logic
type Mapper struct {
	mapping MapToN
	size    int
}

func NewMapper(mapping MapToN, size int) *Mapper {
	return &Mapper{
		mapping: mapping,
		size:    size,
	}
}

func (m *Mapper) GetMapping(key string) int {
	return (m.mapping(key)%m.size + m.size) % m.size
}

func SimpleSumMap(key string) int {
	sum := 0
	for _, char := range key {
		sum += int(char)
	}
	return sum
}

func DJB2Hash(key string) int {
	var h uint32 = 5381
	for i := 0; i < len(key); i++ {
		h = ((h << 5) + h) + uint32(key[i]) // h * 33 + c
	}
	return int(h)
}

// CRC32 raw hash (returns an int)
func CRC32Hash(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key)))
}
