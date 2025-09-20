package config

import (
	"io/fs"
	"os"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Application Application `yaml:"application"`
	Directories Directories `yaml:"directories"`
	Logging     Logging     `yaml:"logging"`
	FileWatcher FileWatcher `yaml:"file_watcher"`
	WorkerPool  WorkerPool  `yaml:"worker_pool"`
}

type Application struct {
	AppName string     `yaml:"app_name"`
	AppMode string     `yaml:"app_mode"`
	Version AppVersion `yaml:"app_version"`
}

type AppVersion struct {
	Major    int    `yaml:"major"`
	Minor    int    `yaml:"minor"`
	Patch    int    `yaml:"patch"`
	Full     string `yaml:"full"`
	Codename string `yaml:"codename"`
}

type Directories struct {
	Plugins       string `yaml:"plugins_dir"`
	PluginConfigs string `yaml:"plugin_configs_dir"`
	Logs          string `yaml:"logs_dir"`
}

type Logging struct {
	Level        string `yaml:"log_level"`
	Filename     string `yaml:"log_filename"`
	MaxSize      int    `yaml:"log_max_size"`
	MaxBackups   int    `yaml:"log_max_backups"`
	MaxAge       int    `yaml:"log_max_age"`
	Compress     bool   `yaml:"log_compress"`
	InclLocation bool   `yaml:"log_include_location"`
	MQ           LogMQ  `yaml:"log_mq"`
}

type LogMQ struct {
	Enabled bool   `yaml:"log_enable_persistent_queue"`
	File    string `yaml:"log_db_file"`
	Queue   string `yaml:"log_queue"`
	Remove  bool   `yaml:"log_remove_on_complete"`
}

type FileWatcher struct {
	Enabled      bool `yaml:"fw_enabled"`
	WatchPlugins bool `yaml:"fw_watch_plugins"`
}

type WorkerPool struct {
	MaxWorkers int `yaml:"wp_max_workers"`
}

// TODO replace with viper

func LoadConfig(root *os.Root, path string) *Config {
	data, err := fs.ReadFile(root.FS(), path)
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return &config
}

// LogLevel determines the logging level based on the configuration, returning a corresponding hclog.Level value.
func (c *Config) LogLevel() hclog.Level {
	return hclog.LevelFromString(c.Logging.Level)
}

// LogFilename returns the filename for the log file as specified in the Logging configuration.
func (c *Config) LogFilename() string {
	return c.Logging.Filename
}

// LogMaxSize retrieves the maximum size, in megabytes, of a log file as specified in the logging configuration.
func (c *Config) LogMaxSize() int {
	return c.Logging.MaxSize
}

// LogMaxBackups retrieves the maximum number of backup log files to retain, as specified in the logging configuration.
func (c *Config) LogMaxBackups() int {
	return c.Logging.MaxBackups
}

// LogMaxAge returns the maximum age (in days) for log retention as specified in the logging configuration.
func (c *Config) LogMaxAge() int {
	return c.Logging.MaxAge
}

// LogCompress checks if log compression is enabled in the configuration and returns true if enabled, otherwise false.
func (c *Config) LogCompress() bool {
	return c.Logging.Compress
}

// AddSource returns a boolean indicating whether the logging configuration includes the source of log entries.
func (c *Config) AddSource() bool {
	return c.Logging.InclLocation
}

//// FullLine returns true if the configuration specifies full-line logging with colors enabled, otherwise false.
//func (c *Config) FullLine() bool {
//	return c.Logging.Colors.FullLine
//}

// PluginsDir returns the configured directory path where plugins are stored.
func (c *Config) PluginsDir() string {
	return c.Directories.Plugins
}

// PluginConfigsDir returns the directory path for plugin configuration files as defined in the Config struct.
func (c *Config) PluginConfigsDir() string {
	return c.Directories.PluginConfigs
}

// LogsDir returns the directory path for storing log files as specified in the configuration.
func (c *Config) LogsDir() string {
	return c.Directories.Logs
}

// FileWatcherEnabled checks if the file watcher functionality is enabled in the configuration and returns true if enabled.
func (c *Config) FileWatcherEnabled() bool {
	return c.FileWatcher.Enabled
}

// FileWatcherWatchPlugins returns true if the file watcher is configured to monitor plugin changes, otherwise false.
func (c *Config) FileWatcherWatchPlugins() bool {
	return c.FileWatcher.WatchPlugins
}

// WorkerPoolMaxWorkers returns the maximum number of workers allowed in the worker pool, as configured.
func (c *Config) WorkerPoolMaxWorkers() int {
	return c.WorkerPool.MaxWorkers
}

// LogMQEnabled checks if the message queue logging is enabled in the configuration and returns true if enabled.
func (c *Config) LogMQEnabled() bool {
	return c.Logging.MQ.Enabled
}

// LogMQFile returns the file path for the persistent message queue log as specified in the configuration.
func (c *Config) LogMQFile() string {
	return c.Logging.MQ.File
}
