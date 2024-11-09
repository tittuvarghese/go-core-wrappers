package config

import (
	"github.com/spf13/viper"
	"github.com/tittuvarghese/core/logger"
)

const DEFAULT_CONFIG_PATH = ".env"

type ConfigManager interface {
	Enable()
	GetString(key string) string
}

type Config struct {
	path string
}

var log = logger.NewLogger("config")

func NewConfigManager(path string) *Config {
	return &Config{path: path}
}

func (conf *Config) Enable() {
	viper.SetConfigFile(conf.path) // specify .env file
	if err := viper.ReadInConfig(); err != nil {
		log.Error("No "+conf.path+" file found, continuing without loading it.", err)
	}
	viper.AutomaticEnv()
}
func (conf *Config) GetString(key string) string {
	return viper.GetString(key)
}
