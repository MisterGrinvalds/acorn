---
description: Help pull and manage Ollama models
argument-hint: [model-name]
allowed-tools: Read, Bash
---

## Task

Help the user pull (download) and manage Ollama models. Recommend appropriate models based on their use case and hardware.

## Quick Pull Commands

```bash
# Using dotfiles function
ollama_pull llama3.2

# Direct ollama command
ollama pull llama3.2

# Pull specific size
ollama pull llama3.2:3b
ollama pull llama3.2:70b
```

## Popular Models by Category

### General Purpose
```bash
# Llama 3.2 (Meta) - Latest, recommended
ollama_pull llama3.2        # Default size
ollama_pull llama3.2:1b     # Smallest, very fast
ollama_pull llama3.2:3b     # Good balance
ollama_pull llama3.2:11b    # More capable
ollama_pull llama3.2:90b    # Most capable

# Mistral (Mistral AI) - Fast and capable
ollama_pull mistral         # 7B, excellent quality

# Mixtral (Mistral AI) - Mixture of experts
ollama_pull mixtral         # 8x7B, very capable
```

### Code Generation
```bash
# CodeLlama (Meta) - Code-specialized
ollama_pull codellama       # 7B default
ollama_pull codellama:7b
ollama_pull codellama:13b
ollama_pull codellama:34b

# Qwen 2.5 Coder (Alibaba) - Latest code model
ollama_pull qwen2.5-coder:0.5b   # Tiny
ollama_pull qwen2.5-coder:1.5b
ollama_pull qwen2.5-coder:7b
ollama_pull qwen2.5-coder:32b

# DeepSeek Coder
ollama_pull deepseek-coder:1.3b
ollama_pull deepseek-coder:6.7b
ollama_pull deepseek-coder:33b
```

### Small & Fast
```bash
# Phi (Microsoft) - Small but capable
ollama_pull phi3            # 3.8B
ollama_pull phi3:mini       # Smaller variant

# Gemma (Google)
ollama_pull gemma2:2b       # Tiny
ollama_pull gemma2:9b       # Medium
ollama_pull gemma2:27b      # Large
```

### Specialized
```bash
# Vision models (image understanding)
ollama_pull llava           # Llama + vision
ollama_pull bakllava        # Vision model

# Embedding models (for RAG/search)
ollama_pull nomic-embed-text
ollama_pull mxbai-embed-large
```

## Model Size Guide

| Model | Size on Disk | RAM Needed | Speed |
|-------|-------------|------------|-------|
| 1-3B | 1-2 GB | 4-8 GB | Very Fast |
| 7B | 4-5 GB | 8-16 GB | Fast |
| 13B | 7-8 GB | 16-24 GB | Medium |
| 34B | 19-20 GB | 32+ GB | Slow |
| 70B+ | 40+ GB | 64+ GB | Very Slow |

## Model Management

### List Installed Models
```bash
ollama_models
# or
ollama list
```

### Check Model Info
```bash
ollama show llama3.2
ollama show llama3.2 --modelfile
```

### Remove Models
```bash
ollama_remove llama3.2
# or
ollama rm llama3.2
```

### Check Storage Usage
```bash
ollama_status
# Shows storage used by models

# Manual check
du -sh ~/.ollama/models
```

## Pulling Strategies

### For Limited Hardware (8GB RAM)
```bash
# Start with these
ollama_pull llama3.2:1b    # General use
ollama_pull phi3:mini      # Alternative
ollama_pull qwen2.5-coder:1.5b  # For code
```

### For Standard Hardware (16GB RAM)
```bash
# Good balance
ollama_pull llama3.2:3b    # General use
ollama_pull codellama:7b   # For code
ollama_pull mistral        # Alternative
```

### For Powerful Hardware (32GB+ RAM)
```bash
# More capable models
ollama_pull llama3.2:11b   # General use
ollama_pull codellama:34b  # For code
ollama_pull mixtral        # Very capable
```

## Quantization Variants

Models often have quantized versions (smaller, slightly less accurate):
```bash
# Q4 - Smallest, fastest
ollama_pull codellama:7b-q4

# Q5 - Good balance
ollama_pull llama3.2:7b-q5

# Q8 - Higher quality
ollama_pull mistral:7b-q8

# Default is usually Q4 for efficiency
```

## Context

@components/ollama/functions.sh
