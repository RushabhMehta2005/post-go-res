package utils

import "sync/atomic"

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

func (counter *AtomicCounter) Read() int64 {
	return counter.c
}
