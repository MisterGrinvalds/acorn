package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	rootconfig "github.com/mistergrinvalds/acorn/config"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	saplingMessage string
	saplingAll     bool
	saplingVerbose bool
	saplingDryRun  bool
)

// saplingCmd represents the sapling command group
var saplingCmd = &cobra.Command{
	Use:   "sapling",
	Short: "Manage the .sapling dotfiles repository",
	Long: `Manage the .sapling directory as a git repository for dotfiles.

The .sapling directory contains your dotfiles configuration sources,
AI tool configurations, and generated outputs. This command helps you
version control and sync these files.

Commands:
  status  - Show git status of .sapling
  commit  - Commit changes in .sapling
  push    - Push commits to remote
  pull    - Pull changes from remote
  sync    - Full sync (commit + push + pull)

Examples:
  acorn sapling status              # Check git status
  acorn sapling commit -m "msg"     # Commit with message
  acorn sapling sync                # Full sync
  acorn sapling sync -m "update"    # Sync with commit message`,
	Aliases: []string{"sp"},
}

// saplingStatusCmd shows git status
var saplingStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show git status of .sapling repository",
	Long: `Display the git status of the .sapling directory.

Shows:
  - Modified files
  - Untracked files
  - Current branch
  - Remote tracking status

Examples:
  acorn sapling status
  acorn sapling status -v    # Verbose output`,
	RunE: runSaplingStatus,
}

// saplingCommitCmd commits changes
var saplingCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit changes in .sapling repository",
	Long: `Commit changes to the .sapling git repository.

By default, commits staged files. Use -a/--all to stage all changes.

Examples:
  acorn sapling commit -m "Update git config"
  acorn sapling commit -a -m "Update all configs"
  acorn sapling commit --dry-run -m "Test commit"`,
	RunE: runSaplingCommit,
}

// saplingPushCmd pushes to remote
var saplingPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push commits to remote repository",
	Long: `Push committed changes to the remote repository.

Examples:
  acorn sapling push
  acorn sapling push --dry-run    # Show what would be pushed`,
	RunE: runSaplingPush,
}

// saplingPullCmd pulls from remote
var saplingPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull changes from remote repository",
	Long: `Pull changes from the remote repository.

Examples:
  acorn sapling pull
  acorn sapling pull --dry-run    # Show what would be pulled`,
	RunE: runSaplingPull,
}

// saplingSyncCmd syncs with remote
var saplingSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Full sync with remote (commit + push + pull)",
	Long: `Perform a full sync with the remote repository:

1. Pull latest changes from remote
2. Stage and commit local changes (if any)
3. Push commits to remote

This ensures your local and remote repositories are in sync.

Examples:
  acorn sapling sync
  acorn sapling sync -m "Update configs"
  acorn sapling sync -a -m "Sync all changes"
  acorn sapling sync --dry-run`,
	RunE: runSaplingSync,
}

func init() {
	rootCmd.AddCommand(saplingCmd)

	// Add subcommands
	saplingCmd.AddCommand(saplingStatusCmd)
	saplingCmd.AddCommand(saplingCommitCmd)
	saplingCmd.AddCommand(saplingPushCmd)
	saplingCmd.AddCommand(saplingPullCmd)
	saplingCmd.AddCommand(saplingSyncCmd)

	// Commit flags
	saplingCommitCmd.Flags().StringVarP(&saplingMessage, "message", "m", "", "Commit message (required)")
	saplingCommitCmd.Flags().BoolVarP(&saplingAll, "all", "a", false, "Stage all changes before committing")
	saplingCommitCmd.MarkFlagRequired("message")

	// Sync flags
	saplingSyncCmd.Flags().StringVarP(&saplingMessage, "message", "m", "", "Commit message (default: auto-generated)")
	saplingSyncCmd.Flags().BoolVarP(&saplingAll, "all", "a", false, "Stage all changes before committing")

	// Global flags
	saplingCmd.PersistentFlags().BoolVarP(&saplingVerbose, "verbose", "v", false, "Show verbose output")
	saplingCmd.PersistentFlags().BoolVar(&saplingDryRun, "dry-run", false, "Show what would be done without executing")
}

// getSaplingRoot returns the .sapling directory path
func getSaplingRoot() (string, error) {
	root, err := rootconfig.SaplingRoot()
	if err != nil {
		return "", fmt.Errorf("failed to find .sapling directory: %w", err)
	}
	return root, nil
}

// runGitCommand executes a git command in the .sapling directory
func runGitCommand(args []string, dryRun bool) error {
	saplingRoot, err := getSaplingRoot()
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Fprintf(os.Stdout, "%s Would run: git %s\n", output.Info("○"), strings.Join(args, " "))
		fmt.Fprintf(os.Stdout, "  In directory: %s\n", saplingRoot)
		return nil
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = saplingRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if saplingVerbose {
		fmt.Fprintf(os.Stdout, "Running: git %s\n", strings.Join(args, " "))
		fmt.Fprintf(os.Stdout, "In: %s\n\n", saplingRoot)
	}

	return cmd.Run()
}

// getGitOutput executes a git command and returns its output
func getGitOutput(args []string) (string, error) {
	saplingRoot, err := getSaplingRoot()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = saplingRoot
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// runSaplingStatus shows git status
func runSaplingStatus(cmd *cobra.Command, args []string) error {
	saplingRoot, err := getSaplingRoot()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Sapling Repository Status\n", output.Info("ℹ"))
	fmt.Fprintf(os.Stdout, "Location: %s\n\n", saplingRoot)

	// Get current branch
	branch, err := getGitOutput([]string{"branch", "--show-current"})
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	branch = strings.TrimSpace(branch)
	fmt.Fprintf(os.Stdout, "Branch: %s\n\n", output.Success(branch))

	// Get status
	if err := runGitCommand([]string{"status"}, false); err != nil {
		return fmt.Errorf("git status failed: %w", err)
	}

	return nil
}

// runSaplingCommit commits changes
func runSaplingCommit(cmd *cobra.Command, args []string) error {
	saplingRoot, err := getSaplingRoot()
	if err != nil {
		return err
	}

	if saplingMessage == "" {
		return fmt.Errorf("commit message is required (use -m flag)")
	}

	fmt.Fprintf(os.Stdout, "%s Committing changes in .sapling\n", output.Info("ℹ"))
	fmt.Fprintf(os.Stdout, "Location: %s\n", saplingRoot)
	fmt.Fprintf(os.Stdout, "Message: %s\n\n", saplingMessage)

	// Stage changes if -a flag is set
	if saplingAll {
		if saplingVerbose || saplingDryRun {
			fmt.Fprintf(os.Stdout, "Staging all changes...\n")
		}
		if err := runGitCommand([]string{"add", "-A"}, saplingDryRun); err != nil {
			return fmt.Errorf("git add failed: %w", err)
		}
	}

	// Commit
	if saplingVerbose || saplingDryRun {
		fmt.Fprintf(os.Stdout, "Creating commit...\n")
	}
	if err := runGitCommand([]string{"commit", "-m", saplingMessage}, saplingDryRun); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	if !saplingDryRun {
		fmt.Fprintf(os.Stdout, "\n%s Committed successfully\n", output.Success("✓"))
	}

	return nil
}

// runSaplingPush pushes to remote
func runSaplingPush(cmd *cobra.Command, args []string) error {
	saplingRoot, err := getSaplingRoot()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Pushing to remote\n", output.Info("ℹ"))
	fmt.Fprintf(os.Stdout, "Location: %s\n\n", saplingRoot)

	if err := runGitCommand([]string{"push"}, saplingDryRun); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	if !saplingDryRun {
		fmt.Fprintf(os.Stdout, "\n%s Pushed successfully\n", output.Success("✓"))
	}

	return nil
}

// runSaplingPull pulls from remote
func runSaplingPull(cmd *cobra.Command, args []string) error {
	saplingRoot, err := getSaplingRoot()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Pulling from remote\n", output.Info("ℹ"))
	fmt.Fprintf(os.Stdout, "Location: %s\n\n", saplingRoot)

	if err := runGitCommand([]string{"pull"}, saplingDryRun); err != nil {
		return fmt.Errorf("git pull failed: %w", err)
	}

	if !saplingDryRun {
		fmt.Fprintf(os.Stdout, "\n%s Pulled successfully\n", output.Success("✓"))
	}

	return nil
}

// runSaplingSync performs full sync
func runSaplingSync(cmd *cobra.Command, args []string) error {
	saplingRoot, err := getSaplingRoot()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Syncing .sapling repository\n", output.Info("ℹ"))
	fmt.Fprintf(os.Stdout, "Location: %s\n\n", saplingRoot)

	// Step 1: Pull from remote
	fmt.Fprintf(os.Stdout, "Step 1: Pulling latest changes from remote\n")
	if err := runGitCommand([]string{"pull"}, saplingDryRun); err != nil {
		fmt.Fprintf(os.Stdout, "%s Pull failed (might be first sync): %v\n\n", output.Warning("!"), err)
	} else {
		fmt.Fprintf(os.Stdout, "%s Pulled successfully\n\n", output.Success("✓"))
	}

	// Step 2: Check for local changes
	if !saplingDryRun {
		statusOut, err := getGitOutput([]string{"status", "--porcelain"})
		if err != nil {
			return fmt.Errorf("failed to check status: %w", err)
		}

		if strings.TrimSpace(statusOut) == "" {
			fmt.Fprintf(os.Stdout, "Step 2: No local changes to commit\n")
			fmt.Fprintf(os.Stdout, "%s Repository is up to date\n", output.Success("✓"))
			return nil
		}
	}

	// Step 3: Stage and commit changes
	fmt.Fprintf(os.Stdout, "Step 2: Staging and committing local changes\n")

	if saplingAll {
		if saplingVerbose || saplingDryRun {
			fmt.Fprintf(os.Stdout, "  Staging all changes...\n")
		}
		if err := runGitCommand([]string{"add", "-A"}, saplingDryRun); err != nil {
			return fmt.Errorf("git add failed: %w", err)
		}
	}

	// Generate commit message if not provided
	commitMsg := saplingMessage
	if commitMsg == "" {
		// Get current timestamp
		cmd := exec.Command("date", "+%Y-%m-%d %H:%M:%S")
		dateBytes, err := cmd.Output()
		if err == nil {
			commitMsg = fmt.Sprintf("Auto-sync dotfiles - %s", strings.TrimSpace(string(dateBytes)))
		} else {
			commitMsg = "Auto-sync dotfiles"
		}
	}

	if saplingVerbose || saplingDryRun {
		fmt.Fprintf(os.Stdout, "  Commit message: %s\n", commitMsg)
	}

	if err := runGitCommand([]string{"commit", "-m", commitMsg}, saplingDryRun); err != nil {
		fmt.Fprintf(os.Stdout, "%s No changes to commit\n\n", output.Info("○"))
	} else {
		fmt.Fprintf(os.Stdout, "%s Committed successfully\n\n", output.Success("✓"))
	}

	// Step 4: Push to remote
	fmt.Fprintf(os.Stdout, "Step 3: Pushing to remote\n")
	if err := runGitCommand([]string{"push"}, saplingDryRun); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	if !saplingDryRun {
		fmt.Fprintf(os.Stdout, "%s Pushed successfully\n\n", output.Success("✓"))
		fmt.Fprintf(os.Stdout, "%s Sync complete!\n", output.Success("✓"))
	}

	return nil
}
