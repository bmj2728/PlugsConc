package logger

import (
	"github.com/hashicorp/go-hclog"
	"gopkg.in/natefinch/lumberjack.v2"
)

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
