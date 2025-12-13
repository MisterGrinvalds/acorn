#!/bin/bash
# AI/ML Module for Automation Framework
# Provides unified interface for Ollama and Hugging Face models

# Source the framework
source "$(dirname "${BASH_SOURCE[0]}")/../framework/core.sh"

# Module information
MODULE_NAME="ai"
MODULE_VERSION="1.0.0"
MODULE_DESCRIPTION="AI/ML model management (Ollama & Hugging Face)"

# Initialize module
ai_init() {
    log_info "Initializing AI/ML module..."
    
    # Check if required tools are available
    if ! command -v python3 >/dev/null 2>&1; then
        log_error "Python 3 is required for AI/ML features"
        return 1
    fi
    
    log_info "AI/ML module initialized"
}

# Main AI command handler
ai_main() {
    local subcommand="$1"
    shift
    
    case "$subcommand" in
        "ollama")
            ai_ollama "$@"
            ;;
        "hf"|"huggingface")
            ai_huggingface "$@"
            ;;
        "setup")
            ai_setup "$@"
            ;;
        "status")
            ai_status "$@"
            ;;
        "models")
            ai_models "$@"
            ;;
        "chat")
            ai_chat "$@"
            ;;
        "generate")
            ai_generate "$@"
            ;;
        "benchmark")
            ai_benchmark "$@"
            ;;
        "examples")
            ai_examples "$@"
            ;;
        "cleanup")
            ai_cleanup "$@"
            ;;
        "--help"|"-h"|"help")
            ai_help
            ;;
        *)
            log_error "Unknown AI command: $subcommand"
            ai_help
            return 1
            ;;
    esac
}

# Ollama management
ai_ollama() {
    local action="$1"
    shift
    
    case "$action" in
        "install")
            log_info "Installing Ollama..."
            if source "$DOTFILES/.bash_tools/ollama.sh" && ollama_install; then
                log_success "Ollama installed successfully"
            else
                log_error "Failed to install Ollama"
                return 1
            fi
            ;;
        "start")
            log_info "Starting Ollama service..."
            source "$DOTFILES/.bash_tools/ollama.sh" && ollama_start
            ;;
        "stop")
            log_info "Stopping Ollama service..."
            source "$DOTFILES/.bash_tools/ollama.sh" && ollama_stop
            ;;
        "status")
            source "$DOTFILES/.bash_tools/ollama.sh" && ollama_status
            ;;
        "models")
            source "$DOTFILES/.bash_tools/ollama.sh" && ollama_models
            ;;
        "pull")
            if [ -z "$1" ]; then
                log_error "Usage: auto ai ollama pull <model_name>"
                return 1
            fi
            log_info "Pulling Ollama model: $1"
            source "$DOTFILES/.bash_tools/ollama.sh" && ollama_pull "$1"
            ;;
        "run")
            if [ -z "$1" ]; then
                log_error "Usage: auto ai ollama run <model_name>"
                return 1
            fi
            source "$DOTFILES/.bash_tools/ollama.sh" && ollama_run "$1"
            ;;
        "chat")
            if [ -z "$1" ] || [ -z "$2" ]; then
                log_error "Usage: auto ai ollama chat <model_name> <prompt>"
                return 1
            fi
            source "$DOTFILES/.bash_tools/ollama.sh" && ollama_chat "$1" "$2"
            ;;
        "setup")
            log_info "Setting up Ollama with recommended models..."
            source "$DOTFILES/.bash_tools/ollama.sh" && ollama_setup
            ;;
        "remove")
            if [ -z "$1" ]; then
                log_error "Usage: auto ai ollama remove <model_name>"
                return 1
            fi
            source "$DOTFILES/.bash_tools/ollama.sh" && ollama_remove "$1"
            ;;
        "benchmark")
            local model="${1:-llama3.2}"
            log_info "Benchmarking Ollama model: $model"
            source "$DOTFILES/.bash_tools/ollama.sh" && ollama_benchmark "$model"
            ;;
        "--help"|"-h"|"help")
            ai_ollama_help
            ;;
        *)
            log_error "Unknown Ollama action: $action"
            ai_ollama_help
            return 1
            ;;
    esac
}

# Hugging Face management
ai_huggingface() {
    local action="$1"
    shift
    
    case "$action" in
        "setup")
            log_info "Setting up Hugging Face environment..."
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_setup
            ;;
        "status")
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_status
            ;;
        "models")
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_models
            ;;
        "download")
            if [ -z "$1" ]; then
                log_error "Usage: auto ai hf download <model_name>"
                return 1
            fi
            log_info "Downloading Hugging Face model: $1"
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_download "$1"
            ;;
        "generate")
            if [ -z "$1" ]; then
                log_error "Usage: auto ai hf generate <prompt> [model_name]"
                return 1
            fi
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_generate "$1" "$2"
            ;;
        "chat")
            local model="${1:-microsoft/DialoGPT-small}"
            log_info "Starting Hugging Face chat with $model"
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_chat "$model"
            ;;
        "summarize")
            if [ -z "$1" ]; then
                log_error "Usage: auto ai hf summarize <text>"
                return 1
            fi
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_summarize "$1"
            ;;
        "sentiment")
            if [ -z "$1" ]; then
                log_error "Usage: auto ai hf sentiment <text>"
                return 1
            fi
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_sentiment "$1"
            ;;
        "qa")
            if [ -z "$1" ] || [ -z "$2" ]; then
                log_error "Usage: auto ai hf qa <question> <context>"
                return 1
            fi
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_qa "$1" "$2"
            ;;
        "code")
            if [ -z "$1" ]; then
                log_error "Usage: auto ai hf code <code_prompt>"
                return 1
            fi
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_code "$1"
            ;;
        "clear-cache")
            log_info "Clearing Hugging Face cache..."
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_clear_cache
            ;;
        "pipelines")
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_pipelines
            ;;
        "examples")
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_examples
            ;;
        "--help"|"-h"|"help")
            ai_huggingface_help
            ;;
        *)
            log_error "Unknown Hugging Face action: $action"
            ai_huggingface_help
            return 1
            ;;
    esac
}

# Setup both AI platforms
ai_setup() {
    local platform="$1"
    
    case "$platform" in
        "ollama")
            ai_ollama setup
            ;;
        "hf"|"huggingface")
            ai_huggingface setup
            ;;
        "all"|"")
            log_info "Setting up complete AI/ML environment..."
            echo ""
            echo "🤖 Installing Ollama..."
            ai_ollama setup
            echo ""
            echo "🤗 Installing Hugging Face..."
            ai_huggingface setup
            echo ""
            log_success "AI/ML environment setup complete!"
            echo ""
            echo "🚀 Quick start commands:"
            echo "  auto ai chat                    # Quick chat with default model"
            echo "  auto ai generate 'Hello world'  # Generate text"
            echo "  auto ai status                  # Check all AI tools"
            echo "  auto ai examples                # Run example commands"
            ;;
        *)
            log_error "Unknown platform: $platform"
            echo "Available platforms: ollama, huggingface, all"
            return 1
            ;;
    esac
}

# Combined status check
ai_status() {
    echo "🤖 AI/ML Environment Status"
    echo "==========================="
    echo ""
    
    # Check Python
    if command -v python3 >/dev/null 2>&1; then
        echo "✅ Python 3: $(python3 --version)"
    else
        echo "❌ Python 3 not found"
    fi
    
    # Check pip
    if command -v pip3 >/dev/null 2>&1; then
        echo "✅ pip3 available"
    else
        echo "❌ pip3 not found"
    fi
    
    echo ""
    echo "📊 Ollama Status:"
    echo "================"
    if command -v ollama >/dev/null 2>&1; then
        source "$DOTFILES/.bash_tools/ollama.sh" && ollama_status
    else
        echo "❌ Ollama not installed"
        echo "   Run: auto ai ollama install"
    fi
    
    echo ""
    echo "🤗 Hugging Face Status:"
    echo "======================"
    if python3 -c "import transformers" 2>/dev/null; then
        source "$DOTFILES/.bash_tools/huggingface.sh" && hf_status
    else
        echo "❌ Hugging Face transformers not installed"
        echo "   Run: auto ai hf setup"
    fi
}

# List all available models
ai_models() {
    echo "🤖 Available AI Models"
    echo "====================="
    echo ""
    
    echo "📊 Ollama Models:"
    echo "================"
    if command -v ollama >/dev/null 2>&1; then
        source "$DOTFILES/.bash_tools/ollama.sh" && ollama_models
    else
        echo "❌ Ollama not installed"
    fi
    
    echo ""
    echo "🤗 Hugging Face Models:"
    echo "======================"
    source "$DOTFILES/.bash_tools/huggingface.sh" && hf_models
}

# Universal chat interface
ai_chat() {
    local platform="$1"
    local model="$2"
    
    case "$platform" in
        "ollama")
            if [ -z "$model" ]; then
                model="llama3.2"
            fi
            log_info "Starting Ollama chat with $model"
            source "$DOTFILES/.bash_tools/ollama.sh" && ollama_run "$model"
            ;;
        "hf"|"huggingface")
            if [ -z "$model" ]; then
                model="microsoft/DialoGPT-small"
            fi
            log_info "Starting Hugging Face chat with $model"
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_chat "$model"
            ;;
        "")
            # Default to Ollama if available, otherwise Hugging Face
            if command -v ollama >/dev/null 2>&1; then
                log_info "Starting chat with Ollama (default)"
                source "$DOTFILES/.bash_tools/ollama.sh" && ollama_run "llama3.2"
            elif python3 -c "import transformers" 2>/dev/null; then
                log_info "Starting chat with Hugging Face (fallback)"
                source "$DOTFILES/.bash_tools/huggingface.sh" && hf_chat
            else
                log_error "No AI platforms available. Run: auto ai setup"
                return 1
            fi
            ;;
        *)
            log_error "Unknown platform: $platform"
            echo "Usage: auto ai chat [ollama|hf] [model_name]"
            return 1
            ;;
    esac
}

# Universal text generation
ai_generate() {
    local prompt="$1"
    local platform="$2"
    local model="$3"
    
    if [ -z "$prompt" ]; then
        log_error "Usage: auto ai generate <prompt> [ollama|hf] [model_name]"
        return 1
    fi
    
    case "$platform" in
        "ollama")
            if [ -z "$model" ]; then
                model="llama3.2"
            fi
            source "$DOTFILES/.bash_tools/ollama.sh" && ollama_chat "$model" "$prompt"
            ;;
        "hf"|"huggingface")
            if [ -z "$model" ]; then
                model="microsoft/DialoGPT-small"
            fi
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_generate "$prompt" "$model"
            ;;
        "")
            # Try Ollama first, then Hugging Face
            if command -v ollama >/dev/null 2>&1; then
                source "$DOTFILES/.bash_tools/ollama.sh" && ollama_chat "llama3.2" "$prompt"
            elif python3 -c "import transformers" 2>/dev/null; then
                source "$DOTFILES/.bash_tools/huggingface.sh" && hf_generate "$prompt"
            else
                log_error "No AI platforms available. Run: auto ai setup"
                return 1
            fi
            ;;
        *)
            log_error "Unknown platform: $platform"
            echo "Usage: auto ai generate <prompt> [ollama|hf] [model_name]"
            return 1
            ;;
    esac
}

# Benchmark AI platforms
ai_benchmark() {
    local platform="$1"
    
    echo "⏱️ AI/ML Performance Benchmark"
    echo "=============================="
    echo ""
    
    case "$platform" in
        "ollama")
            if command -v ollama >/dev/null 2>&1; then
                echo "📊 Benchmarking Ollama..."
                source "$DOTFILES/.bash_tools/ollama.sh" && ollama_benchmark
            else
                echo "❌ Ollama not available"
            fi
            ;;
        "hf"|"huggingface")
            if python3 -c "import transformers, time" 2>/dev/null; then
                echo "🤗 Benchmarking Hugging Face..."
                python3 -c "
import time
import sys
try:
    from transformers import pipeline
    
    print('Loading model...')
    start = time.time()
    generator = pipeline('text-generation', model='distilgpt2')
    load_time = time.time() - start
    
    print('Generating text...')
    start = time.time()
    result = generator('Hello world', max_length=50, num_return_sequences=1)
    gen_time = time.time() - start
    
    print(f'⏱️ Model load time: {load_time:.2f}s')
    print(f'⏱️ Generation time: {gen_time:.2f}s')
    print(f'📝 Generated: {result[0][\"generated_text\"]}')
    
except Exception as e:
    print(f'❌ Benchmark failed: {e}')
    sys.exit(1)
"
            else
                echo "❌ Hugging Face not available"
            fi
            ;;
        ""|"all")
            ai_benchmark "ollama"
            echo ""
            ai_benchmark "hf"
            ;;
        *)
            log_error "Unknown platform: $platform"
            echo "Usage: auto ai benchmark [ollama|hf|all]"
            return 1
            ;;
    esac
}

# Show examples
ai_examples() {
    echo "🎯 AI/ML Usage Examples"
    echo "======================"
    echo ""
    
    echo "🚀 Quick Start:"
    echo "  auto ai setup                    # Install everything"
    echo "  auto ai chat                     # Start interactive chat"
    echo "  auto ai generate 'Hello world'   # Generate text"
    echo ""
    
    echo "📊 Ollama Examples:"
    echo "  auto ai ollama install           # Install Ollama"
    echo "  auto ai ollama pull llama3.2     # Download model"
    echo "  auto ai ollama run llama3.2      # Interactive chat"
    echo "  auto ai ollama chat llama3.2 'Hello'  # Quick question"
    echo ""
    
    echo "🤗 Hugging Face Examples:"
    echo "  auto ai hf setup                 # Setup environment"
    echo "  auto ai hf download gpt2         # Download model"
    echo "  auto ai hf generate 'Once upon'  # Generate text"
    echo "  auto ai hf sentiment 'I love AI' # Sentiment analysis"
    echo "  auto ai hf summarize 'Long text' # Summarize text"
    echo ""
    
    echo "🛠️ Management:"
    echo "  auto ai status                   # Check all systems"
    echo "  auto ai models                   # List all models"
    echo "  auto ai benchmark                # Performance test"
    echo "  auto ai cleanup                  # Clean up resources"
    
    # Run actual examples if requested
    if [ "$1" = "--run" ]; then
        echo ""
        echo "🎯 Running live examples..."
        echo ""
        
        if command -v ollama >/dev/null 2>&1; then
            echo "📊 Ollama example:"
            source "$DOTFILES/.bash_tools/ollama.sh" && ollama_examples
        fi
        
        if python3 -c "import transformers" 2>/dev/null; then
            echo ""
            echo "🤗 Hugging Face example:"
            source "$DOTFILES/.bash_tools/huggingface.sh" && hf_examples
        fi
    fi
}

# Cleanup AI resources
ai_cleanup() {
    echo "🧹 AI/ML Cleanup"
    echo "==============="
    echo ""
    
    # Stop Ollama service
    if command -v ollama >/dev/null 2>&1; then
        echo "🛑 Stopping Ollama service..."
        source "$DOTFILES/.bash_tools/ollama.sh" && ollama_stop
    fi
    
    # Clean Hugging Face cache
    if [ -d ~/.cache/huggingface ]; then
        echo "🗑️ Cleaning Hugging Face cache..."
        source "$DOTFILES/.bash_tools/huggingface.sh" && hf_clear_cache
    fi
    
    # Clean Python cache
    if [ -d ~/.cache/pip ]; then
        echo "🗑️ Cleaning pip cache..."
        pip3 cache purge 2>/dev/null || echo "   pip cache clean skipped"
    fi
    
    echo "✅ Cleanup complete"
}

# Help functions
ai_help() {
    echo "🤖 AI/ML Module Help"
    echo "==================="
    echo ""
    echo "Commands:"
    echo "  setup [platform]       # Setup AI/ML environment"
    echo "  status                  # Check all AI systems"
    echo "  models                  # List available models"
    echo "  chat [platform] [model] # Interactive chat"
    echo "  generate <prompt>       # Generate text"
    echo "  benchmark [platform]    # Performance test"
    echo "  examples [--run]        # Show/run examples"
    echo "  cleanup                 # Clean up resources"
    echo ""
    echo "  ollama <action>         # Ollama management"
    echo "  hf <action>             # Hugging Face management"
    echo ""
    echo "Platforms: ollama, hf (huggingface)"
    echo ""
    echo "Examples:"
    echo "  auto ai setup                    # Install everything"
    echo "  auto ai chat                     # Start chat with best available model"
    echo "  auto ai generate 'Hello world'   # Generate text"
    echo "  auto ai ollama pull llama3.2     # Install Ollama model"
    echo "  auto ai hf sentiment 'Great!'    # Analyze sentiment"
}

ai_ollama_help() {
    echo "📊 Ollama Commands"
    echo "=================="
    echo ""
    echo "Commands:"
    echo "  install                 # Install Ollama"
    echo "  start                   # Start Ollama service"
    echo "  stop                    # Stop Ollama service"
    echo "  status                  # Check Ollama status"
    echo "  models                  # List installed models"
    echo "  pull <model>            # Download model"
    echo "  run <model>             # Interactive chat"
    echo "  chat <model> <prompt>   # Quick question"
    echo "  remove <model>          # Remove model"
    echo "  setup                   # Install with recommended models"
    echo "  benchmark [model]       # Performance test"
    echo ""
    echo "Examples:"
    echo "  auto ai ollama setup           # Complete setup"
    echo "  auto ai ollama pull llama3.2   # Download model"
    echo "  auto ai ollama run llama3.2    # Start chat"
}

ai_huggingface_help() {
    echo "🤗 Hugging Face Commands"
    echo "========================"
    echo ""
    echo "Commands:"
    echo "  setup                   # Setup HF environment"
    echo "  status                  # Check HF status"
    echo "  models                  # List popular models"
    echo "  download <model>        # Download model"
    echo "  generate <prompt>       # Generate text"
    echo "  chat [model]            # Interactive chat"
    echo "  summarize <text>        # Summarize text"
    echo "  sentiment <text>        # Sentiment analysis"
    echo "  qa <question> <context> # Question answering"
    echo "  code <prompt>           # Code generation"
    echo "  clear-cache             # Clear model cache"
    echo "  pipelines               # List available pipelines"
    echo "  examples                # Run examples"
    echo ""
    echo "Examples:"
    echo "  auto ai hf setup                # Setup environment"
    echo "  auto ai hf generate 'Hello'     # Generate text"
    echo "  auto ai hf sentiment 'Great!'   # Analyze sentiment"
}

# Register module
register_module "$MODULE_NAME" ai_main ai_help