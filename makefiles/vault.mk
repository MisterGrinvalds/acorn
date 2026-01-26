# HashiCorp Vault Management
# Secrets management CLI

.PHONY: vault-install vault-setup vault-status vault-test

# Installation
vault-install: ## Install HashiCorp Vault CLI
	@if command -v vault >/dev/null 2>&1; then \
		echo "Vault already installed: $$(vault version | head -1)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew tap hashicorp/tap 2>/dev/null || true; \
		brew install hashicorp/tap/vault; \
	else \
		VAULT_VERSION=$$(curl -s https://api.github.com/repos/hashicorp/vault/releases/latest | grep -oP '"tag_name": "v\K(.*)(?=")'); \
		wget "https://releases.hashicorp.com/vault/$${VAULT_VERSION}/vault_$${VAULT_VERSION}_linux_amd64.zip"; \
		unzip "vault_$${VAULT_VERSION}_linux_amd64.zip"; \
		sudo mv vault /usr/local/bin/; \
		rm "vault_$${VAULT_VERSION}_linux_amd64.zip"; \
	fi

vault-setup: vault-install ## Setup Vault environment
	@if command -v vault >/dev/null 2>&1; then \
		echo "Vault CLI installed successfully"; \
		echo ""; \
		echo "Next steps:"; \
		echo "  1. Set VAULT_ADDR environment variable to your Vault server"; \
		echo "     export VAULT_ADDR=https://vault.example.com:8200"; \
		echo "  2. Authenticate with: vault login"; \
		echo "  3. Verify with: vault status"; \
	fi

# Status
vault-status: ## Check Vault installation status
	@echo "Vault Status"
	@echo "============"
	@echo ""
	@if command -v vault >/dev/null 2>&1; then \
		echo "Version: $$(vault version | head -1 | awk '{print $$2}')"; \
		echo "Path: $$(which vault)"; \
		echo ""; \
		echo "Server: $${VAULT_ADDR:-not set}"; \
		if [ -n "$${VAULT_ADDR}" ]; then \
			if vault status >/dev/null 2>&1; then \
				echo "Connection: OK"; \
				vault status 2>/dev/null | grep -E "(Sealed|Version|Cluster)" || true; \
				echo ""; \
				if vault token lookup >/dev/null 2>&1; then \
					echo "Authentication: Valid token"; \
					vault token lookup 2>/dev/null | grep -E "(display_name|policies|expire_time)" || true; \
				else \
					echo "Authentication: Not logged in"; \
				fi; \
			else \
				echo "Connection: Failed to connect to $${VAULT_ADDR}"; \
			fi; \
		else \
			echo "Set VAULT_ADDR to connect to a Vault server"; \
		fi; \
	else \
		echo "Vault not installed. Run: make vault-install"; \
	fi

# Test
vault-test: ## Test Vault functionality
	@echo "Testing Vault..."
	@command -v vault >/dev/null || (echo "Vault not installed"; exit 1)
	@vault version >/dev/null 2>&1 && echo "Vault test passed" || (echo "Vault test failed"; exit 1)
