package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/benjaminheng/tetrio-metrics/store"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	PollIntervalSeconds         int64 // Poll interval for 40L game mode data
	UserInfoPollIntervalSeconds int64 // Poll interval for user info
	TetrioUserID                string
	DatabaseFilePath            string
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

func (s *Service) checkForNewTetrio40LGames(ctx context.Context) (err error) {
	log.Println("checking for recent tetrio games")
	defer func() {
		if err != nil {
			log.Println(errors.Wrap(err, "error"))
		}
	}()

	// Get recent games from tetrio
	parsedResponse, err := getTetrioRecentUserStreams(ctx, s.config.TetrioUserID)
	if err != nil {
		return errors.Wrap(err, "get tetrio recent user streams")
	}

	// Get the last inserted model
	lastSeen, err := s.store.GetLatestGamemode40L(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(err, "get latest gamemode_40l record")
	}

	// Construct models from tetrio's API response
	var modelsToInsert []store.Gamemode40l
	for _, v := range parsedResponse.Data.Records {
		if v.EndContext.GameType != "40l" {
			continue
		}
		m, err := buildGamemode40LRecord(v)
		if err != nil {
			log.Println(errors.Wrap(err, "build gamemode_40l model"))
		}
		if m.PlayedAt.After(lastSeen) {
			modelsToInsert = append(modelsToInsert, m)
		}
	}

	// Insert models to DB
	for _, m := range modelsToInsert {
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

func (s *Service) poll40LGameMode(ctx context.Context) error {
	s.checkForNewTetrio40LGames(ctx)
	ticker := time.NewTicker(time.Duration(s.config.PollIntervalSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			s.checkForNewTetrio40LGames(ctx)
		}
	}
}

func (s *Service) pollUserInfo(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(s.config.UserInfoPollIntervalSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			// TODO
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
	var wg sync.WaitGroup

	// Poll for 40l game mode data
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = service.poll40LGameMode(ctx)
		if err != nil {
			log.Fatal(errors.Wrap(err, "poll 40l game mode"))
		}
	}()

	// Poll for user info
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = service.pollUserInfo(ctx)
		if err != nil {
			log.Fatal(errors.Wrap(err, "poll user info"))
		}
	}()

	wg.Wait()
}
