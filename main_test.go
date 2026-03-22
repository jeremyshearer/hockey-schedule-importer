package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

func TestParseInput(t *testing.T) {
	input := `Date,Visitor,,Details,,Home,Location
"2025-03-23T19:20:00.000Z","Mid Ice Crisis","Mid Ice Crisis","L 6 - 1","HatTrick Swayzes","HatTrick Swayzes","North"
"2025-04-13T19:00:00.000Z","HatTrick Swayzes","HatTrick Swayzes","W 3 - 2","Puck Me","Puck Me","South"
`
	games, err := parseInput(strings.NewReader(input), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 2 {
		t.Fatalf("expected 2 games, got %d", len(games))
	}
	g := games[0]
	if g.Visitor != "Mid Ice Crisis" {
		t.Errorf("visitor = %q, want %q", g.Visitor, "Mid Ice Crisis")
	}
	if g.Home != "HatTrick Swayzes" {
		t.Errorf("home = %q, want %q", g.Home, "HatTrick Swayzes")
	}
	if g.Location != "North" {
		t.Errorf("location = %q, want %q", g.Location, "North")
	}
}

func TestParseInputSkipsShortRows(t *testing.T) {
	input := `Date,Visitor,,Details,,Home,Location
"2025-03-23T19:20:00.000Z","Only Two Cols"
"2025-04-13T19:00:00.000Z","HatTrick Swayzes","","","Puck Me","Puck Me","South"
`
	var warn bytes.Buffer
	games, err := parseInput(strings.NewReader(input), &warn)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 1 {
		t.Fatalf("expected 1 game (short row skipped), got %d", len(games))
	}
	if !strings.Contains(warn.String(), "skipping row 2") {
		t.Errorf("expected warning about row 2, got %q", warn.String())
	}
}

func TestParseInputSkipsBadDates(t *testing.T) {
	input := `Date,Visitor,,Details,,Home,Location
"not-a-date","Team A","","","Team B","Team B","North"
`
	var warn bytes.Buffer
	games, err := parseInput(strings.NewReader(input), &warn)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 0 {
		t.Fatalf("expected 0 games (bad date skipped), got %d", len(games))
	}
	if !strings.Contains(warn.String(), "invalid date") {
		t.Errorf("expected warning about invalid date, got %q", warn.String())
	}
}

func TestMarshalCSV(t *testing.T) {
	games := []inputGame{{
		Date:     time.Date(2025, 3, 23, 19, 20, 0, 0, time.UTC),
		Visitor:  "Team A",
		Home:     "Team B",
		Location: "North",
	}}
	data, err := marshalCSV(games)
	if err != nil {
		t.Fatal(err)
	}
	csv := string(data)
	if !strings.Contains(csv, "Type,Game Type") {
		t.Error("missing header row")
	}
	if !strings.Contains(csv, "GAME,REGULAR,,Team A,Team B,23/03/2025,7:20 PM,1:30,North,,") {
		t.Errorf("unexpected output row: %s", csv)
	}
}

func TestRoundTrip(t *testing.T) {
	input := `Date,Visitor,,Details,,Home,Location
"2025-03-23T19:20:00.000Z","Mid Ice Crisis","Mid Ice Crisis","L 6 - 1","HatTrick Swayzes","HatTrick Swayzes","North"
"2025-04-13T19:00:00.000Z","HatTrick Swayzes","HatTrick Swayzes","W 3 - 2","Puck Me","Puck Me","South"
`
	games, err := parseInput(strings.NewReader(input), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	data, err := marshalCSV(games)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header + 2 games), got %d", len(lines))
	}
	if !strings.Contains(lines[1], "Mid Ice Crisis") {
		t.Errorf("line 1 missing team name: %s", lines[1])
	}
	if !strings.Contains(lines[2], "Puck Me") {
		t.Errorf("line 2 missing team name: %s", lines[2])
	}
}
