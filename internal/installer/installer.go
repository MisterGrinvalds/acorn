package installer

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mistergrinvalds/acorn/internal/componentconfig"
)

// Installer handles component installation.
type Installer struct {
	loader   *componentconfig.Loader
	platform *Platform
	dryRun   bool
	verbose  bool
	stdout   io.Writer
	stderr   io.Writer
}

// Option configures the Installer.
type Option func(*Installer)

// NewInstaller creates a new Installer.
func NewInstaller(opts ...Option) *Installer {
	i := &Installer{
		loader:   componentconfig.NewLoader(),
		platform: DetectPlatform(),
		stdout:   os.Stdout,
		stderr:   os.Stderr,
	}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

// WithDryRun enables dry-run mode.
func WithDryRun(dryRun bool) Option {
	return func(i *Installer) { i.dryRun = dryRun }
}

// WithVerbose enables verbose output.
func WithVerbose(verbose bool) Option {
	return func(i *Installer) { i.verbose = verbose }
}

// WithOutput sets custom output writers.
func WithOutput(stdout, stderr io.Writer) Option {
	return func(i *Installer) {
		i.stdout = stdout
		i.stderr = stderr
	}
}

// Plan creates an installation plan for a component.
func (i *Installer) Plan(ctx context.Context, component string) (*InstallPlan, error) {
	cfg, err := i.loadInstallConfig(component)
	if err != nil {
		return nil, err
	}

	if len(cfg.Tools) == 0 {
		return nil, fmt.Errorf("component %s has no install configuration", component)
	}

	resolver := NewResolver(i.platform)
	plan, err := resolver.BuildPlan(component, cfg)
	if err != nil {
		return nil, err
	}

	plan.DryRun = i.dryRun
	return plan, nil
}

// Install executes the installation for a component.
func (i *Installer) Install(ctx context.Context, component string) (*InstallResult, error) {
	start := time.Now()

	plan, err := i.Plan(ctx, component)
	if err != nil {
		return nil, err
	}

	result := &InstallResult{
		Component: component,
		Success:   true,
		DryRun:    i.dryRun,
	}

	// Install prerequisites first
	for _, tool := range plan.Prerequisites {
		toolResult := i.installTool(ctx, tool)
		result.Tools = append(result.Tools, toolResult)
		if !toolResult.Success && !toolResult.Skipped {
			result.Success = false
		}
	}

	// Then install direct tools
	for _, tool := range plan.Tools {
		toolResult := i.installTool(ctx, tool)
		result.Tools = append(result.Tools, toolResult)
		if !toolResult.Success && !toolResult.Skipped {
			result.Success = false
		}
	}

	result.Duration = time.Since(start)
	return result, nil
}

// installTool installs a single tool.
func (i *Installer) installTool(ctx context.Context, tool PlannedTool) ToolResult {
	start := time.Now()

	result := ToolResult{
		Name: tool.Name,
	}

	// Skip if already installed
	if tool.AlreadyInstalled {
		result.Skipped = true
		result.SkipReason = "already installed"
		result.Version = tool.Version
		result.Success = true
		return result
	}

	// Dry run - just report what would be done
	if i.dryRun {
		result.Skipped = true
		result.SkipReason = "dry run"
		result.Success = true
		return result
	}

	// Get executor for this method
	executor, err := GetExecutor(tool.Method.Type)
	if err != nil {
		result.Error = err
		result.Success = false
		return result
	}

	// Check if method is available
	if !executor.Available() {
		result.Error = fmt.Errorf("%s not available on this system", tool.Method.Type)
		result.Success = false
		return result
	}

	// Execute installation
	if i.verbose {
		fmt.Fprintf(i.stdout, "Installing %s via %s...\n", tool.Name, tool.Method.Type)
	}

	if err := executor.Execute(ctx, tool, i.stdout, i.stderr); err != nil {
		result.Error = err
		result.Success = false
		return result
	}

	result.Success = true
	result.Duration = time.Since(start)

	// Show post-install message
	if tool.PostInstall.Message != "" && !i.dryRun {
		fmt.Fprintf(i.stdout, "\n%s\n", tool.PostInstall.Message)
	}

	return result
}

// loadInstallConfig loads the install config for a component.
func (i *Installer) loadInstallConfig(component string) (*componentconfig.InstallConfig, error) {
	cfg := &componentconfig.BaseConfig{}
	if err := i.loader.Load(component, cfg); err != nil {
		return nil, fmt.Errorf("failed to load config for %s: %w", component, err)
	}
	return &cfg.Install, nil
}

// GetPlatform returns the detected platform.
func (i *Installer) GetPlatform() *Platform {
	return i.platform
}
