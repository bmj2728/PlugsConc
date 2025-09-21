package mq

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

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

type LogEntry struct {
	Caller    string `json:"@caller"`
	Level     string `json:"@level"`
	Message   string `json:"@message"`
	Module    string `json:"@module"`
	Timestamp string `json:"@timestamp"`
	// You might also want to handle arbitrary additional fields
	Fields map[string]interface{} `json:"-"`
}

func (l *LogEntry) UnmarshalJSON(data []byte) error {
	// First unmarshal into a generic map
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Extract known fields
	l.Caller, _ = raw["@caller"].(string)
	l.Level, _ = raw["@level"].(string)
	l.Message, _ = raw["@message"].(string)
	l.Module, _ = raw["@module"].(string)
	l.Timestamp, _ = raw["@timestamp"].(string)

	// Everything else goes into Fields
	l.Fields = make(map[string]interface{})
	for k, v := range raw {
		if !strings.HasPrefix(k, "@") {
			l.Fields[k] = v
		}
	}

	return nil
}

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
	if err != nil && !strings.Contains(err.Error(), "bad data: undefined type") {
		fmt.Println(err)
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
func LogQueue(conf *config.Config, qLogger hclog.Logger) varmq.PersistentQueue[[]byte] {
	if !conf.LogMQEnabled() {
		hclog.Default().Info("Message queue logging is disabled. Skipping initialization.")
		return nil
	}

	dir := conf.LogsDir()

	aDir, err := filepath.Abs(dir)
	if err != nil {
		hclog.Default().Error("Failed to get absolute path for logs directory", logger.KeyError, err.Error())
		return nil
	}

	sdb := sqliteq.New(filepath.Join(aDir, conf.LogMQFile()))

	persistentQueue, err := sdb.NewQueue(conf.Logging.MQ.Queue, sqliteq.WithRemoveOnComplete(conf.Logging.MQ.Remove))
	if err != nil {
		hclog.Default().Error("Failed to create queue", logger.KeyError, err.Error())
	}

	loggerWorker := varmq.NewWorker(
		func(j varmq.Job[[]byte]) {
			var logEntry LogEntry
			err := logEntry.UnmarshalJSON(j.Data())
			if err != nil {
				hclog.Default().Error("Failed to unmarshal log message", logger.KeyError, err)
			}
			// from here we'll extract the data then use the passed in interceptor to log the message
			lev := hclog.LevelFromString(logEntry.Level)
			msg := logEntry.Message
			var args []any

			args = append(args, "caller", logEntry.Caller)
			args = append(args, "module", logEntry.Module)
			args = append(args, "orig_timestamp", logEntry.Timestamp)

			for k, v := range logEntry.Fields {
				args = append(args, k, v)
			}

			switch lev {
			case hclog.Trace:
				qLogger.Trace(msg, args...)
			case hclog.Debug:
				qLogger.Debug(msg, args...)
			case hclog.Warn:
				qLogger.Warn(msg, args...)
			case hclog.Error:
				qLogger.Error(msg, args...)
			case hclog.Info:
				qLogger.Info(msg, args...)
			default:
				qLogger.Info(msg, args...)

			}
		}, 10,
	)

	// Bind the loggerWorker to the persistent queue
	return loggerWorker.WithPersistentQueue(persistentQueue)
}
