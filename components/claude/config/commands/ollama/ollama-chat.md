---
description: Help with Ollama chat interactions and prompting
argument-hint: [model] [prompt]
allowed-tools: Read, Bash
---

## Task

Help the user interact with Ollama models through chat sessions, one-shot prompts, and code generation.

## Quick Chat Methods

### Interactive Session
```bash
# Using dotfiles function
ollama_run llama3.2

# Direct command
ollama run llama3.2

# With specific size
ollama run llama3.2:3b
```

In interactive mode:
- Type your message and press Enter
- Multi-line: end lines with `\`
- Exit: type `/bye` or Ctrl+D

### One-Shot Query
```bash
# Using dotfiles function
ollama_chat llama3.2 "What is the capital of France?"

# Using echo pipe
echo "Explain recursion" | ollama run llama3.2

# With file input
cat question.txt | ollama run llama3.2
```

### Code Generation
```bash
# Using dotfiles function
ollama_code python "function to check if a number is prime"
ollama_code javascript "async function to fetch JSON from API"
ollama_code bash "script to backup a directory"

# Direct with codellama
echo "Write a Go function to reverse a string" | ollama run codellama
```

## Prompting Techniques

### Clear Instructions
```bash
# Be specific
ollama_chat llama3.2 "List 5 benefits of exercise. Use bullet points."

# Request format
ollama_chat llama3.2 "Explain Docker in 3 sentences."
```

### System Prompts (Custom Models)
```bash
# Create custom model with system prompt
cat > Modelfile << 'EOF'
FROM llama3.2
SYSTEM "You are a helpful Linux system administrator. Give concise answers."
EOF
ollama create linux-admin -f Modelfile
ollama run linux-admin
```

### Code Review
```bash
# Review code from clipboard (macOS)
pbpaste | ollama run codellama "Review this code for bugs:"

# Review a file
cat myfile.py | ollama run codellama "Review this Python code:"
```

### Explain Code
```bash
cat script.sh | ollama run codellama "Explain what this script does:"
```

## Chat Commands (Interactive Mode)

| Command | Description |
|---------|-------------|
| `/bye` | Exit chat |
| `/clear` | Clear context |
| `/set parameter value` | Change parameter |
| `/show info` | Show model info |
| `/load model` | Switch model |
| `"""` | Start/end multi-line |

## API Usage

### Simple Query
```bash
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "prompt": "Hello!",
  "stream": false
}'
```

### Chat API (with history)
```bash
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [
    {"role": "user", "content": "Hello!"},
    {"role": "assistant", "content": "Hi! How can I help?"},
    {"role": "user", "content": "What is 2+2?"}
  ],
  "stream": false
}'
```

### Streaming Response
```bash
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "prompt": "Tell me a story",
  "stream": true
}'
```

## Useful Patterns

### Shell Script Helper
```bash
# Ask about a command
ollama_chat llama3.2 "How do I find files larger than 100MB in Linux?"

# Generate a one-liner
ollama_chat llama3.2 "Write a bash one-liner to count lines in all .py files"
```

### Git Commit Message
```bash
git diff --staged | ollama run llama3.2 "Write a concise git commit message for these changes:"
```

### Documentation
```bash
cat myfunction.py | ollama run codellama "Write docstring for this function:"
```

### Debug Help
```bash
echo "Error: ECONNREFUSED 127.0.0.1:5432" | ollama run llama3.2 "What does this error mean and how to fix it?"
```

## Model Selection for Chat

| Task | Recommended Model |
|------|------------------|
| Quick questions | llama3.2:3b |
| Complex reasoning | llama3.2, mixtral |
| Code questions | codellama, qwen2.5-coder |
| Creative writing | llama3.2, mistral |

## Context

@components/ollama/functions.sh

## Troubleshooting

### Model Not Found
```bash
# Pull the model first
ollama_pull llama3.2
```

### Service Not Running
```bash
ollama_start
# or
ollama serve &
```

### Slow Responses
- Use smaller model (3b instead of 7b)
- Check if GPU is being used
- Close other applications to free RAM
