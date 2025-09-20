package logger

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// MultiLogger creates and returns a new hclog.Logger with the specified name, log level, color setting, location inclusion, and JSON option.
func MultiLogger(name string,
	level hclog.Level,
	color hclog.ColorOption,
	includeLocation bool,
	isJSON bool) hclog.Logger {
	return hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:            name,
		Level:           level,
		Output:          os.Stdout,
		Color:           color,
		IncludeLocation: includeLocation,
		JSONFormat:      isJSON})
}

// DefaultLogger returns a pre-configured logger instance with default parameters for application-level logging.
func DefaultLogger() hclog.Logger {
	return MultiLogger("application", hclog.Info, hclog.AutoColor, true, false)
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
