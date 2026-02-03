// Package golang provides Go development helper functionality.
package golang

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Project represents a Go project.
type Project struct {
	Name   string `json:"name" yaml:"name"`
	Path   string `json:"path" yaml:"path"`
	Module string `json:"module" yaml:"module"`
}

// TestResult contains test run results.
type TestResult struct {
	Package  string `json:"package" yaml:"package"`
	Passed   bool   `json:"passed" yaml:"passed"`
	Output   string `json:"output,omitempty" yaml:"output,omitempty"`
	Coverage string `json:"coverage,omitempty" yaml:"coverage,omitempty"`
}

// BuildTarget represents a build target platform.
type BuildTarget struct {
	OS     string `json:"os" yaml:"os"`
	Arch   string `json:"arch" yaml:"arch"`
	Output string `json:"output" yaml:"output"`
}

// DefaultBuildTargets returns the standard cross-compilation targets.
func DefaultBuildTargets(name string) []BuildTarget {
	return []BuildTarget{
		{OS: "linux", Arch: "amd64", Output: fmt.Sprintf("dist/%s-linux-amd64", name)},
		{OS: "linux", Arch: "arm64", Output: fmt.Sprintf("dist/%s-linux-arm64", name)},
		{OS: "darwin", Arch: "amd64", Output: fmt.Sprintf("dist/%s-darwin-amd64", name)},
		{OS: "darwin", Arch: "arm64", Output: fmt.Sprintf("dist/%s-darwin-arm64", name)},
		{OS: "windows", Arch: "amd64", Output: fmt.Sprintf("dist/%s-windows-amd64.exe", name)},
	}
}

// Helper provides Go development helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// InitProject initializes a new Go project.
func (h *Helper) InitProject(name string) (*Project, error) {
	if name == "" {
		return nil, fmt.Errorf("module name is required")
	}

	// Create directory
	if err := os.MkdirAll(name, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	projectPath, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}

	// Initialize go module
	if err := h.runInDir(projectPath, "go", "mod", "init", name); err != nil {
		return nil, fmt.Errorf("go mod init failed: %w", err)
	}

	// Create main.go
	mainGo := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`
	mainPath := filepath.Join(projectPath, "main.go")
	if err := os.WriteFile(mainPath, []byte(mainGo), 0o644); err != nil {
		return nil, fmt.Errorf("failed to create main.go: %w", err)
	}

	return &Project{
		Name:   name,
		Path:   projectPath,
		Module: name,
	}, nil
}

// InitCobraProject initializes a new Cobra CLI project.
func (h *Helper) InitCobraProject(name string) (*Project, error) {
	if name == "" {
		return nil, fmt.Errorf("app name is required")
	}

	// Check if cobra-cli is installed
	if _, err := exec.LookPath("cobra-cli"); err != nil {
		fmt.Println("cobra-cli not installed. Installing...")
		if err := h.run("go", "install", "github.com/spf13/cobra-cli@latest"); err != nil {
			return nil, fmt.Errorf("failed to install cobra-cli: %w", err)
		}
	}

	// Create directory and initialize
	if err := os.MkdirAll(name, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	projectPath, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}

	// Initialize cobra project
	if err := h.runInDir(projectPath, "cobra-cli", "init", "."); err != nil {
		return nil, fmt.Errorf("cobra-cli init failed: %w", err)
	}

	// Run go mod tidy
	if err := h.runInDir(projectPath, "go", "mod", "tidy"); err != nil {
		return nil, fmt.Errorf("go mod tidy failed: %w", err)
	}

	return &Project{
		Name:   name,
		Path:   projectPath,
		Module: name,
	}, nil
}

// AddCobraCommand adds a command to a Cobra project.
func (h *Helper) AddCobraCommand(cmdName string) error {
	if cmdName == "" {
		return fmt.Errorf("command name is required")
	}

	// Check if cobra-cli is installed
	if _, err := exec.LookPath("cobra-cli"); err != nil {
		return fmt.Errorf("cobra-cli not installed, run: go install github.com/spf13/cobra-cli@latest")
	}

	return h.run("cobra-cli", "add", cmdName)
}

// RunTests runs Go tests with optional pattern filter.
func (h *Helper) RunTests(pattern string) error {
	args := []string{"test", "./..."}
	if pattern != "" {
		args = append(args, "-run", pattern)
	}
	return h.run("go", args...)
}

// RunTestsWithCoverage runs tests with coverage report.
func (h *Helper) RunTestsWithCoverage() error {
	// Run tests with coverage
	if err := h.run("go", "test", "./...", "-coverprofile=coverage.out"); err != nil {
		return err
	}

	// Generate HTML report
	if err := h.run("go", "tool", "cover", "-html=coverage.out", "-o", "coverage.html"); err != nil {
		return err
	}

	// Show summary
	return h.run("go", "tool", "cover", "-func=coverage.out")
}

// RunBenchmarks runs Go benchmarks with optional pattern.
func (h *Helper) RunBenchmarks(pattern string) error {
	benchPattern := "."
	if pattern != "" {
		benchPattern = pattern
	}
	return h.run("go", "test", "-bench="+benchPattern, "./...")
}

// BuildAll builds for multiple platforms.
func (h *Helper) BuildAll(name string) error {
	if name == "" {
		name = "app"
	}

	// Create dist directory
	if err := os.MkdirAll("dist", 0o755); err != nil {
		return fmt.Errorf("failed to create dist directory: %w", err)
	}

	targets := DefaultBuildTargets(name)
	for _, target := range targets {
		fmt.Printf("Building %s/%s -> %s\n", target.OS, target.Arch, target.Output)
		if err := h.buildFor(target); err != nil {
			return fmt.Errorf("build failed for %s/%s: %w", target.OS, target.Arch, err)
		}
	}

	return nil
}

// buildFor builds for a specific target.
func (h *Helper) buildFor(target BuildTarget) error {
	if h.dryRun {
		fmt.Printf("[dry-run] GOOS=%s GOARCH=%s go build -o %s .\n",
			target.OS, target.Arch, target.Output)
		return nil
	}

	cmd := exec.Command("go", "build", "-o", target.Output, ".")
	cmd.Env = append(os.Environ(),
		"GOOS="+target.OS,
		"GOARCH="+target.Arch,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Clean removes build artifacts.
func (h *Helper) Clean() error {
	if err := h.run("go", "clean"); err != nil {
		return err
	}

	// Remove additional artifacts
	artifacts := []string{"dist", "coverage.out", "coverage.html"}
	for _, artifact := range artifacts {
		if _, err := os.Stat(artifact); err == nil {
			if h.dryRun {
				fmt.Printf("[dry-run] would remove: %s\n", artifact)
			} else {
				if err := os.RemoveAll(artifact); err != nil {
					return fmt.Errorf("failed to remove %s: %w", artifact, err)
				}
				if h.verbose {
					fmt.Printf("Removed: %s\n", artifact)
				}
			}
		}
	}

	return nil
}

// GetGoVersion returns the installed Go version.
func (h *Helper) GetGoVersion() (string, error) {
	cmd := exec.Command("go", "version")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// GetGoEnv returns Go environment information.
func (h *Helper) GetGoEnv() map[string]string {
	env := make(map[string]string)

	// Get key Go environment variables
	keys := []string{"GOPATH", "GOROOT", "GOBIN", "GOPROXY", "GOMODCACHE"}
	for _, key := range keys {
		cmd := exec.Command("go", "env", key)
		if out, err := cmd.Output(); err == nil {
			env[key] = strings.TrimSpace(string(out))
		}
	}

	// Add runtime info
	env["GOOS"] = runtime.GOOS
	env["GOARCH"] = runtime.GOARCH

	return env
}

// run executes a command.
func (h *Helper) run(name string, args ...string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: %s %s\n", name, strings.Join(args, " "))
		return nil
	}

	if h.verbose {
		fmt.Printf("Running: %s %s\n", name, strings.Join(args, " "))
	}

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// runInDir executes a command in a specific directory.
func (h *Helper) runInDir(dir, name string, args ...string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run in %s: %s %s\n", dir, name, strings.Join(args, " "))
		return nil
	}

	if h.verbose {
		fmt.Printf("Running in %s: %s %s\n", dir, name, strings.Join(args, " "))
	}

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
