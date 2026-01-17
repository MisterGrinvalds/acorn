package cmd

import "github.com/spf13/cobra"

var cloudCmd = &cobra.Command{
	Use:     "cloud",
	Short:   "Cloud providers",
	Long:    `Commands for managing cloud provider tools and services.`,
	Aliases: []string{"c"},
}

func init() {
	rootCmd.AddCommand(cloudCmd)
}
