package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/components/filesync"
	"github.com/mistergrinvalds/acorn/internal/components/shell"
	"github.com/mistergrinvalds/acorn/internal/utils/config"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	setupDryRun      bool
	setupVerbose     bool
	setupSkipBuild   bool
	setupSaplingRepo string
	setupSaplingPath string
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Full environment setup",
	Long: `Perform complete environment setup in one command.

This command stages everything needed for a working acorn environment:

  0. Setup .sapling repository (clone or link existing)
  1. Build acorn binary (go build)
  2. Generate shell integration scripts
  3. Inject acorn into shell configuration
  4. Create symlinks for generated config files
  5. Sync component configurations (e.g., claude)

This is the recommended way to set up acorn on a new machine or after
pulling changes from the dotfiles repository.

Examples:
  acorn setup                                    # Full setup (interactive)
  acorn setup --sapling-repo git@github.com:user/sapling.git  # Clone sapling repo
  acorn setup --sapling-path ~/my-sapling        # Use existing sapling directory
  acorn setup --dry-run                          # Preview what would be done
  acorn setup --skip-build                       # Skip go build step
  acorn setup -v                                 # Verbose output`,
	RunE: runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.Flags().BoolVar(&setupDryRun, "dry-run", false, "Show what would be done without executing")
	setupCmd.Flags().BoolVarP(&setupVerbose, "verbose", "v", false, "Show verbose output")
	setupCmd.Flags().BoolVar(&setupSkipBuild, "skip-build", false, "Skip the go build step")
	setupCmd.Flags().StringVar(&setupSaplingRepo, "sapling-repo", "", "Git repository URL to clone .sapling from")
	setupCmd.Flags().StringVar(&setupSaplingPath, "sapling-path", "", "Path to existing .sapling directory to link")
}

func runSetup(cmd *cobra.Command, args []string) error {
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return fmt.Errorf("failed to get dotfiles root: %w", err)
	}

	if setupDryRun {
		fmt.Fprintf(os.Stdout, "%s Setup Preview (dry-run)\n", output.Info("ℹ"))
	} else {
		fmt.Fprintf(os.Stdout, "%s Acorn Environment Setup\n", output.Info("ℹ"))
	}
	fmt.Fprintf(os.Stdout, "  Dotfiles: %s\n\n", dotfilesRoot)

	// Step 0: Setup .sapling repository
	if err := setupSapling(dotfilesRoot); err != nil {
		return err
	}

	// Step 1: Build acorn
	if !setupSkipBuild {
		if err := setupBuild(dotfilesRoot); err != nil {
			return err
		}
	} else if setupVerbose {
		fmt.Fprintf(os.Stdout, "%s Skipping build step\n\n", output.Info("○"))
	}

	// Step 2: Generate shell scripts
	if err := setupShellGenerate(); err != nil {
		return err
	}

	// Step 3: Inject into shell rc
	if err := setupShellInject(); err != nil {
		return err
	}

	// Step 4: Create symlinks for generated files
	if err := setupSyncLink(); err != nil {
		return err
	}

	// Step 5: Sync component configurations
	if err := setupComponentSync(dotfilesRoot); err != nil {
		return err
	}

	// Summary
	fmt.Fprintln(os.Stdout)
	if setupDryRun {
		fmt.Fprintf(os.Stdout, "%s Setup preview complete. Run without --dry-run to apply.\n", output.Success("✓"))
	} else {
		fmt.Fprintf(os.Stdout, "%s Setup complete!\n", output.Success("✓"))
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Next steps:")
		fmt.Fprintln(os.Stdout, "  1. Restart your shell or run: source ~/.bashrc (or ~/.zshrc)")
		fmt.Fprintln(os.Stdout, "  2. Verify with: acorn version")
	}

	return nil
}

// setupSapling sets up the .sapling repository
func setupSapling(dotfilesRoot string) error {
	fmt.Fprintf(os.Stdout, "Step 0: Setting up .sapling repository\n")

	saplingDir := filepath.Join(dotfilesRoot, ".sapling")
	home, _ := os.UserHomeDir()
	homeSaplingLink := filepath.Join(home, ".sapling")

	// Check if .sapling is a REAL sapling repo (has config/ or .git)
	// Just having generated/ doesn't count - that's auto-created by shell generate
	if isValidSaplingRepo(saplingDir) {
		fmt.Fprintf(os.Stdout, "  %s .sapling repository exists at %s\n", output.Success("✓"), saplingDir)

		// Ensure ~/.sapling symlink exists
		if err := ensureSaplingSymlink(homeSaplingLink, saplingDir); err != nil {
			fmt.Fprintf(os.Stdout, "  %s Could not create ~/.sapling symlink: %v\n", output.Warning("!"), err)
		}

		fmt.Fprintln(os.Stdout)
		return nil
	}

	// Clean up any auto-created .sapling directory that's not a real repo
	if info, err := os.Stat(saplingDir); err == nil && info.IsDir() {
		fmt.Fprintf(os.Stdout, "  %s Found incomplete .sapling directory (no config/ or .git)\n", output.Warning("!"))
		fmt.Fprintf(os.Stdout, "  %s This was likely auto-created. Removing it...\n", output.Info("○"))
		if !setupDryRun {
			os.RemoveAll(saplingDir)
		}
	}

	// .sapling doesn't exist - need to set it up
	if setupDryRun {
		fmt.Fprintf(os.Stdout, "  %s Would setup .sapling repository\n\n", output.Info("○"))
		return nil
	}

	// Check for command-line flags first
	if setupSaplingRepo != "" {
		return cloneSaplingRepo(saplingDir, setupSaplingRepo, homeSaplingLink)
	}

	if setupSaplingPath != "" {
		return linkSaplingPath(saplingDir, setupSaplingPath, homeSaplingLink)
	}

	// Interactive mode - ask user what to do
	return setupSaplingInteractive(saplingDir, homeSaplingLink)
}

// setupSaplingInteractive prompts the user for sapling setup options
func setupSaplingInteractive(saplingDir, homeSaplingLink string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "  The .sapling directory contains your dotfiles configurations.")
	fmt.Fprintln(os.Stdout, "  How would you like to set it up?")
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "  1) Clone from a git repository")
	fmt.Fprintln(os.Stdout, "  2) Link to an existing .sapling directory")
	fmt.Fprintln(os.Stdout, "  3) Initialize a new empty .sapling")
	fmt.Fprintln(os.Stdout, "  4) Skip (I'll set it up manually)")
	fmt.Fprintln(os.Stdout)
	fmt.Fprint(os.Stdout, "  Choose an option [1-4]: ")

	choice, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		fmt.Fprint(os.Stdout, "  Enter git repository URL: ")
		repoURL, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		repoURL = strings.TrimSpace(repoURL)
		if repoURL == "" {
			return fmt.Errorf("repository URL cannot be empty")
		}
		return cloneSaplingRepo(saplingDir, repoURL, homeSaplingLink)

	case "2":
		fmt.Fprint(os.Stdout, "  Enter path to existing .sapling directory: ")
		existingPath, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		existingPath = strings.TrimSpace(existingPath)
		if existingPath == "" {
			return fmt.Errorf("path cannot be empty")
		}
		// Expand ~ if present
		if strings.HasPrefix(existingPath, "~/") {
			home, _ := os.UserHomeDir()
			existingPath = filepath.Join(home, existingPath[2:])
		}
		return linkSaplingPath(saplingDir, existingPath, homeSaplingLink)

	case "3":
		return initSaplingRepo(saplingDir, homeSaplingLink)

	case "4":
		fmt.Fprintf(os.Stdout, "  %s Skipping .sapling setup\n\n", output.Info("○"))
		return nil

	default:
		return fmt.Errorf("invalid choice: %s", choice)
	}
}

// cloneSaplingRepo clones a git repository to .sapling
func cloneSaplingRepo(saplingDir, repoURL, homeSaplingLink string) error {
	fmt.Fprintf(os.Stdout, "  Cloning %s...\n", repoURL)

	cmd := exec.Command("git", "clone", repoURL, saplingDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	fmt.Fprintf(os.Stdout, "  %s Cloned to %s\n", output.Success("✓"), saplingDir)

	// Create ~/.sapling symlink
	if err := ensureSaplingSymlink(homeSaplingLink, saplingDir); err != nil {
		fmt.Fprintf(os.Stdout, "  %s Could not create ~/.sapling symlink: %v\n", output.Warning("!"), err)
	}

	fmt.Fprintln(os.Stdout)
	return nil
}

// linkSaplingPath creates a symlink from .sapling to an existing directory
func linkSaplingPath(saplingDir, existingPath, homeSaplingLink string) error {
	// Verify the existing path exists
	info, err := os.Stat(existingPath)
	if err != nil {
		return fmt.Errorf("path does not exist: %s", existingPath)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", existingPath)
	}

	// Verify it looks like a sapling directory (has config/ subdirectory)
	configDir := filepath.Join(existingPath, "config")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stdout, "  %s Warning: %s doesn't contain a config/ directory\n", output.Warning("!"), existingPath)
	}

	// Create symlink
	fmt.Fprintf(os.Stdout, "  Creating symlink %s -> %s\n", saplingDir, existingPath)
	if err := os.Symlink(existingPath, saplingDir); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	fmt.Fprintf(os.Stdout, "  %s Linked .sapling to %s\n", output.Success("✓"), existingPath)

	// Create ~/.sapling symlink
	if err := ensureSaplingSymlink(homeSaplingLink, saplingDir); err != nil {
		fmt.Fprintf(os.Stdout, "  %s Could not create ~/.sapling symlink: %v\n", output.Warning("!"), err)
	}

	fmt.Fprintln(os.Stdout)
	return nil
}

// initSaplingRepo initializes a new empty .sapling repository
func initSaplingRepo(saplingDir, homeSaplingLink string) error {
	fmt.Fprintf(os.Stdout, "  Initializing new .sapling repository...\n")

	// Create directory structure
	dirs := []string{
		saplingDir,
		filepath.Join(saplingDir, "config"),
		filepath.Join(saplingDir, "generated"),
		filepath.Join(saplingDir, "ai"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create .gitignore
	gitignore := `# Prevent circular symlinks
.sapling

# Optional: ignore generated files (uncomment to not track them)
# generated/
`
	gitignorePath := filepath.Join(saplingDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignore), 0o644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// Create README
	readme := `# Sapling Dotfiles Configuration

This directory contains your dotfiles configurations managed by Acorn.

## Structure

- config/     - Component configuration files (config.yaml per component)
- generated/  - Generated shell scripts and config files
- ai/         - AI tool configurations (Claude, etc.)

## Usage

Run ` + "`acorn setup`" + ` to generate and link your configurations.
`
	readmePath := filepath.Join(saplingDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readme), 0o644); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = saplingDir
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stdout, "  %s Could not initialize git repo: %v\n", output.Warning("!"), err)
	} else {
		fmt.Fprintf(os.Stdout, "  %s Initialized git repository\n", output.Success("✓"))
	}

	fmt.Fprintf(os.Stdout, "  %s Created .sapling at %s\n", output.Success("✓"), saplingDir)

	// Create ~/.sapling symlink
	if err := ensureSaplingSymlink(homeSaplingLink, saplingDir); err != nil {
		fmt.Fprintf(os.Stdout, "  %s Could not create ~/.sapling symlink: %v\n", output.Warning("!"), err)
	}

	fmt.Fprintln(os.Stdout)
	return nil
}

// isValidSaplingRepo checks if a directory is a real sapling repository
// A valid sapling repo must have either config/ directory or .git directory
// Just having generated/ is not enough (that gets auto-created)
func isValidSaplingRepo(saplingDir string) bool {
	// Check if directory exists
	info, err := os.Stat(saplingDir)
	if err != nil || !info.IsDir() {
		return false
	}

	// Check for config/ directory (primary indicator)
	configDir := filepath.Join(saplingDir, "config")
	if _, err := os.Stat(configDir); err == nil {
		return true
	}

	// Check for .git directory (it's a git repo)
	gitDir := filepath.Join(saplingDir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		return true
	}

	// Check for ai/ directory (another indicator)
	aiDir := filepath.Join(saplingDir, "ai")
	if _, err := os.Stat(aiDir); err == nil {
		return true
	}

	return false
}

// ensureSaplingSymlink creates ~/.sapling -> saplingDir symlink if it doesn't exist
func ensureSaplingSymlink(homeSaplingLink, saplingDir string) error {
	// Check if symlink already exists and points to the right place
	if target, err := os.Readlink(homeSaplingLink); err == nil {
		if target == saplingDir {
			return nil // Already correct
		}
		// Remove incorrect symlink
		os.Remove(homeSaplingLink)
	} else if !os.IsNotExist(err) {
		// It exists but is not a symlink - check if it's a directory
		if info, statErr := os.Stat(homeSaplingLink); statErr == nil {
			if info.IsDir() {
				return fmt.Errorf("~/.sapling exists as a directory, not a symlink")
			}
		}
	}

	// Create the symlink
	if err := os.Symlink(saplingDir, homeSaplingLink); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "  %s Created ~/.sapling -> %s\n", output.Success("✓"), saplingDir)
	return nil
}

// setupBuild builds the acorn binary
func setupBuild(dotfilesRoot string) error {
	fmt.Fprintf(os.Stdout, "Step 1: Building acorn binary\n")

	binDir := filepath.Join(dotfilesRoot, "bin")
	binPath := filepath.Join(binDir, "acorn")

	if setupDryRun {
		fmt.Fprintf(os.Stdout, "  %s Would run: go build -o %s ./cmd/acorn\n", output.Info("○"), binPath)
		fmt.Fprintln(os.Stdout)
		return nil
	}

	// Ensure bin directory exists
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Build acorn
	buildCmd := exec.Command("go", "build", "-o", binPath, "./cmd/acorn")
	buildCmd.Dir = dotfilesRoot
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if setupVerbose {
		fmt.Fprintf(os.Stdout, "  Running: go build -o %s ./cmd/acorn\n", binPath)
	}

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	fmt.Fprintf(os.Stdout, "  %s Built %s\n\n", output.Success("✓"), binPath)
	return nil
}

// setupShellGenerate generates shell integration scripts
func setupShellGenerate() error {
	fmt.Fprintf(os.Stdout, "Step 2: Generating shell scripts\n")

	// Check for valid sapling repo
	if !config.IsValidSaplingRepo() {
		if setupDryRun {
			fmt.Fprintf(os.Stdout, "  %s Would generate shell scripts (after .sapling setup)\n\n", output.Info("○"))
			return nil
		}
		return fmt.Errorf("no valid .sapling repository. Step 0 should have created one")
	}

	cfg := shell.NewConfig(setupVerbose, setupDryRun)
	manager := shell.NewManager(cfg)
	shell.RegisterAllComponents(manager)

	result, err := manager.GenerateAll()
	if err != nil {
		return fmt.Errorf("shell generate failed: %w", err)
	}

	if setupDryRun {
		fmt.Fprintf(os.Stdout, "  %s Would generate %d component scripts\n", output.Info("○"), len(result.Scripts))
		if len(result.ConfigFiles) > 0 {
			fmt.Fprintf(os.Stdout, "  %s Would generate %d config files\n", output.Info("○"), len(result.ConfigFiles))
		}
	} else {
		fmt.Fprintf(os.Stdout, "  %s Generated %d component scripts\n", output.Success("✓"), len(result.Scripts))
		if len(result.ConfigFiles) > 0 {
			fmt.Fprintf(os.Stdout, "  %s Generated %d config files\n", output.Success("✓"), len(result.ConfigFiles))
		}
	}

	if setupVerbose {
		for _, script := range result.Scripts {
			fmt.Fprintf(os.Stdout, "    - %s\n", script.Component)
		}
	}

	fmt.Fprintln(os.Stdout)
	return nil
}

// setupShellInject injects acorn into shell rc
func setupShellInject() error {
	fmt.Fprintf(os.Stdout, "Step 3: Injecting into shell configuration\n")

	config := shell.NewConfig(setupVerbose, setupDryRun)
	manager := shell.NewManager(config)

	result, err := manager.Inject()
	if err != nil {
		return fmt.Errorf("shell inject failed: %w", err)
	}

	switch result.Action {
	case "injected":
		fmt.Fprintf(os.Stdout, "  %s Injected into %s\n", output.Success("✓"), result.RCFile)
	case "already_injected":
		fmt.Fprintf(os.Stdout, "  %s Already injected in %s\n", output.Info("○"), result.RCFile)
	default:
		if setupDryRun {
			fmt.Fprintf(os.Stdout, "  %s Would inject into %s\n", output.Info("○"), result.RCFile)
		}
	}

	fmt.Fprintln(os.Stdout)
	return nil
}

// setupSyncLink creates symlinks for generated config files
func setupSyncLink() error {
	fmt.Fprintf(os.Stdout, "Step 4: Creating config symlinks\n")

	// Use .sapling/generated directory
	generatedDir, err := config.GeneratedDir()
	if err != nil {
		fmt.Fprintf(os.Stdout, "  %s Could not find .sapling/generated directory\n\n", output.Warning("!"))
		return nil
	}

	// Check if generated directory exists
	if _, err := os.Stat(generatedDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stdout, "  %s No generated files to link (directory doesn't exist yet)\n\n", output.Info("○"))
		return nil
	}

	if setupDryRun {
		fmt.Fprintf(os.Stdout, "  %s Would create symlinks from %s\n\n", output.Info("○"), generatedDir)
		return nil
	}

	// Count files that would be linked
	count := 0
	err = filepath.Walk(generatedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		count++
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to scan generated directory: %w", err)
	}

	if count == 0 {
		fmt.Fprintf(os.Stdout, "  %s No generated files to link\n\n", output.Info("○"))
		return nil
	}

	// Run sync link
	if err := runSyncLink(nil, nil); err != nil {
		// Non-fatal - continue with setup
		fmt.Fprintf(os.Stdout, "  %s Symlink creation had issues (continuing)\n\n", output.Warning("!"))
		return nil
	}

	fmt.Fprintln(os.Stdout)
	return nil
}

// setupComponentSync syncs component-specific configurations
func setupComponentSync(dotfilesRoot string) error {
	fmt.Fprintf(os.Stdout, "Step 5: Syncing component configurations\n")

	// Get list of components with sync_files
	components := []string{"claude", "git", "karabiner", "python", "r", "ssh", "tmux", "vscode", "wget"} // Add more as needed

	syncedCount := 0
	for _, component := range components {
		loader := config.NewComponentLoader()
		cfg, err := loader.LoadBase(component)
		if err != nil {
			if setupVerbose {
				fmt.Fprintf(os.Stdout, "  %s Skipping %s (not found)\n", output.Info("○"), component)
			}
			continue
		}

		if !cfg.HasSyncFiles() {
			if setupVerbose {
				fmt.Fprintf(os.Stdout, "  %s %s has no sync files\n", output.Info("○"), component)
			}
			continue
		}

		syncer := filesync.NewSyncer(dotfilesRoot, setupDryRun, setupVerbose)
		result, err := syncer.Sync(cfg.GetSyncFiles())
		if err != nil {
			fmt.Fprintf(os.Stdout, "  %s %s sync failed: %v\n", output.Error("✗"), component, err)
			continue
		}

		if len(result.Synced) > 0 {
			if setupDryRun {
				fmt.Fprintf(os.Stdout, "  %s Would sync %d files for %s\n", output.Info("○"), len(result.Synced), component)
			} else {
				fmt.Fprintf(os.Stdout, "  %s Synced %d files for %s\n", output.Success("✓"), len(result.Synced), component)
			}
			syncedCount += len(result.Synced)
		} else if len(result.Skipped) > 0 {
			fmt.Fprintf(os.Stdout, "  %s %s already synced (%d files)\n", output.Info("○"), component, len(result.Skipped))
		}

		if len(result.Errors) > 0 {
			for _, e := range result.Errors {
				fmt.Fprintf(os.Stdout, "    %s %s: %s\n", output.Error("✗"), e.Source, e.Error)
			}
		}
	}

	if syncedCount == 0 && !setupDryRun {
		fmt.Fprintf(os.Stdout, "  %s All components already synced\n", output.Info("○"))
	}

	return nil
}

// getDotfilesRootOrDefault returns dotfiles root or a default
func getDotfilesRootOrDefault() string {
	root, err := getDotfilesRoot()
	if err != nil {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Repos", "personal", "tools")
	}
	return root
}
