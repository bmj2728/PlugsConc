package config

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
