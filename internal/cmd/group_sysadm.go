package cmd

import "github.com/spf13/cobra"

var sysadmCmd = &cobra.Command{
	Use:     "sysadm",
	Short:   "System administration tools",
	Long:    `Commands for system monitoring, resource management, and administration.`,
	Aliases: []string{"sys", "admin"},
}

func init() {
	rootCmd.AddCommand(sysadmCmd)
}
