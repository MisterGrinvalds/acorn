# Helm Management
# Kubernetes package manager

.PHONY: helm-install helm-setup helm-status helm-test

# Installation
helm-install: ## Install Helm
	@if command -v helm >/dev/null 2>&1; then \
		echo "Helm already installed: $$(helm version --short)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install helm; \
	else \
		curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash; \
	fi

helm-setup: helm-install ## Setup Helm environment
	@if command -v helm >/dev/null 2>&1; then \
		mkdir -p "$${XDG_CONFIG_HOME:-$$HOME/.config}/helm"; \
		mkdir -p "$${XDG_CACHE_HOME:-$$HOME/.cache}/helm"; \
		mkdir -p "$${XDG_DATA_HOME:-$$HOME/.local/share}/helm"; \
		echo "Helm directories created"; \
		helm repo add stable https://charts.helm.sh/stable 2>/dev/null || true; \
		helm repo update >/dev/null 2>&1 || true; \
	fi

# Status
helm-status: ## Check Helm installation status
	@echo "Helm Status"
	@echo "==========="
	@echo ""
	@if command -v helm >/dev/null 2>&1; then \
		echo "Version: $$(helm version --short | cut -d' ' -f1 | cut -d':' -f2)"; \
		echo "Path: $$(which helm)"; \
		echo ""; \
		echo "Repositories: $$(helm repo list 2>/dev/null | tail -n +2 | wc -l | tr -d ' ')"; \
		helm repo list 2>/dev/null | tail -n +2 | head -5 || echo "  (none)"; \
	else \
		echo "Helm not installed. Run: make helm-install"; \
	fi

# Test
helm-test: ## Test Helm functionality
	@echo "Testing Helm..."
	@command -v helm >/dev/null || (echo "Helm not installed"; exit 1)
	@helm version >/dev/null 2>&1 && echo "Helm test passed" || (echo "Helm test failed"; exit 1)
