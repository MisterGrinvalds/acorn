package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/components/data/secrets"
	"github.com/spf13/cobra"
)

var (
	secretsOutputFormat string
	secretsVerbose      bool
)

// secretsCmd represents the secrets command group
var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Secrets management and credential checking",
	Long: `Manage secrets and check credential availability.

Provides commands for managing secrets stored in a secure .env file
and checking the availability of cloud provider credentials.

Examples:
  acorn secrets status           # Check secrets file status
  acorn secrets list             # List configured secret keys
  acorn secrets check            # Check all credentials
  acorn secrets check aws        # Check specific credential`,
}

// secretsStatusCmd shows secrets file status
var secretsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check secrets file status",
	Long: `Check the status of the secrets file.

Shows whether the file exists, is readable, and how many keys are defined.

Examples:
  acorn secrets status
  acorn secrets status -o json`,
	RunE: runSecretsStatus,
}

// secretsListCmd lists secret keys
var secretsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured secret keys",
	Long: `List all secret keys defined in the secrets file.

Only shows key names, not values, for security.

Examples:
  acorn secrets list
  acorn secrets list -o json`,
	RunE: runSecretsList,
}

// secretsLoadCmd loads secrets into environment
var secretsLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load secrets into environment",
	Long: `Load secrets from the secrets file into the current process environment.

Note: This only affects the current process. To load secrets into your shell,
use the shell function 'load_secrets' instead.

Examples:
  acorn secrets load`,
	RunE: runSecretsLoad,
}

// secretsCheckCmd checks credentials
var secretsCheckCmd = &cobra.Command{
	Use:   "check [credential]",
	Short: "Check credential availability",
	Long: `Check if credentials are available in the environment.

Without arguments, checks all known credentials. With an argument,
checks a specific credential.

Supported credentials: aws, azure, github, digitalocean, openai, anthropic, huggingface

Examples:
  acorn secrets check            # Check all
  acorn secrets check aws        # Check AWS only
  acorn secrets check github     # Check GitHub only`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSecretsCheck,
}

// secretsValidateCmd validates secrets
var secretsValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate common credentials",
	Long: `Validate that common credentials are configured.

Checks AWS, Azure, GitHub, DigitalOcean, and AI service credentials.

Examples:
  acorn secrets validate
  acorn secrets validate -o json`,
	RunE: runSecretsValidate,
}

// secretsInitCmd creates secrets file
var secretsInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize secrets file",
	Long: `Create a new secrets file with template.

Creates the secrets directory and file with secure permissions (0600)
and a template showing common secret keys.

Examples:
  acorn secrets init`,
	RunE: runSecretsInit,
}

// secretsPathCmd shows secrets path
var secretsPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show secrets file path",
	Long: `Display the path to the secrets file.

Examples:
  acorn secrets path`,
	RunE: runSecretsPath,
}

func init() {
	rootCmd.AddCommand(secretsCmd)

	// Add subcommands
	secretsCmd.AddCommand(secretsStatusCmd)
	secretsCmd.AddCommand(secretsListCmd)
	secretsCmd.AddCommand(secretsLoadCmd)
	secretsCmd.AddCommand(secretsCheckCmd)
	secretsCmd.AddCommand(secretsValidateCmd)
	secretsCmd.AddCommand(secretsInitCmd)
	secretsCmd.AddCommand(secretsPathCmd)

	// Persistent flags
	secretsCmd.PersistentFlags().StringVarP(&secretsOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	secretsCmd.PersistentFlags().BoolVarP(&secretsVerbose, "verbose", "v", false,
		"Show verbose output")
}

func runSecretsStatus(cmd *cobra.Command, args []string) error {
	helper := secrets.NewHelper(secretsVerbose)
	status, err := helper.GetStatus()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(secretsOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Secrets Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintf(os.Stdout, "File: %s\n", status.FilePath)

	if !status.Exists {
		fmt.Fprintf(os.Stdout, "%s Secrets file not found\n", output.Warning("⚠"))
		fmt.Fprintln(os.Stdout, "  Run: acorn secrets init")
		return nil
	}

	if !status.Readable {
		fmt.Fprintf(os.Stdout, "%s Cannot read secrets file (check permissions)\n", output.Error("✗"))
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s Secrets file exists and is readable\n", output.Success("✓"))
	fmt.Fprintf(os.Stdout, "Keys defined: %d\n", status.KeyCount)

	return nil
}

func runSecretsList(cmd *cobra.Command, args []string) error {
	helper := secrets.NewHelper(secretsVerbose)
	keys, err := helper.ListSecrets()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(secretsOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(map[string][]string{"keys": keys})
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Configured Secrets"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(keys) == 0 {
		fmt.Fprintln(os.Stdout, "No secrets configured")
		return nil
	}

	for _, key := range keys {
		fmt.Fprintf(os.Stdout, "  %s\n", key)
	}
	fmt.Fprintf(os.Stdout, "\nTotal: %d keys\n", len(keys))

	return nil
}

func runSecretsLoad(cmd *cobra.Command, args []string) error {
	helper := secrets.NewHelper(secretsVerbose)
	count, err := helper.LoadSecrets()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Loaded %d secrets into environment\n", output.Success("✓"), count)
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Note: Secrets are only loaded into this process.")
	fmt.Fprintln(os.Stdout, "To load into your shell, use: load_secrets")

	return nil
}

func runSecretsCheck(cmd *cobra.Command, args []string) error {
	helper := secrets.NewHelper(secretsVerbose)

	format, err := output.ParseFormat(secretsOutputFormat)
	if err != nil {
		return err
	}

	// Check specific credential if provided
	if len(args) > 0 {
		var cred *secrets.Credential

		switch args[0] {
		case "aws":
			cred = helper.CheckAWS()
		case "azure":
			cred = helper.CheckAzure()
		case "github":
			cred = helper.CheckGitHub()
		case "digitalocean", "do":
			cred = helper.CheckDigitalOcean()
		case "openai":
			cred = helper.CheckOpenAI()
		case "anthropic":
			cred = helper.CheckAnthropic()
		case "huggingface", "hf":
			cred = helper.CheckHuggingFace()
		default:
			return fmt.Errorf("unknown credential: %s", args[0])
		}

		if format != output.FormatTable {
			printer := output.NewPrinter(os.Stdout, format)
			return printer.Print(cred)
		}

		if cred.Available {
			fmt.Fprintf(os.Stdout, "%s %s credentials: available\n", output.Success("✓"), cred.Name)
		} else {
			fmt.Fprintf(os.Stdout, "%s %s credentials: not found\n", output.Error("✗"), cred.Name)
			fmt.Fprintf(os.Stdout, "  Required: %s\n", joinEnvVars(cred.EnvVars))
		}
		return nil
	}

	// Check all credentials
	check := helper.CheckAllCredentials()

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(check)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Credential Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, cred := range check.Credentials {
		if cred.Available {
			fmt.Fprintf(os.Stdout, "%s %s: available\n", output.Success("✓"), cred.Name)
		} else {
			fmt.Fprintf(os.Stdout, "%s %s: not found\n", output.Error("✗"), cred.Name)
		}
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Available: %d, Missing: %d\n", check.Available, check.Missing)

	return nil
}

func runSecretsValidate(cmd *cobra.Command, args []string) error {
	helper := secrets.NewHelper(secretsVerbose)
	check := helper.ValidateSecrets()

	format, err := output.ParseFormat(secretsOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(check)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Secrets Validation"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, cred := range check.Credentials {
		if cred.Available {
			fmt.Fprintf(os.Stdout, "%s %s: configured\n", output.Success("✓"), cred.Name)
		} else {
			fmt.Fprintf(os.Stdout, "%s %s: missing\n", output.Warning("⚠"), cred.Name)
		}
	}

	fmt.Fprintln(os.Stdout)
	if check.Missing == 0 {
		fmt.Fprintf(os.Stdout, "%s All credentials configured\n", output.Success("✓"))
	} else {
		fmt.Fprintf(os.Stdout, "%s %d credentials missing\n", output.Warning("⚠"), check.Missing)
	}

	return nil
}

func runSecretsInit(cmd *cobra.Command, args []string) error {
	helper := secrets.NewHelper(secretsVerbose)

	if err := helper.CreateSecretsFile(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Secrets file created\n", output.Success("✓"))
	fmt.Fprintf(os.Stdout, "Path: %s\n", helper.GetSecretsFile())
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Edit the file to add your credentials.")
	fmt.Fprintln(os.Stdout, "Then run 'load_secrets' to load them into your shell.")

	return nil
}

func runSecretsPath(cmd *cobra.Command, args []string) error {
	helper := secrets.NewHelper(secretsVerbose)
	fmt.Fprintln(os.Stdout, helper.GetSecretsFile())
	return nil
}

func joinEnvVars(vars []string) string {
	if len(vars) == 0 {
		return ""
	}
	result := vars[0]
	for i := 1; i < len(vars); i++ {
		result += ", " + vars[i]
	}
	return result
}
