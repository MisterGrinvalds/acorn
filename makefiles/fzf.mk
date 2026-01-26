# FZF Management
# Fuzzy finder for command line

.PHONY: fzf-install fzf-setup fzf-status fzf-test

# Installation
fzf-install: ## Install fzf
	@if command -v fzf >/dev/null 2>&1; then \
		echo "fzf already installed: $$(fzf --version)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install fzf; \
	else \
		git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf; \
		~/.fzf/install --bin; \
	fi

fzf-setup: fzf-install ## Setup fzf environment
	@if command -v fzf >/dev/null 2>&1; then \
		mkdir -p "$${XDG_CONFIG_HOME:-$$HOME/.config}/fzf"; \
		echo "fzf config directory created"; \
	fi

# Status
fzf-status: ## Check fzf installation status
	@echo "FZF Status"
	@echo "=========="
	@echo ""
	@if command -v fzf >/dev/null 2>&1; then \
		echo "Version: $$(fzf --version | awk '{print $$1}')"; \
		echo "Path: $$(which fzf)"; \
	else \
		echo "fzf not installed. Run: make fzf-install"; \
	fi

# Test
fzf-test: ## Test fzf functionality
	@echo "Testing fzf..."
	@command -v fzf >/dev/null || (echo "fzf not installed"; exit 1)
	@echo "test" | fzf --filter="test" >/dev/null 2>&1 && echo "fzf test passed" || (echo "fzf test failed"; exit 1)
