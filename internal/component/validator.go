package component

import (
	"fmt"
	"os"
)

// ValidationResult represents the result of component validation.
type ValidationResult struct {
	Component *Component
	Valid     bool
	Errors    []string
}

// Validate performs validation on a component's configuration.
func Validate(comp *Component) *ValidationResult {
	vr := &ValidationResult{
		Component: comp,
		Valid:     true,
		Errors:    []string{},
	}

	// Required fields
	if comp.Name == "" {
		vr.addError("missing required field: name")
	}
	if comp.Version == "" {
		vr.addError("missing required field: version")
	}
	if comp.Description == "" {
		vr.addError("missing required field: description")
	}
	if comp.Category == "" {
		vr.addError("missing required field: category")
	}

	// Validate config files
	for i, cfg := range comp.Config.Files {
		// Check source exists
		sourcePath := comp.Path + "/" + cfg.Source
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			vr.addError(fmt.Sprintf("config.files[%d].source does not exist: %s", i, cfg.Source))
		}

		// Check target is specified
		if cfg.Target == "" {
			vr.addError(fmt.Sprintf("config.files[%d].target is empty", i))
		}

		// Validate method
		if cfg.Method != "" && cfg.Method != "symlink" && cfg.Method != "copy" {
			vr.addError(fmt.Sprintf("config.files[%d].method must be 'symlink' or 'copy', got: %s", i, cfg.Method))
		}
	}

	// Validate platforms
	validPlatforms := map[string]bool{"darwin": true, "linux": true, "windows": true}
	for _, platform := range comp.Platforms {
		if !validPlatforms[platform] {
			vr.addError(fmt.Sprintf("invalid platform: %s (must be darwin, linux, or windows)", platform))
		}
	}

	// Validate shells
	validShells := map[string]bool{"bash": true, "zsh": true, "fish": true}
	for _, shell := range comp.Shells {
		if !validShells[shell] {
			vr.addError(fmt.Sprintf("invalid shell: %s (must be bash, zsh, or fish)", shell))
		}
	}

	// Validate shell files syntax
	for _, shellFile := range comp.ShellFiles() {
		if err := checkShellSyntax(shellFile); err != nil {
			vr.addError(fmt.Sprintf("syntax error in %s: %v", shellFile, err))
		}
	}

	return vr
}

// addError adds a validation error.
func (vr *ValidationResult) addError(err string) {
	vr.Errors = append(vr.Errors, err)
	vr.Valid = false
}
