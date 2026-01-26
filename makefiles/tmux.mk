# Tmux Management
# Terminal multiplexer

.PHONY: tmux-install tmux-setup tmux-status tmux-test

# Installation
tmux-install: ## Install tmux
	@if command -v tmux >/dev/null 2>&1; then \
		echo "tmux already installed: $$(tmux -V)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install tmux; \
	else \
		sudo apt-get update && sudo apt-get install -y tmux || \
		sudo yum install -y tmux || \
		echo "Please install tmux manually"; \
	fi

tmux-setup: tmux-install ## Setup tmux environment
	@if command -v tmux >/dev/null 2>&1; then \
		mkdir -p "$${XDG_CONFIG_HOME:-$$HOME/.config}/tmux"; \
		mkdir -p "$${XDG_CONFIG_HOME:-$$HOME/.config}/tmux/plugins"; \
		echo "tmux config directory created"; \
	fi

# Status
tmux-status: ## Check tmux installation status
	@echo "Tmux Status"
	@echo "==========="
	@echo ""
	@if command -v tmux >/dev/null 2>&1; then \
		echo "Version: $$(tmux -V | awk '{print $$2}')"; \
		echo "Path: $$(which tmux)"; \
		echo ""; \
		echo "Sessions: $$(tmux list-sessions 2>/dev/null | wc -l | tr -d ' ')"; \
		tmux list-sessions 2>/dev/null | head -5 || echo "  (none)"; \
	else \
		echo "tmux not installed. Run: make tmux-install"; \
	fi

# Test
tmux-test: ## Test tmux functionality
	@echo "Testing tmux..."
	@command -v tmux >/dev/null || (echo "tmux not installed"; exit 1)
	@tmux -V >/dev/null 2>&1 && echo "tmux test passed" || (echo "tmux test failed"; exit 1)
