package cmd

import "github.com/spf13/cobra"

var aiCmd = &cobra.Command{
	Use:     "ai",
	Short:   "AI and ML tools",
	Long:    `Commands for managing AI and machine learning tools.`,
	Aliases: []string{"ml"},
}

func init() {
	rootCmd.AddCommand(aiCmd)
}
