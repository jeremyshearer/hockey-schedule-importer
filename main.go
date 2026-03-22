package main

import (
	"fmt"
	"os"

	"github.com/jeremyshearer/hockey-schedule-importer/cmd"
)

func main() {
	if err := cmd.Root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
