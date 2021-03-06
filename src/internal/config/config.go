package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ListenAddress             string
	AccessTokenMinuteLifespan string
	RefreshTokenHourLifespan  string
	ApiSecret                 string
	DB                        struct {
		ConnectionString string
		Name             string
		CollectionName   string
	}
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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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
