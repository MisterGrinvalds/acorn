// Package filesync provides file synchronization with symlink and merge support.
package filesync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/componentconfig"
)

// Syncer handles file synchronization operations.
type Syncer struct {
	dotfilesRoot string
	dryRun       bool
	verbose      bool
}

// NewSyncer creates a new file syncer.
func NewSyncer(dotfilesRoot string, dryRun, verbose bool) *Syncer {
	return &Syncer{
		dotfilesRoot: dotfilesRoot,
		dryRun:       dryRun,
		verbose:      verbose,
	}
}

// SyncResult contains the result of a sync operation.
type SyncResult struct {
	Synced   []SyncedFile `json:"synced"`
	Skipped  []SyncedFile `json:"skipped"`
	Errors   []SyncError  `json:"errors"`
	DryRun   bool         `json:"dry_run"`
}

// SyncedFile represents a single synced file.
type SyncedFile struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Mode   string `json:"mode"`
	Action string `json:"action"` // "created", "updated", "unchanged"
}

// SyncError represents a sync error.
type SyncError struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Error  string `json:"error"`
}

// expandPath expands environment variables and ~ in paths.
func expandPath(path string) string {
	// Handle ${VAR:-default} syntax
	result := os.Expand(path, func(key string) string {
		// Check for :- default syntax
		if idx := strings.Index(key, ":-"); idx != -1 {
			varName := key[:idx]
			defaultVal := key[idx+2:]
			if val := os.Getenv(varName); val != "" {
				return val
			}
			return defaultVal
		}
		return os.Getenv(key)
	})

	// Handle ~ for home directory
	if strings.HasPrefix(result, "~/") {
		home, _ := os.UserHomeDir()
		result = filepath.Join(home, result[2:])
	}

	return result
}

// Sync synchronizes files according to the given configuration.
func (s *Syncer) Sync(files []componentconfig.SyncFileConfig) (*SyncResult, error) {
	result := &SyncResult{
		Synced:  []SyncedFile{},
		Skipped: []SyncedFile{},
		Errors:  []SyncError{},
		DryRun:  s.dryRun,
	}

	for _, fc := range files {
		sourcePath := filepath.Join(s.dotfilesRoot, fc.Source)
		targetPath := expandPath(fc.Target)

		var err error
		var action string

		switch fc.Mode {
		case "symlink":
			action, err = s.syncSymlink(sourcePath, targetPath)
		case "copy":
			action, err = s.syncCopy(sourcePath, targetPath)
		case "merge":
			action, err = s.syncMerge(sourcePath, targetPath, fc.MergeConfig)
		default:
			err = fmt.Errorf("unknown sync mode: %s", fc.Mode)
		}

		if err != nil {
			result.Errors = append(result.Errors, SyncError{
				Source: fc.Source,
				Target: targetPath,
				Error:  err.Error(),
			})
			continue
		}

		if action == "unchanged" {
			result.Skipped = append(result.Skipped, SyncedFile{
				Source: fc.Source,
				Target: targetPath,
				Mode:   fc.Mode,
				Action: action,
			})
		} else {
			result.Synced = append(result.Synced, SyncedFile{
				Source: fc.Source,
				Target: targetPath,
				Mode:   fc.Mode,
				Action: action,
			})
		}
	}

	return result, nil
}

// syncSymlink creates a symlink from target to source.
func (s *Syncer) syncSymlink(source, target string) (string, error) {
	// Check if source exists
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return "", fmt.Errorf("source not found: %s", source)
	}

	// Check current state of target
	linkInfo, err := os.Lstat(target)
	if err == nil {
		// Target exists
		if linkInfo.Mode()&os.ModeSymlink != 0 {
			// It's a symlink - check where it points
			currentDest, _ := os.Readlink(target)
			if currentDest == source {
				return "unchanged", nil
			}
			// Wrong target - remove and recreate
			if !s.dryRun {
				os.Remove(target)
			}
		} else if sourceInfo.IsDir() && linkInfo.IsDir() {
			// Both are directories - backup and replace
			if !s.dryRun {
				backup := target + ".backup"
				os.Rename(target, backup)
			}
		} else {
			// Regular file - backup
			if !s.dryRun {
				backup := target + ".backup"
				os.Rename(target, backup)
			}
		}
	}

	if s.dryRun {
		return "created", nil
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Create symlink
	if err := os.Symlink(source, target); err != nil {
		return "", fmt.Errorf("failed to create symlink: %w", err)
	}

	return "created", nil
}

// syncCopy copies source to target.
func (s *Syncer) syncCopy(source, target string) (string, error) {
	// Check if source exists
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return "", fmt.Errorf("source not found: %s", source)
	}

	if sourceInfo.IsDir() {
		return s.syncCopyDir(source, target)
	}

	return s.syncCopyFile(source, target)
}

// syncCopyFile copies a single file.
func (s *Syncer) syncCopyFile(source, target string) (string, error) {
	sourceData, err := os.ReadFile(source)
	if err != nil {
		return "", fmt.Errorf("failed to read source: %w", err)
	}

	// Check if target is identical
	targetData, err := os.ReadFile(target)
	if err == nil && string(sourceData) == string(targetData) {
		return "unchanged", nil
	}

	if s.dryRun {
		return "created", nil
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", fmt.Errorf("failed to create parent directory: %w", err)
	}

	if err := os.WriteFile(target, sourceData, 0o644); err != nil {
		return "", fmt.Errorf("failed to write target: %w", err)
	}

	return "created", nil
}

// syncCopyDir copies a directory recursively.
func (s *Syncer) syncCopyDir(source, target string) (string, error) {
	if s.dryRun {
		return "created", nil
	}

	// Ensure target directory exists
	if err := os.MkdirAll(target, 0o755); err != nil {
		return "", fmt.Errorf("failed to create target directory: %w", err)
	}

	entries, err := os.ReadDir(source)
	if err != nil {
		return "", fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(source, entry.Name())
		dstPath := filepath.Join(target, entry.Name())

		if entry.IsDir() {
			if _, err := s.syncCopyDir(srcPath, dstPath); err != nil {
				return "", err
			}
		} else {
			if _, err := s.syncCopyFile(srcPath, dstPath); err != nil {
				return "", err
			}
		}
	}

	return "created", nil
}

// syncMerge merges source JSON with user's local JSON.
func (s *Syncer) syncMerge(source, target string, cfg *componentconfig.MergeConfig) (string, error) {
	// Read source (base) file
	sourceData, err := os.ReadFile(source)
	if err != nil {
		return "", fmt.Errorf("failed to read source: %w", err)
	}

	var baseConfig map[string]any
	if err := json.Unmarshal(sourceData, &baseConfig); err != nil {
		return "", fmt.Errorf("failed to parse source JSON: %w", err)
	}

	// Determine user file path
	userFile := target + ".local"
	if cfg != nil && cfg.UserFile != "" {
		userFile = expandPath(cfg.UserFile)
	}

	// Read user's local file (if exists)
	var userConfig map[string]any
	if userData, err := os.ReadFile(userFile); err == nil {
		if err := json.Unmarshal(userData, &userConfig); err != nil {
			return "", fmt.Errorf("failed to parse user JSON at %s: %w", userFile, err)
		}
	}

	// Merge configs: base + user overlay
	mergedConfig := mergeJSON(baseConfig, userConfig, cfg)

	// Check if target is identical
	if targetData, err := os.ReadFile(target); err == nil {
		var existingConfig map[string]any
		if err := json.Unmarshal(targetData, &existingConfig); err == nil {
			if jsonEqual(mergedConfig, existingConfig) {
				return "unchanged", nil
			}
		}
	}

	if s.dryRun {
		return "created", nil
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Remove existing symlink if present (merge mode writes a real file)
	if linkInfo, err := os.Lstat(target); err == nil {
		if linkInfo.Mode()&os.ModeSymlink != 0 {
			// It's a symlink - remove it so we can write a real file
			os.Remove(target)
		}
	}

	// Write merged config
	mergedData, err := json.MarshalIndent(mergedConfig, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal merged config: %w", err)
	}

	if err := os.WriteFile(target, mergedData, 0o644); err != nil {
		return "", fmt.Errorf("failed to write target: %w", err)
	}

	return "created", nil
}

// mergeJSON performs a deep merge of two JSON objects.
// User values override base values. PreserveUserKeys are always taken from user.
func mergeJSON(base, user map[string]any, cfg *componentconfig.MergeConfig) map[string]any {
	if base == nil {
		return user
	}
	if user == nil {
		return base
	}

	result := make(map[string]any)

	// Copy base values
	for k, v := range base {
		result[k] = v
	}

	// Determine merge strategy
	deep := cfg == nil || cfg.Strategy == "" || cfg.Strategy == "deep"

	// Merge user values
	for k, v := range user {
		if deep {
			// Deep merge for nested objects
			if baseMap, ok := result[k].(map[string]any); ok {
				if userMap, ok := v.(map[string]any); ok {
					result[k] = mergeJSON(baseMap, userMap, cfg)
					continue
				}
			}
		}
		result[k] = v
	}

	// Always preserve user keys if specified
	if cfg != nil {
		for _, key := range cfg.PreserveUserKeys {
			if v, ok := user[key]; ok {
				result[key] = v
			}
		}
	}

	return result
}

// jsonEqual checks if two JSON objects are equal.
func jsonEqual(a, b map[string]any) bool {
	aData, err1 := json.Marshal(a)
	bData, err2 := json.Marshal(b)
	if err1 != nil || err2 != nil {
		return false
	}
	return string(aData) == string(bData)
}

// Status checks the sync status of files.
func (s *Syncer) Status(files []componentconfig.SyncFileConfig) (*SyncResult, error) {
	result := &SyncResult{
		Synced:  []SyncedFile{},
		Skipped: []SyncedFile{},
		Errors:  []SyncError{},
	}

	for _, fc := range files {
		sourcePath := filepath.Join(s.dotfilesRoot, fc.Source)
		targetPath := expandPath(fc.Target)

		status := s.checkStatus(sourcePath, targetPath, fc.Mode)

		sf := SyncedFile{
			Source: fc.Source,
			Target: targetPath,
			Mode:   fc.Mode,
			Action: status,
		}

		if status == "synced" {
			result.Synced = append(result.Synced, sf)
		} else {
			result.Skipped = append(result.Skipped, sf)
		}
	}

	return result, nil
}

// checkStatus checks the sync status of a single file.
func (s *Syncer) checkStatus(source, target, mode string) string {
	switch mode {
	case "symlink":
		// Check if target is a symlink pointing to source
		linkInfo, err := os.Lstat(target)
		if err != nil {
			return "missing"
		}
		if linkInfo.Mode()&os.ModeSymlink == 0 {
			return "not_symlink"
		}
		dest, _ := os.Readlink(target)
		if dest == source {
			return "synced"
		}
		return "wrong_target"

	case "copy", "merge":
		// Check if target exists
		if _, err := os.Stat(target); err != nil {
			return "missing"
		}
		return "synced"

	default:
		return "unknown_mode"
	}
}
