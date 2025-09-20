package logger

import (
	"errors"

	"github.com/goptics/varmq"
)

var (
	// ErrFailedToWrite indicates a failure to write to the queue.
	ErrFailedToWrite = errors.New("failed to write to queue")
	// ErrNoQueue indicates that the queue is not present.
	ErrNoQueue = errors.New("queue not present")
	// ErrEmptyMessage indicates that the message is empty.
	ErrEmptyMessage = errors.New("empty message")
)

// AsyncWriter represents a writer that queues messages asynchronously using a persistent queue.
type AsyncWriter struct {
	queue varmq.PersistentQueue[[]byte]
}

// NewAsyncWriter creates and returns a new AsyncWriter initialized with the provided persistent queue.
func NewAsyncWriter(queue varmq.PersistentQueue[[]byte]) *AsyncWriter {
	return &AsyncWriter{
		queue: queue,
	}
}

// Write attempts to enqueue the given byte slice into the queue. Returns the number of bytes written or an error.
func (a AsyncWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, ErrEmptyMessage
	}
	ok := a.queue.Add(p) // try to enqueue the message, returns true if successful, false if not
	if !ok {
		return 0, ErrFailedToWrite
	}
	return len(p), nil
}

// Close safely closes the underlying queue of the AsyncWriter instance if it exists, returning an error if not present.
func (a AsyncWriter) Close() error {
	if a.queue == nil {
		return ErrNoQueue
	}
	return a.queue.Close()
}
