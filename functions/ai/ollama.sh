#!/bin/sh
# Ollama Local AI Model Management

# Ollama model management functions
ollama_install() {
    if command -v ollama >/dev/null 2>&1; then
        echo "Ollama already installed"
        ollama --version
        return 0
    fi

    echo "Installing Ollama..."
    case "$CURRENT_PLATFORM" in
        darwin)
            if command -v brew >/dev/null 2>&1; then
                brew install ollama
            else
                curl -fsSL https://ollama.ai/install.sh | sh
            fi
            ;;
        linux)
            curl -fsSL https://ollama.ai/install.sh | sh
            ;;
        *)
            echo "Unsupported OS. Visit https://ollama.ai"
            return 1
            ;;
    esac
}

# Start Ollama service
ollama_start() {
    if pgrep -f ollama >/dev/null; then
        echo "Ollama service already running"
    else
        echo "Starting Ollama service..."
        ollama serve &
        sleep 3
        echo "Ollama service started"
    fi
}

# Stop Ollama service
ollama_stop() {
    if pgrep -f ollama >/dev/null; then
        echo "Stopping Ollama service..."
        pkill -f ollama
        echo "Ollama service stopped"
    else
        echo "Ollama service not running"
    fi
}

# List available models
ollama_models() {
    if ! command -v ollama >/dev/null 2>&1; then
        echo "Ollama not installed. Run: ollama_install"
        return 1
    fi
    echo "Available Ollama models:"
    ollama list
}

# Pull/download a model
ollama_pull() {
    local model="$1"
    if [ -z "$model" ]; then
        echo "Usage: ollama_pull <model_name>"
        echo "Popular models: llama3.2, codellama, mistral, phi3, gemma2"
        return 1
    fi

    if ! command -v ollama >/dev/null 2>&1; then
        echo "Ollama not installed. Run: ollama_install"
        return 1
    fi

    echo "Pulling model: $model"
    ollama pull "$model"
}

# Run a model interactively
ollama_run() {
    local model="$1"
    if [ -z "$model" ]; then
        echo "Usage: ollama_run <model_name>"
        ollama list
        return 1
    fi

    if ! command -v ollama >/dev/null 2>&1; then
        echo "Ollama not installed. Run: ollama_install"
        return 1
    fi

    echo "Starting interactive session with $model"
    ollama run "$model"
}

# Quick chat with a model
ollama_chat() {
    local model="$1"
    local prompt="$2"

    if [ -z "$model" ] || [ -z "$prompt" ]; then
        echo "Usage: ollama_chat <model_name> <prompt>"
        return 1
    fi

    if ! command -v ollama >/dev/null 2>&1; then
        echo "Ollama not installed. Run: ollama_install"
        return 1
    fi

    echo "Asking $model: $prompt"
    echo ""
    echo "$prompt" | ollama run "$model"
}

# Code generation with models
ollama_code() {
    local language="$1"
    local description="$2"

    if [ -z "$language" ] || [ -z "$description" ]; then
        echo "Usage: ollama_code <language> <description>"
        return 1
    fi

    local model="codellama"
    local prompt="Write a $language $description. Only return the code:"

    if ! ollama list | grep -q "$model"; then
        echo "CodeLlama not found, pulling model..."
        ollama_pull "$model"
    fi

    echo "Generating $language code: $description"
    echo ""
    echo "$prompt" | ollama run "$model"
}

# Quick status check
ollama_status() {
    echo "Ollama Status Check"
    echo "==================="

    if command -v ollama >/dev/null 2>&1; then
        echo "Ollama installed: $(ollama --version)"
    else
        echo "Ollama not installed"
        return 1
    fi

    if pgrep -f ollama >/dev/null; then
        echo "Ollama service running"
    else
        echo "Ollama service not running"
    fi

    echo ""
    echo "Installed models:"
    ollama list

    echo ""
    echo "Storage usage:"
    if [ -d ~/.ollama ]; then
        du -sh ~/.ollama 2>/dev/null || echo "Unable to calculate storage"
    else
        echo "No Ollama data directory found"
    fi
}

# Remove a model
ollama_remove() {
    local model="$1"
    if [ -z "$model" ]; then
        echo "Usage: ollama_remove <model_name>"
        ollama list
        return 1
    fi

    echo "Removing model: $model"
    ollama rm "$model"
}

# Model usage examples
ollama_examples() {
    echo "Ollama Usage Examples"
    echo "====================="
    echo ""
    echo "1. General Chat:"
    echo "   ollama_chat llama3.2 'Explain machine learning'"
    echo ""
    echo "2. Code Generation:"
    echo "   ollama_code python 'function to calculate fibonacci'"
    echo ""
    echo "3. Interactive Session:"
    echo "   ollama_run llama3.2"
    echo ""
    echo "4. Model Management:"
    echo "   ollama_models          # List installed models"
    echo "   ollama_pull mistral    # Install new model"
    echo "   ollama_remove phi3     # Remove model"
    echo ""
    echo "5. Status Check:"
    echo "   ollama_status          # Check installation"
}

# Aliases for convenience
alias ollama-start='ollama_start'
alias ollama-stop='ollama_stop'
alias ollama-status='ollama_status'
alias ollama-models='ollama_models'
alias ollama-examples='ollama_examples'
