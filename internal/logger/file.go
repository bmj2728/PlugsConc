package logger

import (
	"github.com/hashicorp/go-hclog"
	"gopkg.in/natefinch/lumberjack.v2"
)

const DefaultLogFilename = "./logs/app.log"

// DefaultRotator is the default rotator for the application.
// Log files can grow to 2MB before being rotated and compressed with a maximum of 25 backups.
// Log files are retained unless manually deleted.
var DefaultRotator = NewRotator(DefaultLogFilename, 2, 25, 0, true)

func NewRotator(file string, maxSize, maxBackups, maxAge int, compress bool) *lumberjack.Logger {
	if file == "" {
		file = DefaultLogFilename
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
	return &lumberjack.Logger{
		Filename:   file,
		MaxSize:    maxSize,    // megabytes
		MaxBackups: maxBackups, // number of backups
		MaxAge:     maxAge,     // days
		Compress:   compress,
	}
}

// FileLogger creates and returns a new instance of hclog.Logger configured with the specified options.
// Accepts a logger name, logging level, output rotator, color options, location inclusion, and JSON formatting settings.
func FileLogger(name string,
	level hclog.Level,
	rotator *lumberjack.Logger,
	color hclog.ColorOption,
	includeLocation bool,
	isJSON bool) hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:            name,
		Level:           level,
		Output:          rotator,
		Color:           color,
		IncludeLocation: includeLocation,
		JSONFormat:      isJSON,
	})
}

// FileSink creates a new hclog.SinkAdapter for logging to a file with configurable options like level, format, and color.
// It supports log file rotation through the provided lumberjack.Logger instance.
func FileSink(name string,
	level hclog.Level,
	rotator *lumberjack.Logger,
	color hclog.ColorOption,
	includeLocation bool,
	isJSON bool) hclog.SinkAdapter {
	return hclog.NewSinkAdapter(&hclog.LoggerOptions{
		Name:            name,
		Level:           level,
		Output:          rotator,
		Color:           color,
		IncludeLocation: includeLocation,
		JSONFormat:      isJSON,
	})
}
