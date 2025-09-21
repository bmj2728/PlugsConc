package config

import (
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
)

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
