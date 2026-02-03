package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// Settings represents a Claude Code settings file.
type Settings struct {
	Permissions    Permissions            `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	Hooks          map[string]interface{} `json:"hooks,omitempty" yaml:"hooks,omitempty"`
	StatusLine     *StatusLine            `json:"statusLine,omitempty" yaml:"statusLine,omitempty"`
	EnabledPlugins map[string]bool        `json:"enabledPlugins,omitempty" yaml:"enabledPlugins,omitempty"`
	Preferences    *Preferences           `json:"preferences,omitempty" yaml:"preferences,omitempty"`
}

// Permissions represents permission rules.
type Permissions struct {
	Allow []string `json:"allow" yaml:"allow"`
	Deny  []string `json:"deny" yaml:"deny"`
}

// StatusLine represents statusLine configuration.
type StatusLine struct {
	Type    string `json:"type" yaml:"type"`
	Command string `json:"command" yaml:"command"`
	Padding int    `json:"padding,omitempty" yaml:"padding,omitempty"`
}

// Preferences represents user preferences.
type Preferences struct {
	AutoApproveReadOperations bool   `json:"autoApproveReadOperations,omitempty" yaml:"autoApproveReadOperations,omitempty"`
	PreferredShell            string `json:"preferredShell,omitempty" yaml:"preferredShell,omitempty"`
}

// SettingsType represents which settings file to use.
type SettingsType string

const (
	SettingsGlobal SettingsType = "global"
	SettingsLocal  SettingsType = "local"
	SettingsConfig SettingsType = "config"
)

// PermissionsView is the view model for permissions display.
type PermissionsView struct {
	Allow []string `json:"allow" yaml:"allow"`
	Deny  []string `json:"deny" yaml:"deny"`
	File  string   `json:"file" yaml:"file"`
}

// ParseSettingsType parses a settings type string.
func ParseSettingsType(s string) (SettingsType, error) {
	switch s {
	case "", "global", "g":
		return SettingsGlobal, nil
	case "local", "l":
		return SettingsLocal, nil
	case "config", "c":
		return SettingsConfig, nil
	default:
		return "", fmt.Errorf("invalid settings type: %s (use global, local, or config)", s)
	}
}

// GetSettingsPath returns the path for the given settings type.
func (h *Helper) GetSettingsPath(st SettingsType) string {
	switch st {
	case SettingsGlobal:
		return h.paths.Settings
	case SettingsLocal:
		return h.paths.Local
	case SettingsConfig:
		return h.paths.Config
	default:
		return h.paths.Settings
	}
}

// GetSettings reads and returns settings from the specified file.
func (h *Helper) GetSettings(st SettingsType) (*Settings, error) {
	path := h.GetSettingsPath(st)
	if !h.FileExists(path) {
		return nil, fmt.Errorf("settings file not found: %s", path)
	}

	var settings Settings
	if err := h.ReadJSONFile(path, &settings); err != nil {
		return nil, fmt.Errorf("failed to read settings: %w", err)
	}

	return &settings, nil
}

// GetSettingsRaw reads settings as raw JSON for display.
func (h *Helper) GetSettingsRaw(st SettingsType) (map[string]interface{}, error) {
	path := h.GetSettingsPath(st)
	if !h.FileExists(path) {
		return nil, fmt.Errorf("settings file not found: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	return raw, nil
}

// EditSettings opens the settings file in the user's editor.
func (h *Helper) EditSettings(st SettingsType) error {
	path := h.GetSettingsPath(st)
	if !h.FileExists(path) {
		return fmt.Errorf("settings file not found: %s", path)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// GetPermissions reads permissions from settings.local.json.
func (h *Helper) GetPermissions() (*PermissionsView, error) {
	path := h.paths.Local
	if !h.FileExists(path) {
		return &PermissionsView{
			Allow: []string{},
			Deny:  []string{},
			File:  path,
		}, nil
	}

	var settings Settings
	if err := h.ReadJSONFile(path, &settings); err != nil {
		return nil, fmt.Errorf("failed to read settings: %w", err)
	}

	return &PermissionsView{
		Allow: settings.Permissions.Allow,
		Deny:  settings.Permissions.Deny,
		File:  path,
	}, nil
}

// AddPermission adds a permission rule.
func (h *Helper) AddPermission(rule string, permType string) error {
	path := h.paths.Local

	// Initialize file if it doesn't exist
	if !h.FileExists(path) {
		initial := Settings{
			Permissions: Permissions{
				Allow: []string{},
				Deny:  []string{},
			},
		}
		if err := h.WriteJSONFile(path, initial); err != nil {
			return fmt.Errorf("failed to create settings file: %w", err)
		}
	}

	var settings Settings
	if err := h.ReadJSONFile(path, &settings); err != nil {
		return fmt.Errorf("failed to read settings: %w", err)
	}

	// Add rule based on type
	if permType == "deny" {
		settings.Permissions.Deny = addUnique(settings.Permissions.Deny, rule)
	} else {
		settings.Permissions.Allow = addUnique(settings.Permissions.Allow, rule)
	}

	if err := h.WriteJSONFile(path, settings); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	return nil
}

// RemovePermission removes a permission rule.
func (h *Helper) RemovePermission(rule string, permType string) error {
	path := h.paths.Local

	if !h.FileExists(path) {
		return fmt.Errorf("settings file not found: %s", path)
	}

	var settings Settings
	if err := h.ReadJSONFile(path, &settings); err != nil {
		return fmt.Errorf("failed to read settings: %w", err)
	}

	// Remove rule based on type
	if permType == "deny" {
		settings.Permissions.Deny = removeItem(settings.Permissions.Deny, rule)
	} else {
		settings.Permissions.Allow = removeItem(settings.Permissions.Allow, rule)
	}

	if err := h.WriteJSONFile(path, settings); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	return nil
}

// addUnique adds an item to a slice if it doesn't already exist.
func addUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}

// removeItem removes an item from a slice.
func removeItem(slice []string, item string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}
