#!/bin/sh
# Ollama Local AI Model Management
# Provides functions for running and managing Ollama models locally

# Ollama model management functions
ollama_install() {
    if command -v ollama >/dev/null 2>&1; then
        echo "✅ Ollama already installed"
        ollama --version
        return 0
    fi
    
    echo "📦 Installing Ollama..."
    if [ "$OSTYPE" = "darwin"* ]; then
        if command -v brew >/dev/null 2>&1; then
            brew install ollama
        else
            echo "🌐 Installing from ollama.ai..."
            curl -fsSL https://ollama.ai/install.sh | sh
        fi
    elif [ "$OSTYPE" = "linux-gnu" ]; then
        echo "🌐 Installing from ollama.ai..."
        curl -fsSL https://ollama.ai/install.sh | sh
    else
        echo "❌ Unsupported OS for auto-install. Please visit https://ollama.ai"
        return 1
    fi
}

# Start Ollama service
ollama_start() {
    if pgrep -f ollama >/dev/null; then
        echo "✅ Ollama service already running"
    else
        echo "🚀 Starting Ollama service..."
        ollama serve &
        sleep 3
        echo "✅ Ollama service started"
    fi
}

# Stop Ollama service
ollama_stop() {
    if pgrep -f ollama >/dev/null; then
        echo "🛑 Stopping Ollama service..."
        pkill -f ollama
        echo "✅ Ollama service stopped"
    else
        echo "ℹ️ Ollama service not running"
    fi
}

# List available models
ollama_models() {
    if ! command -v ollama >/dev/null 2>&1; then
        echo "❌ Ollama not installed. Run: ollama_install"
        return 1
    fi
    
    echo "📋 Available Ollama models:"
    ollama list
}

# Pull/download a model
ollama_pull() {
    local model="$1"
    if [ -z "$model" ]; then
        echo "❌ Usage: ollama_pull <model_name>"
        echo "Popular models: llama3.2, codellama, mistral, phi3, gemma2"
        return 1
    fi
    
    if ! command -v ollama >/dev/null 2>&1; then
        echo "❌ Ollama not installed. Run: ollama_install"
        return 1
    fi
    
    echo "📥 Pulling model: $model"
    ollama pull "$model"
}

# Run a model interactively
ollama_run() {
    local model="$1"
    if [ -z "$model" ]; then
        echo "❌ Usage: ollama_run <model_name>"
        echo "Available models:"
        ollama list
        return 1
    fi
    
    if ! command -v ollama >/dev/null 2>&1; then
        echo "❌ Ollama not installed. Run: ollama_install"
        return 1
    fi
    
    echo "🤖 Starting interactive session with $model"
    echo "Type 'exit' or press Ctrl+D to quit"
    ollama run "$model"
}

# Quick chat with a model
ollama_chat() {
    local model="$1"
    local prompt="$2"
    
    if [ -z "$model" ] || [ -z "$prompt" ]; then
        echo "❌ Usage: ollama_chat <model_name> <prompt>"
        echo "Example: ollama_chat llama3.2 'Explain quantum computing'"
        return 1
    fi
    
    if ! command -v ollama >/dev/null 2>&1; then
        echo "❌ Ollama not installed. Run: ollama_install"
        return 1
    fi
    
    echo "🤖 Asking $model: $prompt"
    echo ""
    echo "$prompt" | ollama run "$model"
}

# Code generation with models
ollama_code() {
    local language="$1"
    local description="$2"
    
    if [ -z "$language" ] || [ -z "$description" ]; then
        echo "❌ Usage: ollama_code <language> <description>"
        echo "Example: ollama_code python 'function to sort a list'"
        return 1
    fi
    
    local model="codellama"
    local prompt="Write a $language $description. Only return the code without explanation:"
    
    if ! ollama list | grep -q "$model"; then
        echo "📥 CodeLlama not found, pulling model..."
        ollama_pull "$model"
    fi
    
    echo "🧑‍💻 Generating $language code: $description"
    echo ""
    echo "$prompt" | ollama run "$model"
}

# Model benchmarking
ollama_benchmark() {
    local model="${1:-llama3.2}"
    local prompt="Write a simple hello world program in Python"
    
    if ! command -v ollama >/dev/null 2>&1; then
        echo "❌ Ollama not installed. Run: ollama_install"
        return 1
    fi
    
    if ! ollama list | grep -q "$model"; then
        echo "📥 Model $model not found, pulling..."
        ollama_pull "$model"
    fi
    
    echo "⏱️ Benchmarking $model..."
    echo "Prompt: $prompt"
    echo ""
    
    local start_time=$(date +%s)
    echo "$prompt" | ollama run "$model" >/dev/null
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    echo "⏱️ Response time: ${duration}s"
}

# Setup recommended models
ollama_setup() {
    echo "🎯 Setting up recommended Ollama models..."
    
    if ! command -v ollama >/dev/null 2>&1; then
        echo "📦 Installing Ollama first..."
        ollama_install
    fi
    
    echo "🚀 Starting Ollama service..."
    ollama_start
    
    echo "📥 Pulling recommended models..."
    
    # General purpose models
    echo "📥 Pulling Llama 3.2 (3B - good balance of speed/quality)..."
    ollama_pull "llama3.2"
    
    # Code generation
    echo "📥 Pulling CodeLlama (7B - specialized for code)..."
    ollama_pull "codellama"
    
    # Lightweight model
    echo "📥 Pulling Phi-3 (3.8B - Microsoft's efficient model)..."
    ollama_pull "phi3"
    
    echo "✅ Ollama setup complete!"
    echo ""
    echo "🚀 Try these commands:"
    echo "  ollama_run llama3.2           # Interactive chat"
    echo "  ollama_chat llama3.2 'Hello'  # Quick question"
    echo "  ollama_code python 'sort list' # Generate code"
    echo "  ollama_models                 # List all models"
}

# Remove a model
ollama_remove() {
    local model="$1"
    if [ -z "$model" ]; then
        echo "❌ Usage: ollama_remove <model_name>"
        echo "Available models:"
        ollama list
        return 1
    fi
    
    echo "🗑️ Removing model: $model"
    ollama rm "$model"
}

# Show model information
ollama_info() {
    local model="$1"
    if [ -z "$model" ]; then
        echo "❌ Usage: ollama_info <model_name>"
        echo "Available models:"
        ollama list
        return 1
    fi
    
    echo "ℹ️ Model information for: $model"
    ollama show "$model"
}

# Quick status check
ollama_status() {
    echo "🔍 Ollama Status Check"
    echo "====================="
    
    if command -v ollama >/dev/null 2>&1; then
        echo "✅ Ollama installed: $(ollama --version)"
    else
        echo "❌ Ollama not installed"
        return 1
    fi
    
    if pgrep -f ollama >/dev/null; then
        echo "✅ Ollama service running"
    else
        echo "❌ Ollama service not running"
    fi
    
    echo ""
    echo "📋 Installed models:"
    ollama list
    
    echo ""
    echo "💾 Storage usage:"
    if command -v du >/dev/null 2>&1; then
        if [ -d ~/.ollama ]; then
            du -sh ~/.ollama 2>/dev/null || echo "Unable to calculate storage"
        else
            echo "No Ollama data directory found"
        fi
    fi
}

# Model usage examples
ollama_examples() {
    echo "🎯 Ollama Usage Examples"
    echo "======================="
    echo ""
    echo "1. General Chat:"
    echo "   ollama_chat llama3.2 'Explain machine learning in simple terms'"
    echo ""
    echo "2. Code Generation:"
    echo "   ollama_code python 'function to calculate fibonacci sequence'"
    echo "   ollama_code javascript 'async function to fetch API data'"
    echo ""
    echo "3. Interactive Session:"
    echo "   ollama_run llama3.2"
    echo ""
    echo "4. Model Management:"
    echo "   ollama_models          # List installed models"
    echo "   ollama_pull mistral    # Install new model"
    echo "   ollama_remove phi3     # Remove model"
    echo ""
    echo "5. Quick Setup:"
    echo "   ollama_setup           # Install and configure everything"
    echo ""
    echo "6. Status Check:"
    echo "   ollama_status          # Check installation and running models"
}

# Alias for convenience
alias ollama-start='ollama_start'
alias ollama-stop='ollama_stop'
alias ollama-status='ollama_status'
alias ollama-setup='ollama_setup'
alias ollama-models='ollama_models'
alias ollama-examples='ollama_examples'