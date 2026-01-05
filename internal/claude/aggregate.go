package claude

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// AggregateResult holds the result of an aggregation operation.
type AggregateResult struct {
	SearchDir      string          `json:"search_dir" yaml:"search_dir"`
	TargetDir      string          `json:"target_dir" yaml:"target_dir"`
	ReposScanned   int             `json:"repos_scanned" yaml:"repos_scanned"`
	AgentsAdded    int             `json:"agents_added" yaml:"agents_added"`
	CommandsAdded  int             `json:"commands_added" yaml:"commands_added"`
	SubagentsAdded int             `json:"subagents_added" yaml:"subagents_added"`
	Skipped        int             `json:"skipped" yaml:"skipped"`
	Renamed        int             `json:"renamed" yaml:"renamed"`
	Items          []AggregateItem `json:"items,omitempty" yaml:"items,omitempty"`
}

// AggregateItem represents an individual aggregated item.
type AggregateItem struct {
	Type       string `json:"type" yaml:"type"` // "agent", "command", "subagent"
	FileName   string `json:"file_name" yaml:"file_name"`
	SourceRepo string `json:"source_repo" yaml:"source_repo"`
	Action     string `json:"action" yaml:"action"` // "added", "skipped", "renamed"
}

// ListResult holds the list of all aggregated items.
type ListResult struct {
	Agents    []string `json:"agents" yaml:"agents"`
	Commands  []string `json:"commands" yaml:"commands"`
	Subagents []string `json:"subagents" yaml:"subagents"`
}

// ClearResult holds the result of a clear operation.
type ClearResult struct {
	Type    string   `json:"type" yaml:"type"`
	Cleared []string `json:"cleared" yaml:"cleared"`
}

// Aggregate scans repositories for .claude directories and aggregates content.
func (h *Helper) Aggregate(searchDir string) (*AggregateResult, error) {
	// Get dotfiles root from environment
	dotfilesRoot := os.Getenv("DOTFILES_ROOT")
	if dotfilesRoot == "" {
		home, _ := os.UserHomeDir()
		dotfilesRoot = filepath.Join(home, ".config", "dotfiles")
	}

	targetDir := filepath.Join(dotfilesRoot, "components", "claude", "config")

	if !h.DirExists(searchDir) {
		return nil, fmt.Errorf("search directory not found: %s", searchDir)
	}

	if !h.DirExists(targetDir) {
		return nil, fmt.Errorf("target directory not found: %s (ensure DOTFILES_ROOT is set)", targetDir)
	}

	result := &AggregateResult{
		SearchDir: searchDir,
		TargetDir: targetDir,
		Items:     []AggregateItem{},
	}

	// Find all .claude directories
	err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Only process .claude directories up to 3 levels deep
		relPath, _ := filepath.Rel(searchDir, path)
		depth := strings.Count(relPath, string(os.PathSeparator))
		if depth > 3 {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() || info.Name() != ".claude" {
			return nil
		}

		// Skip the dotfiles repo itself
		repoDir := filepath.Dir(path)
		if repoDir == dotfilesRoot {
			return filepath.SkipDir
		}

		repoName := filepath.Base(repoDir)
		foundSomething := false

		// Process agents
		agentsDir := filepath.Join(path, "agents")
		if h.DirExists(agentsDir) {
			items := h.processDirectory(agentsDir, filepath.Join(targetDir, "agents"), repoName, "agent")
			result.Items = append(result.Items, items...)
			for _, item := range items {
				switch item.Action {
				case "added", "renamed":
					result.AgentsAdded++
					foundSomething = true
				case "skipped":
					result.Skipped++
				}
				if item.Action == "renamed" {
					result.Renamed++
				}
			}
		}

		// Process commands
		commandsDir := filepath.Join(path, "commands")
		if h.DirExists(commandsDir) {
			items := h.processDirectory(commandsDir, filepath.Join(targetDir, "commands"), repoName, "command")
			result.Items = append(result.Items, items...)
			for _, item := range items {
				switch item.Action {
				case "added", "renamed":
					result.CommandsAdded++
					foundSomething = true
				case "skipped":
					result.Skipped++
				}
				if item.Action == "renamed" {
					result.Renamed++
				}
			}
		}

		// Process subagents
		subagentsDir := filepath.Join(path, "subagents")
		if h.DirExists(subagentsDir) {
			items := h.processDirectory(subagentsDir, filepath.Join(targetDir, "subagents"), repoName, "subagent")
			result.Items = append(result.Items, items...)
			for _, item := range items {
				switch item.Action {
				case "added", "renamed":
					result.SubagentsAdded++
					foundSomething = true
				case "skipped":
					result.Skipped++
				}
				if item.Action == "renamed" {
					result.Renamed++
				}
			}
		}

		if foundSomething {
			result.ReposScanned++
		}

		return filepath.SkipDir // Don't recurse into .claude directories
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// processDirectory processes files in a source directory and copies to target.
func (h *Helper) processDirectory(sourceDir, targetDir, repoName, itemType string) []AggregateItem {
	var items []AggregateItem

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return items
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		// Skip session files
		name := entry.Name()
		if name == "SESSION_CONTEXT.md" || name == "session-notes.md" || strings.HasSuffix(name, ".local.md") {
			continue
		}

		sourcePath := filepath.Join(sourceDir, name)
		targetPath := filepath.Join(targetDir, name)

		item := AggregateItem{
			Type:       itemType,
			FileName:   name,
			SourceRepo: repoName,
		}

		if h.FileExists(targetPath) {
			// Compare contents
			if h.filesEqual(sourcePath, targetPath) {
				item.Action = "skipped"
			} else {
				// Rename with repo prefix
				newName := repoName + "-" + name
				targetPath = filepath.Join(targetDir, newName)
				item.FileName = newName
				if !h.dryRun {
					if err := h.copyFile(sourcePath, targetPath); err == nil {
						item.Action = "renamed"
					}
				} else {
					item.Action = "renamed"
				}
			}
		} else {
			if !h.dryRun {
				if err := h.copyFile(sourcePath, targetPath); err == nil {
					item.Action = "added"
				}
			} else {
				item.Action = "added"
			}
		}

		items = append(items, item)
	}

	return items
}

// filesEqual compares two files for equality.
func (h *Helper) filesEqual(path1, path2 string) bool {
	data1, err1 := os.ReadFile(path1)
	data2, err2 := os.ReadFile(path2)
	if err1 != nil || err2 != nil {
		return false
	}
	return string(data1) == string(data2)
}

// copyFile copies a file from source to destination.
func (h *Helper) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// List returns all aggregated agents, commands, and subagents.
func (h *Helper) List() (*ListResult, error) {
	dotfilesRoot := os.Getenv("DOTFILES_ROOT")
	if dotfilesRoot == "" {
		home, _ := os.UserHomeDir()
		dotfilesRoot = filepath.Join(home, ".config", "dotfiles")
	}

	targetDir := filepath.Join(dotfilesRoot, "components", "claude", "config")

	result := &ListResult{
		Agents:    []string{},
		Commands:  []string{},
		Subagents: []string{},
	}

	// List agents
	agentsDir := filepath.Join(targetDir, "agents")
	if h.DirExists(agentsDir) {
		entries, _ := os.ReadDir(agentsDir)
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				result.Agents = append(result.Agents, strings.TrimSuffix(entry.Name(), ".md"))
			}
		}
		sort.Strings(result.Agents)
	}

	// List commands
	commandsDir := filepath.Join(targetDir, "commands")
	if h.DirExists(commandsDir) {
		entries, _ := os.ReadDir(commandsDir)
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				result.Commands = append(result.Commands, "/"+strings.TrimSuffix(entry.Name(), ".md"))
			}
		}
		sort.Strings(result.Commands)
	}

	// List subagents
	subagentsDir := filepath.Join(targetDir, "subagents")
	if h.DirExists(subagentsDir) {
		entries, _ := os.ReadDir(subagentsDir)
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				result.Subagents = append(result.Subagents, strings.TrimSuffix(entry.Name(), ".md"))
			}
		}
		sort.Strings(result.Subagents)
	}

	return result, nil
}

// Clear clears cache or stats.
func (h *Helper) Clear(what string) (*ClearResult, error) {
	result := &ClearResult{
		Type:    what,
		Cleared: []string{},
	}

	switch what {
	case "stats":
		if h.FileExists(h.paths.StatsCache) {
			if !h.dryRun {
				if err := os.Remove(h.paths.StatsCache); err != nil {
					return nil, err
				}
			}
			result.Cleared = append(result.Cleared, h.paths.StatsCache)
		}

	case "cache", "":
		// Clear shell-snapshots
		snapshotsDir := filepath.Join(h.paths.ClaudeDir, "shell-snapshots")
		if h.DirExists(snapshotsDir) {
			if !h.dryRun {
				if err := os.RemoveAll(snapshotsDir); err != nil {
					return nil, err
				}
			}
			result.Cleared = append(result.Cleared, snapshotsDir)
		}

		// Clear debug directory
		debugDir := filepath.Join(h.paths.ClaudeDir, "debug")
		if h.DirExists(debugDir) {
			if !h.dryRun {
				if err := os.RemoveAll(debugDir); err != nil {
					return nil, err
				}
			}
			result.Cleared = append(result.Cleared, debugDir)
		}

	default:
		return nil, fmt.Errorf("invalid clear type: %s (use cache or stats)", what)
	}

	return result, nil
}
