# Terraform Management
# Infrastructure as Code tool

.PHONY: terraform-install terraform-setup terraform-status terraform-test

# Installation
terraform-install: ## Install Terraform
	@if command -v terraform >/dev/null 2>&1; then \
		echo "Terraform already installed: $$(terraform version | head -1)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew tap hashicorp/tap; \
		brew install hashicorp/tap/terraform; \
	else \
		wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg; \
		echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $$(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list; \
		sudo apt update && sudo apt install terraform; \
	fi

terraform-setup: terraform-install ## Setup Terraform environment
	@if command -v terraform >/dev/null 2>&1; then \
		mkdir -p "$${XDG_DATA_HOME:-$$HOME/.local/share}/terraform"; \
		mkdir -p "$${XDG_CACHE_HOME:-$$HOME/.cache}/terraform"; \
		echo "Terraform directories created"; \
	fi

# Status
terraform-status: ## Check Terraform installation status
	@echo "Terraform Status"
	@echo "================"
	@echo ""
	@if command -v terraform >/dev/null 2>&1; then \
		echo "Version: $$(terraform version | head -1 | awk '{print $$2}')"; \
		echo "Path: $$(which terraform)"; \
	else \
		echo "Terraform not installed. Run: make terraform-install"; \
	fi

# Test
terraform-test: ## Test Terraform functionality
	@echo "Testing Terraform..."
	@command -v terraform >/dev/null || (echo "Terraform not installed"; exit 1)
	@terraform version >/dev/null 2>&1 && echo "Terraform test passed" || (echo "Terraform test failed"; exit 1)
