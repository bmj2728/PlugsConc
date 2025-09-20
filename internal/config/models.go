package config

import (
	"io/fs"
	"log/slog"
	"os"

	"github.com/bmj2728/PlugsConc/internal/logger"
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
	Level      string        `yaml:"log_level"`
	Filename   string        `yaml:"log_filename"`
	MaxSize    int           `yaml:"log_max_size"`
	MaxBackups int           `yaml:"log_max_backups"`
	MaxAge     int           `yaml:"log_max_age"`
	Compress   bool          `yaml:"log_compress"`
	AddSource  bool          `yaml:"log_add_source"`
	MQ         LogMQ         `yaml:"log_mq"`
	Colors     LoggingColors `yaml:"log_colors"`
}

type LogMQ struct {
	Enabled bool   `yaml:"log_enable_persistent_queue"`
	File    string `yaml:"log_db_file"`
	Queue   string `yaml:"log_queue"`
	Remove  bool   `yaml:"log_remove_on_complete"`
}

type LoggingColors struct {
	FullLine bool   `yaml:"log_full_line"`
	InfoFGC  string `yaml:"log_info_fg"`
	InfoBGC  string `yaml:"log_info_bg"`
	WarnFGC  string `yaml:"log_warn_fg"`
	WarnBGC  string `yaml:"log_warn_bg"`
	ErrorFGC string `yaml:"log_error_fg"`
	ErrorBGC string `yaml:"log_error_bg"`
	DebugFGC string `yaml:"log_debug_fg"`
	DebugBGC string `yaml:"log_debug_bg"`
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

func (c *Config) LoggerColorMap() map[slog.Level]logger.ColorSetting {

	if c.Logging.Colors == (LoggingColors{}) {
		return logger.DefaultColorMap
	}

	iFG := logger.AvailableColorLookup.GetColor(c.Logging.Colors.InfoFGC)
	if iFG == logger.ResetColor {
		slog.Warn("Invalid color configuration: InfoFGC is reset color")
		iFG = logger.Default
	}
	iBG := logger.AvailableColorLookup.GetColor(c.Logging.Colors.InfoBGC)
	if iBG == logger.ResetColor {
		slog.Warn("Invalid color configuration: InfoBGC is reset color")
		iBG = logger.DefaultBackground
	}
	dFG := logger.AvailableColorLookup.GetColor(c.Logging.Colors.DebugFGC)
	if dFG == logger.ResetColor {
		slog.Warn("Invalid color configuration: DebugFGC is reset color")
		dFG = logger.Default
	}
	dBG := logger.AvailableColorLookup.GetColor(c.Logging.Colors.DebugBGC)
	if dBG == logger.ResetColor {
		slog.Warn("Invalid color configuration: DebugBGC is reset color")
		dBG = logger.DefaultBackground
	}
	wFG := logger.AvailableColorLookup.GetColor(c.Logging.Colors.WarnFGC)
	if wFG == logger.ResetColor {
		slog.Warn("Invalid color configuration: WarnFGC is reset color")
		wFG = logger.Default
	}
	wBG := logger.AvailableColorLookup.GetColor(c.Logging.Colors.WarnBGC)
	if wBG == logger.ResetColor {
		slog.Warn("Invalid color configuration: WarnBGC is reset color")
		wBG = logger.DefaultBackground
	}
	eFG := logger.AvailableColorLookup.GetColor(c.Logging.Colors.ErrorFGC)
	if eFG == logger.ResetColor {
		slog.Warn("Invalid color configuration: ErrorFGC is reset color")
		eFG = logger.Default
	}
	eBG := logger.AvailableColorLookup.GetColor(c.Logging.Colors.ErrorBGC)
	if eBG == logger.ResetColor {
		slog.Warn("Invalid color configuration: ErrorBGC is reset color")
		eBG = logger.DefaultBackground
	}

	i := logger.NewColorSettingWithBackground(iFG, iBG)
	d := logger.NewColorSettingWithBackground(dFG, dBG)
	w := logger.NewColorSettingWithBackground(wFG, wBG)
	e := logger.NewColorSettingWithBackground(eFG, eBG)

	return logger.NewColorMap(i, d, w, e)
}

// LogLevel determines the logging level based on the configuration, returning a corresponding slog.Level value.
func (c *Config) LogLevel() slog.Level {
	switch c.Logging.Level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
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
	return c.Logging.AddSource
}

// FullLine returns true if the configuration specifies full-line logging with colors enabled, otherwise false.
func (c *Config) FullLine() bool {
	return c.Logging.Colors.FullLine
}

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
