package cmd

import (
	"github.com/spf13/cobra"
)

// automationCmd represents the automation command group
var automationCmd = &cobra.Command{
	Use:   "automation",
	Short: "Workflow automation tools",
	Long: `Workflow automation and integration tools.

Manage automation platforms like n8n for building workflows.

Commands:
  n8n        - n8n workflow automation

Examples:
  acorn automation n8n status
  acorn automation n8n start`,
	Aliases: []string{"auto"},
}

func init() {
	rootCmd.AddCommand(automationCmd)
}
