package cmd

import "github.com/spf13/cobra"

var vcsCmd = &cobra.Command{
	Use:     "vcs",
	Short:   "Version control",
	Long:    `Commands for managing version control systems.`,
	Aliases: []string{"vc", "v"},
}

func init() {
	rootCmd.AddCommand(vcsCmd)
}
