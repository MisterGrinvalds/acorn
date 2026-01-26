# Git Management
# Version control system

.PHONY: git-install git-setup git-status git-test

# Installation
git-install: ## Install Git
	@if command -v git >/dev/null 2>&1; then \
		echo "Git already installed: $$(git --version)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install git; \
	else \
		sudo apt-get update && sudo apt-get install -y git || \
		sudo yum install -y git || \
		echo "Please install git manually"; \
	fi

git-setup: git-install ## Setup Git environment
	@if command -v git >/dev/null 2>&1; then \
		mkdir -p "$${XDG_CONFIG_HOME:-$$HOME/.config}/git"; \
		if [ ! -f "$${XDG_CONFIG_HOME:-$$HOME/.config}/git/config" ]; then \
			echo "Git config directory created. Configure with:"; \
			echo "  git config --global user.name 'Your Name'"; \
			echo "  git config --global user.email 'your@email.com'"; \
		fi; \
	fi

# Status
git-status: ## Check Git installation status
	@echo "Git Status"
	@echo "=========="
	@echo ""
	@if command -v git >/dev/null 2>&1; then \
		echo "Version: $$(git --version | awk '{print $$3}')"; \
		echo "Path: $$(which git)"; \
		echo ""; \
		if git config user.name >/dev/null 2>&1; then \
			echo "User: $$(git config user.name) <$$(git config user.email)>"; \
		else \
			echo "User: Not configured"; \
		fi; \
		echo "Default Branch: $$(git config init.defaultBranch || echo 'master')"; \
	else \
		echo "Git not installed. Run: make git-install"; \
	fi

# Test
git-test: ## Test Git functionality
	@echo "Testing Git..."
	@command -v git >/dev/null || (echo "Git not installed"; exit 1)
	@git --version >/dev/null 2>&1 && echo "Git test passed" || (echo "Git test failed"; exit 1)
