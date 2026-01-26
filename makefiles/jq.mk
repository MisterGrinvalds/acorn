# JQ Management
# JSON processor

.PHONY: jq-install jq-setup jq-status jq-test

# Installation
jq-install: ## Install jq
	@if command -v jq >/dev/null 2>&1; then \
		echo "jq already installed: $$(jq --version)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install jq; \
	else \
		sudo apt-get update && sudo apt-get install -y jq || \
		sudo yum install -y jq || \
		echo "Please install jq manually"; \
	fi

jq-setup: jq-install ## Setup jq environment
	@echo "jq requires no additional setup"

# Status
jq-status: ## Check jq installation status
	@echo "JQ Status"
	@echo "========="
	@echo ""
	@if command -v jq >/dev/null 2>&1; then \
		echo "Version: $$(jq --version | cut -d'-' -f2)"; \
		echo "Path: $$(which jq)"; \
	else \
		echo "jq not installed. Run: make jq-install"; \
	fi

# Test
jq-test: ## Test jq functionality
	@echo "Testing jq..."
	@command -v jq >/dev/null || (echo "jq not installed"; exit 1)
	@echo '{"test": "value"}' | jq -r '.test' | grep -q "value" && echo "jq test passed" || (echo "jq test failed"; exit 1)
