package logger

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"

	"github.com/goptics/sqliteq"
	"github.com/goptics/varmq"
	"github.com/hashicorp/go-hclog"
)

var (
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

// LogQueue handles the initialization of a persistent log queue, processes jobs, and logs messages based on
// their severity level.
func LogQueue(qLogger hclog.Logger) varmq.PersistentQueue[[]byte] {

	dir := "/home/brian/GolandProjects/PlugsConc/logs"

	aDir, err := filepath.Abs(dir)
	if err != nil {
		hclog.Default().Error("Failed to get absolute path for logs directory", KeyError, err.Error())
		return nil
	}

	sdb := sqliteq.New(filepath.Join(aDir, "logs.db"))

	persistentQueue, err := sdb.NewQueue("log-queue", sqliteq.WithRemoveOnComplete(true))
	if err != nil {
		hclog.Default().Error("Failed to create queue", KeyError, err.Error())
	}

	loggerWorker := varmq.NewWorker(
		func(j varmq.Job[[]byte]) {
			var logEntry LogEntry
			err := logEntry.UnmarshalJSON(j.Data())
			if err != nil {
				hclog.Default().Error("Failed to unmarshal log message", KeyError, errors.Join(ErrLogMsgDecoder, err))
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
