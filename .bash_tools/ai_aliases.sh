#!/bin/sh
# AI/ML Convenience Aliases and Functions
# Quick access to AI tools without make

# Source the AI tools
source_ai_tools() {
    if [ -f "$DOTFILES/.bash_tools/ollama.sh" ]; then
        source "$DOTFILES/.bash_tools/ollama.sh"
    fi
    if [ -f "$DOTFILES/.bash_tools/huggingface.sh" ]; then
        source "$DOTFILES/.bash_tools/huggingface.sh"
    fi
}

# Auto-source on shell load
source_ai_tools

# Quick AI chat - auto-detects best available platform
ai() {
    local command="$1"
    shift
    
    case "$command" in
        "chat"|"c")
            if command -v ollama >/dev/null 2>&1; then
                ollama_run "${1:-llama3.2}"
            elif python3 -c "import transformers" 2>/dev/null; then
                hf_chat "$@"
            else
                echo "❌ No AI platforms available. Run: make ai-setup"
                return 1
            fi
            ;;
        "ask"|"q")
            local prompt="$*"
            if [ -z "$prompt" ]; then
                echo "Usage: ai ask <your question>"
                return 1
            fi
            
            if command -v ollama >/dev/null 2>&1; then
                ollama_chat llama3.2 "$prompt"
            elif python3 -c "import transformers" 2>/dev/null; then
                hf_generate "$prompt"
            else
                echo "❌ No AI platforms available. Run: make ai-setup"
                return 1
            fi
            ;;
        "code"|"dev")
            local lang="$1"
            local desc="$2"
            
            if [ -z "$lang" ] || [ -z "$desc" ]; then
                echo "Usage: ai code <language> <description>"
                echo "Example: ai code python 'function to sort a list'"
                return 1
            fi
            
            if command -v ollama >/dev/null 2>&1; then
                ollama_code "$lang" "$desc"
            else
                echo "Code generation requires Ollama. Run: make ollama-setup"
                return 1
            fi
            ;;
        "status"|"info")
            echo "🤖 Quick AI Status"
            echo "=================="
            
            if command -v ollama >/dev/null 2>&1; then
                echo "✅ Ollama available"
                if pgrep -f ollama >/dev/null; then
                    echo "✅ Ollama service running"
                else
                    echo "❌ Ollama service not running"
                fi
            else
                echo "❌ Ollama not installed"
            fi
            
            if python3 -c "import transformers" 2>/dev/null; then
                echo "✅ Hugging Face available"
            else
                echo "❌ Hugging Face not installed"
            fi
            ;;
        "setup"|"install")
            echo "🚀 Setting up AI tools..."
            make ai-setup
            ;;
        "models"|"list")
            if command -v ollama >/dev/null 2>&1; then
                echo "📊 Ollama Models:"
                ollama_models
            else
                echo "❌ Ollama not installed"
            fi
            ;;
        "help"|"--help"|"-h")
            ai_help
            ;;
        *)
            echo "❌ Unknown command: $command"
            ai_help
            return 1
            ;;
    esac
}

# Help function for ai command
ai_help() {
    echo "🤖 AI Quick Commands"
    echo "==================="
    echo ""
    echo "Usage: ai <command> [options]"
    echo ""
    echo "Commands:"
    echo "  chat, c [model]           # Start interactive chat"
    echo "  ask, q <question>         # Ask a quick question"
    echo "  code, dev <lang> <desc>   # Generate code"
    echo "  status, info              # Check AI status"
    echo "  setup, install            # Setup AI environment"
    echo "  models, list              # List available models"
    echo "  help, --help, -h          # Show this help"
    echo ""
    echo "Examples:"
    echo "  ai chat                   # Start chat with best model"
    echo "  ai ask 'What is Python?'  # Quick question"
    echo "  ai code python 'sort list' # Generate Python code"
    echo "  ai status                 # Check what's installed"
    echo ""
    echo "Full toolkit:"
    echo "  make ai-setup             # Complete setup"
    echo "  make ai-status            # Detailed status"
    echo "  make ai-models            # All available models"
}

# Convenience aliases
alias aichat='ai chat'
alias aiask='ai ask'
alias aicode='ai code'
alias aistatus='ai status'
alias aisetup='ai setup'

# Model-specific aliases
alias ollama-chat='ollama_run'
alias ollama-ask='ollama_chat'
alias hf-chat='hf_chat'
alias hf-ask='hf_generate'

# Quick sentiment analysis
sentiment() {
    local text="$*"
    if [ -z "$text" ]; then
        echo "Usage: sentiment <text to analyze>"
        return 1
    fi
    
    if python3 -c "import transformers" 2>/dev/null; then
        hf_sentiment "$text"
    else
        echo "❌ Sentiment analysis requires Hugging Face. Run: make hf-setup"
        return 1
    fi
}

# Quick text summarization
summarize() {
    local text="$*"
    if [ -z "$text" ]; then
        echo "Usage: summarize <text to summarize>"
        return 1
    fi
    
    if python3 -c "import transformers" 2>/dev/null; then
        hf_summarize "$text"
    else
        echo "❌ Summarization requires Hugging Face. Run: make hf-setup"
        return 1
    fi
}

# Quick model downloading
get_model() {
    local platform="$1"
    local model="$2"
    
    case "$platform" in
        "ollama"|"o")
            if [ -z "$model" ]; then
                echo "Usage: get_model ollama <model_name>"
                echo "Popular models: llama3.2, codellama, phi3, mistral"
                return 1
            fi
            ollama_pull "$model"
            ;;
        "hf"|"huggingface"|"h")
            if [ -z "$model" ]; then
                echo "Usage: get_model hf <model_name>"
                echo "Popular models: microsoft/DialoGPT-small, distilgpt2"
                return 1
            fi
            hf_download "$model"
            ;;
        *)
            echo "Usage: get_model <platform> <model_name>"
            echo "Platforms: ollama, hf"
            return 1
            ;;
    esac
}

# Development helper - code review with AI
code_review() {
    local file="$1"
    
    if [ -z "$file" ] || [ ! -f "$file" ]; then
        echo "Usage: code_review <file.py|file.js|etc>"
        return 1
    fi
    
    local file_content=$(cat "$file")
    local prompt="Please review this code for best practices, potential issues, and suggestions for improvement:\n\n$file_content"
    
    if command -v ollama >/dev/null 2>&1; then
        echo "🧑‍💻 Code Review with Ollama (CodeLlama):"
        ollama_chat codellama "$prompt"
    elif python3 -c "import transformers" 2>/dev/null; then
        echo "🧑‍💻 Code Review with Hugging Face:"
        hf_generate "$prompt"
    else
        echo "❌ Code review requires AI tools. Run: make ai-setup"
        return 1
    fi
}

# AI-powered commit message generator
ai_commit() {
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        echo "❌ Not in a git repository"
        return 1
    fi
    
    local diff=$(git diff --cached)
    if [ -z "$diff" ]; then
        echo "❌ No staged changes found. Stage some changes first with 'git add'"
        return 1
    fi
    
    local prompt="Based on this git diff, write a concise commit message that follows conventional commit format:\n\n$diff"
    
    echo "🤖 Generating commit message..."
    if command -v ollama >/dev/null 2>&1; then
        ollama_chat llama3.2 "$prompt"
    elif python3 -c "import transformers" 2>/dev/null; then
        hf_generate "$prompt"
    else
        echo "❌ AI commit requires AI tools. Run: make ai-setup"
        return 1
    fi
}

# Export functions for use in other scripts
export -f ai ai_help sentiment summarize get_model code_review ai_commit