package main

import (
	"context"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/benjaminheng/tetrio-metrics/store"
	"github.com/pkg/errors"
)

type TetrioStreamRecord struct {
	ID         string `json:"_id"`
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
	// RawResponse contains the full Tetrio API response as a JSON string.
	// This field is not part of tetrio's API response. It is derived and
	// injected later.
	RawResponse string
}

type TetrioStreamResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Records []TetrioStreamRecord `json:"records"`
	} `json:"data"`
}

type TetrioStreamRawResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Records []map[string]interface{} `json:"records"`
	} `json:"data"`
}

func getTetrioRecentUserStreams(ctx context.Context, userID string) (parsedResponse *TetrioStreamResponse, err error) {
	// TODO: add retries
	url := "https://ch.tetr.io/api/streams/any_userrecent_" + userID
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
	parsedResponse = &TetrioStreamResponse{}
	err = json.Unmarshal(b, parsedResponse)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal parsed response")
	}

	rawResponse := &TetrioStreamRawResponse{}
	err = json.Unmarshal(b, rawResponse)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal to raw response")
	}

	// Inject raw response into the parsed response
	if len(rawResponse.Data.Records) == len(parsedResponse.Data.Records) {
		for i, v := range rawResponse.Data.Records {
			b, err := json.Marshal(v)
			if err != nil {
				return nil, errors.Wrap(err, "marshal individual raw response record")
			}
			parsedResponse.Data.Records[i].RawResponse = string(b)
		}
	}

	return parsedResponse, nil
}

func buildGamemode40LRecord(apiRecord TetrioStreamRecord) (record store.Gamemode40l, err error) {
	rawData, err := json.Marshal(apiRecord)
	if err != nil {
		return record, errors.Wrap(err, "marshal tetrio stream record to json")
	}
	record.RawData.Scan(string(rawData))
	// TODO: raw data is wrong
	record.FinesseFaults = apiRecord.EndContext.Finesse.Faults
	record.TotalPieces = apiRecord.EndContext.Finesse.Faults + apiRecord.EndContext.Finesse.PerfectPieces
	finessePercent := 1 - (float64(record.FinesseFaults) / float64(record.TotalPieces))
	finessePercent = math.Round(finessePercent*100.0) / 100 // Round to 2 decimal places
	record.FinessePercent = finessePercent
	playedAt, err := time.Parse(time.RFC3339Nano, apiRecord.Timestamp)
	if err != nil {
		return record, errors.Wrap(err, "parse played_at to time.Time")
	}
	record.PlayedAt = playedAt
	// Tetrio returns timings in granular to 1/3 of a millisecond. We'll
	// just round the time down to the nearest millisecond for simplicity.
	record.TimeMs = int64(apiRecord.EndContext.FinalTime)

	return record, nil
}
