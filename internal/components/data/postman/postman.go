// Package postman provides Postman API development environment helper functionality.
package postman

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Status represents Postman installation status.
type Status struct {
	Installed      bool   `json:"installed" yaml:"installed"`
	AppPath        string `json:"app_path,omitempty" yaml:"app_path,omitempty"`
	Version        string `json:"version,omitempty" yaml:"version,omitempty"`
	NewmanInstalled bool   `json:"newman_installed" yaml:"newman_installed"`
	NewmanVersion  string `json:"newman_version,omitempty" yaml:"newman_version,omitempty"`
	CollectionsDir string `json:"collections_dir,omitempty" yaml:"collections_dir,omitempty"`
}

// Collection represents a Postman collection.
type Collection struct {
	Name     string `json:"name" yaml:"name"`
	Path     string `json:"path" yaml:"path"`
	ID       string `json:"id,omitempty" yaml:"id,omitempty"`
	Requests int    `json:"requests,omitempty" yaml:"requests,omitempty"`
}

// Environment represents a Postman environment.
type Environment struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
	ID   string `json:"id,omitempty" yaml:"id,omitempty"`
}

// NewmanResult represents the result of a Newman run.
type NewmanResult struct {
	Collection  string `json:"collection" yaml:"collection"`
	Environment string `json:"environment,omitempty" yaml:"environment,omitempty"`
	Passed      int    `json:"passed" yaml:"passed"`
	Failed      int    `json:"failed" yaml:"failed"`
	Duration    string `json:"duration,omitempty" yaml:"duration,omitempty"`
}

// Helper provides Postman helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Postman Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetStatus returns Postman status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	// Find Postman installation
	appPath := h.findPostman()
	if appPath != "" {
		status.Installed = true
		status.AppPath = appPath
		status.Version = h.getVersion(appPath)
	}

	// Check Newman
	if version := h.getNewmanVersion(); version != "" {
		status.NewmanInstalled = true
		status.NewmanVersion = version
	}

	// Get collections directory
	status.CollectionsDir = h.getCollectionsDir()

	return status
}

// findPostman finds the Postman application.
func (h *Helper) findPostman() string {
	if runtime.GOOS == "darwin" {
		paths := []string{
			"/Applications/Postman.app",
			filepath.Join(os.Getenv("HOME"), "Applications/Postman.app"),
		}

		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	return ""
}

// getVersion extracts Postman version from the app.
func (h *Helper) getVersion(appPath string) string {
	if runtime.GOOS == "darwin" {
		plistPath := filepath.Join(appPath, "Contents/Info.plist")
		cmd := exec.Command("defaults", "read", plistPath, "CFBundleShortVersionString")
		out, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
	}
	return ""
}

// getNewmanVersion returns the Newman version if installed.
func (h *Helper) getNewmanVersion() string {
	cmd := exec.Command("newman", "--version")
	out, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	// Try npx newman
	cmd = exec.Command("npx", "newman", "--version")
	out, err = cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	return ""
}

// getCollectionsDir returns the default collections directory.
func (h *Helper) getCollectionsDir() string {
	home, _ := os.UserHomeDir()

	// Check XDG config
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		xdgConfig = filepath.Join(home, ".config")
	}

	collectionsDir := filepath.Join(xdgConfig, "postman", "collections")
	if _, err := os.Stat(collectionsDir); err == nil {
		return collectionsDir
	}

	// Check common locations
	paths := []string{
		filepath.Join(home, "Postman/collections"),
		filepath.Join(home, ".postman/collections"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return filepath.Join(xdgConfig, "postman", "collections")
}

// Launch starts Postman.
func (h *Helper) Launch() error {
	if h.dryRun {
		fmt.Println("[dry-run] would launch Postman")
		return nil
	}

	if runtime.GOOS == "darwin" {
		appPath := h.findPostman()
		if appPath == "" {
			return fmt.Errorf("Postman not found")
		}

		cmd := exec.Command("open", "-a", appPath)
		return cmd.Start()
	}

	return fmt.Errorf("Postman not found")
}

// InstallNewman installs Newman CLI globally.
func (h *Helper) InstallNewman() error {
	if h.dryRun {
		fmt.Println("[dry-run] would install Newman via npm")
		return nil
	}

	fmt.Println("Installing Newman...")
	cmd := exec.Command("npm", "install", "-g", "newman")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCollection runs a Postman collection using Newman.
func (h *Helper) RunCollection(collectionPath, envPath string, reporters []string) error {
	if collectionPath == "" {
		return fmt.Errorf("collection path is required")
	}

	args := []string{"run", collectionPath}

	if envPath != "" {
		args = append(args, "-e", envPath)
	}

	if len(reporters) > 0 {
		args = append(args, "-r", strings.Join(reporters, ","))
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: newman %s\n", strings.Join(args, " "))
		return nil
	}

	// Try newman directly
	cmd := exec.Command("newman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		// Try npx newman
		npxArgs := append([]string{"newman"}, args...)
		cmd = exec.Command("npx", npxArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return nil
}

// ListCollections lists available collections.
func (h *Helper) ListCollections() ([]Collection, error) {
	collectionsDir := h.getCollectionsDir()
	if _, err := os.Stat(collectionsDir); os.IsNotExist(err) {
		return []Collection{}, nil
	}

	var collections []Collection

	entries, err := os.ReadDir(collectionsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read collections directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}

		path := filepath.Join(collectionsDir, name)
		collection := Collection{
			Name: strings.TrimSuffix(name, ".json"),
			Path: path,
		}

		// Try to parse collection info
		if info := h.parseCollectionInfo(path); info != nil {
			if info.Name != "" {
				collection.Name = info.Name
			}
			collection.ID = info.ID
			collection.Requests = info.Requests
		}

		collections = append(collections, collection)
	}

	return collections, nil
}

// collectionInfo holds parsed collection metadata.
type collectionInfo struct {
	Name     string
	ID       string
	Requests int
}

// parseCollectionInfo parses collection metadata from a JSON file.
func (h *Helper) parseCollectionInfo(path string) *collectionInfo {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var parsed struct {
		Info struct {
			Name string `json:"name"`
			ID   string `json:"_postman_id"`
		} `json:"info"`
		Item []any `json:"item"`
	}

	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil
	}

	return &collectionInfo{
		Name:     parsed.Info.Name,
		ID:       parsed.Info.ID,
		Requests: len(parsed.Item),
	}
}

// ListEnvironments lists available environments.
func (h *Helper) ListEnvironments() ([]Environment, error) {
	collectionsDir := h.getCollectionsDir()
	envDir := filepath.Join(filepath.Dir(collectionsDir), "environments")
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		return []Environment{}, nil
	}

	var environments []Environment

	entries, err := os.ReadDir(envDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read environments directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}

		path := filepath.Join(envDir, name)
		env := Environment{
			Name: strings.TrimSuffix(name, ".json"),
			Path: path,
		}

		// Try to parse environment info
		if info := h.parseEnvironmentInfo(path); info != nil {
			if info.Name != "" {
				env.Name = info.Name
			}
			env.ID = info.ID
		}

		environments = append(environments, env)
	}

	return environments, nil
}

// environmentInfo holds parsed environment metadata.
type environmentInfo struct {
	Name string
	ID   string
}

// parseEnvironmentInfo parses environment metadata from a JSON file.
func (h *Helper) parseEnvironmentInfo(path string) *environmentInfo {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var parsed struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}

	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil
	}

	return &environmentInfo{
		Name: parsed.Name,
		ID:   parsed.ID,
	}
}

// ExportCollection exports a collection to a file.
func (h *Helper) ExportCollection(name, outputPath string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would export collection %s to %s\n", name, outputPath)
		return nil
	}

	// Postman export is typically done via the GUI
	// This is a placeholder for potential API integration
	return fmt.Errorf("collection export requires Postman API key - use Postman GUI or API directly")
}

// ImportCollection imports a collection from a file.
func (h *Helper) ImportCollection(inputPath string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would import collection from %s\n", inputPath)
		return nil
	}

	// Validate collection file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("collection file not found: %s", inputPath)
	}

	// Copy to collections directory
	collectionsDir := h.getCollectionsDir()
	if err := os.MkdirAll(collectionsDir, 0o755); err != nil {
		return fmt.Errorf("failed to create collections directory: %w", err)
	}

	destPath := filepath.Join(collectionsDir, filepath.Base(inputPath))
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read collection: %w", err)
	}

	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write collection: %w", err)
	}

	fmt.Printf("Imported collection to: %s\n", destPath)
	return nil
}
