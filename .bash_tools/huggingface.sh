#!/bin/sh
# Hugging Face Local Model Management
# Provides functions for running Hugging Face models locally using transformers

# Setup Hugging Face environment
hf_setup() {
    echo "🤗 Setting up Hugging Face environment..."
    
    # Check if we're in a virtual environment
    if [ -z "$VIRTUAL_ENV" ]; then
        echo "📦 Creating virtual environment for Hugging Face..."
        if command -v mkvenv >/dev/null 2>&1; then
            mkvenv huggingface
            venv huggingface
        else
            python3 -m venv ~/.venvs/huggingface
            source ~/.venvs/huggingface/bin/activate
        fi
    fi
    
    echo "📦 Installing Hugging Face packages..."
    pip install --upgrade pip
    pip install transformers torch torchvision torchaudio
    pip install datasets tokenizers accelerate
    pip install huggingface_hub
    
    # Optional: Install optimized packages
    echo "🚀 Installing performance optimizations..."
    pip install optimum[onnxruntime]
    
    echo "✅ Hugging Face environment setup complete!"
    echo ""
    echo "🚀 Try these commands:"
    echo "  hf_models                    # List popular models"
    echo "  hf_download microsoft/DialoGPT-small  # Download model"
    echo "  hf_generate 'Hello world'    # Generate text"
    echo "  hf_chat                      # Interactive chat"
}

# List popular models
hf_models() {
    echo "🤗 Popular Hugging Face Models"
    echo "=============================="
    echo ""
    echo "📝 Text Generation:"
    echo "  microsoft/DialoGPT-small       # Conversational AI (117M)"
    echo "  microsoft/DialoGPT-medium      # Conversational AI (345M)"
    echo "  gpt2                           # GPT-2 (124M)"
    echo "  distilgpt2                     # Distilled GPT-2 (82M)"
    echo "  EleutherAI/gpt-neo-125m        # GPT-Neo (125M)"
    echo ""
    echo "🧠 Language Understanding:"
    echo "  distilbert-base-uncased        # Efficient BERT (66M)"
    echo "  bert-base-uncased              # BERT (110M)"
    echo "  roberta-base                   # RoBERTa (125M)"
    echo ""
    echo "🌍 Multilingual:"
    echo "  distilbert-base-multilingual-cased  # Multilingual BERT"
    echo "  xlm-roberta-base               # Multilingual RoBERTa"
    echo ""
    echo "💡 Specialized:"
    echo "  microsoft/codebert-base        # Code understanding"
    echo "  facebook/bart-base             # Text summarization"
    echo "  t5-small                       # Text-to-text generation"
    echo ""
    echo "📥 Usage: hf_download <model_name>"
}

# Download and cache a model
hf_download() {
    local model="$1"
    if [ -z "$model" ]; then
        echo "❌ Usage: hf_download <model_name>"
        echo "Example: hf_download microsoft/DialoGPT-small"
        return 1
    fi
    
    echo "📥 Downloading model: $model"
    python3 -c "
from transformers import AutoTokenizer, AutoModel
import sys

try:
    print('📦 Downloading tokenizer...')
    tokenizer = AutoTokenizer.from_pretrained('$model')
    print('📦 Downloading model...')
    model = AutoModel.from_pretrained('$model')
    print('✅ Model downloaded and cached successfully!')
    print(f'📍 Model info: {model.config}')
except Exception as e:
    print(f'❌ Error downloading model: {e}')
    sys.exit(1)
"
}

# Generate text with a model
hf_generate() {
    local prompt="$1"
    local model="${2:-microsoft/DialoGPT-small}"
    
    if [ -z "$prompt" ]; then
        echo "❌ Usage: hf_generate <prompt> [model_name]"
        echo "Example: hf_generate 'Hello, how are you?' microsoft/DialoGPT-small"
        return 1
    fi
    
    echo "🤖 Generating text with $model..."
    echo "💬 Prompt: $prompt"
    echo ""
    
    python3 -c "
from transformers import AutoTokenizer, AutoModelForCausalLM
import torch
import sys

try:
    # Load model and tokenizer
    tokenizer = AutoTokenizer.from_pretrained('$model')
    model = AutoModelForCausalLM.from_pretrained('$model')
    
    # Add pad token if it doesn't exist
    if tokenizer.pad_token is None:
        tokenizer.pad_token = tokenizer.eos_token
    
    # Encode input
    inputs = tokenizer.encode('$prompt', return_tensors='pt')
    
    # Generate response
    with torch.no_grad():
        outputs = model.generate(
            inputs, 
            max_length=inputs.shape[1] + 50,
            num_return_sequences=1,
            temperature=0.7,
            do_sample=True,
            pad_token_id=tokenizer.eos_token_id
        )
    
    # Decode response
    response = tokenizer.decode(outputs[0], skip_special_tokens=True)
    
    # Extract just the generated part
    generated = response[len('$prompt'):].strip()
    print('🤖 Response:', generated)
    
except Exception as e:
    print(f'❌ Error generating text: {e}')
    sys.exit(1)
"
}

# Interactive chat with a model
hf_chat() {
    local model="${1:-microsoft/DialoGPT-small}"
    
    echo "🤖 Starting interactive chat with $model"
    echo "Type 'exit' or 'quit' to end the conversation"
    echo "============================================"
    
    python3 -c "
from transformers import AutoTokenizer, AutoModelForCausalLM
import torch
import sys

try:
    print('📦 Loading model and tokenizer...')
    tokenizer = AutoTokenizer.from_pretrained('$model')
    model = AutoModelForCausalLM.from_pretrained('$model')
    
    if tokenizer.pad_token is None:
        tokenizer.pad_token = tokenizer.eos_token
    
    print('✅ Model loaded! Ready to chat.')
    print()
    
    chat_history = []
    
    while True:
        try:
            user_input = input('You: ')
            if user_input.lower() in ['exit', 'quit', 'bye']:
                print('👋 Goodbye!')
                break
            
            # Add user input to history
            chat_history.append(user_input)
            
            # Create context from recent history
            context = ' '.join(chat_history[-5:])  # Last 5 exchanges
            
            # Encode and generate
            inputs = tokenizer.encode(context, return_tensors='pt')
            
            with torch.no_grad():
                outputs = model.generate(
                    inputs,
                    max_length=inputs.shape[1] + 50,
                    num_return_sequences=1,
                    temperature=0.7,
                    do_sample=True,
                    pad_token_id=tokenizer.eos_token_id
                )
            
            response = tokenizer.decode(outputs[0], skip_special_tokens=True)
            generated = response[len(context):].strip()
            
            print(f'Bot: {generated}')
            chat_history.append(generated)
            
        except KeyboardInterrupt:
            print('\n👋 Chat interrupted. Goodbye!')
            break
        except Exception as e:
            print(f'❌ Error: {e}')
            continue

except Exception as e:
    print(f'❌ Error loading model: {e}')
    sys.exit(1)
"
}

# Summarize text
hf_summarize() {
    local text="$1"
    local model="${2:-facebook/bart-large-cnn}"
    
    if [ -z "$text" ]; then
        echo "❌ Usage: hf_summarize <text> [model_name]"
        echo "Example: hf_summarize 'Long text to summarize...'"
        return 1
    fi
    
    echo "📄 Summarizing text with $model..."
    
    python3 -c "
from transformers import pipeline
import sys

try:
    summarizer = pipeline('summarization', model='$model')
    
    text = '''$text'''
    summary = summarizer(text, max_length=130, min_length=30, do_sample=False)
    
    print('📝 Summary:', summary[0]['summary_text'])
    
except Exception as e:
    print(f'❌ Error summarizing: {e}')
    sys.exit(1)
"
}

# Sentiment analysis
hf_sentiment() {
    local text="$1"
    local model="${2:-distilbert-base-uncased-finetuned-sst-2-english}"
    
    if [ -z "$text" ]; then
        echo "❌ Usage: hf_sentiment <text> [model_name]"
        echo "Example: hf_sentiment 'I love this product!'"
        return 1
    fi
    
    echo "😊 Analyzing sentiment with $model..."
    
    python3 -c "
from transformers import pipeline
import sys

try:
    classifier = pipeline('sentiment-analysis', model='$model')
    
    result = classifier('$text')
    
    print(f'😊 Sentiment: {result[0][\"label\"]} (confidence: {result[0][\"score\"]:.2f})')
    
except Exception as e:
    print(f'❌ Error analyzing sentiment: {e}')
    sys.exit(1)
"
}

# Question answering
hf_qa() {
    local question="$1"
    local context="$2"
    local model="${3:-distilbert-base-cased-distilled-squad}"
    
    if [ -z "$question" ] || [ -z "$context" ]; then
        echo "❌ Usage: hf_qa <question> <context> [model_name]"
        echo "Example: hf_qa 'What is Python?' 'Python is a programming language...'"
        return 1
    fi
    
    echo "❓ Answering question with $model..."
    
    python3 -c "
from transformers import pipeline
import sys

try:
    qa_pipeline = pipeline('question-answering', model='$model')
    
    result = qa_pipeline(question='$question', context='$context')
    
    print(f'❓ Question: $question')
    print(f'📝 Answer: {result[\"answer\"]} (confidence: {result[\"score\"]:.2f})')
    
except Exception as e:
    print(f'❌ Error answering question: {e}')
    sys.exit(1)
"
}

# Code completion
hf_code() {
    local code_prompt="$1"
    local model="${2:-microsoft/CodeGPT-small-py}"
    
    if [ -z "$code_prompt" ]; then
        echo "❌ Usage: hf_code <code_prompt> [model_name]"
        echo "Example: hf_code 'def fibonacci(n):'"
        return 1
    fi
    
    echo "🧑‍💻 Generating code with $model..."
    
    python3 -c "
from transformers import AutoTokenizer, AutoModelForCausalLM
import torch
import sys

try:
    tokenizer = AutoTokenizer.from_pretrained('$model')
    model = AutoModelForCausalLM.from_pretrained('$model')
    
    if tokenizer.pad_token is None:
        tokenizer.pad_token = tokenizer.eos_token
    
    inputs = tokenizer.encode('$code_prompt', return_tensors='pt')
    
    with torch.no_grad():
        outputs = model.generate(
            inputs,
            max_length=inputs.shape[1] + 100,
            num_return_sequences=1,
            temperature=0.3,
            do_sample=True,
            pad_token_id=tokenizer.eos_token_id
        )
    
    response = tokenizer.decode(outputs[0], skip_special_tokens=True)
    
    print('🧑‍💻 Generated code:')
    print('```python')
    print(response)
    print('```')
    
except Exception as e:
    print(f'❌ Error generating code: {e}')
    # Try alternative approach
    try:
        from transformers import pipeline
        generator = pipeline('text-generation', model='gpt2')
        result = generator('$code_prompt', max_length=150, temperature=0.3)
        print('🧑‍💻 Generated code (alternative):')
        print(result[0]['generated_text'])
    except:
        print('❌ Code generation failed')
        sys.exit(1)
"
}

# Check model status and cache
hf_status() {
    echo "🤗 Hugging Face Status"
    echo "====================="
    
    # Check if transformers is installed
    if python3 -c "import transformers" 2>/dev/null; then
        echo "✅ Transformers installed"
        python3 -c "import transformers; print(f'   Version: {transformers.__version__}')"
    else
        echo "❌ Transformers not installed"
        echo "   Run: hf_setup"
        return 1
    fi
    
    # Check if torch is installed
    if python3 -c "import torch" 2>/dev/null; then
        echo "✅ PyTorch installed"
        python3 -c "import torch; print(f'   Version: {torch.__version__}')"
    else
        echo "❌ PyTorch not installed"
    fi
    
    # Check cache directory
    echo ""
    echo "📁 Model cache:"
    if [ -d ~/.cache/huggingface ]; then
        echo "   Location: ~/.cache/huggingface"
        if command -v du >/dev/null 2>&1; then
            echo "   Size: $(du -sh ~/.cache/huggingface 2>/dev/null | cut -f1)"
        fi
        echo "   Models:"
        find ~/.cache/huggingface -name "config.json" 2>/dev/null | head -5 | while read -r config; do
            model_dir=$(dirname "$config")
            model_name=$(basename "$(dirname "$model_dir")")/$(basename "$model_dir")
            echo "     - $model_name"
        done
    else
        echo "   No cache directory found"
    fi
    
    # Check virtual environment
    echo ""
    if [ -n "$VIRTUAL_ENV" ]; then
        echo "✅ Virtual environment active: $(basename "$VIRTUAL_ENV")"
    else
        echo "⚠️ No virtual environment active"
        echo "   Consider running: hf_setup"
    fi
}

# Clear model cache
hf_clear_cache() {
    echo "🗑️ Clearing Hugging Face model cache..."
    
    if [ -d ~/.cache/huggingface ]; then
        echo "📁 Cache location: ~/.cache/huggingface"
        if command -v du >/dev/null 2>&1; then
            echo "💾 Current size: $(du -sh ~/.cache/huggingface 2>/dev/null | cut -f1)"
        fi
        
        read -p "Are you sure you want to clear the cache? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf ~/.cache/huggingface
            echo "✅ Cache cleared"
        else
            echo "❌ Cache clear cancelled"
        fi
    else
        echo "ℹ️ No cache directory found"
    fi
}

# List available pipelines
hf_pipelines() {
    echo "🤗 Available Hugging Face Pipelines"
    echo "==================================="
    echo ""
    echo "📝 Text Generation:"
    echo "  hf_generate <prompt>           # Generate text"
    echo "  hf_chat [model]                # Interactive chat"
    echo ""
    echo "📄 Text Analysis:"
    echo "  hf_summarize <text>            # Summarize text"
    echo "  hf_sentiment <text>            # Sentiment analysis"
    echo "  hf_qa <question> <context>     # Question answering"
    echo ""
    echo "🧑‍💻 Code:"
    echo "  hf_code <code_prompt>          # Code completion"
    echo ""
    echo "🛠️ Management:"
    echo "  hf_models                      # List popular models"
    echo "  hf_download <model>            # Download model"
    echo "  hf_status                      # Check installation"
    echo "  hf_clear_cache                 # Clear model cache"
    echo "  hf_setup                       # Setup environment"
}

# Quick example runner
hf_examples() {
    echo "🎯 Hugging Face Examples"
    echo "======================="
    echo ""
    echo "Running quick examples..."
    echo ""
    
    # Check if environment is ready
    if ! python3 -c "import transformers" 2>/dev/null; then
        echo "❌ Transformers not available. Run: hf_setup"
        return 1
    fi
    
    echo "1. 😊 Sentiment Analysis:"
    hf_sentiment "I love working with AI models!"
    echo ""
    
    echo "2. 📝 Text Generation:"
    hf_generate "The future of AI is" "distilgpt2"
    echo ""
    
    echo "3. ❓ Question Answering:"
    hf_qa "What is machine learning?" "Machine learning is a subset of artificial intelligence that focuses on algorithms that can learn from data."
    echo ""
    
    echo "✅ Examples complete!"
}

# Aliases for convenience
alias hf-setup='hf_setup'
alias hf-status='hf_status'
alias hf-models='hf_models'
alias hf-chat='hf_chat'
alias hf-examples='hf_examples'
alias hf-pipelines='hf_pipelines'