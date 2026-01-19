# Testing targets
# Syntax checking, dotfiles tests, security, performance

.PHONY: test test-quick test-all test-syntax test-dotfiles test-dotfiles-basic test-dotfiles-advanced
.PHONY: test-required-tools test-auth-status test-install test-installer-components test-integration
.PHONY: test-security test-performance test-workflows test-ci test-dev-env test-comprehensive
.PHONY: analyze-logs test-report benchmark stress-test
.PHONY: test-components test-component test-component-unit review-components

# Quick tests
test: test-quick ## Alias for test-quick

test-quick: setup test-syntax test-dotfiles-basic ## Run quick tests (syntax, basic functionality)

test-all: setup test-syntax test-dotfiles test-components test-integration ## Run complete test suite

# Syntax checking
test-syntax: ## Check syntax of all shell scripts
	@for file in core/*.sh; do \
		[ -f "$$file" ] && bash -n "$$file" || exit 1; \
	done
	@find components -name "*.sh" -type f 2>/dev/null | while read file; do \
		bash -n "$$file" || exit 1; \
	done
	@[ -f .bash_profile ] && bash -n .bash_profile || true

# Dotfiles tests
test-dotfiles: test-dotfiles-basic test-dotfiles-advanced ## Test all dotfiles functionality

test-dotfiles-basic: ## Test basic dotfiles functionality
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; [ "$$CURRENT_SHELL" = "bash" ]'
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; [ -n "$$DOTFILES_ROOT" ]'
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; alias | grep -q "ll="'

test-dotfiles-advanced: ## Test advanced dotfiles features
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; declare -f git_branch >/dev/null'
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; declare -f git_color >/dev/null'
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; declare -f mkvenv >/dev/null'

# Tool availability checks
test-required-tools: ## Check if required CLI tools are installed
	@command -v git >/dev/null && echo "git: ok" || echo "git: missing"
	@command -v curl >/dev/null && echo "curl: ok" || echo "curl: missing"
	@command -v jq >/dev/null && echo "jq: ok" || echo "jq: missing"
	@command -v aws >/dev/null && echo "aws: ok" || echo "aws: missing"
	@command -v docker >/dev/null && echo "docker: ok" || echo "docker: missing"
	@command -v kubectl >/dev/null && echo "kubectl: ok" || echo "kubectl: missing"

# Authentication status
test-auth-status: ## Check authentication status for cloud services
	@if command -v aws >/dev/null 2>&1 && aws sts get-caller-identity >/dev/null 2>&1; then \
		echo "AWS: authenticated"; else echo "AWS: not authenticated"; fi
	@if command -v az >/dev/null 2>&1 && az account show >/dev/null 2>&1; then \
		echo "Azure: authenticated"; else echo "Azure: not authenticated"; fi
	@if command -v gh >/dev/null 2>&1 && gh auth status >/dev/null 2>&1; then \
		echo "GitHub: authenticated"; else echo "GitHub: not authenticated"; fi
	@if command -v kubectl >/dev/null 2>&1 && kubectl cluster-info >/dev/null 2>&1; then \
		echo "Kubernetes: connected"; else echo "Kubernetes: not connected"; fi

# Installation tests
test-install: ## Test installation process
	@bash -n install.sh
	@bash install.sh --help > $(LOG_DIR)/installer-help.log 2>&1 || true
	@bash install.sh --test > $(LOG_DIR)/installer-test.log 2>&1 || true

test-installer-components: ## Test installer component structure
	@[ -d core ] || (echo "core/ missing" && exit 1)
	@[ -f core/bootstrap.sh ] || (echo "core/bootstrap.sh missing" && exit 1)
	@[ -d components ] || (echo "components/ missing" && exit 1)
	@[ -d config ] || (echo "config/ missing" && exit 1)

# Integration tests
test-integration: ## Test component integration
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source core/bootstrap.sh; declare -f dotfiles_status >/dev/null'

# Component tests using acorn CLI
test-components: build ## Test all sapling components
	@./bin/acorn test --skip-missing

test-component: build ## Test a specific component (COMPONENT=name)
	@if [ -z "$(COMPONENT)" ]; then \
		echo "Usage: make test-component COMPONENT=<name>"; \
		exit 1; \
	fi
	@./bin/acorn test $(COMPONENT)

test-component-unit: ## Run Go unit tests for component package
	@go test -v ./internal/utils/component/...

review-components: build ## Review all components
	@./bin/acorn review --all

# Security tests
test-security: ## Check for security issues
	@! grep -r -i "password\|secret\|key" --include="*.sh" . | grep -v "test\|example\|template" || true
	@find . -name "*.sh" -perm +111 | grep -v "install.sh" | grep -v "setup.sh" || true

# Performance tests
test-performance: ## Measure shell loading time
	@time bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source core/bootstrap.sh; exit 0' > $(LOG_DIR)/perf-load.log 2>&1

test-workflows: ## Test complete workflows
	@cd $(TEST_PROJECT_DIR) && bash -c '\
		export DOTFILES_ROOT="$(DOTFILES_DIR)" AUTO_DRY_RUN=true; \
		source $(DOTFILES_DIR)/shell/init.sh; \
		mkvenv test-workflow-env || true; \
	' > $(LOG_DIR)/workflow-dev.log 2>&1

# CI-friendly tests
test-ci: setup test-syntax test-dotfiles-basic test-security test-install ## Run CI/CD tests

test-dev-env: ## Check development environment
	@command -v git >/dev/null || exit 1
	@command -v curl >/dev/null || exit 1
	@command -v jq >/dev/null || echo "jq recommended"
	@command -v aws >/dev/null && echo "aws: found" || echo "aws: not found"
	@command -v kubectl >/dev/null && echo "kubectl: found" || echo "kubectl: not found"

# Benchmarks and stress tests
benchmark: ## Run performance benchmarks
	@mkdir -p $(LOG_DIR)
	@for i in {1..10}; do \
		time bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source core/bootstrap.sh; exit 0' 2>> $(LOG_DIR)/benchmark-load.log; \
	done

stress-test: ## Run stress tests
	@for i in {1..50}; do \
		bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source core/bootstrap.sh; exit 0' >/dev/null 2>&1 || exit 1; \
	done

# Log analysis
analyze-logs: ## Analyze test logs for issues
	@if [ -d $(LOG_DIR) ]; then \
		ls -la $(LOG_DIR)/; \
		grep -i "error\|fail\|exception" $(LOG_DIR)/*.log 2>/dev/null || echo "No errors found"; \
	fi

# Test report
test-report: ## Generate test report
	@echo "# Test Report" > $(TEST_DIR)/test-report.md
	@echo "Generated: $$(date)" >> $(TEST_DIR)/test-report.md
	@echo "OS: $$(uname -s), Shell: $$SHELL" >> $(TEST_DIR)/test-report.md

# Comprehensive test
test-comprehensive: setup test-syntax test-dotfiles test-components test-auth-status test-required-tools test-install test-installer-components test-integration test-security test-performance test-workflows analyze-logs test-report ## Run all tests
