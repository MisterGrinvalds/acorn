# Kubernetes Management
# kubectl CLI and cluster interaction

.PHONY: kubernetes-install kubectl-install kubernetes-setup kubernetes-status kubernetes-test

# Installation
kubernetes-install kubectl-install: ## Install kubectl
	@if command -v kubectl >/dev/null 2>&1; then \
		echo "kubectl already installed: $$(kubectl version --client --short 2>/dev/null || kubectl version --client)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install kubectl; \
	else \
		curl -LO "https://dl.k8s.io/release/$$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"; \
		chmod +x kubectl; \
		sudo mv kubectl /usr/local/bin/; \
	fi

kubernetes-setup: kubernetes-install ## Setup Kubernetes environment
	@mkdir -p "$${HOME}/.kube"
	@echo "Kubernetes config directory created"
	@if command -v kubectl >/dev/null 2>&1; then \
		echo "Current context: $$(kubectl config current-context 2>/dev/null || echo 'none')"; \
	fi

# Status
kubernetes-status: ## Check Kubernetes installation status
	@echo "Kubernetes Status"
	@echo "================="
	@echo ""
	@if command -v kubectl >/dev/null 2>&1; then \
		echo "Version: $$(kubectl version --client --short 2>/dev/null | awk '{print $$3}' || kubectl version --client -o json 2>/dev/null | jq -r '.clientVersion.gitVersion' || echo 'unknown')"; \
		echo "Path: $$(which kubectl)"; \
		echo ""; \
		if kubectl cluster-info >/dev/null 2>&1; then \
			echo "Cluster: Connected"; \
			echo "Context: $$(kubectl config current-context)"; \
			echo "Server Version: $$(kubectl version --short 2>/dev/null | grep Server | awk '{print $$3}' || kubectl version -o json 2>/dev/null | jq -r '.serverVersion.gitVersion' || echo 'N/A')"; \
		else \
			echo "Cluster: Not connected"; \
			echo "Contexts: $$(kubectl config get-contexts -o name 2>/dev/null | wc -l | tr -d ' ')"; \
		fi; \
	else \
		echo "kubectl not installed. Run: make kubernetes-install"; \
	fi

# Test
kubernetes-test: ## Test kubectl functionality
	@echo "Testing kubectl..."
	@command -v kubectl >/dev/null || (echo "kubectl not installed"; exit 1)
	@kubectl version --client >/dev/null 2>&1 && echo "kubectl test passed" || (echo "kubectl test failed"; exit 1)
