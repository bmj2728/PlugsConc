package mq

import (
	"path/filepath"
	"time"

	"github.com/bmj2728/PlugsConc/internal/config"
	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/goptics/sqliteq"
	"github.com/goptics/varmq"
	"github.com/hashicorp/go-hclog"
)

type LoggerJob struct {
	Level hclog.Level
	Msg   string
	Args  []any
}

func NewLoggerJob(level hclog.Level, msg string, args ...any) LoggerJob {
	return LoggerJob{
		Level: level,
		Msg:   msg,
		Args:  args,
	}
}

// LogQueue handles the initialization of a persistent log queue, processes jobs, and logs messages based on
// their severity level.
func LogQueue(conf *config.Config, log hclog.Logger) varmq.PersistentQueue[LoggerJob] {
	if !conf.LogMQEnabled() {
		log.Info("Message queue logging is disabled. Skipping initialization.")
		return nil
	}

	dir := conf.LogsDir()

	aDir, err := filepath.Abs(dir)
	if err != nil {
		log.Error("Failed to get absolute path for logs directory", logger.KeyError, err.Error())
		return nil
	}

	sdb := sqliteq.New(filepath.Join(aDir, conf.LogMQFile()))

	persistentQueue, err := sdb.NewQueue(conf.Logging.MQ.Queue, sqliteq.WithRemoveOnComplete(conf.Logging.MQ.Remove))
	if err != nil {
		log.Error("Failed to create queue", logger.KeyError, err.Error())
	}

	loggerWorker := varmq.NewWorker(
		func(j varmq.Job[LoggerJob]) {
			time.Sleep(10 * time.Second)
			switch j.Data().Level {
			case hclog.Trace:
				log.Trace(j.Data().Msg, j.Data().Args...)
			case hclog.Debug:
				log.Debug(j.Data().Msg, j.Data().Args...)
			case hclog.Warn:
				log.Warn(j.Data().Msg, j.Data().Args...)
			case hclog.Error:
				log.Error(j.Data().Msg, j.Data().Args...)
			case hclog.Info:
				log.Info(j.Data().Msg, j.Data().Args...)
			default:
				log.Info(j.Data().Msg, j.Data().Args)
			}
		}, 10,
	)

	// Bind the loggerWorker to the persistent queue
	return loggerWorker.WithPersistentQueue(persistentQueue)
}
