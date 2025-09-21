package logger

import (
	"io"
	"os"

	"github.com/goptics/varmq"
	"github.com/hashicorp/go-hclog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewOptions creates and returns a pointer to a configured hclog.LoggerOptions structure.
func NewOptions(name string,
	level hclog.Level,
	output io.Writer,
	color hclog.ColorOption,
	includeLocation bool,
	isJson bool) *hclog.LoggerOptions {
	return &hclog.LoggerOptions{
		Name:            name,
		Level:           level,
		Output:          output,
		Color:           color,
		IncludeLocation: includeLocation,
		JSONFormat:      isJson}
}

// ConsoleOptions creates a LoggerOptions instance configured for console output with the specified parameters.
func ConsoleOptions(name string,
	level hclog.Level,
	color hclog.ColorOption,
	includeLocation bool,
	isJson bool) *hclog.LoggerOptions {
	return NewOptions(name, level, os.Stdout, color, includeLocation, isJson)
}

// FileOptions configures and returns a new hclog.LoggerOptions instance with the provided file and logging parameters.
// It uses a rolling file logger with optional compression, size, backup, and age constraints for log file management.
// Default values are applied for fileName, maxSize, maxBackups, and maxAge when invalid inputs are provided.
func FileOptions(name string,
	fileName string,
	maxSize int,
	maxBackups int,
	maxAge int,
	compress bool,
	level hclog.Level,
	color hclog.ColorOption,
	includeLocation bool,
	isJson bool) *hclog.LoggerOptions {
	if fileName == "" {
		fileName = DefaultLogFilename
	}
	// limit max log file size to 2MB
	if maxSize <= 0 || maxSize > 2 {
		maxSize = 2
	}
	// ensure valid values for max back-ups and max age
	// zeroes are valid values meaning no limit
	if maxBackups < 0 {
		maxBackups = 0
	}
	if maxAge < 0 {
		maxAge = 0
	}
	out := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
	return NewOptions(name, level, out, color, includeLocation, isJson)
}

// AsyncOptions configures and returns a pointer to hclog.LoggerOptions with asynchronous message queuing support.
func AsyncOptions(name string,
	level hclog.Level,
	queue varmq.PersistentQueue[[]byte],
	color hclog.ColorOption,
	includeLocation bool,
	isJson bool) *hclog.LoggerOptions {
	output := NewAsyncWriter(queue)
	return NewOptions(name, level, output, color, includeLocation, isJson)
}
