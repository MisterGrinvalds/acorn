package cmd

import "github.com/spf13/cobra"

var artifactsCmd = &cobra.Command{
	Use:     "artifacts",
	Short:   "Artifact repository tools",
	Long:    `Commands for managing artifact repositories and package registries.`,
	Aliases: []string{"art", "registry"},
}

func init() {
	rootCmd.AddCommand(artifactsCmd)
}
