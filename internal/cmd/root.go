// Package cmd provides all Cobra commands for the acorn CLI.
package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/terminal/ghostty"
	"github.com/mistergrinvalds/acorn/internal/components/terminal/iterm2"
	"github.com/mistergrinvalds/acorn/internal/components/terminal/tmux"
	"github.com/mistergrinvalds/acorn/internal/components/vcs/git"
	"github.com/mistergrinvalds/acorn/internal/utils/config"
	"github.com/mistergrinvalds/acorn/internal/utils/configfile"
	"github.com/mistergrinvalds/acorn/internal/utils/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// Register component-specific config file writers
	configfile.Register(ghostty.NewWriter())
	configfile.Register(tmux.NewWriter())
	configfile.Register(git.NewConfigWriter())
	configfile.Register(iterm2.NewWriter())
}

var (
	cfgFile string
	debug   bool
	cfg     *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "acorn",
	Short: "A powerful CLI for managing development environments",
	Long: `Acorn is a comprehensive CLI tool for managing your development environment,
dotfiles, and developer tooling.

It provides commands for:
  - Dotfiles management and synchronization
  - Development environment setup
  - Tool version management
  - Session and workspace management

Built with love from a collection of battle-tested shell scripts.`,
	// Uncomment the following line if your bare application has an action
	// Run: func(cmd *cobra.Command, args []string) { },
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags available to all commands
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $XDG_CONFIG_HOME/acorn/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false,
		"enable debug output")

	// Bind flags to viper
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	// Override with flag value if set
	if debug {
		cfg.Debug = true
	}
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print detailed version information including build metadata.`,
	Run: func(cmd *cobra.Command, args []string) {
		info := version.Get()
		short, _ := cmd.Flags().GetBool("short")
		if short {
			fmt.Println(info.Short())
		} else {
			fmt.Println(info.String())
		}
	},
}

func init() {
	versionCmd.Flags().BoolP("short", "s", false, "print only the version number")
}

// GetConfig returns the loaded configuration
func GetConfig() *config.Config {
	return cfg
}
