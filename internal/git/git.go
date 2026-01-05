// Package git provides Git helper functionality.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Info represents git repository info.
type Info struct {
	IsRepo     bool   `json:"is_repo" yaml:"is_repo"`
	Branch     string `json:"branch,omitempty" yaml:"branch,omitempty"`
	Remote     string `json:"remote,omitempty" yaml:"remote,omitempty"`
	RemoteURL  string `json:"remote_url,omitempty" yaml:"remote_url,omitempty"`
	Status     string `json:"status,omitempty" yaml:"status,omitempty"`
	Ahead      int    `json:"ahead,omitempty" yaml:"ahead,omitempty"`
	Behind     int    `json:"behind,omitempty" yaml:"behind,omitempty"`
	Staged     int    `json:"staged,omitempty" yaml:"staged,omitempty"`
	Modified   int    `json:"modified,omitempty" yaml:"modified,omitempty"`
	Untracked  int    `json:"untracked,omitempty" yaml:"untracked,omitempty"`
}

// Contributor represents a git contributor.
type Contributor struct {
	Name    string `json:"name" yaml:"name"`
	Commits int    `json:"commits" yaml:"commits"`
}

// Helper provides Git helper operations.
type Helper struct {
	verbose bool
	repoDir string
}

// NewHelper creates a new Git Helper.
func NewHelper(verbose bool) *Helper {
	repoDir := os.Getenv("DEFAULT_REPOS_DIR")
	if repoDir == "" {
		home, _ := os.UserHomeDir()
		repoDir = filepath.Join(home, "Repos")
	}

	return &Helper{
		verbose: verbose,
		repoDir: repoDir,
	}
}

// IsGitRepo checks if current directory is a git repository.
func (h *Helper) IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// GetInfo returns git repository info for current directory.
func (h *Helper) GetInfo() *Info {
	info := &Info{}

	if !h.IsGitRepo() {
		return info
	}

	info.IsRepo = true

	// Get branch
	if out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output(); err == nil {
		info.Branch = strings.TrimSpace(string(out))
	}

	// Get remote
	if out, err := exec.Command("git", "remote").Output(); err == nil {
		remotes := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(remotes) > 0 && remotes[0] != "" {
			info.Remote = remotes[0]
			// Get remote URL
			if urlOut, err := exec.Command("git", "remote", "get-url", info.Remote).Output(); err == nil {
				info.RemoteURL = strings.TrimSpace(string(urlOut))
			}
		}
	}

	// Get status counts
	if out, err := exec.Command("git", "status", "--porcelain").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if len(line) < 2 {
				continue
			}
			switch line[0] {
			case 'A', 'M', 'D', 'R', 'C':
				info.Staged++
			case '?':
				info.Untracked++
			}
			if line[1] == 'M' || line[1] == 'D' {
				info.Modified++
			}
		}
	}

	// Get ahead/behind
	if out, err := exec.Command("git", "rev-list", "--left-right", "--count", "HEAD...@{u}").Output(); err == nil {
		parts := strings.Fields(string(out))
		if len(parts) >= 2 {
			fmt.Sscanf(parts[0], "%d", &info.Ahead)
			fmt.Sscanf(parts[1], "%d", &info.Behind)
		}
	}

	return info
}

// GetContributors returns list of contributors.
func (h *Helper) GetContributors() ([]Contributor, error) {
	if !h.IsGitRepo() {
		return nil, fmt.Errorf("not a git repository")
	}

	cmd := exec.Command("git", "shortlog", "-sn", "--all")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var contributors []Contributor
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) == 2 {
			var commits int
			fmt.Sscanf(strings.TrimSpace(parts[0]), "%d", &commits)
			contributors = append(contributors, Contributor{
				Name:    strings.TrimSpace(parts[1]),
				Commits: commits,
			})
		}
	}

	return contributors, nil
}

// FindCommits finds commits by message.
func (h *Helper) FindCommits(search string) ([]string, error) {
	if !h.IsGitRepo() {
		return nil, fmt.Errorf("not a git repository")
	}

	cmd := exec.Command("git", "log", "--oneline", "--all", "--grep="+search)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var commits []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			commits = append(commits, line)
		}
	}

	return commits, nil
}

// CleanMergedBranches removes merged branches.
func (h *Helper) CleanMergedBranches(dryRun bool) ([]string, error) {
	if !h.IsGitRepo() {
		return nil, fmt.Errorf("not a git repository")
	}

	// Get merged branches
	cmd := exec.Command("git", "branch", "--merged")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var deleted []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		// Skip current branch (*) and main/master
		if branch == "" || strings.HasPrefix(branch, "*") {
			continue
		}
		if branch == "main" || branch == "master" || branch == "develop" {
			continue
		}

		if dryRun {
			deleted = append(deleted, branch)
			continue
		}

		// Delete the branch
		delCmd := exec.Command("git", "branch", "-d", branch)
		if err := delCmd.Run(); err == nil {
			deleted = append(deleted, branch)
		}
	}

	return deleted, nil
}

// GetReposDir returns the default repos directory.
func (h *Helper) GetReposDir() string {
	return h.repoDir
}
