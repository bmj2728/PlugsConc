package config

type Config struct {
	Directories Directories
	Logging     Logging
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
	InfoFGC  string `yaml:"info_fgc"`
	InfoBGC  string `yaml:"info_bgc"`
	WarnFGC  string `yaml:"warn_fgc"`
	WarnBGC  string `yaml:"warn_bgc"`
	ErrorFGC string `yaml:"error_fgc"`
	ErrorBGC string `yaml:"error_bgc"`
	DebugFGC string `yaml:"debug_fgc"`
	DebugBGC string `yaml:"debug_bgc"`
}
