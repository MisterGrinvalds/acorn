// Package github provides GitHub CLI helper functionality.
package github

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Status represents GitHub CLI and repo status.
type Status struct {
	GhInstalled   bool   `json:"gh_installed" yaml:"gh_installed"`
	GhVersion     string `json:"gh_version,omitempty" yaml:"gh_version,omitempty"`
	Authenticated bool   `json:"authenticated" yaml:"authenticated"`
	Username      string `json:"username,omitempty" yaml:"username,omitempty"`
	CurrentBranch string `json:"current_branch,omitempty" yaml:"current_branch,omitempty"`
	HasRemote     bool   `json:"has_remote" yaml:"has_remote"`
}

// PRStatus represents pull request status.
type PRStatus struct {
	Branch    string `json:"branch" yaml:"branch"`
	HasPR     bool   `json:"has_pr" yaml:"has_pr"`
	PRNumber  string `json:"pr_number,omitempty" yaml:"pr_number,omitempty"`
	PRState   string `json:"pr_state,omitempty" yaml:"pr_state,omitempty"`
	Checks    string `json:"checks,omitempty" yaml:"checks,omitempty"`
}

// Helper provides GitHub CLI helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new GitHub Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// IsGhInstalled checks if gh CLI is installed.
func (h *Helper) IsGhInstalled() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// GetStatus returns GitHub CLI and repo status.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	// Check if gh is installed
	if _, err := exec.LookPath("gh"); err != nil {
		return status
	}
	status.GhInstalled = true

	// Get version
	if out, err := exec.Command("gh", "--version").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		if len(lines) > 0 {
			status.GhVersion = strings.TrimSpace(lines[0])
		}
	}

	// Check authentication
	if out, err := exec.Command("gh", "auth", "status").CombinedOutput(); err == nil {
		status.Authenticated = true
		// Try to extract username
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, "Logged in to") {
				parts := strings.Split(line, "as")
				if len(parts) > 1 {
					status.Username = strings.TrimSpace(strings.Split(parts[1], "(")[0])
				}
			}
		}
	}

	// Get current branch
	if out, err := exec.Command("git", "branch", "--show-current").Output(); err == nil {
		status.CurrentBranch = strings.TrimSpace(string(out))
	}

	// Check if has remote
	if out, err := exec.Command("git", "remote").Output(); err == nil {
		status.HasRemote = strings.TrimSpace(string(out)) != ""
	}

	return status
}

// CleanupBranches removes merged branches.
func (h *Helper) CleanupBranches() ([]string, error) {
	var deleted []string

	// Get merged branches
	cmd := exec.Command("git", "branch", "--merged")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get merged branches: %w", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		branch := strings.TrimSpace(line)
		if branch == "" || strings.HasPrefix(branch, "*") {
			continue
		}
		if branch == "main" || branch == "master" || branch == "develop" {
			continue
		}

		if h.dryRun {
			deleted = append(deleted, branch)
			continue
		}

		delCmd := exec.Command("git", "branch", "-d", branch)
		if err := delCmd.Run(); err == nil {
			deleted = append(deleted, branch)
		}
	}

	// Prune remote tracking branches
	if !h.dryRun {
		exec.Command("git", "remote", "prune", "origin").Run()
	}

	return deleted, nil
}

// PushBranch pushes current branch to origin with upstream tracking.
func (h *Helper) PushBranch() error {
	// Get current branch
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	branch := strings.TrimSpace(string(out))

	if h.dryRun {
		fmt.Printf("[dry-run] would run: git push -u origin %s\n", branch)
		return nil
	}

	cmd := exec.Command("git", "push", "-u", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CreatePR creates a pull request (opens in browser).
func (h *Helper) CreatePR() error {
	if !h.IsGhInstalled() {
		return fmt.Errorf("GitHub CLI (gh) is not installed")
	}

	// Get current branch
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	branch := strings.TrimSpace(string(out))

	if branch == "main" || branch == "master" {
		return fmt.Errorf("cannot create PR from main/master branch")
	}

	// Push branch first
	if err := h.PushBranch(); err != nil {
		return fmt.Errorf("failed to push branch: %w", err)
	}

	if h.dryRun {
		fmt.Println("[dry-run] would run: gh pr create --web")
		return nil
	}

	cmd := exec.Command("gh", "pr", "create", "--web")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetPRStatus returns status of current PR.
func (h *Helper) GetPRStatus() error {
	if !h.IsGhInstalled() {
		return fmt.Errorf("GitHub CLI (gh) is not installed")
	}

	cmd := exec.Command("gh", "pr", "status")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetPRChecks returns PR checks status.
func (h *Helper) GetPRChecks() error {
	if !h.IsGhInstalled() {
		return fmt.Errorf("GitHub CLI (gh) is not installed")
	}

	cmd := exec.Command("gh", "pr", "checks")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// WatchRun watches current workflow run.
func (h *Helper) WatchRun() error {
	if !h.IsGhInstalled() {
		return fmt.Errorf("GitHub CLI (gh) is not installed")
	}

	cmd := exec.Command("gh", "run", "watch")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// RerunFailed reruns failed workflow jobs.
func (h *Helper) RerunFailed() error {
	if !h.IsGhInstalled() {
		return fmt.Errorf("GitHub CLI (gh) is not installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would run: gh run rerun --failed")
		return nil
	}

	cmd := exec.Command("gh", "run", "rerun", "--failed")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// QuickCommit adds all changes and commits with message.
func (h *Helper) QuickCommit(message string) error {
	if message == "" {
		return fmt.Errorf("commit message is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: git add -A && git commit -m %q\n", message)
		return nil
	}

	// Add all
	addCmd := exec.Command("git", "add", "-A")
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	// Commit
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	return commitCmd.Run()
}

// NewBranch creates and checks out a new branch.
func (h *Helper) NewBranch(name string) error {
	if name == "" {
		return fmt.Errorf("branch name is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: git checkout -b %s\n", name)
		return nil
	}

	cmd := exec.Command("git", "checkout", "-b", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
