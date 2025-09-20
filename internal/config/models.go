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
	File    string `yaml:"log_persistent_db_file"`
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
