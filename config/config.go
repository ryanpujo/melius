package config

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	BaseRoute string `mapstructure:"BASE_ROUTE"`
	DSN       string `mapstructure:"DSN"`
	Port      int    `mapstructure:"PORT"`
	JWTKey    string `mapstructure:"JWT_KEY"`
	JWTConfig jwt    `mapstructure:"JWT"`
}

type jwt struct {
	EXP uint   `mapstructure:"EXP"`
	AUD any    `mapstructure:"AUD"`
	ISS string `mapstructure:"ISS"`
}

var config *Configuration

func readInConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.SetConfigName("proof") // secondary file name without extension
	if err := viper.MergeInConfig(); err != nil {
		log.Panic("Error merging secondary config file:", err)
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
