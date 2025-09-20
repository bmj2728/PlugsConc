package mq

import (
	"log/slog"
	"path/filepath"

	"github.com/bmj2728/PlugsConc/internal/config"
	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/goptics/sqliteq"
	"github.com/goptics/varmq"
)

type LoggerJob struct {
	Level slog.Level
	Msg   string
	Args  []any
}

func NewLoggerJob(level slog.Level, msg string, args ...any) LoggerJob {
	return LoggerJob{
		Level: level,
		Msg:   msg,
		Args:  args,
	}
}

// LogQueue handles the initialization of a persistent log queue, processes jobs, and logs messages based on
// their severity level.
func LogQueue(conf *config.Config) varmq.PersistentQueue[LoggerJob] {
	if !conf.LogMQEnabled() {
		slog.Info("Message queue logging is disabled. Skipping initialization.")
		return nil
	}

	dir := conf.LogsDir()

	aDir, err := filepath.Abs(dir)
	if err != nil {
		slog.Error("Failed to get absolute path for logs directory", slog.Any(logger.KeyError, err))
		return nil
	}

	sdb := sqliteq.New(filepath.Join(aDir, conf.LogMQFile()))

	persistentQueue, err := sdb.NewQueue(conf.Logging.MQ.Queue, sqliteq.WithRemoveOnComplete(conf.Logging.MQ.Remove))
	if err != nil {
		slog.Error("Failed to create queue", slog.Any(logger.KeyError, err))
	}

	loggerWorker := varmq.NewWorker(
		func(j varmq.Job[LoggerJob]) {
			switch j.Data().Level {
			case slog.LevelInfo:
				slog.Info(j.Data().Msg, j.Data().Args...)
			case slog.LevelDebug:
				slog.Debug(j.Data().Msg, j.Data().Args...)
			case slog.LevelWarn:
				slog.Warn(j.Data().Msg, j.Data().Args...)
			case slog.LevelError:
				slog.Error(j.Data().Msg, j.Data().Args...)
			}
		}, 10,
	)

	// Bind the loggerWorker to the persistent queue
	return loggerWorker.WithPersistentQueue(persistentQueue)
}
