package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jeremyshearer/hockey-schedule-importer/converter"
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

	csvBytes, err := converter.Convert(inFile, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*outputPath, csvBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Conversion complete. Output written to", *outputPath)
}
