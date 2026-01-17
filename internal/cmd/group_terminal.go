package cmd

import "github.com/spf13/cobra"

var terminalCmd = &cobra.Command{
	Use:     "terminal",
	Short:   "Terminal and shell tools",
	Long:    `Commands for managing terminal emulators and shell configuration.`,
	Aliases: []string{"term", "t"},
}

func init() {
	rootCmd.AddCommand(terminalCmd)
}
