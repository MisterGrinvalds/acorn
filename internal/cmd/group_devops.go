package cmd

import "github.com/spf13/cobra"

var devopsCmd = &cobra.Command{
	Use:     "devops",
	Short:   "DevOps and infrastructure",
	Long:    `Commands for managing DevOps tools and infrastructure.`,
	Aliases: []string{"infra", "d"},
}

func init() {
	rootCmd.AddCommand(devopsCmd)
}
