package main

import (
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var config Config

type Config struct {
	PollIntervalSeconds int64
}

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/opt/config/")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "read config")
	}
	err = mapstructure.Decode(viper.AllSettings()["main"], &config)
	if err != nil {
		return errors.Wrap(err, "decode viper config to struct")
	}
	return nil
}

func main() {
	err := initConfig()
	if err != nil {
		log.Fatal(err)
	}
}
