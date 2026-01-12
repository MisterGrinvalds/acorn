---
description: Explain Ollama concepts, models, and local AI setup
argument-hint: [topic]
allowed-tools: Read, Glob, Grep
---

## Task

Explain the requested topic about Ollama and local AI models. If no specific topic provided, give an overview.

## Topics

### Core Concepts
- **ollama** - What is Ollama and how it works
- **models** - Understanding LLM models and sizes
- **quantization** - Model compression (Q4, Q5, Q8, etc.)
- **context-window** - Token limits and memory usage

### Models
- **llama** - Meta's Llama family (3, 3.1, 3.2)
- **codellama** - Code-specialized Llama variant
- **mistral** - Mistral AI models
- **mixtral** - Mixture of experts architecture
- **phi** - Microsoft's small models
- **gemma** - Google's open models

### Operations
- **installation** - Installing Ollama on different platforms
- **service** - Running Ollama as a service
- **api** - REST API for integrations
- **gpu** - GPU acceleration (CUDA, Metal)

### Usage
- **prompting** - Effective prompt techniques
- **chat** - Interactive conversations
- **code-gen** - Code generation workflows
- **automation** - Scripting with Ollama

## Context

@components/ollama/component.yaml
@components/ollama/functions.sh
@components/ollama/aliases.sh

## Response Format

When explaining a topic:

1. **Definition** - What it is in simple terms
2. **How it works** - Technical details
3. **Examples** - Practical usage
4. **Dotfiles integration** - Available functions/aliases
5. **Recommendations** - Best practices

## Quick Reference

### Dotfiles Functions
- `ollama_install` - Install Ollama
- `ollama_start` / `ollama_stop` - Service management
- `ollama_status` - Check installation and models
- `ollama_models` - List installed models
- `ollama_pull <model>` - Download model
- `ollama_run <model>` - Interactive session
- `ollama_chat <model> <prompt>` - Quick query
- `ollama_code <lang> <desc>` - Code generation
- `ollama_remove <model>` - Delete model
- `ollama_examples` - Show usage examples

### Model Naming Convention
```
model:tag
  │    └── Version/variant (optional)
  └── Model name

Examples:
  llama3.2         # Latest default
  llama3.2:3b      # 3 billion parameter version
  llama3.2:70b     # 70 billion parameter version
  codellama:7b-q4  # 7B with Q4 quantization
```

### Storage Location
- macOS: `~/.ollama/models`
- Linux: `~/.ollama/models`
- Custom: Set `OLLAMA_MODELS` environment variable
