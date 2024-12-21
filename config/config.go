package config

import "github.com/spf13/viper"

type Configuration struct {
	DSN  string `mapstructure:"DSN"`
	Port int    `mapstructure:"PORT"`
}

var config *Configuration

func readInConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}
}

func Config() *Configuration {
	if config == nil {
		readInConfig()
	}
	return config
}
