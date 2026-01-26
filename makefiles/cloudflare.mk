# Cloudflare Wrangler Management
# Cloudflare Workers CLI

.PHONY: cloudflare-install wrangler-install cloudflare-setup cloudflare-status cloudflare-test

# Installation
cloudflare-install wrangler-install: ## Install Cloudflare Wrangler
	@if command -v wrangler >/dev/null 2>&1; then \
		echo "Wrangler already installed: $$(wrangler --version)"; \
	elif command -v npm >/dev/null 2>&1; then \
		npm install -g wrangler; \
	else \
		echo "npm not found. Install Node.js first: make node-install"; \
		exit 1; \
	fi

cloudflare-setup: cloudflare-install ## Setup Wrangler environment
	@if command -v wrangler >/dev/null 2>&1; then \
		mkdir -p "$${XDG_CONFIG_HOME:-$$HOME/.config}/wrangler"; \
		echo "Wrangler config directory created"; \
		if ! wrangler whoami >/dev/null 2>&1; then \
			echo "Run 'wrangler login' to authenticate with Cloudflare"; \
		fi; \
	fi

# Status
cloudflare-status: ## Check Wrangler installation status
	@echo "Cloudflare Wrangler Status"
	@echo "=========================="
	@echo ""
	@if command -v wrangler >/dev/null 2>&1; then \
		echo "Version: $$(wrangler --version | awk '{print $$2}')"; \
		echo "Path: $$(which wrangler)"; \
		echo ""; \
		if wrangler whoami >/dev/null 2>&1; then \
			echo "Authentication: Logged in"; \
		else \
			echo "Authentication: Not logged in"; \
		fi; \
	else \
		echo "Wrangler not installed. Run: make cloudflare-install"; \
	fi

# Test
cloudflare-test: ## Test Wrangler functionality
	@echo "Testing Wrangler..."
	@command -v wrangler >/dev/null || (echo "Wrangler not installed"; exit 1)
	@wrangler --version >/dev/null 2>&1 && echo "Wrangler test passed" || (echo "Wrangler test failed"; exit 1)
