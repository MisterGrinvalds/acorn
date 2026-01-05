// Package huggingface provides Hugging Face model management functionality.
package huggingface

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents Hugging Face installation status.
type Status struct {
	TransformersInstalled bool   `json:"transformers_installed" yaml:"transformers_installed"`
	TransformersVersion   string `json:"transformers_version,omitempty" yaml:"transformers_version,omitempty"`
	PyTorchInstalled      bool   `json:"pytorch_installed" yaml:"pytorch_installed"`
	PyTorchVersion        string `json:"pytorch_version,omitempty" yaml:"pytorch_version,omitempty"`
	CacheDir              string `json:"cache_dir" yaml:"cache_dir"`
	CacheSize             string `json:"cache_size,omitempty" yaml:"cache_size,omitempty"`
	VirtualEnv            string `json:"virtual_env,omitempty" yaml:"virtual_env,omitempty"`
}

// Model represents a popular model reference.
type Model struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Category    string `json:"category" yaml:"category"`
}

// Pipeline represents a pipeline task.
type Pipeline struct {
	Task        string `json:"task" yaml:"task"`
	Description string `json:"description" yaml:"description"`
}

// Helper provides Hugging Face helper operations.
type Helper struct {
	verbose bool
}

// NewHelper creates a new Hugging Face Helper.
func NewHelper(verbose bool) *Helper {
	return &Helper{
		verbose: verbose,
	}
}

// GetCacheDir returns the Hugging Face cache directory.
func (h *Helper) GetCacheDir() string {
	if hfHome := os.Getenv("HF_HOME"); hfHome != "" {
		return hfHome
	}
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home, _ := os.UserHomeDir()
		cacheHome = filepath.Join(home, ".cache")
	}
	return filepath.Join(cacheHome, "huggingface")
}

// GetStatus returns Hugging Face installation status.
func (h *Helper) GetStatus() *Status {
	status := &Status{
		CacheDir: h.GetCacheDir(),
	}

	// Check transformers
	cmd := exec.Command("python3", "-c", "import transformers; print(transformers.__version__)")
	if out, err := cmd.Output(); err == nil {
		status.TransformersInstalled = true
		status.TransformersVersion = strings.TrimSpace(string(out))
	}

	// Check PyTorch
	cmd = exec.Command("python3", "-c", "import torch; print(torch.__version__)")
	if out, err := cmd.Output(); err == nil {
		status.PyTorchInstalled = true
		status.PyTorchVersion = strings.TrimSpace(string(out))
	}

	// Check cache size
	if _, err := os.Stat(status.CacheDir); err == nil {
		cmd = exec.Command("du", "-sh", status.CacheDir)
		if out, err := cmd.Output(); err == nil {
			parts := strings.Fields(string(out))
			if len(parts) > 0 {
				status.CacheSize = parts[0]
			}
		}
	}

	// Check virtual environment
	if venv := os.Getenv("VIRTUAL_ENV"); venv != "" {
		status.VirtualEnv = filepath.Base(venv)
	}

	return status
}

// GetModels returns list of popular models.
func (h *Helper) GetModels() []Model {
	return []Model{
		// Text Generation
		{Name: "microsoft/DialoGPT-small", Description: "Conversational AI (117M)", Category: "Text Generation"},
		{Name: "gpt2", Description: "GPT-2 (124M)", Category: "Text Generation"},
		{Name: "distilgpt2", Description: "Distilled GPT-2 (82M)", Category: "Text Generation"},
		// Language Understanding
		{Name: "distilbert-base-uncased", Description: "Efficient BERT (66M)", Category: "Language Understanding"},
		{Name: "bert-base-uncased", Description: "BERT (110M)", Category: "Language Understanding"},
		// Specialized
		{Name: "microsoft/codebert-base", Description: "Code understanding", Category: "Specialized"},
		{Name: "facebook/bart-base", Description: "Text summarization", Category: "Specialized"},
	}
}

// GetPipelines returns list of available pipelines.
func (h *Helper) GetPipelines() []Pipeline {
	return []Pipeline{
		{Task: "text-generation", Description: "Generate text continuations"},
		{Task: "summarization", Description: "Summarize long text"},
		{Task: "sentiment-analysis", Description: "Analyze text sentiment"},
		{Task: "question-answering", Description: "Answer questions about text"},
		{Task: "translation", Description: "Translate between languages"},
		{Task: "fill-mask", Description: "Fill in masked words"},
		{Task: "text-classification", Description: "Classify text categories"},
	}
}

// ClearCache removes the model cache.
func (h *Helper) ClearCache(force bool) error {
	cacheDir := h.GetCacheDir()

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return fmt.Errorf("no cache directory found at %s", cacheDir)
	}

	if !force {
		return fmt.Errorf("use --force to actually clear the cache")
	}

	if err := os.RemoveAll(cacheDir); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	return nil
}

// GetCacheInfo returns cache directory info.
func (h *Helper) GetCacheInfo() (string, string, error) {
	cacheDir := h.GetCacheDir()

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return cacheDir, "", nil
	}

	cmd := exec.Command("du", "-sh", cacheDir)
	out, err := cmd.Output()
	if err != nil {
		return cacheDir, "unknown", nil
	}

	parts := strings.Fields(string(out))
	if len(parts) > 0 {
		return cacheDir, parts[0], nil
	}

	return cacheDir, "unknown", nil
}
