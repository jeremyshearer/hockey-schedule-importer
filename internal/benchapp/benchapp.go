package benchapp

import (
	"bytes"
	"encoding/csv"
)

type Game struct {
	Type     string
	GameType string
	Title    string
	Away     string
	Home     string
	Date     string
	Time     string
	Duration string
	Location string
	Address  string
	Notes    string
}

// WriteCSV accepts a slice of Game structs and returns the data as a csv []byte
func WriteCSV(games []Game) ([]byte, error) {
	var out []byte
	buffer := bytes.NewBuffer(out)
	writer := csv.NewWriter(buffer)
	header := []string{"Type", "Game Type", "Title (Optional)", "Away", "Home", "Date", "Time", "Duration", "Location (Optional)", "Address (Optional)", "Notes (Optional)"}
	if err := writer.Write(header); err != nil {
		return out, err
	}
	for _, g := range games {
		row := []string{
			g.Type,
			g.GameType,
			g.Title,
			g.Away,
			g.Home,
			g.Date,
			g.Time,
			g.Duration,
			g.Location,
			g.Address,
			g.Notes,
		}
		if err := writer.Write(row); err != nil {
			return out, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return out, err
	}
	return buffer.Bytes(), nil
}
