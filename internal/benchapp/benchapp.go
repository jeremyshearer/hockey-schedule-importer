package benchapp

import (
	"encoding/csv"
	"io"
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

// WriteCSV writes a slice of Game to the given writer as CSV, including the header row.
func WriteCSV(w io.Writer, games []Game) error {
	writer := csv.NewWriter(w)
	header := []string{"Type", "Game Type", "Title (Optional)", "Away", "Home", "Date", "Time", "Duration", "Location (Optional)", "Address (Optional)", "Notes (Optional)"}
	if err := writer.Write(header); err != nil {
		return err
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
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
