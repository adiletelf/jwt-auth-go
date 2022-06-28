package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ListenAddress string
}

var Cfg *Config

func New() (*Config, error) {
	if Cfg != nil {
		return Cfg, nil
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	Cfg = config
	return Cfg, nil
}