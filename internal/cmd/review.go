package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mistergrinvalds/acorn/internal/utils/component"
	"github.com/mistergrinvalds/acorn/internal/utils/config"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	reviewAll bool
)

// ReviewStatus represents the review status of all components.
type ReviewStatus struct {
	LastUpdated time.Time                   `yaml:"last_updated"`
	Components  map[string]*ComponentReview `yaml:"components"`
}

// ComponentReview represents the review status of a single component.
type ComponentReview struct {
	Reviewed    bool      `yaml:"reviewed"`
	ReviewedAt  time.Time `yaml:"reviewed_at,omitempty"`
	ReviewedBy  string    `yaml:"reviewed_by,omitempty"`
	Notes       string    `yaml:"notes,omitempty"`
	Checklist   Checklist `yaml:"checklist"`
	TestsPassed bool      `yaml:"tests_passed"`
}

// Checklist represents the review checklist for a component.
type Checklist struct {
	HasNameDescVersion bool `yaml:"has_name_desc_version"`
	HasFeatures        bool `yaml:"has_features"`
	HasInstallCheck    bool `yaml:"has_install_check"`
	TestsPass          bool `yaml:"tests_pass"`
}

// reviewCmd represents the review command
var reviewCmd = &cobra.Command{
	Use:   "review [component]",
	Short: "Review component configurations",
	Long: `Interactive review of component configurations.

Reviews check:
  - config.yaml has name, description, version
  - Has at least one of: env, aliases, shell_functions, wrappers
  - install.tools has check command defined
  - All tests pass

Examples:
  acorn review python     # Review specific component
  acorn review --all      # Review all components
  acorn review --status   # Show review status`,
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: completeComponentNames,
	RunE:              runReview,
}

func init() {
	rootCmd.AddCommand(reviewCmd)

	reviewCmd.Flags().BoolVar(&reviewAll, "all", false,
		"Review all components")
	// Output format is inherited from root command
}

func runReview(cmd *cobra.Command, args []string) error {
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return err
	}

	statusPath := filepath.Join(dotfilesRoot, ".sapling", "review-status.yaml")
	status, err := loadReviewStatus(statusPath)
	if err != nil {
		// Create new status if not found
		status = &ReviewStatus{
			LastUpdated: time.Now(),
			Components:  make(map[string]*ComponentReview),
		}
	}

	if len(args) == 1 {
		// Review specific component
		return reviewComponent(dotfilesRoot, args[0], status, statusPath)
	}

	if reviewAll {
		// Review all components
		return reviewAllComponents(dotfilesRoot, status, statusPath)
	}

	// Show current review status
	return showReviewStatus(dotfilesRoot, status)
}

func reviewComponent(dotfilesRoot, name string, status *ReviewStatus, statusPath string) error {
	configPath := filepath.Join(dotfilesRoot, ".sapling", "config", name)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("component not found: %s", name)
	}

	// Load config
	cfg, err := loadComponentConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create review
	review := &ComponentReview{
		Checklist: Checklist{},
	}

	// Check name, description, version
	review.Checklist.HasNameDescVersion = cfg.Name != "" && cfg.Description != "" && cfg.Version != ""

	// Check for features
	review.Checklist.HasFeatures = len(cfg.Env) > 0 ||
		len(cfg.Aliases) > 0 ||
		len(cfg.ShellFunctions) > 0 ||
		len(cfg.Wrappers) > 0 ||
		len(cfg.Files) > 0

	// Check install config
	if cfg.HasInstall() {
		hasCheck := true
		for _, tool := range cfg.Install.Tools {
			if tool.Check == "" {
				hasCheck = false
				break
			}
		}
		review.Checklist.HasInstallCheck = hasCheck
	} else {
		review.Checklist.HasInstallCheck = true // No install config is ok
	}

	// Run tests
	tester := component.NewTester(dotfilesRoot)
	testResult, err := tester.TestComponent(name)
	if err != nil {
		review.TestsPassed = false
	} else {
		review.TestsPassed = !testResult.HasFailures()
		review.Checklist.TestsPass = review.TestsPassed
	}

	// Print review results
	fmt.Fprintf(os.Stdout, "Review: %s\n", output.Info(name))
	fmt.Fprintln(os.Stdout, "─────────────────────────────────")

	printChecklistItem("name, description, version", review.Checklist.HasNameDescVersion)
	printChecklistItem("has features (env/aliases/functions)", review.Checklist.HasFeatures)
	printChecklistItem("install has check command", review.Checklist.HasInstallCheck)
	printChecklistItem("all tests pass", review.Checklist.TestsPass)

	fmt.Fprintln(os.Stdout)

	// Check if all items pass
	allPass := review.Checklist.HasNameDescVersion &&
		review.Checklist.HasFeatures &&
		review.Checklist.HasInstallCheck &&
		review.Checklist.TestsPass

	if allPass {
		review.Reviewed = true
		review.ReviewedAt = time.Now()
		fmt.Fprintln(os.Stdout, output.Success("Component passed all review checks"))
	} else {
		review.Reviewed = false
		fmt.Fprintln(os.Stdout, output.Warning("Component needs attention"))
	}

	// Save status
	status.Components[name] = review
	status.LastUpdated = time.Now()

	return saveReviewStatus(statusPath, status)
}

func reviewAllComponents(dotfilesRoot string, status *ReviewStatus, statusPath string) error {
	configDir := filepath.Join(dotfilesRoot, ".sapling", "config")

	entries, err := os.ReadDir(configDir)
	if err != nil {
		return fmt.Errorf("failed to read config directory: %w", err)
	}

	passedCount := 0
	failedCount := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Review this component
		err := reviewComponent(dotfilesRoot, name, status, statusPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reviewing %s: %v\n", name, err)
			failedCount++
			continue
		}

		if review, ok := status.Components[name]; ok && review.Reviewed {
			passedCount++
		} else {
			failedCount++
		}

		fmt.Fprintln(os.Stdout)
	}

	// Summary
	fmt.Fprintln(os.Stdout, "═══════════════════════════════════")
	fmt.Fprintf(os.Stdout, "Summary: %s passed, %s need attention\n",
		output.Success(fmt.Sprintf("%d", passedCount)),
		output.Warning(fmt.Sprintf("%d", failedCount)))

	return nil
}

func showReviewStatus(dotfilesRoot string, status *ReviewStatus) error {
	configDir := filepath.Join(dotfilesRoot, ".sapling", "config")

	entries, err := os.ReadDir(configDir)
	if err != nil {
		return fmt.Errorf("failed to read config directory: %w", err)
	}

	fmt.Fprintln(os.Stdout, "Component Review Status")
	fmt.Fprintln(os.Stdout, "═══════════════════════════════════")
	fmt.Fprintln(os.Stdout)

	reviewed := 0
	pending := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		symbol := output.Colorize("○", output.ColorYellow)
		statusText := "pending"

		if review, ok := status.Components[name]; ok && review.Reviewed {
			symbol = output.Colorize("✓", output.ColorGreen)
			statusText = "reviewed"
			reviewed++
		} else {
			pending++
		}

		fmt.Fprintf(os.Stdout, "%s %s (%s)\n", symbol, name, statusText)
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Total: %d components (%d reviewed, %d pending)\n",
		reviewed+pending, reviewed, pending)

	if status.LastUpdated.IsZero() {
		fmt.Fprintln(os.Stdout, "Last updated: never")
	} else {
		fmt.Fprintf(os.Stdout, "Last updated: %s\n", status.LastUpdated.Format(time.RFC3339))
	}

	return nil
}

func printChecklistItem(label string, passed bool) {
	symbol := output.Colorize("✓", output.ColorGreen)
	if !passed {
		symbol = output.Colorize("✗", output.ColorRed)
	}
	fmt.Fprintf(os.Stdout, "  %s %s\n", symbol, label)
}

func loadComponentConfig(configPath string) (*config.BaseConfig, error) {
	yamlPath := filepath.Join(configPath, "config.yaml")

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	var cfg config.BaseConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadReviewStatus(path string) (*ReviewStatus, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var status ReviewStatus
	if err := yaml.Unmarshal(data, &status); err != nil {
		return nil, err
	}

	if status.Components == nil {
		status.Components = make(map[string]*ComponentReview)
	}

	return &status, nil
}

func saveReviewStatus(path string, status *ReviewStatus) error {
	data, err := yaml.Marshal(status)
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
