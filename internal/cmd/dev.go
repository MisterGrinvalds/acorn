package cmd

import "github.com/spf13/cobra"

// devCmd represents the dev command group for development and tooling commands.
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Development and tooling commands",
	Long: `Development commands for building and maintaining the acorn CLI.

These commands help with:
  - Auditing commands for JSON-first compliance
  - Validating command implementations
  - Code generation and scaffolding`,
	Aliases: []string{"develop"},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
