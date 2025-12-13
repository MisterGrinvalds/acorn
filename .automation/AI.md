# AI/ML Integration Guide

This guide covers the comprehensive AI/ML capabilities integrated into the bash profile automation framework, including Ollama and Hugging Face model management.

## 🚀 Quick Start

```bash
# Setup complete AI/ML environment
make ai-setup

# Or use automation framework directly
auto ai setup all

# Start chatting immediately
make ai-chat
# or
auto ai chat
```

## 🤖 Available Platforms

### 📊 Ollama
- **Local LLM hosting** - Run models locally without internet
- **Multiple models** - Llama 3.2, CodeLlama, Phi-3, Mistral, and more
- **Fast inference** - Optimized for local hardware
- **No API costs** - Completely free to use

### 🤗 Hugging Face
- **Transformers library** - Access to thousands of models
- **Specialized tasks** - Text generation, sentiment analysis, QA, summarization
- **Fine-tuned models** - Pre-trained models for specific use cases
- **Python ecosystem** - Full ML/AI Python stack

## 📋 Installation Options

### Option 1: Complete Setup (Recommended)
```bash
# Install both platforms with recommended models
make ai-setup
```

### Option 2: Platform-Specific Setup
```bash
# Ollama only
make ollama-setup
# or
auto ai ollama setup

# Hugging Face only
make hf-setup
# or
auto ai hf setup
```

### Option 3: Manual Installation
```bash
# Install Ollama
auto ai ollama install

# Setup Python environment for Hugging Face
auto ai hf setup
```

## 🚀 Common Usage Patterns

### Interactive Chat
```bash
# Auto-detect best available model
auto ai chat

# Specific platform
auto ai chat ollama llama3.2
auto ai chat hf microsoft/DialoGPT-small

# Make targets
make ai-chat           # Auto-detect
make ai-chat-ollama    # Ollama with Llama 3.2
make ai-chat-hf        # Hugging Face
```

### Text Generation
```bash
# Quick generation
auto ai generate "Explain quantum computing"

# Platform-specific
auto ai generate "Write a Python function" ollama codellama
auto ai hf generate "Once upon a time" gpt2
```

### Specialized Tasks
```bash
# Sentiment analysis
auto ai hf sentiment "I love this product!"

# Text summarization
auto ai hf summarize "Long article text here..."

# Question answering
auto ai hf qa "What is Python?" "Python is a programming language..."

# Code generation
auto ai ollama chat codellama "Write a REST API in Python"
auto ai hf code "def fibonacci(n):"
```

## 🔧 Management Commands

### Status and Information
```bash
# Check all AI systems
make ai-status
auto ai status

# List available models
make ai-models
auto ai models

# Platform-specific status
make ollama-status
make hf-status
```

### Model Management
```bash
# Download/pull models
auto ai ollama pull llama3.2
auto ai hf download microsoft/DialoGPT-small

# List installed models
auto ai ollama models
auto ai hf models

# Remove models
auto ai ollama remove phi3
```

### Service Management
```bash
# Start/stop Ollama service
make ollama-start
make ollama-stop

# Service status
make ollama-status
```

### Cache and Cleanup
```bash
# Clean up all AI resources
make ai-cleanup
auto ai cleanup

# Clear specific caches
make hf-clear-cache
auto ai hf clear-cache
```

## 🎯 Available Models

### 📊 Ollama Models
| Model | Size | Use Case | Command |
|-------|------|----------|---------|
| llama3.2 | 3B | General chat | `auto ai ollama pull llama3.2` |
| codellama | 7B | Code generation | `auto ai ollama pull codellama` |
| phi3 | 3.8B | Efficient general purpose | `auto ai ollama pull phi3` |
| mistral | 7B | Instruction following | `auto ai ollama pull mistral` |
| gemma2 | 9B | Google's model | `auto ai ollama pull gemma2` |

### 🤗 Hugging Face Models
| Model | Use Case | Command |
|-------|----------|---------|
| microsoft/DialoGPT-small | Conversational AI | `auto ai hf download microsoft/DialoGPT-small` |
| distilgpt2 | Fast text generation | `auto ai hf download distilgpt2` |
| distilbert-base-uncased | Text classification | `auto ai hf download distilbert-base-uncased` |
| facebook/bart-large-cnn | Summarization | `auto ai hf download facebook/bart-large-cnn` |
| microsoft/codebert-base | Code understanding | `auto ai hf download microsoft/codebert-base` |

## 🧪 Testing and Validation

### Basic Tests
```bash
# Test AI functionality
make ai-test

# Check if AI tools are available
make test-ai-tools
```

### Performance Benchmarking
```bash
# Benchmark all platforms
make ai-benchmark
auto ai benchmark

# Platform-specific benchmarks
auto ai benchmark ollama
auto ai benchmark hf
```

### Example Workflows
```bash
# Show usage examples
make ai-examples
auto ai examples

# Run live examples (requires models)
make ai-examples-run
auto ai examples --run
```

## 📁 File Structure

```
.bash_tools/
├── ollama.sh              # Ollama shell functions
└── huggingface.sh         # Hugging Face shell functions

.automation/
├── modules/
│   └── ai.sh              # AI automation module
└── AI.md                  # This documentation

Makefile                   # AI/ML make targets
```

## ⚙️ Configuration

### Environment Variables
```bash
# Ollama
export OLLAMA_MODELS="/path/to/models"  # Custom model location

# Hugging Face
export TRANSFORMERS_CACHE="/path/to/cache"  # Custom cache location
export HF_HOME="/path/to/hf"  # Hugging Face home directory
```

### Virtual Environments
The Hugging Face setup automatically creates a virtual environment:
```bash
# Activate manually if needed
source ~/.venvs/huggingface/bin/activate

# Or use the enhanced mkvenv function
mkvenv huggingface
venv huggingface
```

## 🚨 Troubleshooting

### Common Issues

**Ollama not starting:**
```bash
# Check if port 11434 is in use
lsof -i :11434

# Restart service
make ollama-stop
make ollama-start
```

**Hugging Face import errors:**
```bash
# Reinstall in fresh environment
auto ai hf setup

# Check Python environment
python3 -c "import transformers; print(transformers.__version__)"
```

**Model download failures:**
```bash
# Check internet connection and disk space
df -h
ping huggingface.co

# Clear cache and retry
auto ai hf clear-cache
auto ai hf download model-name
```

**Memory issues:**
```bash
# Use smaller models for testing
auto ai ollama pull llama3.2:1b  # 1B parameter version
auto ai hf download distilgpt2   # Distilled model
```

### Debug Commands
```bash
# Verbose status check
auto ai status

# Check logs
ls tests/logs/ai-*.log

# Test syntax
bash -n .bash_tools/ollama.sh
bash -n .bash_tools/huggingface.sh
bash -n .automation/modules/ai.sh
```

## 🔒 Security Considerations

- **Local models**: Ollama runs entirely locally with no external API calls
- **HF API keys**: Store in secrets management system if using Hugging Face Hub
- **Model cache**: Consider disk space and cleanup policies
- **Network usage**: Model downloads can be large (GB sizes)

## 🚀 Advanced Usage

### Custom Model Integration
```bash
# Use custom Ollama model
ollama pull custom-model:latest
auto ai ollama run custom-model

# Use local Hugging Face model
auto ai hf download ./path/to/local/model
```

### Automation Workflows
```bash
# Generate code and save to file
auto ai ollama chat codellama "Python REST API" > generated_api.py

# Batch sentiment analysis
echo "Great product!" | auto ai hf sentiment
echo "Terrible service!" | auto ai hf sentiment
```

### Integration with Other Tools
```bash
# Use with tmux for persistent sessions
tmux new-session -d -s ai-chat 'auto ai chat'

# Combine with file processing
cat README.md | auto ai hf summarize
```

## 📊 Performance Characteristics

### Ollama
- **Cold start**: 2-5 seconds (model loading)
- **Warm inference**: 100-500ms per response
- **Memory usage**: 2-8GB depending on model size
- **Disk usage**: 2-20GB per model

### Hugging Face
- **Model download**: 1-10 minutes (depending on size)
- **First inference**: 5-30 seconds (model loading)
- **Subsequent inference**: 500ms-5s per response
- **Memory usage**: 1-4GB for typical models

## 🎯 Best Practices

1. **Start small**: Begin with smaller models (llama3.2, distilgpt2)
2. **Monitor resources**: Check disk space and memory usage
3. **Use appropriate models**: CodeLlama for code, BART for summarization
4. **Cache management**: Regularly clean up unused models
5. **Virtual environments**: Keep AI dependencies isolated
6. **Backup important models**: Some models take time to download

## 🔗 Useful Resources

- [Ollama Model Library](https://ollama.ai/library)
- [Hugging Face Model Hub](https://huggingface.co/models)
- [Transformers Documentation](https://huggingface.co/docs/transformers)
- [Ollama GitHub](https://github.com/ollama/ollama)

## 📞 Support

If you encounter issues:
1. Check the troubleshooting section above
2. Run `make ai-test` to validate your setup
3. Check logs in `tests/logs/ai-*.log`
4. Review the automation framework documentation