package config

func DefaultConfig() *Config {
	app := Application{
		AppName: "app",
		AppMode: "dev",
		Version: AppVersion{
			Major:    0,
			Minor:    0,
			Patch:    0,
			Full:     "0.0.0",
			Codename: "alpha",
		},
	}

	dir := Directories{
		Plugins:       "./plugins",
		PluginConfigs: "./plugin_configs",
		Logs:          "./logs",
	}

	Logging := Logging{
		Level:      "info",
		Filename:   "app.log",
		MaxSize:    0,
		MaxBackups: 0,
		MaxAge:     0,
		Compress:   false,
		AddSource:  true,
		MQ: LogMQ{
			Enabled: false,
			File:    "",
		},
		Colors: LoggingColors{
			FullLine: false,
			InfoFGC:  "BrightBlue",
			InfoBGC:  "DefaultBackground",
			WarnFGC:  "BrightYellow",
			WarnBGC:  "DefaultBackground",
			ErrorFGC: "Red",
			ErrorBGC: "DefaultBackground",
			DebugFGC: "BrightGreen",
			DebugBGC: "DefaultBackground",
		},
	}
	fw := FileWatcher{
		Enabled:      false,
		WatchPlugins: false,
	}
	wp := WorkerPool{
		MaxWorkers: 100,
	}
	return &Config{
		Application: app,
		Directories: dir,
		Logging:     Logging,
		FileWatcher: fw,
		WorkerPool:  wp,
	}
}
