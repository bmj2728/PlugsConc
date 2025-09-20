package logger

import (
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

const DefaultLogFilename = "./logs/app.log"

var DefaultRotator = NewRotator(DefaultLogFilename, 5, 3, 7, true)

func NewRotator(file string, maxSize, maxBackups, maxAge int, compress bool) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   file,
		MaxSize:    maxSize,    // megabytes
		MaxBackups: maxBackups, // number of backups
		MaxAge:     maxAge,     // days
		Compress:   compress,
	}
}

var DefaultFileLogHandlerOptions = NewFileLogHandlerOptions(true, slog.LevelInfo)

func NewFileLogHandlerOptions(addSource bool, level slog.Level) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		AddSource: addSource,
		Level:     level,
	}
}

func NewFileLogHandler(file *lumberjack.Logger, opts *slog.HandlerOptions) *slog.JSONHandler {
	// set unset options
	if opts == nil {
		opts = DefaultFileLogHandlerOptions
	} else if opts.Level == nil {
		opts.Level = slog.LevelInfo
	}

	// resolve rotator config issues
	if file == nil {
		file = DefaultRotator
	} else {
		if file.Filename == "" {
			file.Filename = DefaultRotator.Filename
		}
		if file.MaxAge <= 0 {
			file.MaxAge = 7
		}
		if file.MaxBackups <= 0 {
			file.MaxBackups = 3
		}
		if file.MaxSize <= 0 {
			file.MaxSize = 5
		}
	}
	// validate/fix log file path
	path, err := filepath.Abs(file.Filename)
	if err != nil {
		file.Filename = DefaultLogFilename
	}
	dir := filepath.Dir(path)
	// MkDirAll returns nil on success OR if the dir already exists
	if err := os.MkdirAll(dir, 0750); err != nil {
		// Handle error or fall back to the default if the directory can't be created
		file.Filename = DefaultLogFilename
	} else {
		file.Filename = path // Use the verified, absolute path
	}
	return slog.NewJSONHandler(file, opts)
}
