package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
)

// InfoColor defines the ANSI color code for bright blue, typically used for informational messages.
// DebugColor defines the ANSI color code for bright green, typically used for debug messages.
// WarnColor defines the ANSI color code for bright yellow, typically used for warning messages.
// ErrorColor defines the ANSI color code for red, typically used for error messages.
const (
	InfoColor  = "\033[1;94m" // Bright Blue
	DebugColor = "\033[1;92m" // Bright Green
	WarnColor  = "\033[1;93m" // Bright Yellow
	ErrorColor = "\033[1;31m" // Red
	ResetColor = "\033[0m"
)

var (
	DefaultColorMap = map[slog.Level]string{
		slog.LevelInfo:  InfoColor,
		slog.LevelDebug: DebugColor,
		slog.LevelWarn:  WarnColor,
		slog.LevelError: ErrorColor,
	}

	// DefaultOptions defines the default logging configuration with an info log level and a predefined color map.
	DefaultOptions = Options{
		Level:    slog.LevelInfo,
		ColorMap: DefaultColorMap,
	}
)

// NewColorMap creates a new color map with the provided color codes for each log level.
// The color codes should be formatted as ANSI escape sequences e.g. for black "\033[1;30m"
// The following color codes are supported:
// 30-37(FG): Black, Red, Green, Yellow, Blue, Magenta, Cyan, White
// 90-97(BG): Bright Black, Bright Red, Bright Green, Bright Yellow,
// Bright Blue, Bright Magenta, Bright Cyan, Bright White
// 40-47(FG): Black, Red, Green, Yellow, Blue, Magenta, Cyan, White
// 100-107(BG): Bright Black, Bright Red, Bright Green, Bright Yellow,
// Bright Blue, Bright Magenta, Bright Cyan, Bright White
func NewColorMap(info, debug, warn, error string) map[slog.Level]string {
	return map[slog.Level]string{
		slog.LevelInfo:  info,
		slog.LevelDebug: debug,
		slog.LevelWarn:  warn,
		slog.LevelError: error,
	}
}

// ColorHandler is a log handler that formats and outputs log records with optional colorized levels.
// It supports customization of log levels, color mappings, and formatting through the Options struct.
// ColorHandler is thread-safe, leveraging a mutex for synchronization when handling concurrent log output.
type ColorHandler struct {
	opts          Options
	preformatted  []byte // data from WithGroup and WithAttrs
	unopenedGroup string // group from WithGroup that is pending Attrs and has not been opened
	mu            *sync.Mutex
	out           io.Writer
}

// Options defines configuration for customizing log level handling and mapping levels to color codes in a logger.
// Level specifies the log level filter, determining which log messages are processed.
// ColorMap maps log levels to their respective color codes for enhanced readability in console output.
type Options struct {
	Level    slog.Leveler
	ColorMap map[slog.Level]string
}

// New creates a new ColorHandler instance with the provided output writer and options.
// Defaults are applied if opts is nil.
func New(out io.Writer, opts *Options) *ColorHandler {
	ch := &ColorHandler{
		out: out,
		mu:  &sync.Mutex{},
	}
	if opts != nil {
		ch.opts = *opts
	} else {
		ch.opts = DefaultOptions
	}
	// apply defaults if passed options are nil
	if ch.opts.Level == nil {
		ch.opts.Level = slog.LevelInfo
	}
	if ch.opts.ColorMap == nil {
		ch.opts.ColorMap = DefaultColorMap
	}
	return ch
}

// NewDefault initializes and returns a new ColorHandler instance with default options and output set to stderr.
func NewDefault() *ColorHandler {
	return &ColorHandler{
		out:  os.Stderr,
		mu:   &sync.Mutex{},
		opts: DefaultOptions,
	}
}

// Enabled determines if a log level is enabled based on the handler's configuration and the provided context.
func (c *ColorHandler) Enabled(_ context.Context, level slog.Level) bool {
	// in a real-world implementation, this would check the context for a log level filter or something similar
	return level >= c.opts.Level.Level()
}

// Handle writes a log record to the output with optional color formatting based on the log level.
// It ensures thread-safety using a mutex and skips processing if the log level is disabled.
func (c *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	if !c.Enabled(ctx, r.Level) {
		return nil
	}
	color := c.getColor(r.Level)
	// get a buffer from the sync pool
	buf := make([]byte, 0, 1024)
	// set the color for the log level
	buf = append(buf, color...)
	if !r.Time.IsZero() {
		buf = append(buf, r.Time.Format("2006-01-02 15:04:05")...)
		buf = append(buf, ' ')
	}
	buf = append(buf, r.Level.String()...)
	buf = append(buf, ' ')
	buf = append(buf, r.Message...)
	buf = append(buf, ' ')
	buf = append(buf, c.preformatted...)
	r.Attrs(func(a slog.Attr) bool {
		entry := fmt.Sprintf("%s:%s ", a.Key, a.Value)
		buf = append(buf, entry...)
		return true
	})
	// reset the color
	buf = append(buf, ResetColor...)
	buf = append(buf, '\n')
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.out.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// WithAttrs returns a new handler with the specified attributes added to the preformatted log output.
func (c *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return c
	}
	c2 := *c
	c2.preformatted = make([]byte, len(c.preformatted))
	copy(c2.preformatted, c.preformatted)
	if c2.unopenedGroup == "" {
		for _, attr := range attrs {
			entry := fmt.Sprintf("%s:%s ", attr.Key, attr.Value)
			c2.preformatted = append(c2.preformatted, entry...)
		}
	} else if c2.unopenedGroup != "" {
		for _, attr := range attrs {
			entry := fmt.Sprintf("%s.%s=%s ", c2.unopenedGroup, attr.Key, attr.Value)
			c2.preformatted = append(c2.preformatted, entry...)
		}
	}
	c2.unopenedGroup = ""
	return &c2
}

// WithGroup returns a new ColorHandler with the specified group name, enabling group-based attribute organization.
func (c *ColorHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return c
	}
	c2 := *c
	if c2.unopenedGroup != "" {
		c2.unopenedGroup = fmt.Sprintf("%s.%s", c2.unopenedGroup, name)
	} else {
		c2.unopenedGroup = name
	}
	return &c2
}

// getColor retrieves the ANSI color code for the provided log level based on the configured or default color map.
func (c *ColorHandler) getColor(level slog.Level) string {
	if c.opts.ColorMap == nil {
		color, ok := DefaultColorMap[level]
		if !ok {
			return getDefaultColor(level)
		}
		return color
	} else {
		color, ok := c.opts.ColorMap[level]
		if !ok {
			return getDefaultColor(level)
		}
		return color
	}
}

// getDefaultColor returns the default ANSI color code for the given log level using DefaultColorMap.
// If the level is not present in DefaultColorMap, it returns the ResetColor.
func getDefaultColor(level slog.Level) string {
	color, ok := DefaultColorMap[level]
	if !ok {
		return ResetColor
	}
	return color
}
