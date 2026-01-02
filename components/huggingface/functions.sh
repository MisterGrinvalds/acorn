#!/bin/sh
# components/huggingface/functions.sh - Hugging Face local model management

# =============================================================================
# Setup
# =============================================================================

# Setup Hugging Face environment
hf_setup() {
    echo "Setting up Hugging Face environment..."

    if [ -z "$VIRTUAL_ENV" ]; then
        echo "Creating virtual environment for Hugging Face..."
        if command -v mkvenv >/dev/null 2>&1; then
            mkvenv huggingface
            venv huggingface
        else
            python3 -m venv ~/.venvs/huggingface
            . ~/.venvs/huggingface/bin/activate
        fi
    fi

    echo "Installing Hugging Face packages..."
    pip install --upgrade pip
    pip install transformers torch torchvision torchaudio
    pip install datasets tokenizers accelerate
    pip install huggingface_hub

    echo "Hugging Face environment setup complete!"
}

# =============================================================================
# Status and Info
# =============================================================================

# Check model status and cache
hf_status() {
    echo "Hugging Face Status"
    echo "==================="

    if python3 -c "import transformers" 2>/dev/null; then
        echo "Transformers: installed"
        python3 -c "import transformers; print(f'  Version: {transformers.__version__}')"
    else
        echo "Transformers: not installed"
        echo "  Run: hf_setup"
        return 1
    fi

    if python3 -c "import torch" 2>/dev/null; then
        echo "PyTorch: installed"
        python3 -c "import torch; print(f'  Version: {torch.__version__}')"
    else
        echo "PyTorch: not installed"
    fi

    echo ""
    echo "Model cache:"
    if [ -d "$HF_HOME" ]; then
        echo "  Location: $HF_HOME"
        du -sh "$HF_HOME" 2>/dev/null | head -1
    else
        echo "  No cache directory found"
    fi

    echo ""
    if [ -n "$VIRTUAL_ENV" ]; then
        echo "Virtual environment: $(basename "$VIRTUAL_ENV")"
    else
        echo "Virtual environment: none active"
    fi
}

# List popular models
hf_models() {
    echo "Popular Hugging Face Models"
    echo "==========================="
    echo ""
    echo "Text Generation:"
    echo "  microsoft/DialoGPT-small       # Conversational AI (117M)"
    echo "  gpt2                           # GPT-2 (124M)"
    echo "  distilgpt2                     # Distilled GPT-2 (82M)"
    echo ""
    echo "Language Understanding:"
    echo "  distilbert-base-uncased        # Efficient BERT (66M)"
    echo "  bert-base-uncased              # BERT (110M)"
    echo ""
    echo "Specialized:"
    echo "  microsoft/codebert-base        # Code understanding"
    echo "  facebook/bart-base             # Text summarization"
    echo ""
    echo "Usage: Download models using transformers library"
}

# =============================================================================
# Cache Management
# =============================================================================

# Clear model cache
hf_clear_cache() {
    echo "Hugging Face Model Cache"
    echo "========================"

    if [ -d "$HF_HOME" ]; then
        echo "Cache location: $HF_HOME"
        du -sh "$HF_HOME" 2>/dev/null | head -1

        printf "Clear cache? (y/N): "
        read -r reply
        case "$reply" in
            [Yy]*)
                rm -rf "$HF_HOME"
                echo "Cache cleared"
                ;;
            *)
                echo "Cache clear cancelled"
                ;;
        esac
    else
        echo "No cache directory found"
    fi
}

# =============================================================================
# Pipeline Help
# =============================================================================

# List available pipelines
hf_pipelines() {
    echo "Hugging Face Pipelines"
    echo "======================"
    echo ""
    echo "Common pipeline tasks:"
    echo "  text-generation     # Generate text continuations"
    echo "  summarization       # Summarize long text"
    echo "  sentiment-analysis  # Analyze text sentiment"
    echo "  question-answering  # Answer questions about text"
    echo "  translation         # Translate between languages"
    echo "  fill-mask          # Fill in masked words"
    echo "  text-classification # Classify text categories"
    echo ""
    echo "Example Python usage:"
    echo "  from transformers import pipeline"
    echo "  gen = pipeline('text-generation', model='gpt2')"
    echo "  gen('Hello, I am')"
    echo ""
    echo "Management:"
    echo "  hf_models       - List popular models"
    echo "  hf_status       - Check installation"
    echo "  hf_clear_cache  - Clear model cache"
    echo "  hf_setup        - Setup environment"
}
