# AI/ML Management
# Ollama and Hugging Face tooling via acorn CLI

.PHONY: ai-setup ai-status ai-test ai-models ai-chat ai-chat-ollama ai-chat-hf
.PHONY: ai-benchmark ai-examples ai-cleanup
.PHONY: ollama-install ollama-status ollama-models ollama-start ollama-stop
.PHONY: hf-setup hf-status hf-models hf-clear-cache
.PHONY: test-ai test-ai-tools test-ai-go test-ai-components test-ai-json
.PHONY: test-ai-component-claude test-ai-component-ollama test-ai-component-huggingface

# ─────────────────────────────────────────────────────────────────────────────
# Setup and Status
# ─────────────────────────────────────────────────────────────────────────────

ai-setup: ollama-install hf-setup ## Setup AI/ML environment (Ollama + Hugging Face)

ai-status: ## Check AI/ML tools and models status
	@echo "AI/ML Status"
	@echo "============"
	@echo ""
	@echo "Python:"
	@command -v python3 >/dev/null && echo "  Version: $$(python3 --version)" || echo "  Not installed"
	@echo ""
	@echo "Ollama:"
	@command -v ollama >/dev/null && echo "  Installed: $$(ollama --version 2>/dev/null || echo 'yes')" || echo "  Not installed"
	@command -v ollama >/dev/null && pgrep -x ollama >/dev/null && echo "  Service: running" || echo "  Service: not running"
	@echo ""
	@echo "Hugging Face:"
	@python3 -c "import transformers; print('  transformers:', transformers.__version__)" 2>/dev/null || echo "  transformers: not installed"
	@python3 -c "import torch; print('  torch:', torch.__version__)" 2>/dev/null || echo "  torch: not installed"

ai-models: ## List all available AI models
	@echo "Ollama Models:"
	@command -v ollama >/dev/null && ollama list 2>/dev/null || echo "  Ollama not installed"
	@echo ""
	@echo "Popular models to pull:"
	@echo "  llama3.2, codellama, mistral, phi3, gemma2"

# ─────────────────────────────────────────────────────────────────────────────
# Chat Interfaces
# ─────────────────────────────────────────────────────────────────────────────

ai-chat: ## Start interactive AI chat (auto-detects best available)
	@if command -v ollama >/dev/null 2>&1; then \
		model=$$(ollama list 2>/dev/null | awk 'NR==2 {print $$1}'); \
		if [ -n "$$model" ]; then \
			echo "Starting chat with $$model..."; \
			ollama run "$$model"; \
		else \
			echo "No models installed. Run: make ollama-models"; \
		fi; \
	else \
		echo "Ollama not installed. Run: make ollama-install"; \
	fi

ai-chat-ollama: ## Start Ollama chat (specify MODEL=name or uses first available)
	@if [ -n "$(MODEL)" ]; then \
		ollama run $(MODEL); \
	else \
		model=$$(ollama list 2>/dev/null | awk 'NR==2 {print $$1}'); \
		if [ -n "$$model" ]; then ollama run "$$model"; \
		else echo "No models. Run: ollama pull llama3.2"; fi; \
	fi

ai-chat-hf: ## Start Hugging Face chat (requires transformers)
	@python3 -c "import transformers" 2>/dev/null || (echo "transformers not installed" && exit 1)
	@echo "Hugging Face chat requires a Python script. See: acorn ai huggingface examples"

# ─────────────────────────────────────────────────────────────────────────────
# Benchmarks and Examples
# ─────────────────────────────────────────────────────────────────────────────

ai-benchmark: ## Run AI performance benchmarks
	@mkdir -p $(LOG_DIR)
	@echo "AI Benchmark - $$(date)" > $(LOG_DIR)/ai-benchmark.log
	@echo "" >> $(LOG_DIR)/ai-benchmark.log
	@echo "System:" >> $(LOG_DIR)/ai-benchmark.log
	@uname -a >> $(LOG_DIR)/ai-benchmark.log
	@echo "" >> $(LOG_DIR)/ai-benchmark.log
	@echo "Ollama Models:" >> $(LOG_DIR)/ai-benchmark.log
	@command -v ollama >/dev/null && ollama list >> $(LOG_DIR)/ai-benchmark.log 2>&1 || echo "  not installed" >> $(LOG_DIR)/ai-benchmark.log
	@echo "Results: $(LOG_DIR)/ai-benchmark.log"

ai-examples: ## Show AI usage examples
	@echo "Ollama Examples:"
	@echo "  ollama pull llama3.2      # Download model"
	@echo "  ollama run llama3.2       # Interactive chat"
	@echo "  ollama list               # List models"
	@echo ""
	@echo "Acorn Commands:"
	@echo "  acorn ai ollama status    # Status"
	@echo "  acorn ai ollama models    # List models"
	@echo "  acorn ai ollama chat llama3.2 'Hello'  # Quick chat"
	@echo ""
	@echo "Shell Functions (after sourcing generated/shell/ollama.sh):"
	@echo "  ollama_status, ollama_models, ollama_pull, ollama_run"

ai-cleanup: ## Clean AI model caches and stop services
	@command -v ollama >/dev/null && pkill ollama 2>/dev/null && echo "Stopped Ollama" || true
	@[ -d ~/.cache/huggingface ] && echo "HF cache: ~/.cache/huggingface ($$(du -sh ~/.cache/huggingface 2>/dev/null | cut -f1))" || true

# ─────────────────────────────────────────────────────────────────────────────
# Ollama Targets
# ─────────────────────────────────────────────────────────────────────────────

ollama-install: ## Install Ollama
	@if command -v ollama >/dev/null 2>&1; then \
		echo "Ollama already installed"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install ollama; \
	else \
		curl -fsSL https://ollama.com/install.sh | sh; \
	fi

ollama-status: ## Check Ollama status
	@./bin/acorn ai ollama status 2>/dev/null || acorn ai ollama status 2>/dev/null || (command -v ollama >/dev/null && ollama list || echo "Ollama not installed")

ollama-models: ## List Ollama models
	@command -v ollama >/dev/null && ollama list || echo "Ollama not installed"

ollama-start: ## Start Ollama service
	@if command -v ollama >/dev/null; then \
		ollama serve > /dev/null 2>&1 & \
		echo "Ollama service started"; \
	else \
		echo "Ollama not installed"; \
	fi

ollama-stop: ## Stop Ollama service
	@pkill ollama 2>/dev/null && echo "Ollama stopped" || echo "Ollama not running"

# ─────────────────────────────────────────────────────────────────────────────
# Hugging Face Targets
# ─────────────────────────────────────────────────────────────────────────────

hf-setup: ## Setup Hugging Face environment
	@if python3 -c "import transformers" 2>/dev/null; then \
		echo "Hugging Face already installed"; \
	else \
		echo "Installing transformers..."; \
		pip3 install transformers torch --quiet || echo "Install failed - try: pip3 install transformers torch"; \
	fi

hf-status: ## Check Hugging Face status
	@echo "Hugging Face Status"
	@echo "==================="
	@python3 -c "import transformers; print('transformers:', transformers.__version__)" 2>/dev/null || echo "transformers: not installed"
	@python3 -c "import torch; print('torch:', torch.__version__)" 2>/dev/null || echo "torch: not installed"
	@[ -d ~/.cache/huggingface ] && echo "Cache: $$(du -sh ~/.cache/huggingface 2>/dev/null | cut -f1)" || echo "Cache: empty"

hf-models: ## List popular Hugging Face models
	@echo "Popular Hugging Face Models:"
	@echo "  Text Generation: gpt2, facebook/opt-350m, EleutherAI/gpt-neo-125m"
	@echo "  Sentiment: distilbert-base-uncased-finetuned-sst-2-english"
	@echo "  Summarization: facebook/bart-large-cnn"
	@echo "  Translation: Helsinki-NLP/opus-mt-en-de"

hf-clear-cache: ## Clear Hugging Face model cache
	@if [ -d ~/.cache/huggingface ]; then \
		size=$$(du -sh ~/.cache/huggingface 2>/dev/null | cut -f1); \
		read -p "Delete $$size of cached models? [y/N] " confirm; \
		if [ "$$confirm" = "y" ]; then rm -rf ~/.cache/huggingface && echo "Cache cleared"; fi; \
	else \
		echo "No cache to clear"; \
	fi

# ─────────────────────────────────────────────────────────────────────────────
# Testing
# ─────────────────────────────────────────────────────────────────────────────

AI_COMPONENTS := claude ollama huggingface

test-ai: test-ai-tools test-ai-go test-ai-components ## Run all AI tests

test-ai-tools: ## Check if AI/ML tools are available
	@echo "AI Tool Availability:"
	@printf "  %-15s" "ollama:"; command -v ollama >/dev/null && echo "installed ($$(ollama --version 2>/dev/null || echo 'ok'))" || echo "not installed"
	@printf "  %-15s" "python3:"; command -v python3 >/dev/null && echo "installed ($$(python3 --version 2>&1 | awk '{print $$2}'))" || echo "not installed"
	@printf "  %-15s" "transformers:"; python3 -c "import transformers; print('installed (' + transformers.__version__ + ')')" 2>/dev/null || echo "not installed"
	@printf "  %-15s" "torch:"; python3 -c "import torch; print('installed (' + torch.__version__ + ')')" 2>/dev/null || echo "not installed"

test-ai-go: ## Test AI Go code compiles
	@echo "Testing AI Go packages..."
	@go build -o /dev/null ./internal/components/ai/ollama/... 2>/dev/null && echo "  ollama package: ok" || echo "  ollama package: FAIL"
	@go build -o /dev/null ./internal/components/ai/huggingface/... 2>/dev/null && echo "  huggingface package: ok" || echo "  huggingface package: FAIL"
	@go build -o /dev/null ./internal/cmd/... 2>/dev/null && echo "  cmd package: ok" || echo "  cmd package: FAIL"

test-ai-components: ## Test all AI components with JSON output
	@echo "AI Component Configuration"
	@echo "=========================="
	@for comp in $(AI_COMPONENTS); do \
		$(MAKE) -s test-ai-component-$$comp; \
		echo ""; \
	done

test-ai-component-claude: ## Test claude component config
	@echo ""
	@echo "Component: claude"
	@echo "-----------------"
	@config=".sapling/config/claude/config.yaml"; \
	if [ ! -f "$$config" ]; then \
		echo '{"error": "config not found"}'; \
	else \
		echo "Config: $$config"; \
		echo ""; \
		yq -o json '{ \
			"name": .name, \
			"description": .description, \
			"env": (.env | keys), \
			"aliases": (.aliases | keys), \
			"wrappers": [.wrappers[].name], \
			"shell_functions": (.shell_functions | keys), \
			"sync_files": [.sync_files[] | {"source": .source, "target": .target, "mode": .mode}] \
		}' "$$config" 2>/dev/null || echo '{"error": "failed to parse"}'; \
	fi
	@echo ""
	@echo "Generated shell script:"
	@if [ -f "generated/shell/claude.sh" ]; then \
		size=$$(wc -c < "generated/shell/claude.sh" | tr -d ' '); \
		lines=$$(wc -l < "generated/shell/claude.sh" | tr -d ' '); \
		echo "  generated/shell/claude.sh: $$size bytes, $$lines lines"; \
		bash -n "generated/shell/claude.sh" 2>/dev/null && echo "  Syntax: valid" || echo "  Syntax: INVALID"; \
	else \
		echo "  generated/shell/claude.sh: NOT GENERATED"; \
	fi
	@echo ""
	@echo "Sync status:"
	@./bin/acorn ai claude sync --dry-run -o json 2>/dev/null | jq -c '{synced: (.synced | length), unchanged: (.skipped | length), errors: (.errors | length)}' 2>/dev/null || echo '{"error": "sync check failed"}'

test-ai-component-ollama: ## Test ollama component config
	@echo ""
	@echo "Component: ollama"
	@echo "-----------------"
	@config=".sapling/config/ollama/config.yaml"; \
	if [ ! -f "$$config" ]; then \
		echo '{"error": "config not found"}'; \
	else \
		echo "Config: $$config"; \
		echo ""; \
		yq -o json '{ \
			"name": .name, \
			"description": .description, \
			"env": (.env | keys), \
			"aliases": (.aliases | keys), \
			"wrappers": [.wrappers[].name], \
			"shell_functions": (.shell_functions | keys) \
		}' "$$config" 2>/dev/null || echo '{"error": "failed to parse"}'; \
	fi
	@echo ""
	@echo "Generated shell script:"
	@if [ -f "generated/shell/ollama.sh" ]; then \
		size=$$(wc -c < "generated/shell/ollama.sh" | tr -d ' '); \
		lines=$$(wc -l < "generated/shell/ollama.sh" | tr -d ' '); \
		echo "  generated/shell/ollama.sh: $$size bytes, $$lines lines"; \
		bash -n "generated/shell/ollama.sh" 2>/dev/null && echo "  Syntax: valid" || echo "  Syntax: INVALID"; \
	else \
		echo "  generated/shell/ollama.sh: NOT GENERATED"; \
	fi

test-ai-component-huggingface: ## Test huggingface component config
	@echo ""
	@echo "Component: huggingface"
	@echo "----------------------"
	@config=".sapling/config/huggingface/config.yaml"; \
	if [ ! -f "$$config" ]; then \
		echo '{"error": "config not found"}'; \
	else \
		echo "Config: $$config"; \
		echo ""; \
		yq -o json '{ \
			"name": .name, \
			"description": .description, \
			"env": (.env | keys), \
			"aliases": (.aliases | keys), \
			"wrappers": [.wrappers[].name], \
			"shell_functions": (.shell_functions | keys) \
		}' "$$config" 2>/dev/null || echo '{"error": "failed to parse"}'; \
	fi
	@echo ""
	@echo "Generated shell script:"
	@if [ -f "generated/shell/huggingface.sh" ]; then \
		size=$$(wc -c < "generated/shell/huggingface.sh" | tr -d ' '); \
		lines=$$(wc -l < "generated/shell/huggingface.sh" | tr -d ' '); \
		echo "  generated/shell/huggingface.sh: $$size bytes, $$lines lines"; \
		bash -n "generated/shell/huggingface.sh" 2>/dev/null && echo "  Syntax: valid" || echo "  Syntax: INVALID"; \
	else \
		echo "  generated/shell/huggingface.sh: NOT GENERATED"; \
	fi

test-ai-json: ## Output all AI components as single JSON document
	@echo '{'
	@echo '  "components": {'
	@for comp in $(AI_COMPONENTS); do \
		config=".sapling/config/$$comp/config.yaml"; \
		if [ -f "$$config" ]; then \
			printf '    "%s": ' "$$comp"; \
			yq -o json '{ \
				"description": .description, \
				"env": (.env | keys // []), \
				"aliases": (.aliases | keys // []), \
				"wrappers": ([.wrappers[]?.name] // []), \
				"shell_functions": (.shell_functions | keys // []), \
				"sync_files": ([.sync_files[]? | .mode] // []) \
			}' "$$config" 2>/dev/null | jq -c '.' || echo '{}'; \
			if [ "$$comp" != "huggingface" ]; then echo ","; else echo ""; fi; \
		fi; \
	done
	@echo '  },'
	@echo '  "generated_scripts": {'
	@for comp in $(AI_COMPONENTS); do \
		script="generated/shell/$$comp.sh"; \
		if [ -f "$$script" ]; then \
			size=$$(wc -c < "$$script" | tr -d ' '); \
			lines=$$(wc -l < "$$script" | tr -d ' '); \
			valid=$$(bash -n "$$script" 2>/dev/null && echo "true" || echo "false"); \
			printf '    "%s": {"bytes": %s, "lines": %s, "valid": %s}' "$$comp" "$$size" "$$lines" "$$valid"; \
		else \
			printf '    "%s": null' "$$comp"; \
		fi; \
		if [ "$$comp" != "huggingface" ]; then echo ","; else echo ""; fi; \
	done
	@echo '  }'
	@echo '}'
