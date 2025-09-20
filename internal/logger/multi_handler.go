package logger

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

// MultiHandler is a composite handler that forwards logging records to multiple slog.Handlers.
type MultiHandler struct {
	mu       sync.RWMutex
	handlers []slog.Handler
}

// NewMultiHandler creates a new MultiHandler with a slice of slog.Handler instances as its handlers.
func NewMultiHandler(handlers []slog.Handler) *MultiHandler {
	return &MultiHandler{
		mu:       sync.RWMutex{},
		handlers: handlers,
	}
}

// Enabled checks if any of the underlying handlers is enabled for the given context and log level.
func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle processes a log record with all handlers in the MultiHandler and aggregates any errors encountered.
func (m *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	var joinedErrors error
	for _, h := range m.handlers {
		err := h.Handle(ctx, record)
		if err != nil {
			if joinedErrors == nil {
				joinedErrors = err
			} else {
				joinedErrors = errors.Join(joinedErrors, err)
			}
		}
	}
	return joinedErrors
}

// WithAttrs returns a new MultiHandler with the provided attributes added to each underlying handler.
func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return m
	}
	resSlice := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		res := h.WithAttrs(attrs)
		resSlice[i] = res
	}
	return NewMultiHandler(resSlice)
}

// WithGroup returns a new MultiHandler with the specified group name applied to all underlying handlers.
func (m *MultiHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return m
	}
	resSlice := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		res := h.WithGroup(name)
		resSlice[i] = res
	}
	return NewMultiHandler(resSlice)
}
