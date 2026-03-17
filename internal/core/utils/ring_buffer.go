package utils

import (
	"sync"
)

// https://medium.com/checker-engineering/a-practical-guide-to-implementing-a-generic-ring-buffer-in-go-866d27ec1a05

type RingBuffer[T any] struct {
	buffer []T
	size   int
	mu     sync.Mutex
	write  int
	count  int
}

// NewRingBuffer creates a new ring buffer with a fixed size.
func NewRingBuffer[T any](size int) *RingBuffer[T] {
	if size <= 0 {
		panic("ring buffer size must be > 0")
	}

	return &RingBuffer[T]{
		buffer: make([]T, size),
		size:   size,
	}
}

// Add inserts a new element into the buffer, overwriting the oldest if full.
func (rb *RingBuffer[T]) Add(value T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.buffer[rb.write] = value
	rb.write = (rb.write + 1) % rb.size

	if rb.count < rb.size {
		rb.count++
	}
}

// Get returns the contents of the buffer in FIFO order.
func (rb *RingBuffer[T]) Get() []T {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	result := make([]T, 0, rb.count)

	for i := 0; i < rb.count; i++ {
		index := (rb.write + rb.size - rb.count + i) % rb.size
		result = append(result, rb.buffer[index])
	}

	return result
}

// Len returns the current number of elements in the buffer.
func (rb *RingBuffer[T]) Len() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.count
}
