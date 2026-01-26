# Neovim Management
# Modern Vim-based text editor

.PHONY: neovim-install nvim-install neovim-setup neovim-status neovim-test

# Installation
neovim-install nvim-install: ## Install Neovim
	@if command -v nvim >/dev/null 2>&1; then \
		echo "Neovim already installed: $$(nvim --version | head -1)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install neovim; \
	else \
		sudo apt-get update && sudo apt-get install -y neovim || \
		sudo yum install -y neovim || \
		echo "Please install neovim manually"; \
	fi

neovim-setup: neovim-install ## Setup Neovim environment
	@if command -v nvim >/dev/null 2>&1; then \
		mkdir -p "$${XDG_CONFIG_HOME:-$$HOME/.config}/nvim"; \
		mkdir -p "$${XDG_DATA_HOME:-$$HOME/.local/share}/nvim"; \
		mkdir -p "$${XDG_STATE_HOME:-$$HOME/.local/state}/nvim"; \
		mkdir -p "$${XDG_CACHE_HOME:-$$HOME/.cache}/nvim"; \
		echo "Neovim directories created"; \
	fi

# Status
neovim-status: ## Check Neovim installation status
	@echo "Neovim Status"
	@echo "============="
	@echo ""
	@if command -v nvim >/dev/null 2>&1; then \
		echo "Version: $$(nvim --version | head -1 | awk '{print $$2}')"; \
		echo "Path: $$(which nvim)"; \
		echo ""; \
		echo "Config: $${XDG_CONFIG_HOME:-$$HOME/.config}/nvim"; \
		if [ -f "$${XDG_CONFIG_HOME:-$$HOME/.config}/nvim/init.lua" ] || [ -f "$${XDG_CONFIG_HOME:-$$HOME/.config}/nvim/init.vim" ]; then \
			echo "Config file: Found"; \
		else \
			echo "Config file: Not found"; \
		fi; \
	else \
		echo "Neovim not installed. Run: make neovim-install"; \
	fi

# Test
neovim-test: ## Test Neovim functionality
	@echo "Testing Neovim..."
	@command -v nvim >/dev/null || (echo "Neovim not installed"; exit 1)
	@nvim --version >/dev/null 2>&1 && echo "Neovim test passed" || (echo "Neovim test failed"; exit 1)
