package cmd

import "github.com/spf13/cobra"

var programmingCmd = &cobra.Command{
	Use:     "programming",
	Short:   "Programming languages",
	Long:    `Commands for managing programming language tools and environments.`,
	Aliases: []string{"lang", "p"},
}

func init() {
	rootCmd.AddCommand(programmingCmd)
}
