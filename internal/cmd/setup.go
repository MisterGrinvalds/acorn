package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mistergrinvalds/acorn/internal/componentconfig"
	"github.com/mistergrinvalds/acorn/internal/filesync"
	"github.com/mistergrinvalds/acorn/internal/output"
	"github.com/mistergrinvalds/acorn/internal/shell"
	"github.com/spf13/cobra"
)

var (
	setupDryRun  bool
	setupVerbose bool
	setupSkipBuild bool
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Full environment setup",
	Long: `Perform complete environment setup in one command.

This command stages everything needed for a working acorn environment:

  1. Build acorn binary (go build)
  2. Generate shell integration scripts
  3. Inject acorn into shell configuration
  4. Create symlinks for generated config files
  5. Sync component configurations (e.g., claude)

This is the recommended way to set up acorn on a new machine or after
pulling changes from the dotfiles repository.

Examples:
  acorn setup                    # Full setup
  acorn setup --dry-run          # Preview what would be done
  acorn setup --skip-build       # Skip go build step
  acorn setup -v                 # Verbose output`,
	RunE: runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.Flags().BoolVar(&setupDryRun, "dry-run", false, "Show what would be done without executing")
	setupCmd.Flags().BoolVarP(&setupVerbose, "verbose", "v", false, "Show verbose output")
	setupCmd.Flags().BoolVar(&setupSkipBuild, "skip-build", false, "Skip the go build step")
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

	config := shell.NewConfig(setupVerbose, setupDryRun)
	manager := shell.NewManager(config)
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

	generatedDir := filepath.Join(getDotfilesRootOrDefault(), "generated")

	// Check if generated directory exists
	if _, err := os.Stat(generatedDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stdout, "  %s No generated files to link\n\n", output.Info("○"))
		return nil
	}

	if setupDryRun {
		fmt.Fprintf(os.Stdout, "  %s Would create symlinks from %s\n\n", output.Info("○"), generatedDir)
		return nil
	}

	// Count files that would be linked
	count := 0
	err := filepath.Walk(generatedDir, func(path string, info os.FileInfo, err error) error {
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
	components := []string{"claude"} // Add more as needed

	syncedCount := 0
	for _, component := range components {
		loader := componentconfig.NewLoader()
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
