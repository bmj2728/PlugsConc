package logger

import (
	"errors"
	"io"

	"github.com/goptics/varmq"
	"github.com/hashicorp/go-hclog"
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

// AsyncSink creates and returns a SinkAdapter for asynchronous logging using a persistent message queue.
// The queue must be initialized with an AsyncInterceptLogger. This sync can then be passed to the multi-writer.
func AsyncSink(name string,
	queue varmq.PersistentQueue[[]byte],
	level hclog.Level,
	color hclog.ColorOption,
	includeLocation bool,
	isJSON bool) hclog.SinkAdapter {
	return hclog.NewSinkAdapter(AsyncOptions(name, level, queue, color, includeLocation, isJSON))
}

// AsyncInterceptLogger creates the base logger for asynchronous logging using a persistent message queue.
// This logger must be passed to LogQueue. The defined logger will use the queue to write messages asynchronously.
// Additional async destinations can be added to the logger by calling RegisterSink on the returned logger.
func AsyncInterceptLogger(name string,
	level hclog.Level,
	output io.Writer,
	color hclog.ColorOption,
	includeLocation bool,
	isJSON bool) hclog.InterceptLogger {
	return hclog.NewInterceptLogger(NewOptions(name, level, output, color, includeLocation, isJSON))
}

/*
Logging Setup
**imports from config should create logger options - the type can be used to create loggers, intercept loggers or sinks
1. default logger until config is loaded
2. setup console log from config if configured
3. we need:
	a. Multi-Write Intercept Logger - done
		i. this takes the console logger as the base logger - done
	b. If any other loggers are set for async, we create:
		i. an Async Sink - register it to the multi-writer
		ii. An Async Intercept Logger + the first async logger is used to create it
		iii. additional sinks created for any loggers using async and registered to async sink
		iv. Async Intercept is passed to LogQueue
		v. flow becomes multi -> async sink -> async writer -> queue -> async logger -> fanned to each sink
	c. If any loggers aare not set for async, we create:
		i. sinks for each logger and then register it to the multi-writer: the current implementation works fine
*/
