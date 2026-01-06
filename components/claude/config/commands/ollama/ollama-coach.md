---
description: Interactive coaching session to learn local AI with Ollama
argument-hint: [skill-level: beginner|intermediate|advanced]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Guide the user through learning Ollama and local AI model management interactively.

## Approach

1. **Assess level** - Ask about AI/LLM experience
2. **Check setup** - Verify Ollama installation
3. **Start small** - Begin with lightweight models
4. **Progressive exercises** - Build from basic to advanced
5. **Real-time practice** - Have them run actual commands

## Skill Levels

### Beginner
- What is local AI?
- Installing Ollama
- Pulling your first model
- Basic chat interaction
- Understanding model sizes

### Intermediate
- Different model types
- Code generation with CodeLlama
- Using the API
- Service management
- Multiple models workflow

### Advanced
- Custom modelfiles
- Fine-tuning considerations
- GPU optimization
- Embedding models
- Production deployment

## Interactive Exercises

### Beginner Exercises
```bash
# Exercise 1: Check if Ollama is installed
ollama_status

# Exercise 2: Install Ollama (if needed)
ollama_install

# Exercise 3: Start the service
ollama_start

# Exercise 4: Pull a small model
ollama_pull llama3.2:3b

# Exercise 5: Try a simple chat
ollama_chat llama3.2:3b "What is 2+2?"

# Exercise 6: Interactive session
ollama_run llama3.2:3b
# Type: "Hello, who are you?"
# Type: /bye to exit
```

### Intermediate Exercises
```bash
# Exercise 7: Pull a code model
ollama_pull codellama:7b

# Exercise 8: Generate code
ollama_code python "function to reverse a string"

# Exercise 9: Use the API
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2:3b",
  "prompt": "Explain recursion",
  "stream": false
}'

# Exercise 10: List and manage models
ollama_models
ollama_remove llama3.2:3b  # Clean up
```

### Advanced Exercises
```bash
# Exercise 11: Create a custom modelfile
cat > Modelfile << 'EOF'
FROM llama3.2:3b
PARAMETER temperature 0.7
PARAMETER top_p 0.9
SYSTEM "You are a helpful coding assistant."
EOF
ollama create coding-assistant -f Modelfile

# Exercise 12: Try an embedding model
ollama_pull nomic-embed-text
curl http://localhost:11434/api/embeddings -d '{
  "model": "nomic-embed-text",
  "prompt": "Hello world"
}'

# Exercise 13: Compare model responses
echo "Explain Docker" | ollama run llama3.2:3b
echo "Explain Docker" | ollama run mistral
```

## Model Selection Guide

| Use Case | Recommended Model | Size |
|----------|------------------|------|
| Quick tasks | llama3.2:3b, phi3 | ~2GB |
| General use | llama3.2, mistral | ~4GB |
| Coding | codellama:7b, qwen2.5-coder | ~4GB |
| Complex tasks | mixtral, llama3.2:70b | 8-40GB |

## Context

@components/ollama/functions.sh
@components/ollama/aliases.sh

## Coaching Style

- Start with status check (`ollama_status`)
- Use small models first (3b) for fast feedback
- Show dotfiles functions for convenience
- Explain resource usage as you go
- Build toward practical automation
- Encourage experimentation with different models
