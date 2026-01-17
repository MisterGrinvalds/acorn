# AI/ML Management
# Ollama, Hugging Face, and related tooling

.PHONY: ai-setup ai-status ai-test ai-models ai-chat ai-chat-ollama ai-chat-hf
.PHONY: ai-benchmark ai-examples ai-examples-run ai-cleanup
.PHONY: ollama-install ollama-setup ollama-start ollama-stop ollama-status ollama-models
.PHONY: hf-setup hf-status hf-models hf-examples hf-clear-cache
.PHONY: test-ai-module test-ai-tools

# Setup and status
ai-setup: ## Setup AI/ML environment (Ollama + Hugging Face)
	@bash -c 'source functions/ai/ollama.sh && ollama_setup'
	@bash -c 'source functions/ai/huggingface.sh && hf_setup'

ai-status: ## Check AI/ML tools and models status
	@echo "AI/ML Status"
	@echo "============"
	@command -v python3 >/dev/null && echo "Python 3: $$(python3 --version)" || echo "Python 3: not found"
	@command -v pip3 >/dev/null && echo "pip3: available" || echo "pip3: not found"
	@echo ""
	@echo "Ollama:"
	@bash -c 'source functions/ai/ollama.sh && ollama_status' 2>/dev/null || echo "  not available"
	@echo ""
	@echo "Hugging Face:"
	@bash -c 'source functions/ai/huggingface.sh && hf_status' 2>/dev/null || echo "  not available"

ai-test: ## Test AI/ML functionality
	@mkdir -p $(LOG_DIR)
	@find components/ollama components/huggingface -name "*.sh" -type f 2>/dev/null | while read file; do \
		bash -n "$$file" || exit 1; \
	done

# Model listing
ai-models: ## List all available AI models
	@echo "Ollama Models:"
	@bash -c 'source functions/ai/ollama.sh && ollama_models' 2>/dev/null || echo "  Ollama not installed"
	@echo ""
	@echo "Hugging Face Models:"
	@bash -c 'source functions/ai/huggingface.sh && hf_models'

# Chat interfaces
ai-chat: ## Start interactive AI chat (auto-detects best available)
	@if command -v ollama >/dev/null 2>&1; then \
		bash -c 'source functions/ai/ollama.sh && ollama_run llama3.2'; \
	elif python3 -c "import transformers" 2>/dev/null; then \
		bash -c 'source functions/ai/huggingface.sh && hf_chat'; \
	else \
		echo "No AI platforms available. Run: make ai-setup"; exit 1; \
	fi

ai-chat-ollama: ## Start Ollama chat with Llama 3.2
	@bash -c 'source functions/ai/ollama.sh && ollama_run llama3.2'

ai-chat-hf: ## Start Hugging Face chat
	@bash -c 'source functions/ai/huggingface.sh && hf_chat'

# Benchmarks and examples
ai-benchmark: ## Run AI performance benchmarks
	@mkdir -p $(LOG_DIR)
	@echo "AI Benchmark - $$(date)" > $(LOG_DIR)/ai-benchmark.log
	@command -v ollama >/dev/null && ollama list >> $(LOG_DIR)/ai-benchmark.log 2>&1 || true
	@echo "Results saved to: $(LOG_DIR)/ai-benchmark.log"

ai-examples: ## Show AI usage examples
	@echo "Quick Start:"
	@echo "  make ai-setup       # Install everything"
	@echo "  make ai-chat        # Start interactive chat"
	@echo ""
	@echo "Ollama:"
	@echo "  ollama_install, ollama_pull llama3.2, ollama_run llama3.2"
	@echo ""
	@echo "Hugging Face:"
	@echo "  hf_setup, hf_generate 'text', hf_sentiment 'text'"

ai-examples-run: ## Run live AI examples
	@command -v ollama >/dev/null 2>&1 && bash -c 'source functions/ai/ollama.sh && ollama_examples' || true
	@python3 -c "import transformers" 2>/dev/null && bash -c 'source functions/ai/huggingface.sh && hf_examples' || true

ai-cleanup: ## Clean AI model caches and stop services
	@command -v ollama >/dev/null 2>&1 && bash -c 'source functions/ai/ollama.sh && ollama_stop' || true
	@[ -d ~/.cache/huggingface ] && bash -c 'source functions/ai/huggingface.sh && hf_clear_cache' || true

# Ollama targets
ollama-install: ## Install Ollama
	@bash -c 'source functions/ai/ollama.sh && ollama_install'

ollama-setup: ## Setup Ollama with recommended models
	@bash -c 'source functions/ai/ollama.sh && ollama_setup'

ollama-start: ## Start Ollama service
	@bash -c 'source functions/ai/ollama.sh && ollama_start'

ollama-stop: ## Stop Ollama service
	@bash -c 'source functions/ai/ollama.sh && ollama_stop'

ollama-status: ## Check Ollama status
	@bash -c 'source functions/ai/ollama.sh && ollama_status'

ollama-models: ## List Ollama models
	@bash -c 'source functions/ai/ollama.sh && ollama_models'

# Hugging Face targets
hf-setup: ## Setup Hugging Face environment
	@bash -c 'source functions/ai/huggingface.sh && hf_setup'

hf-status: ## Check Hugging Face status
	@bash -c 'source functions/ai/huggingface.sh && hf_status'

hf-models: ## List popular Hugging Face models
	@bash -c 'source functions/ai/huggingface.sh && hf_models'

hf-examples: ## Run Hugging Face examples
	@bash -c 'source functions/ai/huggingface.sh && hf_examples'

hf-clear-cache: ## Clear Hugging Face model cache
	@bash -c 'source functions/ai/huggingface.sh && hf_clear_cache'

# Testing
test-ai-module: ## Test AI components
	@find components/ollama components/huggingface -name "*.sh" -type f 2>/dev/null | while read file; do \
		bash -n "$$file" || exit 1; \
	done

test-ai-tools: ## Check if AI/ML tools are available
	@command -v ollama >/dev/null && echo "ollama: installed" || echo "ollama: not installed"
	@python3 -c "import transformers" 2>/dev/null && echo "transformers: installed" || echo "transformers: not installed"
	@python3 -c "import torch" 2>/dev/null && echo "torch: installed" || echo "torch: not installed"
