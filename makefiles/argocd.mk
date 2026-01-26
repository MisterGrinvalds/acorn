# ArgoCD Management
# GitOps continuous delivery tool for Kubernetes

.PHONY: argocd-install argocd-setup argocd-status argocd-test

# Installation
argocd-install: ## Install ArgoCD CLI
	@if command -v argocd >/dev/null 2>&1; then \
		echo "ArgoCD CLI already installed: $$(argocd version --client --short 2>/dev/null)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install argocd; \
	else \
		ARGOCD_VERSION=$$(curl -s https://api.github.com/repos/argoproj/argo-cd/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")'); \
		curl -sSL -o argocd https://github.com/argoproj/argo-cd/releases/download/$$ARGOCD_VERSION/argocd-linux-amd64; \
		chmod +x argocd; \
		sudo mv argocd /usr/local/bin/; \
	fi

argocd-setup: argocd-install ## Setup ArgoCD environment
	@if command -v argocd >/dev/null 2>&1; then \
		mkdir -p "$${XDG_CONFIG_HOME:-$$HOME/.config}/argocd"; \
		echo "ArgoCD config directory created"; \
	fi

# Status
argocd-status: ## Check ArgoCD installation status
	@echo "ArgoCD Status"
	@echo "============="
	@echo ""
	@if command -v argocd >/dev/null 2>&1; then \
		echo "Version: $$(argocd version --client --short 2>/dev/null | awk '{print $$2}')"; \
		echo "Path: $$(which argocd)"; \
	else \
		echo "ArgoCD CLI not installed. Run: make argocd-install"; \
	fi

# Test
argocd-test: ## Test ArgoCD functionality
	@echo "Testing ArgoCD..."
	@command -v argocd >/dev/null || (echo "ArgoCD not installed"; exit 1)
	@argocd version --client >/dev/null 2>&1 && echo "ArgoCD test passed" || (echo "ArgoCD test failed"; exit 1)
