package cmd

import "github.com/spf13/cobra"

var Root = &cobra.Command{
	Use:   "hockey-schedule-importer",
	Short: "Convert hockey schedule CSVs between formats",
}
