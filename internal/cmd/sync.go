package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/utils/configfile"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/components/terminal/shell"
	"github.com/spf13/cobra"
)

var (
	syncQuiet bool
)

// syncCmd represents the sync command group
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize dotfiles repository",
	Long: `Synchronize dotfiles repository with remote and manage config symlinks.

The sync commands help you keep your dotfiles in sync with the remote repository
and manage symlinks from XDG directories to generated config files in the repo.

Commands:
  status  - Show dotfiles repository status
  pull    - Pull latest changes from remote
  push    - Commit and push local changes
  drift   - Quick drift check (commits ahead/behind)
  audit   - Full audit of all changes
  link    - Create symlinks from XDG to generated configs
  unlink  - Remove symlinks (safely)
  update  - Pull and reload shell configuration

Examples:
  acorn sync status
  acorn sync pull
  acorn sync push "Update tmux config"
  acorn sync link`,
	Aliases: []string{"s"},
}

// syncStatusCmd shows repository status
var syncStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show dotfiles repository status",
	Long: `Display the current status of the dotfiles repository.

Shows:
  - Current branch
  - Commits ahead/behind remote
  - Modified and untracked files
  - Symlink status for config files`,
	Aliases: []string{"st"},
	RunE:    runSyncStatus,
}

// syncPullCmd pulls latest changes
var syncPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull latest changes from remote",
	Long: `Pull the latest changes from the remote repository.

Uses git pull --rebase to maintain a clean history.`,
	RunE: runSyncPull,
}

// syncPushCmd commits and pushes changes
var syncPushCmd = &cobra.Command{
	Use:   "push [message]",
	Short: "Commit and push local changes",
	Long: `Commit all local changes and push to the remote repository.

If no message is provided, a default message with timestamp is used.

Examples:
  acorn sync push
  acorn sync push "Update tmux config"`,
	RunE: runSyncPush,
}

// syncDriftCmd checks for drift
var syncDriftCmd = &cobra.Command{
	Use:   "drift",
	Short: "Check for drift from remote",
	Long: `Quickly check if local and remote are in sync.

Shows commits ahead/behind without detailed file changes.
Use --quiet for use in shell startup scripts.`,
	RunE: runSyncDrift,
}

// syncAuditCmd performs full audit
var syncAuditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Full audit of repository state",
	Long: `Perform a comprehensive audit of the dotfiles repository.

Shows detailed information about:
  - Repository status
  - All modified files
  - Untracked files
  - Symlink health`,
	RunE: runSyncAudit,
}

// syncLinkCmd creates symlinks
var syncLinkCmd = &cobra.Command{
	Use:   "link",
	Short: "Create symlinks for config files",
	Long: `Create symlinks from XDG directories to generated config files.

For each component with generated config files, creates a symlink from
the XDG path (e.g., ~/.config/tmux/tmux.conf) to the generated file
in the repository (e.g., $DOTFILES_ROOT/generated/tmux/tmux.conf).

This allows you to:
  - Keep config files version controlled
  - See changes in git diff
  - Easily rollback configurations`,
	RunE: runSyncLink,
}

// syncUnlinkCmd removes symlinks
var syncUnlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Remove config symlinks",
	Long: `Remove symlinks that point to the generated config files.

Only removes symlinks that point to files in $DOTFILES_ROOT/generated/.
Regular files are left untouched to prevent data loss.`,
	RunE: runSyncUnlink,
}

// syncUpdateCmd pulls and reloads
var syncUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull changes and reload configuration",
	Long: `Pull latest changes from remote and reload shell configuration.

Equivalent to running:
  acorn sync pull
  acorn shell generate
  source $ACORN_CONFIG_DIR/shell.sh (in your current shell)`,
	RunE: runSyncUpdate,
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Add subcommands
	syncCmd.AddCommand(syncStatusCmd)
	syncCmd.AddCommand(syncPullCmd)
	syncCmd.AddCommand(syncPushCmd)
	syncCmd.AddCommand(syncDriftCmd)
	syncCmd.AddCommand(syncAuditCmd)
	syncCmd.AddCommand(syncLinkCmd)
	syncCmd.AddCommand(syncUnlinkCmd)
	syncCmd.AddCommand(syncUpdateCmd)

	// Flags
	syncDriftCmd.Flags().BoolVarP(&syncQuiet, "quiet", "q", false, "Minimal output (for shell startup)")
}

// getSyncRoot returns the dotfiles repository root for sync operations
func getSyncRoot() string {
	root, err := getDotfilesRoot()
	if err != nil {
		// Fallback
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Repos", "personal", "tools")
	}
	return root
}

// getGeneratedDir returns the generated directory path
func getGeneratedDir() string {
	return filepath.Join(getSyncRoot(), "generated")
}

// isSyncGitRepo checks if the directory is a git repository
func isSyncGitRepo(dir string) bool {
	gitDir := filepath.Join(dir, ".git")
	info, err := os.Stat(gitDir)
	return err == nil && info.IsDir()
}

// syncGitCmd runs a git command in the dotfiles directory
func syncGitCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	cmd.Dir = getSyncRoot()
	return cmd
}

// runSyncStatus shows repository status
func runSyncStatus(cmd *cobra.Command, args []string) error {
	root := getSyncRoot()

	if !isSyncGitRepo(root) {
		return fmt.Errorf("not a git repository: %s", root)
	}

	fmt.Fprintf(os.Stdout, "%s Dotfiles Status\n", output.Info("ℹ"))
	fmt.Fprintf(os.Stdout, "  Repository: %s\n", root)

	// Current branch
	out, err := syncGitCmd("branch", "--show-current").Output()
	if err == nil {
		fmt.Fprintf(os.Stdout, "  Branch:     %s\n", strings.TrimSpace(string(out)))
	}

	// Commits ahead/behind
	ahead, behind := getCommitCounts()
	if ahead > 0 || behind > 0 {
		fmt.Fprintf(os.Stdout, "  Remote:     %d ahead, %d behind\n", ahead, behind)
	} else {
		fmt.Fprintf(os.Stdout, "  Remote:     up to date\n")
	}

	fmt.Fprintln(os.Stdout)

	// Show git status
	statusCmd := syncGitCmd("status", "--short")
	statusCmd.Stdout = os.Stdout
	statusCmd.Stderr = os.Stderr
	if err := statusCmd.Run(); err != nil {
		return fmt.Errorf("git status failed: %w", err)
	}

	return nil
}

// getCommitCounts returns commits ahead and behind remote
func getCommitCounts() (ahead, behind int) {
	// Fetch first (silently)
	syncGitCmd("fetch", "-q").Run()

	// Get counts
	out, err := syncGitCmd("rev-list", "--left-right", "--count", "@{u}...HEAD").Output()
	if err != nil {
		return 0, 0
	}

	parts := strings.Fields(string(out))
	if len(parts) == 2 {
		fmt.Sscanf(parts[0], "%d", &behind)
		fmt.Sscanf(parts[1], "%d", &ahead)
	}
	return
}

// runSyncPull pulls latest changes
func runSyncPull(cmd *cobra.Command, args []string) error {
	root := getSyncRoot()

	if !isSyncGitRepo(root) {
		return fmt.Errorf("not a git repository: %s", root)
	}

	fmt.Fprintf(os.Stdout, "%s Pulling latest changes...\n", output.Info("→"))

	pullCmd := syncGitCmd("pull", "--rebase")
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr
	if err := pullCmd.Run(); err != nil {
		return fmt.Errorf("git pull failed: %w", err)
	}

	fmt.Fprintf(os.Stdout, "%s Pull complete\n", output.Success("✓"))
	return nil
}

// runSyncPush commits and pushes changes
func runSyncPush(cmd *cobra.Command, args []string) error {
	root := getSyncRoot()

	if !isSyncGitRepo(root) {
		return fmt.Errorf("not a git repository: %s", root)
	}

	// Check for changes
	statusOut, _ := syncGitCmd("status", "--porcelain").Output()
	if len(statusOut) == 0 {
		fmt.Fprintf(os.Stdout, "%s No changes to commit\n", output.Info("ℹ"))
		return nil
	}

	// Determine commit message
	message := "Update dotfiles"
	if len(args) > 0 {
		message = args[0]
	}

	fmt.Fprintf(os.Stdout, "%s Committing changes...\n", output.Info("→"))

	// Add all changes
	addCmd := syncGitCmd("add", "-A")
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	// Commit
	commitCmd := syncGitCmd("commit", "-m", message)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	fmt.Fprintf(os.Stdout, "%s Pushing to remote...\n", output.Info("→"))

	// Push
	pushCmd := syncGitCmd("push")
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	fmt.Fprintf(os.Stdout, "%s Push complete\n", output.Success("✓"))
	return nil
}

// runSyncDrift checks for drift
func runSyncDrift(cmd *cobra.Command, args []string) error {
	root := getSyncRoot()

	if !isSyncGitRepo(root) {
		if !syncQuiet {
			fmt.Fprintf(os.Stderr, "not a git repository: %s\n", root)
		}
		return nil // Don't error for quiet mode
	}

	ahead, behind := getCommitCounts()

	if syncQuiet {
		// Minimal output for shell startup
		if ahead > 0 || behind > 0 {
			fmt.Fprintf(os.Stdout, "[dotfiles] %d ahead, %d behind\n", ahead, behind)
		}
	} else {
		if ahead == 0 && behind == 0 {
			fmt.Fprintf(os.Stdout, "%s Dotfiles in sync with remote\n", output.Success("✓"))
		} else {
			if ahead > 0 {
				fmt.Fprintf(os.Stdout, "%s %d commit(s) ahead of remote\n", output.Warning("!"), ahead)
			}
			if behind > 0 {
				fmt.Fprintf(os.Stdout, "%s %d commit(s) behind remote\n", output.Warning("!"), behind)
			}
		}
	}

	return nil
}

// runSyncAudit performs full audit
func runSyncAudit(cmd *cobra.Command, args []string) error {
	root := getSyncRoot()

	if !isSyncGitRepo(root) {
		return fmt.Errorf("not a git repository: %s", root)
	}

	fmt.Fprintf(os.Stdout, "%s Dotfiles Audit\n", output.Info("ℹ"))
	fmt.Fprintf(os.Stdout, "═══════════════════════════════════════════════════\n")

	// Repository info
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Repository:")
	fmt.Fprintf(os.Stdout, "  Path:   %s\n", root)

	out, err := syncGitCmd("branch", "--show-current").Output()
	if err == nil {
		fmt.Fprintf(os.Stdout, "  Branch: %s\n", strings.TrimSpace(string(out)))
	}

	// Remote status
	ahead, behind := getCommitCounts()
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Remote Status:")
	if ahead == 0 && behind == 0 {
		fmt.Fprintf(os.Stdout, "  %s In sync with remote\n", output.Success("✓"))
	} else {
		if ahead > 0 {
			fmt.Fprintf(os.Stdout, "  %s %d commit(s) to push\n", output.Warning("→"), ahead)
		}
		if behind > 0 {
			fmt.Fprintf(os.Stdout, "  %s %d commit(s) to pull\n", output.Warning("←"), behind)
		}
	}

	// Changed files
	statusOut, _ := syncGitCmd("status", "--porcelain").Output()
	if len(statusOut) > 0 {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Changes:")
		lines := strings.Split(string(statusOut), "\n")
		for _, line := range lines {
			if line != "" {
				fmt.Fprintf(os.Stdout, "  %s\n", line)
			}
		}
	} else {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "  %s Working tree clean\n", output.Success("✓"))
	}

	// Symlink status
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Config Symlinks:")
	if err := checkSymlinks(); err != nil {
		fmt.Fprintf(os.Stdout, "  %s Error checking symlinks: %v\n", output.Error("✗"), err)
	}

	return nil
}

// checkSymlinks verifies symlink status
func checkSymlinks() error {
	generatedDir := getGeneratedDir()

	// Check if generated directory exists
	if _, err := os.Stat(generatedDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stdout, "  %s Generated directory not found: %s\n", output.Warning("!"), generatedDir)
		fmt.Fprintf(os.Stdout, "    Run 'acorn shell generate' to create config files\n")
		return nil
	}

	// Walk through generated directory
	found := false
	err := filepath.Walk(generatedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if info.IsDir() {
			return nil
		}

		found = true

		// Get relative path from generated dir
		relPath, _ := filepath.Rel(generatedDir, path)
		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) < 2 {
			return nil
		}

		component := parts[0]
		filename := parts[len(parts)-1]

		// Determine XDG target
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig == "" {
			home, _ := os.UserHomeDir()
			xdgConfig = filepath.Join(home, ".config")
		}

		// Special case: shell scripts go to acorn/ directory, not shell/
		targetComponent := component
		if component == "shell" {
			targetComponent = "acorn"
		}
		target := filepath.Join(xdgConfig, targetComponent, filename)

		// Check symlink
		linkInfo, err := os.Lstat(target)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stdout, "  %s %s → not linked\n", output.Warning("○"), target)
		} else if linkInfo.Mode()&os.ModeSymlink != 0 {
			linkDest, _ := os.Readlink(target)
			if linkDest == path {
				fmt.Fprintf(os.Stdout, "  %s %s → %s\n", output.Success("✓"), target, path)
			} else {
				fmt.Fprintf(os.Stdout, "  %s %s → %s (wrong target)\n", output.Warning("!"), target, linkDest)
			}
		} else {
			fmt.Fprintf(os.Stdout, "  %s %s (regular file, not symlink)\n", output.Warning("!"), target)
		}

		return nil
	})

	if !found {
		fmt.Fprintf(os.Stdout, "  %s No generated config files found\n", output.Info("ℹ"))
	}

	return err
}

// runSyncLink creates symlinks
func runSyncLink(cmd *cobra.Command, args []string) error {
	generatedDir := getGeneratedDir()

	// Check if generated directory exists
	if _, err := os.Stat(generatedDir); os.IsNotExist(err) {
		return fmt.Errorf("generated directory not found: %s\nRun 'acorn shell generate' first", generatedDir)
	}

	fmt.Fprintf(os.Stdout, "%s Creating symlinks...\n", output.Info("→"))

	// Walk through generated directory
	count := 0
	err := filepath.Walk(generatedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		// Get relative path from generated dir
		relPath, _ := filepath.Rel(generatedDir, path)
		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) < 2 {
			return nil
		}

		component := parts[0]
		filename := parts[len(parts)-1]

		// Determine XDG target
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig == "" {
			home, _ := os.UserHomeDir()
			xdgConfig = filepath.Join(home, ".config")
		}

		// Special case: shell scripts go to acorn/ directory, not shell/
		// This is because shell scripts (bootstrap.sh, etc.) are sourced from
		// $ACORN_CONFIG_DIR which is ~/.config/acorn/
		targetComponent := component
		if component == "shell" {
			targetComponent = "acorn"
		}
		target := filepath.Join(xdgConfig, targetComponent, filename)

		// Create parent directory
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "  %s Failed to create directory: %v\n", output.Error("✗"), err)
			return nil
		}

		// Check if target exists
		if linkInfo, err := os.Lstat(target); err == nil {
			if linkInfo.Mode()&os.ModeSymlink != 0 {
				// Remove existing symlink
				os.Remove(target)
			} else {
				// Regular file - backup
				backup := target + ".backup"
				fmt.Fprintf(os.Stdout, "  %s Backing up %s to %s\n", output.Warning("!"), target, backup)
				os.Rename(target, backup)
			}
		}

		// Create symlink
		if err := os.Symlink(path, target); err != nil {
			fmt.Fprintf(os.Stderr, "  %s Failed to create symlink %s: %v\n", output.Error("✗"), target, err)
			return nil
		}

		fmt.Fprintf(os.Stdout, "  %s %s → %s\n", output.Success("✓"), target, path)
		count++

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking generated directory: %w", err)
	}

	if count == 0 {
		fmt.Fprintf(os.Stdout, "%s No config files to link\n", output.Info("ℹ"))
	} else {
		fmt.Fprintf(os.Stdout, "\n%s Created %d symlink(s)\n", output.Success("✓"), count)
	}

	return nil
}

// runSyncUnlink removes symlinks
func runSyncUnlink(cmd *cobra.Command, args []string) error {
	generatedDir := getGeneratedDir()

	fmt.Fprintf(os.Stdout, "%s Removing symlinks...\n", output.Info("→"))

	// Walk through generated directory to find what should be linked
	count := 0
	err := filepath.Walk(generatedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		// Get relative path from generated dir
		relPath, _ := filepath.Rel(generatedDir, path)
		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) < 2 {
			return nil
		}

		component := parts[0]
		filename := parts[len(parts)-1]

		// Determine XDG target
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig == "" {
			home, _ := os.UserHomeDir()
			xdgConfig = filepath.Join(home, ".config")
		}

		// Special case: shell scripts go to acorn/ directory, not shell/
		targetComponent := component
		if component == "shell" {
			targetComponent = "acorn"
		}
		target := filepath.Join(xdgConfig, targetComponent, filename)

		// Check if it's a symlink pointing to our generated file
		linkInfo, err := os.Lstat(target)
		if os.IsNotExist(err) {
			return nil // Nothing to unlink
		}

		if linkInfo.Mode()&os.ModeSymlink == 0 {
			fmt.Fprintf(os.Stdout, "  %s %s is not a symlink (skipping)\n", output.Warning("!"), target)
			return nil
		}

		linkDest, err := os.Readlink(target)
		if err != nil {
			return nil
		}

		// Only remove if it points to our generated directory
		if strings.HasPrefix(linkDest, generatedDir) {
			if err := os.Remove(target); err != nil {
				fmt.Fprintf(os.Stderr, "  %s Failed to remove %s: %v\n", output.Error("✗"), target, err)
			} else {
				fmt.Fprintf(os.Stdout, "  %s Removed %s\n", output.Success("✓"), target)
				count++
			}
		} else {
			fmt.Fprintf(os.Stdout, "  %s %s points elsewhere (skipping)\n", output.Warning("!"), target)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking generated directory: %w", err)
	}

	if count == 0 {
		fmt.Fprintf(os.Stdout, "%s No symlinks to remove\n", output.Info("ℹ"))
	} else {
		fmt.Fprintf(os.Stdout, "\n%s Removed %d symlink(s)\n", output.Success("✓"), count)
	}

	return nil
}

// runSyncUpdate pulls and regenerates
func runSyncUpdate(cmd *cobra.Command, args []string) error {
	// Pull
	fmt.Fprintln(os.Stdout, "Step 1/2: Pulling latest changes...")
	if err := runSyncPull(cmd, args); err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout)

	// Regenerate shell config
	fmt.Fprintln(os.Stdout, "Step 2/2: Regenerating shell configuration...")
	config := shell.NewConfig(false, false)
	manager := shell.NewManager(config)
	shell.RegisterAllComponents(manager)

	result, err := manager.GenerateAll()
	if err != nil {
		return fmt.Errorf("failed to generate shell config: %w", err)
	}

	fmt.Fprintf(os.Stdout, "%s Generated %d component scripts\n", output.Success("✓"), len(result.Scripts))
	if len(result.ConfigFiles) > 0 {
		fmt.Fprintf(os.Stdout, "%s Generated %d config files\n", output.Success("✓"), len(result.ConfigFiles))
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "%s Update complete. Run 'source %s/shell.sh' to reload.\n",
		output.Success("✓"), config.AcornDir)

	return nil
}

// Ensure configfile is used (prevents import error)
var _ = configfile.NewManager
