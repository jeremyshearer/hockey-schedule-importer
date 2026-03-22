package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type inputGame struct {
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

func parseInput(r io.Reader, warn io.Writer) ([]inputGame, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}
	var games []inputGame
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
		games = append(games, inputGame{
			Date:     t,
			Visitor:  row[colVisitor],
			Home:     row[colHome],
			Location: row[colLocation],
		})
	}
	return games, nil
}

func marshalCSV(games []inputGame) ([]byte, error) {
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

func main() {
	inputPath := flag.String("in", "data/input.csv", "Input CSV file path")
	outputPath := flag.String("out", "data/output.csv", "Output CSV file path")
	flag.Parse()

	inFile, err := os.Open(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input: %v\n", err)
		os.Exit(1)
	}
	defer inFile.Close()

	inputGames, err := parseInput(inFile, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing input: %v\n", err)
		os.Exit(1)
	}

	csvBytes, err := marshalCSV(inputGames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*outputPath, csvBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Conversion complete. Output written to", *outputPath)
}
