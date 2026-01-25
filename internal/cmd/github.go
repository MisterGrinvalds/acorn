package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/vcs/github"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	ghVerbose bool
	ghDryRun  bool
)

// ghCmd represents the github command group
var ghCmd = &cobra.Command{
	Use:   "gh",
	Short: "GitHub CLI workflow helpers",
	Long: `GitHub CLI workflow helpers for common operations.

Provides commands for branch management, PR creation,
and workflow monitoring.

Examples:
  acorn gh status              # Show gh CLI and repo status
  acorn gh pr create           # Push branch and create PR
  acorn gh pr status           # Show PR status
  acorn gh cleanup             # Clean merged branches`,
	Aliases: []string{"github"},
}

// ghStatusCmd shows status
var ghStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show GitHub CLI and repo status",
	Long: `Display GitHub CLI installation, authentication, and repo status.

Examples:
  acorn gh status
  acorn gh status -o json`,
	RunE: runGhStatus,
}

// ghCleanupCmd cleans merged branches
var ghCleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean merged branches",
	Long: `Remove local branches that have been merged and prune remote tracking.

Skips main, master, and develop branches.

Examples:
  acorn gh cleanup
  acorn gh cleanup --dry-run`,
	RunE: runGhCleanup,
}

// ghPRCmd is the parent for PR subcommands
var ghPRCmd = &cobra.Command{
	Use:   "pr",
	Short: "Pull request commands",
	Long: `Commands for managing pull requests.

Examples:
  acorn gh pr create     # Push and create PR
  acorn gh pr status     # Show PR status
  acorn gh pr checks     # Show PR checks`,
}

// ghPRCreateCmd creates a PR
var ghPRCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Push branch and create PR",
	Long: `Push current branch to origin and create a pull request.

Opens the PR creation page in your browser.

Examples:
  acorn gh pr create`,
	Aliases: []string{"new"},
	RunE:    runGhPRCreate,
}

// ghPRStatusCmd shows PR status
var ghPRStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show PR status",
	Long: `Display pull request status for current branch.

Examples:
  acorn gh pr status`,
	RunE: runGhPRStatus,
}

// ghPRChecksCmd shows PR checks
var ghPRChecksCmd = &cobra.Command{
	Use:   "checks",
	Short: "Show PR checks",
	Long: `Display CI/CD check status for current PR.

Examples:
  acorn gh pr checks`,
	RunE: runGhPRChecks,
}

// ghRunCmd is the parent for workflow run subcommands
var ghRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Workflow run commands",
	Long: `Commands for managing GitHub Actions workflow runs.

Examples:
  acorn gh run watch     # Watch current run
  acorn gh run rerun     # Rerun failed jobs`,
}

// ghRunWatchCmd watches a run
var ghRunWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch current workflow run",
	Long: `Watch the current workflow run in real-time.

Examples:
  acorn gh run watch`,
	RunE: runGhRunWatch,
}

// ghRunRerunCmd reruns failed jobs
var ghRunRerunCmd = &cobra.Command{
	Use:   "rerun",
	Short: "Rerun failed workflow jobs",
	Long: `Rerun failed jobs from the most recent workflow run.

Examples:
  acorn gh run rerun`,
	RunE: runGhRunRerun,
}

// ghCommitCmd does quick commit
var ghCommitCmd = &cobra.Command{
	Use:   "commit <message>",
	Short: "Quick commit all changes",
	Long: `Add all changes and commit with the given message.

Examples:
  acorn gh commit "Fix bug in login"`,
	Args: cobra.ExactArgs(1),
	RunE: runGhCommit,
}

// ghBranchCmd creates new branch
var ghBranchCmd = &cobra.Command{
	Use:   "branch <name>",
	Short: "Create and checkout new branch",
	Long: `Create a new branch and switch to it.

Examples:
  acorn gh branch feature/new-feature`,
	Args: cobra.ExactArgs(1),
	RunE: runGhBranch,
}

// ghPushCmd pushes current branch
var ghPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push current branch to origin",
	Long: `Push current branch to origin with upstream tracking.

Examples:
  acorn gh push`,
	RunE: runGhPush,
}

func init() {
	vcsCmd.AddCommand(ghCmd)

	// Add subcommands
	ghCmd.AddCommand(ghStatusCmd)
	ghCmd.AddCommand(ghCleanupCmd)
	ghCmd.AddCommand(ghPRCmd)
	ghCmd.AddCommand(ghRunCmd)
	ghCmd.AddCommand(ghCommitCmd)
	ghCmd.AddCommand(ghBranchCmd)
	ghCmd.AddCommand(ghPushCmd)

	// PR subcommands
	ghPRCmd.AddCommand(ghPRCreateCmd)
	ghPRCmd.AddCommand(ghPRStatusCmd)
	ghPRCmd.AddCommand(ghPRChecksCmd)

	// Run subcommands
	ghRunCmd.AddCommand(ghRunWatchCmd)
	ghRunCmd.AddCommand(ghRunRerunCmd)

	// Persistent flags
	ghCmd.PersistentFlags().BoolVarP(&ghVerbose, "verbose", "v", false,
		"Show verbose output")
	ghCmd.PersistentFlags().BoolVar(&ghDryRun, "dry-run", false,
		"Show what would be done without executing")
}

func runGhStatus(cmd *cobra.Command, args []string) error {
	helper := github.NewHelper(ghVerbose, ghDryRun)
	status := helper.GetStatus()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("GitHub CLI Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if !status.GhInstalled {
		fmt.Fprintf(os.Stdout, "%s GitHub CLI (gh) not installed\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: brew install gh")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s %s\n", output.Success("✓"), status.GhVersion)

	if status.Authenticated {
		fmt.Fprintf(os.Stdout, "%s Authenticated as: %s\n", output.Success("✓"), status.Username)
	} else {
		fmt.Fprintf(os.Stdout, "%s Not authenticated\n", output.Warning("⚠"))
		fmt.Fprintln(os.Stdout, "  Run: gh auth login")
	}

	fmt.Fprintln(os.Stdout)
	if status.CurrentBranch != "" {
		fmt.Fprintf(os.Stdout, "Branch: %s\n", status.CurrentBranch)
	}
	if status.HasRemote {
		fmt.Fprintf(os.Stdout, "%s Remote configured\n", output.Success("✓"))
	}

	return nil
}

func runGhCleanup(cmd *cobra.Command, args []string) error {
	helper := github.NewHelper(ghVerbose, ghDryRun)

	deleted, err := helper.CleanupBranches()
	if err != nil {
		return err
	}

	if len(deleted) == 0 {
		fmt.Fprintln(os.Stdout, "No merged branches to clean")
		return nil
	}

	if ghDryRun {
		fmt.Fprintf(os.Stdout, "%s\n", output.Info("Branches that would be deleted:"))
	} else {
		fmt.Fprintf(os.Stdout, "%s\n", output.Info("Deleted branches:"))
	}

	for _, b := range deleted {
		fmt.Fprintf(os.Stdout, "  %s %s\n", output.Success("✓"), b)
	}

	if !ghDryRun {
		fmt.Fprintln(os.Stdout, "\nPruned remote tracking branches")
	}

	return nil
}

func runGhPRCreate(cmd *cobra.Command, args []string) error {
	helper := github.NewHelper(ghVerbose, ghDryRun)
	return helper.CreatePR()
}

func runGhPRStatus(cmd *cobra.Command, args []string) error {
	helper := github.NewHelper(ghVerbose, ghDryRun)
	return helper.GetPRStatus()
}

func runGhPRChecks(cmd *cobra.Command, args []string) error {
	helper := github.NewHelper(ghVerbose, ghDryRun)
	return helper.GetPRChecks()
}

func runGhRunWatch(cmd *cobra.Command, args []string) error {
	helper := github.NewHelper(ghVerbose, ghDryRun)
	return helper.WatchRun()
}

func runGhRunRerun(cmd *cobra.Command, args []string) error {
	helper := github.NewHelper(ghVerbose, ghDryRun)
	return helper.RerunFailed()
}

func runGhCommit(cmd *cobra.Command, args []string) error {
	helper := github.NewHelper(ghVerbose, ghDryRun)
	return helper.QuickCommit(args[0])
}

func runGhBranch(cmd *cobra.Command, args []string) error {
	helper := github.NewHelper(ghVerbose, ghDryRun)
	if err := helper.NewBranch(args[0]); err != nil {
		return err
	}
	if !ghDryRun {
		fmt.Fprintf(os.Stdout, "%s Created and switched to branch: %s\n", output.Success("✓"), args[0])
	}
	return nil
}

func runGhPush(cmd *cobra.Command, args []string) error {
	helper := github.NewHelper(ghVerbose, ghDryRun)
	return helper.PushBranch()
}
