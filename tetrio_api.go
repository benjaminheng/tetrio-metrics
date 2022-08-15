package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type TetrioStreamResponse struct {
	Success bool
	Data    struct {
		Records []struct {
			ID         string `json:"_id"`
			ReplayID   string `json:"replayid"`
			Timestamp  string `json:"ts"`
			EndContext struct {
				PiecesPlaced int64   `json:"piecesplaced"`
				GameType     string  `json:"gametype"`
				FinalTime    float64 `json:"finalTime"`
				Finesse      struct {
					Faults        int64 `json:"faults"`
					PerfectPieces int64 `json:"perfectpieces"`
				}
			}
		}
	}
}

func getTetrioRecentUserStreams(ctx context.Context, userID string) (parsedResponse *TetrioStreamResponse, rawResponse string, err error) {
	// TODO: add retries
	url := "https://ch.tetr.io/api/streams/any_userrecent_" + userID
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, "", errors.Wrap(err, "get request")
	}
	req.Header.Add("User-Agent", "Bot to archive my personal replays (repo: github.com/benjaminheng/tetrio-metrics)")

	client := &http.Client{
		Timeout: 3000 * time.Millisecond,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", errors.Wrap(err, "call tetrio api")
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", errors.Wrap(err, "read body")
	}
	parsedResponse = &TetrioStreamResponse{}
	err = json.Unmarshal(b, parsedResponse)
	if err != nil {
		return nil, "", errors.Wrap(err, "unmarshal response")
	}
	return parsedResponse, string(b), nil
}
