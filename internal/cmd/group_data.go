package cmd

import "github.com/spf13/cobra"

var dataCmd = &cobra.Command{
	Use:     "data",
	Short:   "Database and secrets",
	Long:    `Commands for managing databases and secrets.`,
	Aliases: []string{"db"},
}

func init() {
	rootCmd.AddCommand(dataCmd)
}
