package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/devops/docker"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	dockerVerbose bool
	dockerDryRun  bool
	dockerAll     bool
)

// dockerCmd represents the docker command group
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Docker helper commands",
	Long: `Docker helper commands for container management.

Provides container, image, volume, and network operations.

Examples:
  acorn docker status         # Show Docker status
  acorn docker ps             # List containers
  acorn docker images         # List images
  acorn docker clean          # Clean unused resources
  acorn docker compose up     # Docker compose up`,
	Aliases: []string{"dk", "d"},
}

// dockerStatusCmd shows Docker status
var dockerStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Docker daemon status",
	Long: `Display Docker installation and daemon status.

Shows version, running containers, and image counts.

Examples:
  acorn docker status
  acorn docker status -o json`,
	RunE: runDockerStatus,
}

// dockerPsCmd lists containers
var dockerPsCmd = &cobra.Command{
	Use:   "ps",
	Short: "List containers",
	Long: `List Docker containers.

Use -a to show all containers including stopped ones.

Examples:
  acorn docker ps
  acorn docker ps -a`,
	Aliases: []string{"containers"},
	RunE:    runDockerPs,
}

// dockerImagesCmd lists images
var dockerImagesCmd = &cobra.Command{
	Use:   "images",
	Short: "List images",
	Long: `List Docker images.

Examples:
  acorn docker images
  acorn docker images -o json`,
	RunE: runDockerImages,
}

// dockerVolumesCmd lists volumes
var dockerVolumesCmd = &cobra.Command{
	Use:   "volumes",
	Short: "List volumes",
	Long: `List Docker volumes.

Examples:
  acorn docker volumes`,
	Aliases: []string{"vol"},
	RunE:    runDockerVolumes,
}

// dockerNetworksCmd lists networks
var dockerNetworksCmd = &cobra.Command{
	Use:   "networks",
	Short: "List networks",
	Long: `List Docker networks.

Examples:
  acorn docker networks`,
	Aliases: []string{"net"},
	RunE:    runDockerNetworks,
}

// dockerLogsCmd shows container logs
var dockerLogsCmd = &cobra.Command{
	Use:   "logs [container]",
	Short: "Show container logs",
	Long: `Display logs from a container.

Examples:
  acorn docker logs mycontainer
  acorn docker logs mycontainer -f`,
	Args: cobra.ExactArgs(1),
	RunE: runDockerLogs,
}

// dockerExecCmd executes command in container
var dockerExecCmd = &cobra.Command{
	Use:   "exec [container] [command...]",
	Short: "Execute command in container",
	Long: `Execute a command inside a running container.

Examples:
  acorn docker exec mycontainer /bin/sh
  acorn docker exec mycontainer ls -la`,
	Args: cobra.MinimumNArgs(2),
	RunE: runDockerExec,
}

// dockerStopCmd stops a container
var dockerStopCmd = &cobra.Command{
	Use:   "stop [container]",
	Short: "Stop a container",
	Long: `Stop a running container.

Examples:
  acorn docker stop mycontainer`,
	Args: cobra.ExactArgs(1),
	RunE: runDockerStop,
}

// dockerRmCmd removes a container
var dockerRmCmd = &cobra.Command{
	Use:   "rm [container]",
	Short: "Remove a container",
	Long: `Remove a container.

Use -f to force remove a running container.

Examples:
  acorn docker rm mycontainer
  acorn docker rm -f mycontainer`,
	Args: cobra.ExactArgs(1),
	RunE: runDockerRm,
}

// dockerCleanCmd cleans unused resources
var dockerCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean unused resources",
	Long: `Remove unused containers, images, volumes, and networks.

Use -a to remove all unused images, not just dangling ones.

Examples:
  acorn docker clean
  acorn docker clean -a
  acorn docker clean --dry-run`,
	Aliases: []string{"prune"},
	RunE:    runDockerClean,
}

// dockerComposeCmd is the compose subcommand group
var dockerComposeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Docker Compose commands",
	Long: `Docker Compose helper commands.

Examples:
  acorn docker compose up
  acorn docker compose down
  acorn docker compose logs`,
	Aliases: []string{"dc"},
}

// dockerComposeUpCmd runs compose up
var dockerComposeUpCmd = &cobra.Command{
	Use:   "up [services...]",
	Short: "Start compose services",
	Long: `Start Docker Compose services.

Examples:
  acorn docker compose up
  acorn docker compose up -d
  acorn docker compose up web api`,
	RunE: runDockerComposeUp,
}

// dockerComposeDownCmd runs compose down
var dockerComposeDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop compose services",
	Long: `Stop and remove Docker Compose services.

Examples:
  acorn docker compose down
  acorn docker compose down -v`,
	RunE: runDockerComposeDown,
}

// dockerComposeLogsCmd shows compose logs
var dockerComposeLogsCmd = &cobra.Command{
	Use:   "logs [services...]",
	Short: "Show compose logs",
	Long: `Display logs from Docker Compose services.

Examples:
  acorn docker compose logs
  acorn docker compose logs -f
  acorn docker compose logs web`,
	RunE: runDockerComposeLogs,
}

var (
	dockerFollow        bool
	dockerTail          int
	dockerForce         bool
	dockerDetach        bool
	dockerBuild         bool
	dockerRemoveVolumes bool
	dockerRemoveOrphans bool
	dockerComposeFile   string
)

func init() {
	devopsCmd.AddCommand(dockerCmd)

	// Add subcommands
	dockerCmd.AddCommand(dockerStatusCmd)
	dockerCmd.AddCommand(dockerPsCmd)
	dockerCmd.AddCommand(dockerImagesCmd)
	dockerCmd.AddCommand(dockerVolumesCmd)
	dockerCmd.AddCommand(dockerNetworksCmd)
	dockerCmd.AddCommand(dockerLogsCmd)
	dockerCmd.AddCommand(dockerExecCmd)
	dockerCmd.AddCommand(dockerStopCmd)
	dockerCmd.AddCommand(dockerRmCmd)
	dockerCmd.AddCommand(dockerCleanCmd)
	dockerCmd.AddCommand(dockerComposeCmd)

	// Compose subcommands
	dockerComposeCmd.AddCommand(dockerComposeUpCmd)
	dockerComposeCmd.AddCommand(dockerComposeDownCmd)
	dockerComposeCmd.AddCommand(dockerComposeLogsCmd)

	// Persistent flags
	dockerCmd.PersistentFlags().BoolVarP(&dockerVerbose, "verbose", "v", false,
		"Show verbose output")
	dockerCmd.PersistentFlags().BoolVar(&dockerDryRun, "dry-run", false,
		"Show what would be done without executing")

	// Command-specific flags
	dockerPsCmd.Flags().BoolVarP(&dockerAll, "all", "a", false, "Show all containers")
	dockerLogsCmd.Flags().BoolVarP(&dockerFollow, "follow", "f", false, "Follow log output")
	dockerLogsCmd.Flags().IntVar(&dockerTail, "tail", 0, "Number of lines to show")
	dockerRmCmd.Flags().BoolVarP(&dockerForce, "force", "f", false, "Force remove")
	dockerCleanCmd.Flags().BoolVarP(&dockerAll, "all", "a", false, "Remove all unused images")

	// Compose flags
	dockerComposeCmd.PersistentFlags().StringVarP(&dockerComposeFile, "file", "f", "",
		"Compose file path")
	dockerComposeUpCmd.Flags().BoolVarP(&dockerDetach, "detach", "d", false, "Run in background")
	dockerComposeUpCmd.Flags().BoolVar(&dockerBuild, "build", false, "Build images before starting")
	dockerComposeDownCmd.Flags().BoolVarP(&dockerRemoveVolumes, "volumes", "v", false, "Remove volumes")
	dockerComposeDownCmd.Flags().BoolVar(&dockerRemoveOrphans, "remove-orphans", false, "Remove orphan containers")
	dockerComposeLogsCmd.Flags().BoolVarP(&dockerFollow, "follow", "f", false, "Follow log output")
}

func runDockerStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := docker.NewHelper(dockerVerbose, dockerDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	status := helper.GetStatus()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "Installed:      %v\n", status.Installed)
	fmt.Fprintf(os.Stdout, "Running:        %v\n", status.Running)
	if status.Version != "" {
		fmt.Fprintf(os.Stdout, "Version:        %s\n", status.Version)
	}
	fmt.Fprintf(os.Stdout, "Containers:     %d (%d running)\n", status.Containers, status.RunningCount)
	fmt.Fprintf(os.Stdout, "Images:         %d\n", status.Images)
	fmt.Fprintf(os.Stdout, "Compose:        %v\n", status.ComposeExists)

	return nil
}

func runDockerPs(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := docker.NewHelper(dockerVerbose, dockerDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	containers, err := helper.GetContainers(dockerAll)
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"containers": containers})
	}

	if len(containers) == 0 {
		fmt.Fprintln(os.Stdout, "No containers found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-12s %-20s %-25s %-20s\n", "ID", "NAME", "IMAGE", "STATUS")
	for _, c := range containers {
		fmt.Fprintf(os.Stdout, "%-12s %-20s %-25s %-20s\n", c.ID[:12], truncateStr(c.Names, 20), truncateStr(c.Image, 25), truncateStr(c.Status, 20))
	}

	return nil
}

func runDockerImages(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := docker.NewHelper(dockerVerbose, dockerDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	images, err := helper.GetImages()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"images": images})
	}

	if len(images) == 0 {
		fmt.Fprintln(os.Stdout, "No images found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-40s %-15s %-12s %-10s\n", "REPOSITORY", "TAG", "IMAGE ID", "SIZE")
	for _, img := range images {
		fmt.Fprintf(os.Stdout, "%-40s %-15s %-12s %-10s\n", truncateStr(img.Repository, 40), truncateStr(img.Tag, 15), img.ImageID[:12], img.Size)
	}

	return nil
}

func runDockerVolumes(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := docker.NewHelper(dockerVerbose, dockerDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	volumes, err := helper.GetVolumes()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"volumes": volumes})
	}

	if len(volumes) == 0 {
		fmt.Fprintln(os.Stdout, "No volumes found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-50s %-10s\n", "NAME", "DRIVER")
	for _, v := range volumes {
		fmt.Fprintf(os.Stdout, "%-50s %-10s\n", truncateStr(v.Name, 50), v.Driver)
	}

	return nil
}

func runDockerNetworks(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := docker.NewHelper(dockerVerbose, dockerDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	networks, err := helper.GetNetworks()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"networks": networks})
	}

	if len(networks) == 0 {
		fmt.Fprintln(os.Stdout, "No networks found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-20s %-12s %-10s %-10s\n", "NAME", "ID", "DRIVER", "SCOPE")
	for _, n := range networks {
		fmt.Fprintf(os.Stdout, "%-20s %-12s %-10s %-10s\n", truncateStr(n.Name, 20), n.ID[:12], n.Driver, n.Scope)
	}

	return nil
}

func runDockerLogs(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(dockerVerbose, dockerDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	return helper.GetLogs(args[0], dockerFollow, dockerTail)
}

func runDockerExec(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(dockerVerbose, dockerDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	return helper.ExecInContainer(args[0], args[1:], true)
}

func runDockerStop(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(dockerVerbose, dockerDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	if err := helper.StopContainer(args[0]); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Stopped container %s\n", output.Success("✓"), args[0])
	return nil
}

func runDockerRm(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(dockerVerbose, dockerDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	if err := helper.RemoveContainer(args[0], dockerForce); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Removed container %s\n", output.Success("✓"), args[0])
	return nil
}

func runDockerClean(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(dockerVerbose, dockerDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	result, err := helper.Cleanup(dockerAll)
	if err != nil {
		return err
	}

	if dockerDryRun {
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s Cleanup complete:\n", output.Success("✓"))
	fmt.Fprintf(os.Stdout, "  Containers: %d removed\n", result.ContainersRemoved)
	fmt.Fprintf(os.Stdout, "  Images:     %d removed\n", result.ImagesRemoved)
	fmt.Fprintf(os.Stdout, "  Volumes:    %d removed\n", result.VolumesRemoved)
	fmt.Fprintf(os.Stdout, "  Networks:   %d removed\n", result.NetworksRemoved)

	return nil
}

func runDockerComposeUp(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(dockerVerbose, dockerDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	return helper.ComposeUp(dockerComposeFile, dockerDetach, dockerBuild, args)
}

func runDockerComposeDown(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(dockerVerbose, dockerDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	return helper.ComposeDown(dockerComposeFile, dockerRemoveVolumes, dockerRemoveOrphans)
}

func runDockerComposeLogs(cmd *cobra.Command, args []string) error {
	helper := docker.NewHelper(dockerVerbose, dockerDryRun)

	if !helper.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	return helper.GetComposeLogs(dockerComposeFile, dockerFollow, args)
}

// truncateStr truncates a string to max length
func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
