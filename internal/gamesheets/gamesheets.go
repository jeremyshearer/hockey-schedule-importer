package gamesheets

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"
)

type Game struct {
	Date     time.Time
	Visitor  string
	Details  string
	Home     string
	Home2    string
	Location string
}

// ParseFromCSV parses games from a CSV reader and returns a slice of Game.
func ParseFromCSV(r io.Reader) ([]Game, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}
	games := make([]Game, 0)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 7 {
			continue
		}
		gameTime, err := time.Parse(time.RFC3339, strings.Trim(row[0], "\""))
		if err != nil {
			continue
		}
		games = append(games, Game{
			Date:     gameTime,
			Visitor:  strings.Trim(row[1], "\""),
			Details:  strings.Trim(row[3], "\""),
			Home:     strings.Trim(row[4], "\""),
			Home2:    strings.Trim(row[5], "\""),
			Location: strings.Trim(row[6], "\""),
		})
	}
	return games, nil
}
