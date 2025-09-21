package mq

import (
	"bytes"
	"encoding/gob"
	"errors"
	"path/filepath"

	"github.com/bmj2728/PlugsConc/internal/config"
	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/goptics/sqliteq"
	"github.com/goptics/varmq"
	"github.com/hashicorp/go-hclog"
)

var (
	ErrLogMsgEncoder = errors.New("error encoding log message")
	ErrLogMsgDecoder = errors.New("error decoding log message")
)

type LoggerJob struct {
	Level hclog.Level
	Msg   string
	Args  []any
}

func (j LoggerJob) Encode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(j)
	if err != nil {
		err = errors.Join(ErrLogMsgEncoder, err)
		hclog.Default().Error("error encoding log message", "error", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeLoggerJob(b []byte) (LoggerJob, error) {
	var j LoggerJob

	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&j)
	if err != nil {
		err = errors.Join(ErrLogMsgDecoder, err)
		hclog.Default().Error("error decoding log message", "error", err)
		return j, err
	}
	return j, nil
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
func LogQueue(conf *config.Config, log hclog.Logger) varmq.PersistentQueue[[]byte] {
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
		func(j varmq.Job[[]byte]) {
			lj, err := DecodeLoggerJob(j.Data())
			if err != nil {
				log.Error("Failed to decode log message", logger.KeyError, err)
			}
			switch lj.Level {
			case hclog.Trace:
				log.Trace(lj.Msg, lj.Args...)
			case hclog.Debug:
				log.Debug(lj.Msg, lj.Args...)
			case hclog.Warn:
				log.Warn(lj.Msg, lj.Args...)
			case hclog.Error:
				log.Error(lj.Msg, lj.Args...)
			case hclog.Info:
				log.Info(lj.Msg, lj.Args...)
			default:
				log.Info(lj.Msg, lj.Args)
			}
		}, 10,
	)

	// Bind the loggerWorker to the persistent queue
	return loggerWorker.WithPersistentQueue(persistentQueue)
}
