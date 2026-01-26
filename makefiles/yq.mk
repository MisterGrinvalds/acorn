# YQ Management
# YAML processor

.PHONY: yq-install yq-setup yq-status yq-test

# Installation
yq-install: ## Install yq
	@if command -v yq >/dev/null 2>&1; then \
		echo "yq already installed: $$(yq --version)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install yq; \
	else \
		YQ_VERSION=$$(curl -s https://api.github.com/repos/mikefarah/yq/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")'); \
		wget https://github.com/mikefarah/yq/releases/download/$$YQ_VERSION/yq_linux_amd64 -O yq; \
		chmod +x yq; \
		sudo mv yq /usr/local/bin/; \
	fi

yq-setup: yq-install ## Setup yq environment
	@echo "yq requires no additional setup"

# Status
yq-status: ## Check yq installation status
	@echo "YQ Status"
	@echo "========="
	@echo ""
	@if command -v yq >/dev/null 2>&1; then \
		echo "Version: $$(yq --version | awk '{print $$NF}')"; \
		echo "Path: $$(which yq)"; \
	else \
		echo "yq not installed. Run: make yq-install"; \
	fi

# Test
yq-test: ## Test yq functionality
	@echo "Testing yq..."
	@command -v yq >/dev/null || (echo "yq not installed"; exit 1)
	@echo 'test: value' | yq '.test' | grep -q "value" && echo "yq test passed" || (echo "yq test failed"; exit 1)
