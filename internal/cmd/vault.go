package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/components/data/vault"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	vaultVerbose bool
	vaultDryRun  bool
)

// vaultCmd represents the vault command group
var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "HashiCorp Vault secrets management",
	Long: `HashiCorp Vault CLI commands for secrets management.

Provides commands for managing secrets, authentication, and server status.

Examples:
  acorn vault status                    # Show Vault status
  acorn vault read secret/myapp/config  # Read a secret
  acorn vault list secret/              # List secrets`,
	Aliases: []string{"v"},
}

// vaultStatusCmd shows vault status
var vaultStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Vault installation and connection status",
	Long: `Display Vault CLI installation status, server connection, and authentication status.

Examples:
  acorn vault status
  acorn vault status -o json`,
	RunE: runVaultStatus,
}

// vaultReadCmd reads a secret
var vaultReadCmd = &cobra.Command{
	Use:   "read <path>",
	Short: "Read a secret from Vault",
	Long: `Read a secret from the specified path in Vault.

Examples:
  acorn vault read secret/myapp/database
  acorn vault read secret/myapp/api-keys`,
	Aliases: []string{"get"},
	Args:    cobra.ExactArgs(1),
	RunE:    runVaultRead,
}

// vaultWriteCmd writes a secret
var vaultWriteCmd = &cobra.Command{
	Use:   "write <path> <key=value> [key=value...]",
	Short: "Write a secret to Vault",
	Long: `Write a secret to the specified path in Vault.

Examples:
  acorn vault write secret/myapp/db password=secret123
  acorn vault write secret/myapp/api key=abc123 url=https://api.example.com`,
	Aliases: []string{"put"},
	Args:    cobra.MinimumNArgs(2),
	RunE:    runVaultWrite,
}

// vaultDeleteCmd deletes a secret
var vaultDeleteCmd = &cobra.Command{
	Use:   "delete <path>",
	Short: "Delete a secret from Vault",
	Long: `Delete a secret at the specified path in Vault.

Examples:
  acorn vault delete secret/myapp/old-key`,
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    runVaultDelete,
}

// vaultListCmd lists secrets
var vaultListCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "List secrets at a path",
	Long: `List all secrets at the specified path in Vault.

Defaults to listing secrets at 'secret/' if no path is provided.

Examples:
  acorn vault list
  acorn vault list secret/myapp`,
	Aliases: []string{"ls"},
	Args:    cobra.MaximumNArgs(1),
	RunE:    runVaultList,
}

// vaultLoginCmd authenticates to Vault
var vaultLoginCmd = &cobra.Command{
	Use:   "login <token>",
	Short: "Authenticate to Vault with a token",
	Long: `Authenticate to Vault using a token.

Examples:
  acorn vault login s.abc123xyz`,
	Args: cobra.ExactArgs(1),
	RunE: runVaultLogin,
}

// vaultRenewCmd renews the token
var vaultRenewCmd = &cobra.Command{
	Use:   "renew [increment]",
	Short: "Renew Vault token",
	Long: `Renew the current Vault authentication token.

Increment duration defaults to 1h if not specified.

Examples:
  acorn vault renew
  acorn vault renew 24h`,
	Args: cobra.MaximumNArgs(1),
	RunE: runVaultRenew,
}

// vaultEnvCmd sets the Vault environment
var vaultEnvCmd = &cobra.Command{
	Use:   "env <environment>",
	Short: "Set Vault server environment",
	Long: `Set the VAULT_ADDR environment variable for a specific environment.

Environments: prod, staging, dev

Examples:
  acorn vault env prod
  acorn vault env staging`,
	Args: cobra.ExactArgs(1),
	RunE: runVaultEnv,
}

func init() {
	dataCmd.AddCommand(vaultCmd)

	// Add subcommands
	vaultCmd.AddCommand(vaultStatusCmd)
	vaultCmd.AddCommand(vaultReadCmd)
	vaultCmd.AddCommand(vaultWriteCmd)
	vaultCmd.AddCommand(vaultDeleteCmd)
	vaultCmd.AddCommand(vaultListCmd)
	vaultCmd.AddCommand(vaultLoginCmd)
	vaultCmd.AddCommand(vaultRenewCmd)
	vaultCmd.AddCommand(vaultEnvCmd)

	// Persistent flags
	vaultCmd.PersistentFlags().BoolVarP(&vaultVerbose, "verbose", "v", false,
		"Show verbose output")
	vaultCmd.PersistentFlags().BoolVar(&vaultDryRun, "dry-run", false,
		"Show what would be done without executing")
}

func runVaultStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := vault.NewHelper(vaultVerbose, vaultDryRun)

	status := helper.GetStatus()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintln(os.Stdout, output.Info("Vault Status"))
	fmt.Fprintln(os.Stdout, strings.Repeat("=", 12))
	fmt.Fprintln(os.Stdout)

	if status.Installed {
		fmt.Fprintf(os.Stdout, "Version:   %s\n", status.Version)
		fmt.Fprintf(os.Stdout, "Location:  %s\n", status.Location)
		fmt.Fprintln(os.Stdout)

		if status.ServerAddress != "" {
			fmt.Fprintf(os.Stdout, "Server:    %s\n", status.ServerAddress)
			if status.Connected {
				fmt.Fprintf(os.Stdout, "Connected: %s\n", output.Success("Yes"))
				if status.Sealed {
					fmt.Fprintf(os.Stdout, "Sealed:    %s\n", output.Warning("Yes"))
				} else {
					fmt.Fprintf(os.Stdout, "Sealed:    %s\n", output.Success("No"))
				}
			} else {
				fmt.Fprintf(os.Stdout, "Connected: %s\n", output.Error("No"))
			}
			fmt.Fprintln(os.Stdout)

			if status.Authenticated {
				fmt.Fprintf(os.Stdout, "Auth:      %s\n", output.Success("Authenticated"))
			} else {
				fmt.Fprintf(os.Stdout, "Auth:      %s (run: acorn vault login)\n", output.Warning("Not authenticated"))
			}
		} else {
			fmt.Fprintf(os.Stdout, "Server:    %s\n", output.Warning("VAULT_ADDR not set"))
			fmt.Fprintln(os.Stdout)
			fmt.Fprintln(os.Stdout, "Set with: export VAULT_ADDR=https://vault.example.com:8200")
			fmt.Fprintln(os.Stdout, "Or use:   acorn vault env <prod|staging|dev>")
		}
	} else {
		fmt.Fprintf(os.Stdout, "Vault CLI: %s\n", output.Error("Not installed"))
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Install with: make vault-install")
	}

	return nil
}

func runVaultRead(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := vault.NewHelper(vaultVerbose, vaultDryRun)

	path := args[0]
	secret, err := helper.Read(path)
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(secret)
	}

	fmt.Fprintf(os.Stdout, "Secret at: %s\n", path)
	if len(secret.Data) > 0 {
		table := output.NewTable("KEY", "VALUE")
		for k, v := range secret.Data {
			table.AddRow(k, v)
		}
		table.Render(os.Stdout)
	} else {
		fmt.Fprintln(os.Stdout, "(no data)")
	}

	return nil
}

func runVaultWrite(cmd *cobra.Command, args []string) error {
	helper := vault.NewHelper(vaultVerbose, vaultDryRun)

	path := args[0]
	data := make(map[string]string)

	// Parse key=value pairs
	for _, arg := range args[1:] {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid key=value pair: %s", arg)
		}
		data[parts[0]] = parts[1]
	}

	if err := helper.Write(path, data); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Secret written to: %s\n", output.Success("✓"), path)
	return nil
}

func runVaultDelete(cmd *cobra.Command, args []string) error {
	helper := vault.NewHelper(vaultVerbose, vaultDryRun)

	path := args[0]
	if err := helper.Delete(path); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Secret deleted at: %s\n", output.Success("✓"), path)
	return nil
}

func runVaultList(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := vault.NewHelper(vaultVerbose, vaultDryRun)

	path := "secret/"
	if len(args) > 0 {
		path = args[0]
	}

	secrets, err := helper.List(path)
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(secrets)
	}

	fmt.Fprintf(os.Stdout, "Secrets at: %s\n\n", path)
	if len(secrets) > 0 {
		for _, s := range secrets {
			fmt.Fprintf(os.Stdout, "  %s\n", s)
		}
		fmt.Fprintf(os.Stdout, "\nTotal: %d secrets\n", len(secrets))
	} else {
		fmt.Fprintln(os.Stdout, "(no secrets)")
	}

	return nil
}

func runVaultLogin(cmd *cobra.Command, args []string) error {
	helper := vault.NewHelper(vaultVerbose, vaultDryRun)

	token := args[0]
	if err := helper.Login(token); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Authenticated to Vault\n", output.Success("✓"))
	return nil
}

func runVaultRenew(cmd *cobra.Command, args []string) error {
	helper := vault.NewHelper(vaultVerbose, vaultDryRun)

	increment := "1h"
	if len(args) > 0 {
		increment = args[0]
	}

	if err := helper.RenewToken(increment); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Token renewed for: %s\n", output.Success("✓"), increment)
	return nil
}

func runVaultEnv(cmd *cobra.Command, args []string) error {
	helper := vault.NewHelper(vaultVerbose, vaultDryRun)

	env := args[0]
	addr, err := helper.SetEnvironment(env)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Vault environment set to: %s\n", output.Success("✓"), env)
	fmt.Fprintf(os.Stdout, "\nRun the following command to apply:\n")
	fmt.Fprintf(os.Stdout, "  export VAULT_ADDR=%s\n", addr)

	return nil
}
