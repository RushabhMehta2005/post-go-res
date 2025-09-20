package utils

import "sync/atomic"

// AtomicCounter is a simple wrapper around an int64 counter
// that provides atomic increment, decrement, and addition.
type AtomicCounter struct {
	c int64
}

func NewCounter() *AtomicCounter {
	return &AtomicCounter{}
}

func (counter *AtomicCounter) Add(delta int64) {
	atomic.AddInt64(&counter.c, delta)
}

func (counter *AtomicCounter) Inc() {
	atomic.AddInt64(&counter.c, 1)
}

func (counter *AtomicCounter) Dec() {
	atomic.AddInt64(&counter.c, -1)
}

// Read returns the current value of the counter atomically.
func (counter *AtomicCounter) Read() int64 {
	return counter.c
}
