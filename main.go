package main

import (
	"context"
	"log"
	"time"

	"github.com/benjaminheng/tetrio-metrics/store"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	PollIntervalSeconds int64
	TetrioUserID        string
	DatabaseFilePath    string
}

type Service struct {
	store  *store.Store
	config Config
}

func initConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/opt/config/")
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, errors.Wrap(err, "read config")
	}
	config := Config{}
	err = mapstructure.Decode(viper.AllSettings()["main"], &config)
	if err != nil {
		return Config{}, errors.Wrap(err, "decode viper config to struct")
	}
	return config, nil
}

func (s *Service) checkForNewTetrioGames(ctx context.Context) (err error) {
	log.Println("checking for recent tetrio games")
	defer func() {
		if err != nil {
			log.Println(errors.Wrap(err, "error"))
		}
	}()
	parsedResponse, err := getTetrioRecentUserStreams(ctx, s.config.TetrioUserID)
	if err != nil {
		return errors.Wrap(err, "get tetrio recent user streams")
	}
	for _, record := range parsedResponse.Data.Records {
		m, err := buildGamemode40LRecord(record)
		if err != nil {
			log.Println(errors.Wrap(err, "build gamemode_40l model"))
		}
		// TODO: insert only if it doesn't exist
		log.Printf("saving game: ts=%v time=%vs\n", m.PlayedAt.Format(time.RFC3339Nano), float64(m.TimeMs)/1000.0)
		_, err = s.store.InsertGamemode40L(ctx, store.InsertGamemode40LParams{
			PlayedAt:        m.PlayedAt,
			TimeMs:          m.TimeMs,
			FinessePercent:  m.FinessePercent,
			TotalPieces:     m.TotalPieces,
			PiecesPerSecond: m.PiecesPerSecond,
			RawData:         m.RawData,
		})
		if err != nil {
			log.Println(errors.Wrap(err, "insert gamemode_40l record to DB"))
		}
	}
	return nil
}

func (s *Service) poll(ctx context.Context) error {
	s.checkForNewTetrioGames(ctx)
	ticker := time.NewTicker(time.Duration(s.config.PollIntervalSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			s.checkForNewTetrioGames(ctx)
		}
	}
}

func NewService(config Config) (*Service, error) {
	storeConfig := store.Config{
		DatabaseFilePath: config.DatabaseFilePath,
	}
	store, err := store.NewStore(storeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "initialize storage")
	}
	s := &Service{
		store:  store,
		config: config,
	}
	return s, nil
}

func main() {
	config, err := initConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "init config"))
	}

	service, err := NewService(config)
	if err != nil {
		log.Fatal(errors.Wrap(err, "init service"))
	}

	ctx := context.Background()
	err = service.poll(ctx)
	if err != nil {
		log.Fatal(errors.Wrap(err, "poll"))
	}
}
