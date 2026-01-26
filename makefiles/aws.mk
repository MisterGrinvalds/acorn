# AWS CLI Management
# Amazon Web Services command line interface

.PHONY: aws-install aws-setup aws-status aws-test

# Installation
aws-install: ## Install AWS CLI
	@if command -v aws >/dev/null 2>&1; then \
		echo "AWS CLI already installed: $$(aws --version)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install awscli; \
	else \
		curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"; \
		unzip -q awscliv2.zip; \
		sudo ./aws/install; \
		rm -rf aws awscliv2.zip; \
	fi

aws-setup: aws-install ## Setup AWS CLI environment
	@if command -v aws >/dev/null 2>&1; then \
		mkdir -p "$${HOME}/.aws"; \
		if [ ! -f "$${HOME}/.aws/config" ]; then \
			echo "AWS CLI installed. Configure with: aws configure"; \
		fi; \
	fi

# Status
aws-status: ## Check AWS CLI installation status
	@echo "AWS CLI Status"
	@echo "=============="
	@echo ""
	@if command -v aws >/dev/null 2>&1; then \
		echo "Version: $$(aws --version | awk '{print $$1}' | cut -d'/' -f2)"; \
		echo "Path: $$(which aws)"; \
		echo ""; \
		if [ -f "$${HOME}/.aws/credentials" ]; then \
			echo "Credentials: Configured"; \
			echo "Profiles: $$(grep -c '^\[' "$${HOME}/.aws/credentials" 2>/dev/null || echo 0)"; \
		else \
			echo "Credentials: Not configured"; \
		fi; \
	else \
		echo "AWS CLI not installed. Run: make aws-install"; \
	fi

# Test
aws-test: ## Test AWS CLI functionality
	@echo "Testing AWS CLI..."
	@command -v aws >/dev/null || (echo "AWS CLI not installed"; exit 1)
	@aws --version >/dev/null 2>&1 && echo "AWS CLI test passed" || (echo "AWS CLI test failed"; exit 1)
