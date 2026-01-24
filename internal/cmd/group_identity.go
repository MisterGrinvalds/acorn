package cmd

import (
	"github.com/spf13/cobra"
)

// identityCmd represents the identity command group
var identityCmd = &cobra.Command{
	Use:   "identity",
	Short: "Identity and access management tools",
	Long: `Identity and access management (IAM) tools.

Manage authentication providers, identity systems, and access control.

Commands:
  keycloak   - Keycloak IAM platform

Examples:
  acorn identity keycloak status
  acorn identity keycloak start`,
	Aliases: []string{"iam", "auth"},
}

func init() {
	rootCmd.AddCommand(identityCmd)
}
