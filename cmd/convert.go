package cmd

import (
	"fmt"
	"os"

	"github.com/jeremyshearer/hockey-schedule-importer/converter"
	"github.com/spf13/cobra"
)

var (
	inputPath  string
	outputPath string
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert a gamesheets CSV to an importable schedule CSV",
	RunE: func(cmd *cobra.Command, args []string) error {
		inFile, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("opening input: %w", err)
		}
		defer inFile.Close()

		csvBytes, err := converter.Convert(inFile, os.Stderr)
		if err != nil {
			return fmt.Errorf("converting: %w", err)
		}

		if err := os.WriteFile(outputPath, csvBytes, 0644); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}
		fmt.Println("Conversion complete. Output written to", outputPath)
		return nil
	},
}

func init() {
	convertCmd.Flags().StringVar(&inputPath, "in", "data/input.csv", "Input CSV file path")
	convertCmd.Flags().StringVar(&outputPath, "out", "data/output.csv", "Output CSV file path")
	Root.AddCommand(convertCmd)
}
