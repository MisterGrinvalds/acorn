// Package ollama provides Ollama local AI model management functionality.
package ollama

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Status represents Ollama installation status.
type Status struct {
	Installed      bool     `json:"installed" yaml:"installed"`
	Version        string   `json:"version,omitempty" yaml:"version,omitempty"`
	ServiceRunning bool     `json:"service_running" yaml:"service_running"`
	Models         []Model  `json:"models,omitempty" yaml:"models,omitempty"`
	HomeDir        string   `json:"home_dir" yaml:"home_dir"`
	StorageSize    string   `json:"storage_size,omitempty" yaml:"storage_size,omitempty"`
}

// Model represents an Ollama model.
type Model struct {
	Name       string `json:"name" yaml:"name"`
	Size       string `json:"size,omitempty" yaml:"size,omitempty"`
	Modified   string `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Helper provides Ollama helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Ollama Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetHomeDir returns the Ollama home directory.
func (h *Helper) GetHomeDir() string {
	if home := os.Getenv("OLLAMA_HOME"); home != "" {
		return home
	}
	userHome, _ := os.UserHomeDir()
	return filepath.Join(userHome, ".ollama")
}

// IsInstalled checks if Ollama is installed.
func (h *Helper) IsInstalled() bool {
	_, err := exec.LookPath("ollama")
	return err == nil
}

// GetVersion returns the Ollama version.
func (h *Helper) GetVersion() string {
	out, err := exec.Command("ollama", "--version").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// IsServiceRunning checks if Ollama service is running.
func (h *Helper) IsServiceRunning() bool {
	cmd := exec.Command("pgrep", "-f", "ollama serve")
	return cmd.Run() == nil
}

// GetStatus returns Ollama status.
func (h *Helper) GetStatus() *Status {
	status := &Status{
		HomeDir: h.GetHomeDir(),
	}

	if h.IsInstalled() {
		status.Installed = true
		status.Version = h.GetVersion()
	}

	status.ServiceRunning = h.IsServiceRunning()

	// Get models
	if status.Installed {
		models, _ := h.ListModels()
		status.Models = models
	}

	// Get storage size
	homeDir := status.HomeDir
	if _, err := os.Stat(homeDir); err == nil {
		cmd := exec.Command("du", "-sh", homeDir)
		if out, err := cmd.Output(); err == nil {
			parts := strings.Fields(string(out))
			if len(parts) > 0 {
				status.StorageSize = parts[0]
			}
		}
	}

	return status
}

// ListModels returns list of installed models.
func (h *Helper) ListModels() ([]Model, error) {
	if !h.IsInstalled() {
		return nil, fmt.Errorf("Ollama not installed")
	}

	cmd := exec.Command("ollama", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var models []Model
	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		if i == 0 || line == "" { // Skip header
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			model := Model{Name: fields[0]}
			if len(fields) >= 3 {
				model.Size = fields[2]
			}
			if len(fields) >= 4 {
				model.Modified = strings.Join(fields[3:], " ")
			}
			models = append(models, model)
		}
	}

	return models, nil
}

// Install installs Ollama.
func (h *Helper) Install() error {
	if h.IsInstalled() {
		return fmt.Errorf("Ollama already installed: %s", h.GetVersion())
	}

	if h.dryRun {
		fmt.Println("[dry-run] would install Ollama")
		return nil
	}

	fmt.Println("Installing Ollama...")

	switch runtime.GOOS {
	case "darwin":
		// Try brew first
		if _, err := exec.LookPath("brew"); err == nil {
			cmd := exec.Command("brew", "install", "ollama")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
		// Fall back to curl
		cmd := exec.Command("bash", "-c", "curl -fsSL https://ollama.ai/install.sh | sh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	case "linux":
		cmd := exec.Command("bash", "-c", "curl -fsSL https://ollama.ai/install.sh | sh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	default:
		return fmt.Errorf("unsupported OS: %s. Visit https://ollama.ai", runtime.GOOS)
	}
}

// Start starts the Ollama service.
func (h *Helper) Start() error {
	if !h.IsInstalled() {
		return fmt.Errorf("Ollama not installed")
	}

	if h.IsServiceRunning() {
		return fmt.Errorf("Ollama service already running")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would start Ollama service")
		return nil
	}

	fmt.Println("Starting Ollama service...")
	cmd := exec.Command("ollama", "serve")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return err
	}

	fmt.Println("Ollama service started in background")
	return nil
}

// Stop stops the Ollama service.
func (h *Helper) Stop() error {
	if !h.IsServiceRunning() {
		return fmt.Errorf("Ollama service not running")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would stop Ollama service")
		return nil
	}

	fmt.Println("Stopping Ollama service...")
	cmd := exec.Command("pkill", "-f", "ollama serve")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}

	fmt.Println("Ollama service stopped")
	return nil
}

// Pull downloads a model.
func (h *Helper) Pull(model string) error {
	if !h.IsInstalled() {
		return fmt.Errorf("Ollama not installed")
	}

	if model == "" {
		return fmt.Errorf("model name required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would pull model: %s\n", model)
		return nil
	}

	fmt.Printf("Pulling model: %s\n", model)
	cmd := exec.Command("ollama", "pull", model)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Remove removes a model.
func (h *Helper) Remove(model string) error {
	if !h.IsInstalled() {
		return fmt.Errorf("Ollama not installed")
	}

	if model == "" {
		return fmt.Errorf("model name required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would remove model: %s\n", model)
		return nil
	}

	fmt.Printf("Removing model: %s\n", model)
	cmd := exec.Command("ollama", "rm", model)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Chat sends a prompt to a model.
func (h *Helper) Chat(model, prompt string) error {
	if !h.IsInstalled() {
		return fmt.Errorf("Ollama not installed")
	}

	if model == "" || prompt == "" {
		return fmt.Errorf("model and prompt required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would chat with %s: %s\n", model, prompt)
		return nil
	}

	fmt.Printf("Asking %s: %s\n\n", model, prompt)

	cmd := exec.Command("ollama", "run", model)
	cmd.Stdin = strings.NewReader(prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Code generates code using a model.
func (h *Helper) Code(language, description string) error {
	if !h.IsInstalled() {
		return fmt.Errorf("Ollama not installed")
	}

	if language == "" || description == "" {
		return fmt.Errorf("language and description required")
	}

	model := "codellama"

	// Check if codellama is installed
	models, _ := h.ListModels()
	hasCodeLlama := false
	for _, m := range models {
		if strings.Contains(m.Name, "codellama") {
			hasCodeLlama = true
			model = m.Name
			break
		}
	}

	if !hasCodeLlama && !h.dryRun {
		fmt.Println("CodeLlama not found, pulling model...")
		if err := h.Pull("codellama"); err != nil {
			return fmt.Errorf("failed to pull codellama: %w", err)
		}
	}

	prompt := fmt.Sprintf("Write a %s %s. Only return the code:", language, description)

	if h.dryRun {
		fmt.Printf("[dry-run] would generate %s code: %s\n", language, description)
		return nil
	}

	fmt.Printf("Generating %s code: %s\n\n", language, description)

	cmd := exec.Command("ollama", "run", model)
	cmd.Stdin = strings.NewReader(prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetExamples returns usage examples.
func (h *Helper) GetExamples() string {
	return `Ollama Usage Examples
=====================

1. General Chat:
   acorn ollama chat llama3.2 "Explain machine learning"

2. Code Generation:
   acorn ollama code python "function to calculate fibonacci"

3. Interactive Session:
   ollama run llama3.2

4. Model Management:
   acorn ollama models         # List installed models
   acorn ollama pull mistral   # Install new model
   acorn ollama rm phi3        # Remove model

5. Service Management:
   acorn ollama start          # Start service
   acorn ollama stop           # Stop service
   acorn ollama status         # Check status

Popular Models:
   llama3.2    - Meta's latest Llama model
   codellama   - Code generation model
   mistral     - Mistral AI model
   phi3        - Microsoft's Phi-3 model
   gemma2      - Google's Gemma 2 model
`
}
