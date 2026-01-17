package cmd

import "github.com/spf13/cobra"

var ideCmd = &cobra.Command{
	Use:     "ide",
	Short:   "IDEs and editors",
	Long:    `Commands for managing IDEs and text editors.`,
	Aliases: []string{"editor", "e"},
}

func init() {
	rootCmd.AddCommand(ideCmd)
}
