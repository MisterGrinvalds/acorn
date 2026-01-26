# K9s Management
# Kubernetes CLI UI

.PHONY: k9s-install k9s-setup k9s-status k9s-test

# Installation
k9s-install: ## Install k9s
	@if command -v k9s >/dev/null 2>&1; then \
		echo "k9s already installed: $$(k9s version --short 2>/dev/null | head -1)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install k9s; \
	else \
		K9S_VERSION=$$(curl -s https://api.github.com/repos/derailed/k9s/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")'); \
		curl -sL https://github.com/derailed/k9s/releases/download/$$K9S_VERSION/k9s_Linux_amd64.tar.gz | tar xz; \
		chmod +x k9s; \
		sudo mv k9s /usr/local/bin/; \
	fi

k9s-setup: k9s-install ## Setup k9s environment
	@if command -v k9s >/dev/null 2>&1; then \
		mkdir -p "$${XDG_CONFIG_HOME:-$$HOME/.config}/k9s"; \
		echo "k9s config directory created"; \
	fi

# Status
k9s-status: ## Check k9s installation status
	@echo "K9s Status"
	@echo "=========="
	@echo ""
	@if command -v k9s >/dev/null 2>&1; then \
		echo "Version: $$(k9s version --short 2>/dev/null | head -1 | awk '{print $$2}')"; \
		echo "Path: $$(which k9s)"; \
	else \
		echo "k9s not installed. Run: make k9s-install"; \
	fi

# Test
k9s-test: ## Test k9s functionality
	@echo "Testing k9s..."
	@command -v k9s >/dev/null || (echo "k9s not installed"; exit 1)
	@k9s version >/dev/null 2>&1 && echo "k9s test passed" || (echo "k9s test failed"; exit 1)
