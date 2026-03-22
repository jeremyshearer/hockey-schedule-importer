package converter

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

type Game struct {
	Date     time.Time
	Visitor  string
	Home     string
	Location string
}

// Column indices in the gamesheets CSV.
const (
	colDate     = 0
	colVisitor  = 1
	colHome     = 4
	colLocation = 6
)

func ParseInput(r io.Reader, warn io.Writer) ([]Game, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}
	var games []Game
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 7 {
			fmt.Fprintf(warn, "warning: skipping row %d: expected 7+ columns, got %d\n", i+1, len(row))
			continue
		}
		t, err := time.Parse(time.RFC3339, row[colDate])
		if err != nil {
			fmt.Fprintf(warn, "warning: skipping row %d: invalid date %q\n", i+1, row[colDate])
			continue
		}
		games = append(games, Game{
			Date:     t,
			Visitor:  row[colVisitor],
			Home:     row[colHome],
			Location: row[colLocation],
		})
	}
	return games, nil
}

func MarshalCSV(games []Game) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	header := []string{"Type", "Game Type", "Title (Optional)", "Away", "Home", "Date", "Time", "Duration", "Location (Optional)", "Address (Optional)", "Notes (Optional)"}
	if err := w.Write(header); err != nil {
		return nil, err
	}
	for _, g := range games {
		row := []string{
			"GAME", "REGULAR", "",
			g.Visitor, g.Home,
			g.Date.Format("02/01/2006"), g.Date.Format("3:04 PM"),
			"1:30", g.Location, "", "",
		}
		if err := w.Write(row); err != nil {
			return nil, err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Convert(r io.Reader, warn io.Writer) ([]byte, error) {
	games, err := ParseInput(r, warn)
	if err != nil {
		return nil, err
	}
	return MarshalCSV(games)
}
