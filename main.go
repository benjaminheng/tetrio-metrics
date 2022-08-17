package main

import (
	"context"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var config Config

type Config struct {
	PollIntervalSeconds int64
	TetrioUserID        string
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

func checkForNewTetrioGames(ctx context.Context) (err error) {
	log.Println("checking for recent tetrio games")
	defer func() {
		if err != nil {
			log.Println(errors.Wrap(err, "error"))
		}
	}()
	parsedResponse, rawResponse, err := getTetrioRecentUserStreams(ctx, config.TetrioUserID)
	if err != nil {
		return errors.Wrap(err, "get tetrio recent user streams")
	}
	_ = parsedResponse
	_ = rawResponse
	// TODO: save to DB
	return nil
}

func poll(ctx context.Context) error {
	checkForNewTetrioGames(ctx)
	ticker := time.NewTicker(time.Duration(config.PollIntervalSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			checkForNewTetrioGames(ctx)
		}
	}
}

func main() {
	err := initConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "init config"))
	}
	ctx := context.Background()
	err = poll(ctx)
	if err != nil {
		log.Fatal(errors.Wrap(err, "poll"))
	}
}
