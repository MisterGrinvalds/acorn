# GitHub CLI Management
# Official GitHub command line tool

.PHONY: github-install gh-install github-setup github-status github-test

# Installation
github-install gh-install: ## Install GitHub CLI (gh)
	@if command -v gh >/dev/null 2>&1; then \
		echo "GitHub CLI already installed: $$(gh --version | head -1)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install gh; \
	else \
		curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg; \
		sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg; \
		echo "deb [arch=$$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null; \
		sudo apt update; \
		sudo apt install gh -y; \
	fi

github-setup: github-install ## Setup GitHub CLI
	@if command -v gh >/dev/null 2>&1; then \
		mkdir -p "$${XDG_CONFIG_HOME:-$$HOME/.config}/gh"; \
		echo "GitHub CLI config directory created"; \
		if ! gh auth status >/dev/null 2>&1; then \
			echo "Run 'gh auth login' to authenticate with GitHub"; \
		fi; \
	fi

# Status
github-status: ## Check GitHub CLI installation status
	@echo "GitHub CLI Status"
	@echo "================="
	@echo ""
	@if command -v gh >/dev/null 2>&1; then \
		echo "Version: $$(gh --version | head -1 | awk '{print $$3}')"; \
		echo "Path: $$(which gh)"; \
		echo ""; \
		if gh auth status >/dev/null 2>&1; then \
			echo "Authentication: Logged in"; \
			gh auth status 2>&1 | grep "Logged in" | head -1; \
		else \
			echo "Authentication: Not logged in"; \
		fi; \
	else \
		echo "GitHub CLI not installed. Run: make github-install"; \
	fi

# Test
github-test: ## Test GitHub CLI functionality
	@echo "Testing GitHub CLI..."
	@command -v gh >/dev/null || (echo "GitHub CLI not installed"; exit 1)
	@gh --version >/dev/null 2>&1 && echo "GitHub CLI test passed" || (echo "GitHub CLI test failed"; exit 1)
