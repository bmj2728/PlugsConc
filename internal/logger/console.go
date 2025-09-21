package logger

import (
	"github.com/hashicorp/go-hclog"
)

// ConsoleLogger creates a new logger instance with specified name, level, color, location inclusion,
// and JSON formatting options.
func ConsoleLogger(name string,
	level hclog.Level,
	color hclog.ColorOption,
	includeLocation bool,
	isJson bool) hclog.Logger {
	return hclog.New(ConsoleOptions(name, level, color, includeLocation, isJson))
}
