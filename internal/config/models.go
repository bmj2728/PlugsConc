package config

import (
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Directories Directories
	Logging     Logging
	FileWatcher FileWatcher
	WorkerPool  WorkerPool
}

type Directories struct {
	Plugins       string `yaml:"plugins"`
	PluginConfigs string `yaml:"plugin_configs"`
}

type Logging struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
	Colors     LoggingColors
}

type LoggingColors struct {
	FullLine bool   `yaml:"full_line"`
	InfoFGC  string `yaml:"info_fg"`
	InfoBGC  string `yaml:"info_bg"`
	WarnFGC  string `yaml:"warn_fg"`
	WarnBGC  string `yaml:"warn_bg"`
	ErrorFGC string `yaml:"error_fg"`
	ErrorBGC string `yaml:"error_bg"`
	DebugFGC string `yaml:"debug_fg"`
	DebugBGC string `yaml:"debug_bg"`
}

type FileWatcher struct {
	Enabled      bool `yaml:"enabled"`
	WatchPlugins bool `yaml:"watch_plugins"`
}

type WorkerPool struct {
	MaxWorkers int `yaml:"max_workers"`
}

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
