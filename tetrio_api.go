package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/benjaminheng/tetrio-metrics/store"
	"github.com/pkg/errors"
)

type TetrioResponse_Entry struct {
	ID        string `json:"_id"`
	GameMode  string `json:"gamemode"`
	Timestamp string `json:"ts"`
	Results   struct {
		Stats struct {
			PiecesPlaced int64   `json:"piecesplaced"`
			FinalTime    float64 `json:"finaltime"`
			Finesse      struct {
				Faults        int64 `json:"faults"`
				PerfectPieces int64 `json:"perfectpieces"`
			} `json:"finesse"`
		} `json:"stats"`
	} `json:"results"`
}

type TetrioResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Entries []TetrioResponse_Entry `json:"entries"`
	} `json:"data"`
}

func getTetrioRecentUserStreams(ctx context.Context, username string) (parsedResponse *TetrioResponse, err error) {
	// TODO: add retries
	url := fmt.Sprintf("https://ch.tetr.io/api/users/%s/records/40l/recent?limit=50", username)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "get request")
	}
	req.Header.Add("User-Agent", "Bot to archive my personal replays (repo: github.com/benjaminheng/tetrio-metrics)")

	client := &http.Client{
		Timeout: 3000 * time.Millisecond,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "call tetrio api")
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body")
	}
	parsedResponse = &TetrioResponse{}
	err = json.Unmarshal(b, parsedResponse)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal parsed response")
	}

	return parsedResponse, nil
}

func buildGamemode40LRecord(entry TetrioResponse_Entry) (record store.Gamemode40l, err error) {
	rawData, err := json.Marshal(entry)
	if err != nil {
		return record, errors.Wrap(err, "marshal tetrio stream record to json")
	}
	record.RawData.Scan(string(rawData))
	record.TimeMs = int64(math.Round(entry.Results.Stats.FinalTime))
	record.TotalPieces = entry.Results.Stats.PiecesPlaced
	finessePercent := (float64(entry.Results.Stats.Finesse.PerfectPieces) / float64(record.TotalPieces)) * 100
	finessePercent = math.Round(finessePercent*100.0) / 100 // Round to 2 decimal places
	record.FinessePercent = finessePercent
	playedAt, err := time.Parse(time.RFC3339Nano, entry.Timestamp)
	if err != nil {
		return record, errors.Wrap(err, "parse played_at to time.Time")
	}
	pps := (float64(record.TotalPieces) * 1000) / float64(record.TimeMs)
	pps = math.Round(pps*100.0) / 100 // Round to 2 decimal places
	record.PiecesPerSecond = pps
	record.PlayedAt = playedAt
	return record, nil
}
