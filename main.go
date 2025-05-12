package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jeremyshearer/hockey-schedule-importer/internal/benchapp"
	"github.com/jeremyshearer/hockey-schedule-importer/internal/converter"
	gs "github.com/jeremyshearer/hockey-schedule-importer/internal/gamesheets"
)

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

	inputGames, err := gs.ParseFromCSV(inFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing input CSV: %v\n", err)
		os.Exit(1)
	}

	outputGames, err := converter.Convert(inputGames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting games: %v\n", err)
		os.Exit(1)
	}

	// Use output path as provided, relative to current working directory
	outFile, err := os.Create(*outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	err = benchapp.WriteCSV(outFile, outputGames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Conversion complete. Output written to", *outputPath)
}
