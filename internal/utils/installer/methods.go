package installer

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/mistergrinvalds/acorn/internal/utils/config"
)

// MethodExecutor executes installation for a specific method type.
type MethodExecutor interface {
	// Type returns the method type identifier.
	Type() string

	// Available returns true if this method can be used on the current system.
	Available() bool

	// Execute installs a tool using this method.
	Execute(ctx context.Context, tool PlannedTool, stdout, stderr io.Writer) error
}

// GetExecutor returns an executor for the given install method type.
func GetExecutor(methodType string) (MethodExecutor, error) {
	switch methodType {
	case InstallTypeBrew:
		return &BrewExecutor{}, nil
	case InstallTypeApt:
		return &AptExecutor{}, nil
	case InstallTypeNpm:
		return &NpmExecutor{}, nil
	case InstallTypeGo:
		return &GoExecutor{}, nil
	case InstallTypeCurl:
		return &CurlExecutor{}, nil
	default:
		return nil, fmt.Errorf("unknown install method: %s", methodType)
	}
}

// BrewExecutor handles Homebrew installations.
type BrewExecutor struct{}

func (e *BrewExecutor) Type() string { return InstallTypeBrew }

func (e *BrewExecutor) Available() bool {
	return commandExists("brew")
}

func (e *BrewExecutor) Execute(ctx context.Context, tool PlannedTool, stdout, stderr io.Writer) error {
	pkg := tool.Method.Package
	if pkg == "" {
		pkg = tool.Name
	}

	args := []string{"install", pkg}
	args = append(args, tool.Method.Args...)

	return runCommand(ctx, "brew", args, stdout, stderr)
}

// AptExecutor handles APT installations (Debian/Ubuntu).
type AptExecutor struct{}

func (e *AptExecutor) Type() string { return InstallTypeApt }

func (e *AptExecutor) Available() bool {
	return commandExists("apt")
}

func (e *AptExecutor) Execute(ctx context.Context, tool PlannedTool, stdout, stderr io.Writer) error {
	pkg := tool.Method.Package
	if pkg == "" {
		pkg = tool.Name
	}

	// apt install requires sudo
	args := []string{"apt", "install", "-y", pkg}
	args = append(args, tool.Method.Args...)

	return runCommand(ctx, "sudo", args, stdout, stderr)
}

// NpmExecutor handles npm installations.
type NpmExecutor struct{}

func (e *NpmExecutor) Type() string { return InstallTypeNpm }

func (e *NpmExecutor) Available() bool {
	return commandExists("npm")
}

func (e *NpmExecutor) Execute(ctx context.Context, tool PlannedTool, stdout, stderr io.Writer) error {
	pkg := tool.Method.Package
	if pkg == "" {
		pkg = tool.Name
	}

	args := []string{"install"}
	if tool.Method.Global {
		args = append(args, "-g")
	}
	args = append(args, pkg)
	args = append(args, tool.Method.Args...)

	return runCommand(ctx, "npm", args, stdout, stderr)
}

// GoExecutor handles go install.
type GoExecutor struct{}

func (e *GoExecutor) Type() string { return InstallTypeGo }

func (e *GoExecutor) Available() bool {
	return commandExists("go")
}

func (e *GoExecutor) Execute(ctx context.Context, tool PlannedTool, stdout, stderr io.Writer) error {
	pkg := tool.Method.Package
	if pkg == "" {
		return fmt.Errorf("go install requires package path")
	}

	args := []string{"install", pkg}
	args = append(args, tool.Method.Args...)

	return runCommand(ctx, "go", args, stdout, stderr)
}

// CurlExecutor handles curl-based script installations.
type CurlExecutor struct{}

func (e *CurlExecutor) Type() string { return InstallTypeCurl }

func (e *CurlExecutor) Available() bool {
	return commandExists("curl")
}

func (e *CurlExecutor) Execute(ctx context.Context, tool PlannedTool, stdout, stderr io.Writer) error {
	if tool.Method.URL == "" {
		return fmt.Errorf("curl install requires URL")
	}

	// curl -fsSL <url> | sh
	curlCmd := exec.CommandContext(ctx, "curl", "-fsSL", tool.Method.URL)
	shCmd := exec.CommandContext(ctx, "sh")

	// Pipe curl output to sh
	pipe, err := curlCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
	}

	shCmd.Stdin = pipe
	shCmd.Stdout = stdout
	shCmd.Stderr = stderr

	if err := curlCmd.Start(); err != nil {
		return fmt.Errorf("failed to start curl: %w", err)
	}

	if err := shCmd.Start(); err != nil {
		curlCmd.Process.Kill()
		return fmt.Errorf("failed to start sh: %w", err)
	}

	if err := curlCmd.Wait(); err != nil {
		return fmt.Errorf("curl failed: %w", err)
	}

	if err := shCmd.Wait(); err != nil {
		return fmt.Errorf("install script failed: %w", err)
	}

	return nil
}

// runCommand executes a command with the given arguments.
func runCommand(ctx context.Context, name string, args []string, stdout, stderr io.Writer) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = os.Stdin // Allow interactive prompts

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", name, err)
	}

	return nil
}

// SelectMethod selects the appropriate install method for a platform.
func SelectMethod(methods map[string]config.InstallMethod, platform *Platform) (config.InstallMethod, bool) {
	// Try platform-specific keys in order of specificity
	for _, key := range platform.GetMethodKeys() {
		if method, ok := methods[key]; ok {
			return method, true
		}
	}

	return config.InstallMethod{}, false
}
