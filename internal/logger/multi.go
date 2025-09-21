package logger

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

// MultiLogger creates and returns a new hclog.Logger with the specified name, log level, color setting,
// location inclusion, and JSON option. The primary logger must write to stdout.
// Additional locations can be added to the logger by calling RegisterSink on the returned logger.
// A sink is created by passing hclog.LoggerOptions to NewSinkAdapter - no changes are needed to the options.
func MultiLogger(name string,
	level hclog.Level,
	color hclog.ColorOption,
	includeLocation bool,
	isJSON bool) hclog.InterceptLogger {
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
	return MultiLogger("application", hclog.Info, hclog.ForceColor, true, false)
}
