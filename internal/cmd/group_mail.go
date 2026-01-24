package cmd

import "github.com/spf13/cobra"

var mailCmd = &cobra.Command{
	Use:     "mail",
	Short:   "Email and communication tools",
	Long:    `Commands for managing email clients and mail configuration.`,
	Aliases: []string{"email"},
}

func init() {
	rootCmd.AddCommand(mailCmd)
}
