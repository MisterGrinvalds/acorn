package cmd

import "github.com/spf13/cobra"

var networkCmd = &cobra.Command{
	Use:     "network",
	Short:   "Network and VPN tools",
	Long:    `Commands for managing network and VPN tools.`,
	Aliases: []string{"net", "n"},
}

func init() {
	rootCmd.AddCommand(networkCmd)
}
