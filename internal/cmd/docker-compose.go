package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/docker"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	composeVerbose       bool
	composeDryRun        bool
	composeFile          string
	composeFollow        bool
	composeDetach        bool
	composeBuild         bool
	composeRemoveVolumes bool
	composeRemoveOrphans bool
	composeScale         int
	composeProfile       string
)

// composeCmd represents the docker-compose command group
var composeCmd = &cobra.Command{
	Use:   "docker-compose",
	Short: "Docker Compose commands",
	Long: `Docker Compose commands for multi-container applications.

Provides service lifecycle management, logging, and debugging.

Examples:
  acorn docker-compose up           # Start services
  acorn docker-compose down         # Stop services
  acorn docker-compose logs         # View service logs
  acorn docker-compose ps           # List services
  acorn docker-compose status       # Show compose status`,
	Aliases: []string{"compose", "dco"},
}

// composeStatusCmd shows compose status
var composeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Docker Compose status",
	Long: `Display Docker Compose installation and project status.

Shows compose version and current project services.

Examples:
  acorn docker-compose status
  acorn docker-compose status -o json`,
	RunE: runComposeStatus,
}

// composePsCmd lists compose services
var composePsCmd = &cobra.Command{
	Use:   "ps",
	Short: "List compose services",
	Long: `List Docker Compose services and their status.

Examples:
  acorn docker-compose ps
  acorn docker-compose ps -f docker-compose.dev.yml`,
	Aliases: []string{"services"},
	RunE:    runComposePs,
}

// composeUpCmd starts compose services
var composeUpCmd = &cobra.Command{
	Use:   "up [services...]",
	Short: "Start compose services",
	Long: `Start Docker Compose services.

Examples:
  acorn docker-compose up
  acorn docker-compose up -d
  acorn docker-compose up --build
  acorn docker-compose up web api`,
	Aliases: []string{"start"},
	RunE:    runComposeUp,
}

// composeDownCmd stops compose services
var composeDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop compose services",
	Long: `Stop and remove Docker Compose services.

Examples:
  acorn docker-compose down
  acorn docker-compose down -v
  acorn docker-compose down --remove-orphans`,
	Aliases: []string{"stop"},
	RunE:    runComposeDown,
}

// composeLogsCmd shows compose logs
var composeLogsCmd = &cobra.Command{
	Use:   "logs [services...]",
	Short: "Show compose logs",
	Long: `Display logs from Docker Compose services.

Examples:
  acorn docker-compose logs
  acorn docker-compose logs -f
  acorn docker-compose logs web api`,
	RunE: runComposeLogs,
}

// composeRestartCmd restarts compose services
var composeRestartCmd = &cobra.Command{
	Use:   "restart [services...]",
	Short: "Restart compose services",
	Long: `Restart Docker Compose services.

Examples:
  acorn docker-compose restart
  acorn docker-compose restart web`,
	RunE: runComposeRestart,
}

// composeBuildCmd builds compose services
var composeBuildCmd = &cobra.Command{
	Use:   "build [services...]",
	Short: "Build compose images",
	Long: `Build Docker Compose service images.

Examples:
  acorn docker-compose build
  acorn docker-compose build web api`,
	RunE: runComposeBuild,
}

// composePullCmd pulls compose images
var composePullCmd = &cobra.Command{
	Use:   "pull [services...]",
	Short: "Pull compose images",
	Long: `Pull Docker Compose service images.

Examples:
  acorn docker-compose pull
  acorn docker-compose pull web`,
	RunE: runComposePull,
}

// composeExecCmd executes command in service
var composeExecCmd = &cobra.Command{
	Use:   "exec [service] [command...]",
	Short: "Execute command in service",
	Long: `Execute a command inside a running compose service.

Examples:
  acorn docker-compose exec web /bin/sh
  acorn docker-compose exec db psql -U postgres`,
	Args: cobra.MinimumNArgs(2),
	RunE: runComposeExec,
}

// composeValidateCmd validates compose config (passthrough to docker compose config)
var composeValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate compose configuration",
	Long: `Validate and view resolved Docker Compose configuration.

Runs 'docker compose config' to resolve and validate the compose file.

Examples:
  acorn docker-compose validate
  acorn docker-compose validate -f docker-compose.dev.yml`,
	RunE: runComposeConfig,
}

func init() {

	// Add subcommands
	composeCmd.AddCommand(composeStatusCmd)
	composeCmd.AddCommand(composePsCmd)
	composeCmd.AddCommand(composeUpCmd)
	composeCmd.AddCommand(composeDownCmd)
	composeCmd.AddCommand(composeLogsCmd)
	composeCmd.AddCommand(composeRestartCmd)
	composeCmd.AddCommand(composeBuildCmd)
	composeCmd.AddCommand(composePullCmd)
	composeCmd.AddCommand(composeExecCmd)
	composeCmd.AddCommand(composeValidateCmd)
	composeCmd.AddCommand(configcmd.NewConfigRouter("docker-compose"))

	// Persistent flags
	composeCmd.PersistentFlags().BoolVarP(&composeVerbose, "verbose", "v", false,
		"Show verbose output")
	composeCmd.PersistentFlags().BoolVar(&composeDryRun, "dry-run", false,
		"Show what would be done without executing")
	composeCmd.PersistentFlags().StringVarP(&composeFile, "file", "f", "",
		"Compose file path")
	composeCmd.PersistentFlags().StringVar(&composeProfile, "profile", "",
		"Specify a profile to enable")

	// Command-specific flags
	composeUpCmd.Flags().BoolVarP(&composeDetach, "detach", "d", false, "Run in background")
	composeUpCmd.Flags().BoolVar(&composeBuild, "build", false, "Build images before starting")
	composeUpCmd.Flags().IntVar(&composeScale, "scale", 0, "Scale service to N instances")

	composeDownCmd.Flags().BoolVarP(&composeRemoveVolumes, "volumes", "v", false, "Remove volumes")
	composeDownCmd.Flags().BoolVar(&composeRemoveOrphans, "remove-orphans", false, "Remove orphan containers")

	composeLogsCmd.Flags().BoolVarP(&composeFollow, "follow", "f", false, "Follow log output")
}

func runComposeStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := docker.NewHelper(composeVerbose, composeDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	status := helper.GetComposeStatus(composeFile)

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	fmt.Fprintf(os.Stdout, "Compose Installed: %v\n", status.Installed)
	if status.Version != "" {
		fmt.Fprintf(os.Stdout, "Version:           %s\n", status.Version)
	}
	if status.ProjectName != "" {
		fmt.Fprintf(os.Stdout, "Project:           %s\n", status.ProjectName)
	}
	if status.ComposeFile != "" {
		fmt.Fprintf(os.Stdout, "File:              %s\n", status.ComposeFile)
	}
	fmt.Fprintf(os.Stdout, "Services:          %d\n", status.ServiceCount)

	return nil
}

func runComposePs(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := docker.NewHelper(composeVerbose, composeDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	services, err := helper.GetComposeServices(composeFile)
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"services": services})
	}

	if len(services) == 0 {
		fmt.Fprintln(os.Stdout, "No services found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-20s %-15s %-30s\n", "SERVICE", "STATUS", "PORTS")
	for _, s := range services {
		fmt.Fprintf(os.Stdout, "%-20s %-15s %-30s\n", s.Name, s.Status, s.Ports)
	}

	return nil
}

func runComposeUp(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(composeVerbose, composeDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	if err := helper.ComposeUp(composeFile, composeDetach, composeBuild, args); err != nil {
		return err
	}

	if !composeDetach {
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s Services started\n", output.Success("✓"))
	return nil
}

func runComposeDown(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(composeVerbose, composeDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	if err := helper.ComposeDown(composeFile, composeRemoveVolumes, composeRemoveOrphans); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Services stopped\n", output.Success("✓"))
	return nil
}

func runComposeLogs(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(composeVerbose, composeDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	return helper.GetComposeLogs(composeFile, composeFollow, args)
}

func runComposeRestart(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(composeVerbose, composeDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	if err := helper.ComposeRestart(composeFile, args); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Services restarted\n", output.Success("✓"))
	return nil
}

func runComposeBuild(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(composeVerbose, composeDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	if err := helper.ComposeBuild(composeFile, args); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Build complete\n", output.Success("✓"))
	return nil
}

func runComposePull(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(composeVerbose, composeDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	if err := helper.ComposePull(composeFile, args); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Pull complete\n", output.Success("✓"))
	return nil
}

func runComposeExec(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(composeVerbose, composeDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	return helper.ComposeExec(composeFile, args[0], args[1:])
}

func runComposeConfig(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(composeVerbose, composeDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	return helper.ComposeConfig(composeFile)
}

func init() {
	components.Register(&components.Registration{
		Name: "docker-compose",
		RegisterCmd: func() *cobra.Command { return composeCmd },
	})
}
