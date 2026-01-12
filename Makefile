# Makefile for Testing Component-Based Dotfiles System
# Comprehensive test suite for all functionality

.PHONY: help test test-all test-quick test-dotfiles test-automation test-cloud test-modules test-syntax test-security test-install clean setup ai-setup ai-status ai-test ai-models ai-chat ai-benchmark ai-cleanup ai-examples nvm-install nvm-setup nvm-status node-install node-lts pnpm-install pnpm-setup node-status shell-status shell-test-discovery shell-test-xdg shell-test-theme shell-test-env shell-test-options shell-test-aliases shell-test-functions shell-test-prompt shell-test-all dotfiles-install dotfiles-inject dotfiles-eject dotfiles-link dotfiles-unlink dotfiles-status dotfiles-reload dotfiles-update uv-install uv-setup uv-status python-status venv-create venv-status go-install go-setup go-status go-tools status brew-status brew-update brew-install-devops brew-install-dev brew-install-db brew-install-all db-install-mysql db-install-mongo db-install-redis db-install-neo4j db-install-kafka component-list component-status component-new component-validate test-components test-component-loader test-component-deps

# Default target
help: ## Show this help message
	@echo "Component-Based Dotfiles Test Suite"
	@echo "===================================="
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Test Categories:"
	@echo "  test-quick        - Quick syntax and basic functionality tests"
	@echo "  test-all          - Complete test suite"
	@echo "  test-dotfiles     - Test dotfiles configuration only"
	@echo "  test-auth-status  - Test authentication status for all services"
	@echo "  test-required-tools - Test if required CLI tools are installed"
	@echo ""
	@echo "AI/ML Management:"
	@echo "  ai-setup          - Setup AI/ML environment (Ollama + Hugging Face)"
	@echo "  ai-status         - Check AI/ML tools status"
	@echo "  ai-test           - Test AI/ML functionality"
	@echo "  ai-models         - List available AI models"
	@echo "  ai-chat           - Start interactive AI chat"
	@echo "  ai-benchmark      - Run AI performance benchmarks"
	@echo "  ai-examples       - Show AI usage examples"
	@echo "  ai-cleanup        - Clean AI model caches and stop services"
	@echo ""
	@echo "Node.js/NVM Management:"
	@echo "  nvm-install       - Install NVM (Node Version Manager)"
	@echo "  nvm-setup         - Install NVM + latest LTS Node + pnpm"
	@echo "  nvm-status        - Check NVM and Node.js status"
	@echo "  node-install      - Install latest Node.js via NVM"
	@echo "  node-lts          - Install latest LTS Node.js via NVM"
	@echo "  pnpm-install      - Install pnpm globally"
	@echo "  pnpm-setup        - Setup pnpm with corepack"
	@echo "  node-status       - Check Node.js ecosystem status"
	@echo ""
	@echo "Shell Layer Testing (DAG Parity):"
	@echo "  shell-status      - Show all shell module status"
	@echo "  shell-test-all    - Test complete shell loading sequence"
	@echo "  shell-test-*      - Test individual modules (discovery, xdg, theme, etc.)"
	@echo ""
	@echo "Component Management:"
	@echo "  component-list    - List all available components"
	@echo "  component-status  - Show component health and loading status"
	@echo "  component-new     - Create new component from template"
	@echo "  component-validate - Validate all component.yaml files"
	@echo "  test-components   - Test component discovery and loading"
	@echo ""
	@echo "Dotfiles Management:"
	@echo "  dotfiles-install  - Run full dotfiles installation"
	@echo "  dotfiles-inject   - Create shell bootstrap files"
	@echo "  dotfiles-eject    - Remove shell bootstrap files"
	@echo "  dotfiles-link     - Link app configurations"
	@echo "  dotfiles-unlink   - Remove app configuration links"
	@echo "  dotfiles-status   - Show installation status"
	@echo "  dotfiles-update   - Git pull and reload"
	@echo ""
	@echo "Python/UV Management:"
	@echo "  uv-install        - Install UV package manager"
	@echo "  uv-setup          - Complete UV setup with Python"
	@echo "  uv-status         - Check UV and Python status"
	@echo "  venv-create       - Create virtual environment"
	@echo "  venv-status       - Show active venv info"
	@echo ""
	@echo "Go Management:"
	@echo "  go-install        - Install Go"
	@echo "  go-setup          - Setup Go environment"
	@echo "  go-status         - Check Go installation"
	@echo "  go-tools          - Install common Go tools"
	@echo ""
	@echo "Brew Package Management:"
	@echo "  brew-status       - Show status of all managed brew packages"
	@echo "  brew-update       - Update all brew packages"
	@echo "  brew-install-devops - Install all DevOps tools"
	@echo "  brew-install-dev  - Install all dev tools"
	@echo "  brew-install-db   - Install all database tools"
	@echo "  brew-install-all  - Install all managed packages"
	@echo ""
	@echo "Database Tool Installs:"
	@echo "  db-install-mysql  - Install MySQL client + mycli"
	@echo "  db-install-mongo  - Install MongoDB shell + mongocli"
	@echo "  db-install-redis  - Install Redis + iredis"
	@echo "  db-install-neo4j  - Install Neo4j"
	@echo "  db-install-kafka  - Install Kafka"
	@echo ""
	@echo "Comprehensive:"
	@echo "  status            - Show complete environment status"

# Variables
SHELL := /bin/bash
DOTFILES_DIR := $(PWD)
TEST_DIR := $(DOTFILES_DIR)/tests
LOG_DIR := $(TEST_DIR)/logs
BACKUP_DIR := $(TEST_DIR)/backups

# Test configuration
TEST_PROJECT_DIR := $(TEST_DIR)/test-projects
TEST_TIMEOUT := 30

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

# Setup test environment
setup: ## Setup test environment
	@echo -e "$(BLUE)Setting up test environment...$(NC)"
	@mkdir -p $(LOG_DIR) $(BACKUP_DIR) $(TEST_PROJECT_DIR)
	@echo "Test directories created"
	@# Backup existing configurations
	@if [ -f ~/.bash_profile ]; then cp ~/.bash_profile $(BACKUP_DIR)/bash_profile.backup; fi
	@if [ -d ~/.config/shell ]; then cp -r ~/.config/shell $(BACKUP_DIR)/; fi
	@echo -e "$(GREEN)Test environment setup complete$(NC)"

# Quick tests - syntax and basic functionality
test-quick: setup test-syntax test-dotfiles-basic ## Run quick tests (syntax, basic functionality)
	@echo -e "$(GREEN)✅ Quick tests completed successfully$(NC)"

# Complete test suite
test-all: setup test-syntax test-dotfiles test-components test-integration ## Run complete test suite
	@echo -e "$(GREEN)✅ All tests completed successfully$(NC)"

# Individual test categories
test: test-quick ## Alias for test-quick

# Test syntax of all shell files
test-syntax: ## Test syntax of all shell scripts
	@echo -e "$(BLUE)Testing shell script syntax...$(NC)"
	@echo "Testing core modules..."
	@for file in core/*.sh; do \
		if [ -f "$$file" ]; then \
			echo "  Testing $$file..."; \
			bash -n "$$file" || (echo -e "$(RED)❌ Syntax error in $$file$(NC)" && exit 1); \
		fi; \
	done
	@echo "Testing component scripts..."
	@find components -name "*.sh" -type f 2>/dev/null | while read file; do \
		echo "  Testing $$file..."; \
		bash -n "$$file" || (echo -e "$(RED)❌ Syntax error in $$file$(NC)" && exit 1); \
	done
	@# Legacy files (if they exist)
	@if [ -f .bash_profile ]; then bash -n .bash_profile; fi
	@echo -e "$(GREEN)✅ All syntax tests passed$(NC)"

# Test dotfiles functionality
test-dotfiles: test-dotfiles-basic test-dotfiles-advanced ## Test all dotfiles functionality

test-dotfiles-basic: ## Test basic dotfiles functionality
	@echo -e "$(BLUE)Testing basic dotfiles functionality...$(NC)"
	@# Note: We set IS_INTERACTIVE=true to bypass early exit for non-interactive shells
	@# Test shell detection
	@echo "Testing shell detection..."
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; [ "$$CURRENT_SHELL" = "bash" ]' || \
		(echo -e "$(RED)❌ Shell detection failed$(NC)" && exit 1)
	@# Test environment loading
	@echo "Testing environment loading..."
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; [ -n "$$DOTFILES_ROOT" ]' || \
		(echo -e "$(RED)❌ DOTFILES_ROOT variable not set$(NC)" && exit 1)
	@# Test aliases
	@echo "Testing aliases..."
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; alias | grep -q "ll="' || \
		(echo -e "$(RED)❌ Basic aliases not loaded$(NC)" && exit 1)
	@echo -e "$(GREEN)✅ Basic dotfiles tests passed$(NC)"

test-dotfiles-advanced: ## Test advanced dotfiles features
	@echo -e "$(BLUE)Testing advanced dotfiles features...$(NC)"
	@# Test git prompt functions
	@echo "Testing git prompt functions..."
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; declare -f git_branch >/dev/null' || \
		(echo -e "$(RED)❌ git_branch function not defined$(NC)" && exit 1)
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; declare -f git_color >/dev/null' || \
		(echo -e "$(RED)❌ git_color function not defined$(NC)" && exit 1)
	@# Test custom functions
	@echo "Testing custom functions..."
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; declare -f mkvenv >/dev/null' || \
		(echo -e "$(RED)❌ mkvenv function not defined$(NC)" && exit 1)
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; declare -f h >/dev/null' || \
		(echo -e "$(RED)❌ h function not defined$(NC)" && exit 1)
	@echo -e "$(GREEN)✅ Advanced dotfiles tests passed$(NC)"

test-required-tools: ## Test if required CLI tools are installed
	@echo -e "$(BLUE)Checking required CLI tools...$(NC)"
	@echo "=== CLI Tools Status ===" > $(LOG_DIR)/cli-tools-status.log
	@echo "Checking development tools..." >> $(LOG_DIR)/cli-tools-status.log
	@command -v git >/dev/null && echo "✅ git installed" || echo "❌ git not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@command -v curl >/dev/null && echo "✅ curl installed" || echo "❌ curl not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@command -v jq >/dev/null && echo "✅ jq installed" || echo "❌ jq not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@echo "" >> $(LOG_DIR)/cli-tools-status.log
	@echo "Checking cloud CLI tools..." >> $(LOG_DIR)/cli-tools-status.log
	@command -v aws >/dev/null && echo "✅ AWS CLI installed" || echo "❌ AWS CLI not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@command -v az >/dev/null && echo "✅ Azure CLI installed" || echo "❌ Azure CLI not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@command -v doctl >/dev/null && echo "✅ doctl installed" || echo "❌ doctl not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@echo "" >> $(LOG_DIR)/cli-tools-status.log
	@echo "Checking container tools..." >> $(LOG_DIR)/cli-tools-status.log
	@command -v docker >/dev/null && echo "✅ Docker installed" || echo "❌ Docker not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@command -v kubectl >/dev/null && echo "✅ kubectl installed" || echo "❌ kubectl not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@command -v helm >/dev/null && echo "✅ Helm installed" || echo "❌ Helm not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@echo "" >> $(LOG_DIR)/cli-tools-status.log
	@echo "Checking development tools..." >> $(LOG_DIR)/cli-tools-status.log
	@command -v gh >/dev/null && echo "✅ GitHub CLI installed" || echo "❌ GitHub CLI not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@command -v python3 >/dev/null && echo "✅ Python 3 installed" || echo "❌ Python 3 not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@command -v go >/dev/null && echo "✅ Go installed" || echo "❌ Go not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@command -v node >/dev/null && echo "✅ Node.js installed" || echo "❌ Node.js not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@command -v npm >/dev/null && echo "✅ npm installed" || echo "❌ npm not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@echo -e "$(GREEN)✅ CLI tools check completed$(NC)"

test-auth-status: ## Test authentication status for all services
	@echo -e "$(BLUE)Testing authentication status...$(NC)"
	@echo "=== Authentication Status Report ===" > $(LOG_DIR)/auth-status.log
	@echo "Generated: $$(date)" >> $(LOG_DIR)/auth-status.log
	@echo "" >> $(LOG_DIR)/auth-status.log
	@# AWS Authentication
	@echo "AWS Authentication:" >> $(LOG_DIR)/auth-status.log
	@if command -v aws >/dev/null 2>&1; then \
		if aws sts get-caller-identity >/dev/null 2>&1; then \
			echo "✅ AWS: Authenticated" | tee -a $(LOG_DIR)/auth-status.log; \
			aws sts get-caller-identity --query 'Account' --output text 2>/dev/null | sed 's/^/   Account: /' >> $(LOG_DIR)/auth-status.log; \
		else echo "❌ AWS: Not authenticated" | tee -a $(LOG_DIR)/auth-status.log; fi; \
	else echo "❌ AWS CLI not installed" | tee -a $(LOG_DIR)/auth-status.log; fi
	@echo "" >> $(LOG_DIR)/auth-status.log
	@# Azure Authentication
	@echo "Azure Authentication:" >> $(LOG_DIR)/auth-status.log
	@if command -v az >/dev/null 2>&1; then \
		if az account show >/dev/null 2>&1; then \
			echo "✅ Azure: Authenticated" | tee -a $(LOG_DIR)/auth-status.log; \
			az account show --query 'name' --output tsv 2>/dev/null | sed 's/^/   Subscription: /' >> $(LOG_DIR)/auth-status.log; \
		else echo "❌ Azure: Not authenticated" | tee -a $(LOG_DIR)/auth-status.log; fi; \
	else echo "❌ Azure CLI not installed" | tee -a $(LOG_DIR)/auth-status.log; fi
	@echo "" >> $(LOG_DIR)/auth-status.log
	@# DigitalOcean Authentication
	@echo "DigitalOcean Authentication:" >> $(LOG_DIR)/auth-status.log
	@if command -v doctl >/dev/null 2>&1; then \
		if doctl account get >/dev/null 2>&1; then \
			echo "✅ DigitalOcean: Authenticated" | tee -a $(LOG_DIR)/auth-status.log; \
			doctl account get --format Email --no-header 2>/dev/null | sed 's/^/   Account: /' >> $(LOG_DIR)/auth-status.log; \
		else echo "❌ DigitalOcean: Not authenticated" | tee -a $(LOG_DIR)/auth-status.log; fi; \
	else echo "❌ doctl not installed" | tee -a $(LOG_DIR)/auth-status.log; fi
	@echo "" >> $(LOG_DIR)/auth-status.log
	@# GitHub Authentication
	@echo "GitHub Authentication:" >> $(LOG_DIR)/auth-status.log
	@if command -v gh >/dev/null 2>&1; then \
		if gh auth status >/dev/null 2>&1; then \
			echo "✅ GitHub: Authenticated" | tee -a $(LOG_DIR)/auth-status.log; \
			gh api user --jq '.login' 2>/dev/null | sed 's/^/   User: /' >> $(LOG_DIR)/auth-status.log; \
		else echo "❌ GitHub: Not authenticated" | tee -a $(LOG_DIR)/auth-status.log; fi; \
	else echo "❌ GitHub CLI not installed" | tee -a $(LOG_DIR)/auth-status.log; fi
	@echo "" >> $(LOG_DIR)/auth-status.log
	@# Docker Authentication
	@echo "Docker Authentication:" >> $(LOG_DIR)/auth-status.log
	@if command -v docker >/dev/null 2>&1; then \
		if docker info >/dev/null 2>&1; then \
			echo "✅ Docker: Running" | tee -a $(LOG_DIR)/auth-status.log; \
		else echo "❌ Docker: Not running" | tee -a $(LOG_DIR)/auth-status.log; fi; \
	else echo "❌ Docker not installed" | tee -a $(LOG_DIR)/auth-status.log; fi
	@echo "" >> $(LOG_DIR)/auth-status.log
	@# Kubernetes Authentication
	@echo "Kubernetes Authentication:" >> $(LOG_DIR)/auth-status.log
	@if command -v kubectl >/dev/null 2>&1; then \
		if kubectl cluster-info >/dev/null 2>&1; then \
			echo "✅ Kubernetes: Connected" | tee -a $(LOG_DIR)/auth-status.log; \
			kubectl config current-context 2>/dev/null | sed 's/^/   Context: /' >> $(LOG_DIR)/auth-status.log; \
		else echo "❌ Kubernetes: No cluster access" | tee -a $(LOG_DIR)/auth-status.log; fi; \
	else echo "❌ kubectl not installed" | tee -a $(LOG_DIR)/auth-status.log; fi
	@echo -e "$(GREEN)✅ Authentication status check completed$(NC)"
	@echo "📋 Full report saved to: $(LOG_DIR)/auth-status.log"

# Test installation and setup
test-install: ## Test installation process
	@echo -e "$(BLUE)Testing installation process...$(NC)"
	@# Test install.sh script
	@echo "Testing install.sh script..."
	@bash -n install.sh || (echo -e "$(RED)❌ Install script syntax error$(NC)" && exit 1)
	@# Test installer help
	@echo "Testing installer help..."
	@bash install.sh --help > $(LOG_DIR)/installer-help.log 2>&1 || \
		(echo -e "$(RED)❌ Installer help failed$(NC)" && exit 1)
	@# Test installer options validation
	@echo "Testing installer argument parsing..."
	@bash install.sh --test > $(LOG_DIR)/installer-test.log 2>&1 || \
		echo -e "$(YELLOW)⚠️ Installer test mode completed$(NC)"
	@echo -e "$(GREEN)✅ Installation tests passed$(NC)"

test-installer-components: ## Test individual installer components
	@echo -e "$(BLUE)Testing installer components...$(NC)"
	@# Test backup functionality
	@echo "Testing backup creation..."
	@mkdir -p $(TEST_DIR)/installer-test
	@touch $(TEST_DIR)/installer-test/.bash_profile
	@# Test dotfiles structure
	@echo "Testing dotfiles structure..."
	@[ -d core ] || (echo -e "$(RED)❌ Core directory missing$(NC)" && exit 1)
	@[ -f core/bootstrap.sh ] || (echo -e "$(RED)❌ Core bootstrap.sh missing$(NC)" && exit 1)
	@[ -d components ] || (echo -e "$(RED)❌ Components directory missing$(NC)" && exit 1)
	@[ -d config ] || (echo -e "$(RED)❌ Config directory missing$(NC)" && exit 1)
	@echo -e "$(GREEN)✅ Installer components tests passed$(NC)"

# Integration tests
test-integration: ## Test integration between components
	@echo -e "$(BLUE)Testing component integration...$(NC)"
	@# Test dotfiles integration
	@echo "Testing dotfiles integration..."
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source core/bootstrap.sh; declare -f dotfiles_status >/dev/null' || \
		(echo -e "$(RED)❌ Dotfiles integration not loaded$(NC)" && exit 1)
	@echo -e "$(GREEN)✅ Integration tests passed$(NC)"

# Security tests
test-security: ## Test for security issues
	@echo -e "$(BLUE)Testing for security issues...$(NC)"
	@# Check for hardcoded secrets
	@echo "Checking for hardcoded secrets..."
	@! grep -r -i "password\|secret\|key" --include="*.sh" . | grep -v "test\|example\|template" || \
		(echo -e "$(YELLOW)⚠️ Potential hardcoded secrets found$(NC)")
	@# Check file permissions
	@echo "Checking file permissions..."
	@find . -name "*.sh" -perm +111 | grep -v "install.sh" | grep -v "setup.sh" && \
		echo -e "$(YELLOW)⚠️ Unexpected executable permissions found$(NC)" || true
	@echo -e "$(GREEN)✅ Security tests passed$(NC)"

# Performance tests
test-performance: ## Test performance of shell loading
	@echo -e "$(BLUE)Testing shell loading performance...$(NC)"
	@# Time shell loading
	@echo "Testing bootstrap load time..."
	@time bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source core/bootstrap.sh; exit 0' > $(LOG_DIR)/perf-load.log 2>&1 || \
		(echo -e "$(RED)❌ Performance test failed$(NC)" && exit 1)
	@echo -e "$(GREEN)✅ Performance tests passed$(NC)"

# Functional tests for specific workflows
test-workflows: ## Test complete workflows
	@echo -e "$(BLUE)Testing complete workflows...$(NC)"
	@# Test development workflow
	@echo "Testing development workflow..."
	@cd $(TEST_PROJECT_DIR) && bash -c '\
		export DOTFILES_ROOT="$(DOTFILES_DIR)"; \
		source $(DOTFILES_DIR)/shell/init.sh; \
		export AUTO_DRY_RUN=true; \
		echo "Testing python environment creation..."; \
		mkvenv test-workflow-env || echo "Workflow test informational only"; \
	' > $(LOG_DIR)/workflow-dev.log 2>&1
	@echo -e "$(GREEN)✅ Workflow tests completed$(NC)"

# Test log analysis
analyze-logs: ## Analyze test logs for issues
	@echo -e "$(BLUE)Analyzing test logs...$(NC)"
	@if [ -d $(LOG_DIR) ]; then \
		echo "Log files created:"; \
		ls -la $(LOG_DIR)/; \
		echo ""; \
		echo "Checking for errors in logs..."; \
		grep -i "error\|fail\|exception" $(LOG_DIR)/*.log || echo "No errors found in logs"; \
	else \
		echo "No log directory found"; \
	fi

# Generate test report
test-report: ## Generate comprehensive test report
	@echo -e "$(BLUE)Generating test report...$(NC)"
	@echo "# Bash Profile & Automation Framework Test Report" > $(TEST_DIR)/test-report.md
	@echo "Generated: $$(date)" >> $(TEST_DIR)/test-report.md
	@echo "" >> $(TEST_DIR)/test-report.md
	@echo "## Test Environment" >> $(TEST_DIR)/test-report.md
	@echo "- OS: $$(uname -s)" >> $(TEST_DIR)/test-report.md
	@echo "- Shell: $$SHELL" >> $(TEST_DIR)/test-report.md
	@echo "- Directory: $$(pwd)" >> $(TEST_DIR)/test-report.md
	@echo "" >> $(TEST_DIR)/test-report.md
	@echo "## File Structure" >> $(TEST_DIR)/test-report.md
	@echo "\`\`\`" >> $(TEST_DIR)/test-report.md
	@tree -I 'node_modules|.git|tests' -L 3 >> $(TEST_DIR)/test-report.md 2>/dev/null || \
		find . -type d -not -path '*/.*' -not -path '*/node_modules*' | head -20 >> $(TEST_DIR)/test-report.md
	@echo "\`\`\`" >> $(TEST_DIR)/test-report.md
	@echo "" >> $(TEST_DIR)/test-report.md
	@if [ -d $(LOG_DIR) ]; then \
		echo "## Test Logs" >> $(TEST_DIR)/test-report.md; \
		for log in $(LOG_DIR)/*.log; do \
			echo "### $$(basename $$log)" >> $(TEST_DIR)/test-report.md; \
			echo "\`\`\`" >> $(TEST_DIR)/test-report.md; \
			head -20 "$$log" >> $(TEST_DIR)/test-report.md; \
			echo "\`\`\`" >> $(TEST_DIR)/test-report.md; \
			echo "" >> $(TEST_DIR)/test-report.md; \
		done; \
	fi
	@echo -e "$(GREEN)Test report generated: $(TEST_DIR)/test-report.md$(NC)"

# Clean test artifacts
clean: ## Clean test artifacts and restore backups
	@echo -e "$(BLUE)Cleaning test artifacts...$(NC)"
	@# Remove test directories
	@rm -rf $(TEST_DIR)/test-projects
	@rm -rf $(LOG_DIR)
	@# Restore backups if they exist
	@if [ -f $(BACKUP_DIR)/bash_profile.backup ]; then \
		cp $(BACKUP_DIR)/bash_profile.backup ~/.bash_profile; \
		echo "Restored ~/.bash_profile"; \
	fi
	@if [ -d $(BACKUP_DIR)/shell ]; then \
		cp -r $(BACKUP_DIR)/shell ~/.config/; \
		echo "Restored ~/.config/shell"; \
	fi
	@echo -e "$(GREEN)✅ Cleanup completed$(NC)"

# CI/CD friendly test
test-ci: ## Run CI-friendly tests (no interactive components)
	@echo -e "$(BLUE)Running CI/CD tests...$(NC)"
	@$(MAKE) setup
	@$(MAKE) test-syntax
	@$(MAKE) test-dotfiles-basic
	@$(MAKE) test-automation-basic
	@$(MAKE) test-security
	@$(MAKE) test-install
	@echo -e "$(GREEN)✅ CI/CD tests completed$(NC)"

# Development tests
test-dev-env: ## Test development environment setup
	@echo -e "$(BLUE)Testing development environment...$(NC)"
	@# Check required tools
	@echo "Checking required development tools..."
	@command -v git >/dev/null || (echo -e "$(RED)❌ git not found$(NC)" && exit 1)
	@command -v jq >/dev/null || echo -e "$(YELLOW)⚠️ jq not found (recommended)$(NC)"
	@command -v curl >/dev/null || (echo -e "$(RED)❌ curl not found$(NC)" && exit 1)
	@# Check optional cloud tools
	@echo "Checking optional cloud tools..."
	@command -v aws >/dev/null && echo "✅ AWS CLI found" || echo "❌ AWS CLI not found"
	@command -v az >/dev/null && echo "✅ Azure CLI found" || echo "❌ Azure CLI not found"
	@command -v doctl >/dev/null && echo "✅ doctl found" || echo "❌ doctl not found"
	@command -v kubectl >/dev/null && echo "✅ kubectl found" || echo "❌ kubectl not found"
	@command -v helm >/dev/null && echo "✅ helm found" || echo "❌ helm not found"
	@command -v gh >/dev/null && echo "✅ GitHub CLI found" || echo "❌ GitHub CLI not found"
	@echo -e "$(GREEN)✅ Development environment check completed$(NC)"

# Benchmark tests
benchmark: ## Run performance benchmarks
	@echo -e "$(BLUE)Running performance benchmarks...$(NC)"
	@mkdir -p $(LOG_DIR)
	@echo "Benchmarking bootstrap loading..."
	@for i in {1..10}; do \
		time bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source core/bootstrap.sh; exit 0' 2>> $(LOG_DIR)/benchmark-load.log; \
	done
	@echo "Benchmark results saved to $(LOG_DIR)/benchmark-*.log"
	@echo -e "$(GREEN)✅ Benchmarks completed$(NC)"

# Stress tests
stress-test: ## Run stress tests
	@echo -e "$(BLUE)Running stress tests...$(NC)"
	@# Test rapid shell loading
	@echo "Testing rapid shell loading..."
	@for i in {1..50}; do \
		bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source core/bootstrap.sh; exit 0' >/dev/null 2>&1 || \
		(echo -e "$(RED)❌ Stress test failed at iteration $$i$(NC)" && exit 1); \
	done
	@echo -e "$(GREEN)✅ Stress tests passed$(NC)"

# Final comprehensive test
test-comprehensive: setup test-syntax test-dotfiles test-components test-auth-status test-required-tools test-install test-installer-components test-integration test-security test-performance test-workflows analyze-logs test-report ## Run all tests and generate report
	@echo -e "$(GREEN)🎉 Comprehensive test suite completed successfully!$(NC)"
	@echo -e "$(BLUE)📊 Test report available at: $(TEST_DIR)/test-report.md$(NC)"

# AI/ML Management Targets
# ========================

ai-setup: ## Setup complete AI/ML environment (Ollama + Hugging Face)
	@echo -e "$(BLUE)Setting up AI/ML environment...$(NC)"
	@echo "📊 Installing Ollama..."
	@bash -c 'source functions/ai/ollama.sh && ollama_setup'
	@echo ""
	@echo "🤗 Installing Hugging Face..."
	@bash -c 'source functions/ai/huggingface.sh && hf_setup'
	@echo -e "$(GREEN)✅ AI/ML environment setup complete$(NC)"

ai-status: ## Check AI/ML tools and models status
	@echo -e "$(BLUE)Checking AI/ML status...$(NC)"
	@echo "🤖 AI/ML Environment Status"
	@echo "==========================="
	@echo ""
	@# Check Python
	@command -v python3 >/dev/null && echo "✅ Python 3: $$(python3 --version)" || echo "❌ Python 3 not found"
	@command -v pip3 >/dev/null && echo "✅ pip3 available" || echo "❌ pip3 not found"
	@echo ""
	@echo "📊 Ollama Status:"
	@echo "================"
	@bash -c 'source functions/ai/ollama.sh && ollama_status' 2>/dev/null || echo "❌ Ollama not available"
	@echo ""
	@echo "🤗 Hugging Face Status:"
	@echo "======================"
	@bash -c 'source functions/ai/huggingface.sh && hf_status' 2>/dev/null || echo "❌ Hugging Face not available"
	@echo -e "$(GREEN)✅ AI/ML status check complete$(NC)"

ai-test: ## Test AI/ML functionality
	@echo -e "$(BLUE)Testing AI/ML functionality...$(NC)"
	@mkdir -p $(LOG_DIR)
	@# Test component syntax
	@echo "Testing AI component syntax..."
	@find components/ollama components/huggingface -name "*.sh" -type f 2>/dev/null | while read file; do \
		bash -n "$$file" || (echo -e "$(RED)❌ Syntax error in $$file$(NC)" && exit 1); \
	done
	@echo -e "$(GREEN)✅ AI/ML tests completed$(NC)"

ai-models: ## List all available AI models
	@echo -e "$(BLUE)Listing AI models...$(NC)"
	@echo "🤖 Available AI Models"
	@echo "====================="
	@echo ""
	@echo "📊 Ollama Models:"
	@echo "================"
	@bash -c 'source functions/ai/ollama.sh && ollama_models' 2>/dev/null || echo "❌ Ollama not installed"
	@echo ""
	@echo "🤗 Hugging Face Models:"
	@echo "======================"
	@bash -c 'source functions/ai/huggingface.sh && hf_models'

ai-chat: ## Start interactive AI chat (auto-detects best available model)
	@echo -e "$(BLUE)Starting AI chat...$(NC)"
	@if command -v ollama >/dev/null 2>&1; then \
		echo "🤖 Starting Ollama chat (llama3.2)..."; \
		bash -c 'source functions/ai/ollama.sh && ollama_run llama3.2'; \
	elif python3 -c "import transformers" 2>/dev/null; then \
		echo "🤗 Starting Hugging Face chat..."; \
		bash -c 'source functions/ai/huggingface.sh && hf_chat'; \
	else \
		echo "❌ No AI platforms available. Run: make ai-setup"; \
		exit 1; \
	fi

ai-chat-ollama: ## Start Ollama chat with Llama 3.2
	@echo -e "$(BLUE)Starting Ollama chat...$(NC)"
	@bash -c 'source functions/ai/ollama.sh && ollama_run llama3.2'

ai-chat-hf: ## Start Hugging Face chat
	@echo -e "$(BLUE)Starting Hugging Face chat...$(NC)"
	@bash -c 'source functions/ai/huggingface.sh && hf_chat'

ai-benchmark: ## Run AI performance benchmarks
	@echo -e "$(BLUE)Running AI benchmarks...$(NC)"
	@mkdir -p $(LOG_DIR)
	@echo "=== AI/ML Performance Benchmark ===" > $(LOG_DIR)/ai-benchmark.log
	@echo "Generated: $$(date)" >> $(LOG_DIR)/ai-benchmark.log
	@echo "" >> $(LOG_DIR)/ai-benchmark.log
	@command -v ollama >/dev/null && ollama list >> $(LOG_DIR)/ai-benchmark.log 2>&1 || \
		echo "Ollama not available" >> $(LOG_DIR)/ai-benchmark.log
	@echo "📊 Benchmark results saved to: $(LOG_DIR)/ai-benchmark.log"
	@echo -e "$(GREEN)✅ AI benchmarks completed$(NC)"

ai-examples: ## Show AI usage examples
	@echo -e "$(BLUE)AI/ML Usage Examples$(NC)"
	@echo "🎯 AI/ML Usage Examples"
	@echo "======================"
	@echo ""
	@echo "🚀 Quick Start:"
	@echo "  make ai-setup                    # Install everything"
	@echo "  make ai-chat                     # Start interactive chat"
	@echo "  ollama_chat llama3.2 'Hello'     # Quick Ollama question"
	@echo ""
	@echo "📊 Ollama Examples:"
	@echo "  ollama_install                   # Install Ollama"
	@echo "  ollama_pull llama3.2             # Download model"
	@echo "  ollama_run llama3.2              # Interactive chat"
	@echo "  ollama_code python 'sort list'   # Generate code"
	@echo ""
	@echo "🤗 Hugging Face Examples:"
	@echo "  hf_setup                         # Setup environment"
	@echo "  hf_generate 'Once upon'          # Generate text"
	@echo "  hf_sentiment 'I love AI'         # Sentiment analysis"
	@echo "  hf_summarize 'Long text'         # Summarize text"
	@echo ""
	@echo "🛠️ Management:"
	@echo "  make ai-status                   # Check all systems"
	@echo "  make ai-models                   # List all models"
	@echo "  make ai-cleanup                  # Clean up resources"

ai-examples-run: ## Run live AI examples (requires models)
	@echo -e "$(BLUE)Running live AI examples...$(NC)"
	@if command -v ollama >/dev/null 2>&1; then \
		echo "📊 Ollama example:"; \
		bash -c 'source functions/ai/ollama.sh && ollama_examples'; \
	fi
	@if python3 -c "import transformers" 2>/dev/null; then \
		echo "🤗 Hugging Face example:"; \
		bash -c 'source functions/ai/huggingface.sh && hf_examples'; \
	fi

ai-cleanup: ## Clean AI model caches and stop services
	@echo -e "$(BLUE)Cleaning up AI/ML resources...$(NC)"
	@# Stop Ollama service
	@if command -v ollama >/dev/null 2>&1; then \
		echo "🛑 Stopping Ollama service..."; \
		bash -c 'source functions/ai/ollama.sh && ollama_stop'; \
	fi
	@# Clean Hugging Face cache
	@if [ -d ~/.cache/huggingface ]; then \
		echo "🗑️ Cleaning Hugging Face cache..."; \
		bash -c 'source functions/ai/huggingface.sh && hf_clear_cache'; \
	fi
	@echo -e "$(GREEN)✅ AI/ML cleanup complete$(NC)"

# Ollama specific targets
ollama-install: ## Install Ollama
	@echo -e "$(BLUE)Installing Ollama...$(NC)"
	@bash -c 'source functions/ai/ollama.sh && ollama_install'

ollama-setup: ## Setup Ollama with recommended models
	@echo -e "$(BLUE)Setting up Ollama...$(NC)"
	@bash -c 'source functions/ai/ollama.sh && ollama_setup'

ollama-start: ## Start Ollama service
	@echo -e "$(BLUE)Starting Ollama service...$(NC)"
	@bash -c 'source functions/ai/ollama.sh && ollama_start'

ollama-stop: ## Stop Ollama service
	@echo -e "$(BLUE)Stopping Ollama service...$(NC)"
	@bash -c 'source functions/ai/ollama.sh && ollama_stop'

ollama-status: ## Check Ollama status
	@bash -c 'source functions/ai/ollama.sh && ollama_status'

ollama-models: ## List Ollama models
	@bash -c 'source functions/ai/ollama.sh && ollama_models'

# Hugging Face specific targets
hf-setup: ## Setup Hugging Face environment
	@echo -e "$(BLUE)Setting up Hugging Face...$(NC)"
	@bash -c 'source functions/ai/huggingface.sh && hf_setup'

hf-status: ## Check Hugging Face status
	@bash -c 'source functions/ai/huggingface.sh && hf_status'

hf-models: ## List popular Hugging Face models
	@bash -c 'source functions/ai/huggingface.sh && hf_models'

hf-examples: ## Run Hugging Face examples
	@bash -c 'source functions/ai/huggingface.sh && hf_examples'

hf-clear-cache: ## Clear Hugging Face model cache
	@echo -e "$(BLUE)Clearing Hugging Face cache...$(NC)"
	@bash -c 'source functions/ai/huggingface.sh && hf_clear_cache'

# AI testing in main test suite
test-ai-module: ## Test AI components
	@echo -e "$(BLUE)Testing AI/ML components...$(NC)"
	@# Test AI component syntax
	@find components/ollama components/huggingface -name "*.sh" -type f 2>/dev/null | while read file; do \
		bash -n "$$file" || (echo -e "$(RED)❌ Syntax error in $$file$(NC)" && exit 1); \
	done
	@echo -e "$(GREEN)✅ AI component tests passed$(NC)"

# AI tools status in required tools check
test-ai-tools: ## Test if AI/ML tools are available
	@echo -e "$(BLUE)Checking AI/ML tools...$(NC)"
	@echo "=== AI/ML Tools Status ===" >> $(LOG_DIR)/cli-tools-status.log
	@command -v python3 >/dev/null && echo "✅ Python 3 installed" || echo "❌ Python 3 not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@python3 -c "import pip" 2>/dev/null && echo "✅ pip available" || echo "❌ pip not available" | tee -a $(LOG_DIR)/cli-tools-status.log
	@command -v ollama >/dev/null && echo "✅ Ollama installed" || echo "❌ Ollama not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@python3 -c "import transformers" 2>/dev/null && echo "✅ Hugging Face transformers installed" || echo "❌ Hugging Face transformers not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@python3 -c "import torch" 2>/dev/null && echo "✅ PyTorch installed" || echo "❌ PyTorch not installed" | tee -a $(LOG_DIR)/cli-tools-status.log
	@echo "" >> $(LOG_DIR)/cli-tools-status.log
	@echo -e "$(GREEN)✅ AI/ML tools check completed$(NC)"

# =============================================================================
# Node.js/NVM Management Targets
# =============================================================================

# NVM directory (XDG-compliant)
NVM_DIR := $(XDG_DATA_HOME)/nvm
ifeq ($(NVM_DIR),/nvm)
NVM_DIR := $(HOME)/.local/share/nvm
endif

nvm-install: ## Install NVM (Node Version Manager)
	@echo -e "$(BLUE)Installing NVM...$(NC)"
	@if [ -d "$(NVM_DIR)" ] && [ -f "$(NVM_DIR)/nvm.sh" ]; then \
		echo -e "$(GREEN)✅ NVM already installed at $(NVM_DIR)$(NC)"; \
	else \
		echo "Downloading NVM installer..."; \
		mkdir -p "$(NVM_DIR)"; \
		export NVM_DIR="$(NVM_DIR)"; \
		curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash; \
		echo -e "$(GREEN)✅ NVM installed successfully$(NC)"; \
		echo ""; \
		echo "⚠️  Restart your shell or run:"; \
		echo "   source ~/.bashrc"; \
	fi

nvm-setup: nvm-install node-lts pnpm-install ## Complete NVM setup: install NVM + LTS Node + pnpm
	@echo -e "$(GREEN)✅ NVM setup complete!$(NC)"
	@echo ""
	@echo "🎉 Node.js environment ready:"
	@echo "   - NVM installed at: $(NVM_DIR)"
	@echo "   - Node.js LTS installed"
	@echo "   - pnpm available globally"
	@echo ""
	@echo "Restart your shell to apply changes."

nvm-status: ## Check NVM installation status
	@echo -e "$(BLUE)Checking NVM status...$(NC)"
	@echo "📦 NVM Status"
	@echo "============="
	@if [ -f "$(NVM_DIR)/nvm.sh" ]; then \
		echo "✅ NVM installed at: $(NVM_DIR)"; \
		bash -c 'export NVM_DIR="$(NVM_DIR)"; source "$(NVM_DIR)/nvm.sh"; echo "   Version: $$(nvm --version)"'; \
		bash -c 'export NVM_DIR="$(NVM_DIR)"; source "$(NVM_DIR)/nvm.sh"; echo "   Current Node: $$(nvm current)"'; \
		bash -c 'export NVM_DIR="$(NVM_DIR)"; source "$(NVM_DIR)/nvm.sh"; echo "   Installed versions:"; nvm ls --no-colors 2>/dev/null | head -10'; \
	else \
		echo "❌ NVM not installed"; \
		echo "   Run: make nvm-install"; \
	fi

node-install: ## Install latest Node.js via NVM
	@echo -e "$(BLUE)Installing latest Node.js...$(NC)"
	@if [ ! -f "$(NVM_DIR)/nvm.sh" ]; then \
		echo -e "$(RED)❌ NVM not installed. Run: make nvm-install$(NC)"; \
		exit 1; \
	fi
	@bash -c 'export NVM_DIR="$(NVM_DIR)"; source "$(NVM_DIR)/nvm.sh"; \
		echo "Installing latest Node.js..."; \
		nvm install node; \
		nvm use node; \
		nvm alias default node; \
		echo ""; \
		echo "✅ Node.js installed:"; \
		node --version; \
		npm --version'

node-lts: ## Install latest LTS Node.js via NVM
	@echo -e "$(BLUE)Installing Node.js LTS...$(NC)"
	@if [ ! -f "$(NVM_DIR)/nvm.sh" ]; then \
		echo -e "$(RED)❌ NVM not installed. Run: make nvm-install$(NC)"; \
		exit 1; \
	fi
	@bash -c 'export NVM_DIR="$(NVM_DIR)"; source "$(NVM_DIR)/nvm.sh"; \
		echo "Installing Node.js LTS..."; \
		nvm install --lts; \
		nvm use --lts; \
		nvm alias default lts/*; \
		echo ""; \
		echo "✅ Node.js LTS installed:"; \
		node --version; \
		npm --version'

pnpm-install: ## Install pnpm globally via npm
	@echo -e "$(BLUE)Installing pnpm...$(NC)"
	@if ! command -v node >/dev/null 2>&1; then \
		if [ -f "$(NVM_DIR)/nvm.sh" ]; then \
			bash -c 'export NVM_DIR="$(NVM_DIR)"; source "$(NVM_DIR)/nvm.sh"; \
				if ! command -v pnpm >/dev/null 2>&1; then \
					echo "Installing pnpm via npm..."; \
					npm install -g pnpm; \
					echo "✅ pnpm installed: $$(pnpm --version)"; \
				else \
					echo "✅ pnpm already installed: $$(pnpm --version)"; \
				fi'; \
		else \
			echo -e "$(RED)❌ Node.js not available. Run: make nvm-setup$(NC)"; \
			exit 1; \
		fi; \
	else \
		if ! command -v pnpm >/dev/null 2>&1; then \
			echo "Installing pnpm via npm..."; \
			npm install -g pnpm; \
			echo -e "$(GREEN)✅ pnpm installed: $$(pnpm --version)$(NC)"; \
		else \
			echo -e "$(GREEN)✅ pnpm already installed: $$(pnpm --version)$(NC)"; \
		fi; \
	fi

pnpm-setup: ## Setup pnpm with corepack (Node 16.13+)
	@echo -e "$(BLUE)Setting up pnpm with corepack...$(NC)"
	@if ! command -v node >/dev/null 2>&1; then \
		if [ -f "$(NVM_DIR)/nvm.sh" ]; then \
			bash -c 'export NVM_DIR="$(NVM_DIR)"; source "$(NVM_DIR)/nvm.sh"; \
				echo "Enabling corepack..."; \
				corepack enable 2>/dev/null || npm install -g corepack; \
				corepack prepare pnpm@latest --activate; \
				echo "✅ pnpm setup via corepack: $$(pnpm --version)"'; \
		else \
			echo -e "$(RED)❌ Node.js not available. Run: make nvm-setup$(NC)"; \
			exit 1; \
		fi; \
	else \
		echo "Enabling corepack..."; \
		corepack enable 2>/dev/null || npm install -g corepack; \
		corepack prepare pnpm@latest --activate 2>/dev/null || npm install -g pnpm; \
		echo -e "$(GREEN)✅ pnpm ready: $$(pnpm --version)$(NC)"; \
	fi

node-status: ## Check complete Node.js ecosystem status
	@echo -e "$(BLUE)Node.js Ecosystem Status$(NC)"
	@echo "========================="
	@echo ""
	@echo "📦 NVM:"
	@if [ -f "$(NVM_DIR)/nvm.sh" ]; then \
		echo "   ✅ Installed at: $(NVM_DIR)"; \
		bash -c 'export NVM_DIR="$(NVM_DIR)"; source "$(NVM_DIR)/nvm.sh" 2>/dev/null; echo "   Version: $$(nvm --version 2>/dev/null || echo "unknown")"'; \
	else \
		echo "   ❌ Not installed"; \
	fi
	@echo ""
	@echo "🟢 Node.js:"
	@if command -v node >/dev/null 2>&1; then \
		echo "   ✅ Version: $$(node --version)"; \
		echo "   Path: $$(which node)"; \
	elif [ -f "$(NVM_DIR)/nvm.sh" ]; then \
		bash -c 'export NVM_DIR="$(NVM_DIR)"; source "$(NVM_DIR)/nvm.sh"; \
			if command -v node >/dev/null 2>&1; then \
				echo "   ✅ Version: $$(node --version)"; \
				echo "   Path: $$(which node)"; \
			else \
				echo "   ❌ Not installed (run: make node-lts)"; \
			fi'; \
	else \
		echo "   ❌ Not installed"; \
	fi
	@echo ""
	@echo "📦 npm:"
	@if command -v npm >/dev/null 2>&1; then \
		echo "   ✅ Version: $$(npm --version)"; \
	elif [ -f "$(NVM_DIR)/nvm.sh" ]; then \
		bash -c 'export NVM_DIR="$(NVM_DIR)"; source "$(NVM_DIR)/nvm.sh"; \
			if command -v npm >/dev/null 2>&1; then \
				echo "   ✅ Version: $$(npm --version)"; \
			else \
				echo "   ❌ Not available"; \
			fi'; \
	else \
		echo "   ❌ Not available"; \
	fi
	@echo ""
	@echo "🚀 pnpm:"
	@if command -v pnpm >/dev/null 2>&1; then \
		echo "   ✅ Version: $$(pnpm --version)"; \
		echo "   Path: $$(which pnpm)"; \
	elif [ -f "$(NVM_DIR)/nvm.sh" ]; then \
		bash -c 'export NVM_DIR="$(NVM_DIR)"; source "$(NVM_DIR)/nvm.sh"; \
			if command -v pnpm >/dev/null 2>&1; then \
				echo "   ✅ Version: $$(pnpm --version)"; \
			else \
				echo "   ❌ Not installed (run: make pnpm-install)"; \
			fi'; \
	else \
		echo "   ❌ Not installed"; \
	fi
	@echo ""
	@echo "📦 Global packages:"
	@if command -v npm >/dev/null 2>&1; then \
		npm list -g --depth=0 2>/dev/null | tail -n +2 | head -10 || echo "   (none)"; \
	elif [ -f "$(NVM_DIR)/nvm.sh" ]; then \
		bash -c 'export NVM_DIR="$(NVM_DIR)"; source "$(NVM_DIR)/nvm.sh"; \
			npm list -g --depth=0 2>/dev/null | tail -n +2 | head -10 || echo "   (none)"'; \
	else \
		echo "   (npm not available)"; \
	fi

node-update: ## Update Node.js to latest LTS and reinstall globals
	@echo -e "$(BLUE)Updating Node.js...$(NC)"
	@if [ ! -f "$(NVM_DIR)/nvm.sh" ]; then \
		echo -e "$(RED)❌ NVM not installed$(NC)"; \
		exit 1; \
	fi
	@bash -c 'export NVM_DIR="$(NVM_DIR)"; source "$(NVM_DIR)/nvm.sh"; \
		echo "Current: $$(node --version 2>/dev/null || echo "none")"; \
		echo "Installing latest LTS..."; \
		nvm install --lts --reinstall-packages-from=current; \
		nvm alias default lts/*; \
		echo ""; \
		echo "✅ Updated to: $$(node --version)"'

# =============================================================================
# Shell Layer Targets (DAG Parity)
# =============================================================================
# These targets correspond to the shell loading DAG:
# discovery -> xdg -> theme -> env/secrets/options -> aliases -> functions -> prompt
# =============================================================================

shell-status: ## Show status of all shell modules
	@echo -e "$(BLUE)Shell Module Status$(NC)"
	@echo "==================="
	@echo ""
	@echo "📂 Shell Modules:"
	@for file in shell/discovery.sh shell/xdg.sh shell/theme.sh shell/environment.sh shell/secrets.sh shell/options.sh shell/aliases.sh shell/prompt.sh shell/completions.sh shell/init.sh; do \
		if [ -f "$$file" ]; then \
			echo "   ✅ $$file"; \
		else \
			echo "   ❌ $$file (missing)"; \
		fi; \
	done
	@echo ""
	@echo "📂 Function Modules:"
	@for dir in functions/core functions/dev functions/cloud functions/ai; do \
		if [ -d "$$dir" ]; then \
			count=$$(find "$$dir" -name "*.sh" 2>/dev/null | wc -l | tr -d ' '); \
			echo "   ✅ $$dir ($$count files)"; \
		else \
			echo "   ❌ $$dir (missing)"; \
		fi; \
	done
	@echo ""
	@echo "🔧 Current Environment:"
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh 2>/dev/null; \
		echo "   CURRENT_SHELL: $${CURRENT_SHELL:-not set}"; \
		echo "   CURRENT_PLATFORM: $${CURRENT_PLATFORM:-not set}"; \
		echo "   DOTFILES_ROOT: $${DOTFILES_ROOT:-not set}"; \
		echo "   XDG_CONFIG_HOME: $${XDG_CONFIG_HOME:-not set}"; \
		echo "   XDG_DATA_HOME: $${XDG_DATA_HOME:-not set}"'

shell-test-discovery: ## Test shell/discovery.sh module
	@echo -e "$(BLUE)Testing discovery.sh...$(NC)"
	@# Test CURRENT_SHELL detection
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source shell/discovery.sh; \
		[ "$$CURRENT_SHELL" = "bash" ] || (echo "❌ CURRENT_SHELL not set to bash" && exit 1); \
		echo "✅ CURRENT_SHELL=$$CURRENT_SHELL"'
	@# Test CURRENT_PLATFORM detection
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source shell/discovery.sh; \
		[ -n "$$CURRENT_PLATFORM" ] || (echo "❌ CURRENT_PLATFORM not set" && exit 1); \
		echo "✅ CURRENT_PLATFORM=$$CURRENT_PLATFORM"'
	@# Test IS_INTERACTIVE override
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/discovery.sh; \
		[ "$$IS_INTERACTIVE" = "true" ] || (echo "❌ IS_INTERACTIVE override failed" && exit 1); \
		echo "✅ IS_INTERACTIVE override works"'
	@echo -e "$(GREEN)✅ discovery.sh tests passed$(NC)"

shell-test-xdg: ## Test shell/xdg.sh module
	@echo -e "$(BLUE)Testing xdg.sh...$(NC)"
	@# Test XDG variables are set
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/discovery.sh; source shell/xdg.sh; \
		[ -n "$$XDG_CONFIG_HOME" ] || (echo "❌ XDG_CONFIG_HOME not set" && exit 1); \
		echo "✅ XDG_CONFIG_HOME=$$XDG_CONFIG_HOME"'
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/discovery.sh; source shell/xdg.sh; \
		[ -n "$$XDG_DATA_HOME" ] || (echo "❌ XDG_DATA_HOME not set" && exit 1); \
		echo "✅ XDG_DATA_HOME=$$XDG_DATA_HOME"'
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/discovery.sh; source shell/xdg.sh; \
		[ -n "$$XDG_CACHE_HOME" ] || (echo "❌ XDG_CACHE_HOME not set" && exit 1); \
		echo "✅ XDG_CACHE_HOME=$$XDG_CACHE_HOME"'
	@echo -e "$(GREEN)✅ xdg.sh tests passed$(NC)"

shell-test-theme: ## Test shell/theme.sh module
	@echo -e "$(BLUE)Testing theme.sh...$(NC)"
	@# Test theme color variables are set
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/discovery.sh; source shell/xdg.sh; source shell/theme.sh; \
		[ -n "$$THEME_GREEN" ] || (echo "❌ THEME_GREEN not set" && exit 1); \
		echo "✅ THEME_GREEN defined"'
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/discovery.sh; source shell/xdg.sh; source shell/theme.sh; \
		[ -n "$$THEME_GIT_CLEAN" ] || (echo "❌ THEME_GIT_CLEAN not set" && exit 1); \
		echo "✅ THEME_GIT_CLEAN defined"'
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/discovery.sh; source shell/xdg.sh; source shell/theme.sh; \
		[ -n "$$THEME_RESET" ] || (echo "❌ THEME_RESET not set" && exit 1); \
		echo "✅ THEME_RESET defined"'
	@echo -e "$(GREEN)✅ theme.sh tests passed$(NC)"

shell-test-env: ## Test shell/environment.sh module
	@echo -e "$(BLUE)Testing environment.sh...$(NC)"
	@# Test environment variables are set
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source shell/discovery.sh; source shell/xdg.sh; source shell/environment.sh; \
		[ -n "$$PYTHONSTARTUP" ] || (echo "❌ PYTHONSTARTUP not set" && exit 1); \
		echo "✅ PYTHONSTARTUP=$$PYTHONSTARTUP"'
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source shell/discovery.sh; source shell/xdg.sh; source shell/environment.sh; \
		[ -n "$$NVM_DIR" ] || (echo "❌ NVM_DIR not set" && exit 1); \
		echo "✅ NVM_DIR=$$NVM_DIR"'
	@echo -e "$(GREEN)✅ environment.sh tests passed$(NC)"

shell-test-options: ## Test shell/options.sh module
	@echo -e "$(BLUE)Testing options.sh...$(NC)"
	@# Test options file sources without error
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source shell/discovery.sh; source shell/xdg.sh; source shell/options.sh' || \
		(echo -e "$(RED)❌ options.sh failed to source$(NC)" && exit 1)
	@echo "✅ options.sh sources without error"
	@echo -e "$(GREEN)✅ options.sh tests passed$(NC)"

shell-test-aliases: ## Test shell/aliases.sh module
	@echo -e "$(BLUE)Testing aliases.sh...$(NC)"
	@# Test key aliases are defined
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source shell/discovery.sh; source shell/xdg.sh; source shell/aliases.sh; \
		alias | grep -q "ll=" || (echo "❌ ll alias not defined" && exit 1); \
		echo "✅ ll alias defined"'
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source shell/discovery.sh; source shell/xdg.sh; source shell/aliases.sh; \
		alias | grep -q "la=" || (echo "❌ la alias not defined" && exit 1); \
		echo "✅ la alias defined"'
	@echo -e "$(GREEN)✅ aliases.sh tests passed$(NC)"

shell-test-functions: ## Test all function modules
	@echo -e "$(BLUE)Testing function modules...$(NC)"
	@# Test core functions
	@echo "Testing core functions..."
	@for file in functions/core/*.sh; do \
		bash -n "$$file" || (echo -e "$(RED)❌ Syntax error in $$file$(NC)" && exit 1); \
	done
	@echo "✅ Core functions syntax OK"
	@# Test dev functions
	@echo "Testing dev functions..."
	@for file in functions/dev/*.sh; do \
		bash -n "$$file" || (echo -e "$(RED)❌ Syntax error in $$file$(NC)" && exit 1); \
	done
	@echo "✅ Dev functions syntax OK"
	@# Test cloud functions
	@echo "Testing cloud functions..."
	@for file in functions/cloud/*.sh; do \
		bash -n "$$file" || (echo -e "$(RED)❌ Syntax error in $$file$(NC)" && exit 1); \
	done
	@echo "✅ Cloud functions syntax OK"
	@# Test ai functions
	@echo "Testing AI functions..."
	@for file in functions/ai/*.sh; do \
		bash -n "$$file" || (echo -e "$(RED)❌ Syntax error in $$file$(NC)" && exit 1); \
	done
	@echo "✅ AI functions syntax OK"
	@# Test functions are loaded
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; \
		declare -f mkvenv >/dev/null || (echo "❌ mkvenv not loaded" && exit 1); \
		echo "✅ mkvenv function loaded"'
	@echo -e "$(GREEN)✅ Function modules tests passed$(NC)"

shell-test-prompt: ## Test shell/prompt.sh module
	@echo -e "$(BLUE)Testing prompt.sh...$(NC)"
	@# Test prompt functions exist
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; \
		declare -f git_branch >/dev/null || (echo "❌ git_branch not defined" && exit 1); \
		echo "✅ git_branch function defined"'
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; \
		declare -f git_color >/dev/null || (echo "❌ git_color not defined" && exit 1); \
		echo "✅ git_color function defined"'
	@# Test PS1 is set (bash)
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh; \
		[ -n "$$PS1" ] || (echo "❌ PS1 not set" && exit 1); \
		echo "✅ PS1 is set"'
	@echo -e "$(GREEN)✅ prompt.sh tests passed$(NC)"

shell-test-all: shell-test-discovery shell-test-xdg shell-test-theme shell-test-env shell-test-options shell-test-aliases shell-test-functions shell-test-prompt ## Test complete shell loading sequence
	@echo -e "$(GREEN)✅ All shell layer tests passed$(NC)"

# =============================================================================
# Component Management Targets
# =============================================================================
# Component-based architecture for modular dotfiles configuration
# Components are self-contained units with: component.yaml, env.sh, aliases.sh,
# functions.sh, completions.sh, and setup.sh
# =============================================================================

COMPONENTS_DIR := $(DOTFILES_DIR)/components
CORE_DIR := $(DOTFILES_DIR)/core

component-list: ## List all available components
	@echo -e "$(BLUE)Available Components$(NC)"
	@echo "===================="
	@echo ""
	@if [ -d "$(COMPONENTS_DIR)" ]; then \
		for dir in $(COMPONENTS_DIR)/*/; do \
			name=$$(basename "$$dir"); \
			if [ "$$name" != "_template" ] && [ -f "$$dir/component.yaml" ]; then \
				if command -v yq >/dev/null 2>&1; then \
					desc=$$(yq -r '.description // "No description"' "$$dir/component.yaml" 2>/dev/null); \
					version=$$(yq -r '.version // "0.0.0"' "$$dir/component.yaml" 2>/dev/null); \
					echo "  📦 $$name (v$$version)"; \
					echo "     $$desc"; \
				else \
					echo "  📦 $$name"; \
				fi; \
			fi; \
		done; \
	else \
		echo "  ❌ Components directory not found"; \
	fi
	@echo ""
	@echo "Template available at: $(COMPONENTS_DIR)/_template/"

component-status: ## Show component health and loading status
	@echo -e "$(BLUE)Component Status$(NC)"
	@echo "================"
	@echo ""
	@echo "📂 Core Framework:"
	@for file in bootstrap.sh discovery.sh xdg.sh theme.sh loader.sh sync.sh; do \
		if [ -f "$(CORE_DIR)/$$file" ]; then \
			echo "   ✅ core/$$file"; \
		else \
			echo "   ❌ core/$$file (missing)"; \
		fi; \
	done
	@echo ""
	@echo "📦 Components:"
	@if [ -d "$(COMPONENTS_DIR)" ]; then \
		for dir in $(COMPONENTS_DIR)/*/; do \
			name=$$(basename "$$dir"); \
			if [ "$$name" != "_template" ]; then \
				if [ -f "$$dir/component.yaml" ]; then \
					echo "   ✅ $$name"; \
					for f in env.sh aliases.sh functions.sh completions.sh; do \
						if [ -f "$$dir/$$f" ]; then \
							echo "      ├── $$f"; \
						fi; \
					done; \
				else \
					echo "   ❌ $$name (missing component.yaml)"; \
				fi; \
			fi; \
		done; \
	fi
	@echo ""
	@echo "🔧 Dependencies:"
	@if command -v yq >/dev/null 2>&1; then \
		echo "   ✅ yq installed (required for YAML parsing)"; \
	else \
		echo "   ❌ yq not installed (install with: brew install yq)"; \
	fi

component-new: ## Create new component from template (usage: make component-new NAME=mycomponent)
	@if [ -z "$(NAME)" ]; then \
		echo -e "$(RED)Error: NAME is required$(NC)"; \
		echo "Usage: make component-new NAME=mycomponent"; \
		exit 1; \
	fi
	@if [ -d "$(COMPONENTS_DIR)/$(NAME)" ]; then \
		echo -e "$(RED)Error: Component '$(NAME)' already exists$(NC)"; \
		exit 1; \
	fi
	@echo -e "$(BLUE)Creating component: $(NAME)$(NC)"
	@cp -r "$(COMPONENTS_DIR)/_template" "$(COMPONENTS_DIR)/$(NAME)"
	@# Update component.yaml with the new name
	@if command -v sed >/dev/null 2>&1; then \
		sed -i.bak "s/name: template/name: $(NAME)/" "$(COMPONENTS_DIR)/$(NAME)/component.yaml" && \
		rm -f "$(COMPONENTS_DIR)/$(NAME)/component.yaml.bak"; \
	fi
	@echo -e "$(GREEN)✅ Created component at: $(COMPONENTS_DIR)/$(NAME)/$(NC)"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Edit component.yaml with metadata"
	@echo "  2. Add environment variables to env.sh"
	@echo "  3. Add aliases to aliases.sh"
	@echo "  4. Add functions to functions.sh"
	@echo "  5. Add completions to completions.sh"

component-validate: ## Validate all component.yaml files
	@echo -e "$(BLUE)Validating component.yaml files...$(NC)"
	@if ! command -v yq >/dev/null 2>&1; then \
		echo -e "$(RED)Error: yq is required for validation$(NC)"; \
		echo "Install with: brew install yq"; \
		exit 1; \
	fi
	@errors=0; \
	for dir in $(COMPONENTS_DIR)/*/; do \
		name=$$(basename "$$dir"); \
		if [ "$$name" != "_template" ] && [ -f "$$dir/component.yaml" ]; then \
			if yq eval '.' "$$dir/component.yaml" >/dev/null 2>&1; then \
				comp_name=$$(yq -r '.name // ""' "$$dir/component.yaml"); \
				if [ -z "$$comp_name" ]; then \
					echo "   ❌ $$name: missing 'name' field"; \
					errors=$$((errors + 1)); \
				else \
					echo "   ✅ $$name"; \
				fi; \
			else \
				echo "   ❌ $$name: invalid YAML syntax"; \
				errors=$$((errors + 1)); \
			fi; \
		fi; \
	done; \
	if [ $$errors -gt 0 ]; then \
		echo -e "$(RED)❌ $$errors validation errors found$(NC)"; \
		exit 1; \
	fi
	@echo -e "$(GREEN)✅ All component.yaml files are valid$(NC)"

test-components: test-component-loader test-component-deps ## Test component system
	@echo -e "$(GREEN)✅ All component tests passed$(NC)"

test-component-loader: ## Test component loader functionality
	@echo -e "$(BLUE)Testing component loader...$(NC)"
	@# Test core files exist
	@[ -f "$(CORE_DIR)/bootstrap.sh" ] || (echo -e "$(RED)❌ core/bootstrap.sh not found$(NC)" && exit 1)
	@[ -f "$(CORE_DIR)/loader.sh" ] || (echo -e "$(RED)❌ core/loader.sh not found$(NC)" && exit 1)
	@echo "✅ Core files found"
	@# Test bootstrap can be sourced
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source core/bootstrap.sh 2>/dev/null; \
		[ -n "$$DOTFILES_COMPONENTS_LOADED" ] || (echo "❌ Components not loaded" && exit 1)' && \
		echo "✅ Bootstrap sources successfully" || \
		echo -e "$(YELLOW)⚠️ Bootstrap may have issues$(NC)"
	@# Test component discovery
	@if [ -d "$(COMPONENTS_DIR)" ]; then \
		count=$$(find $(COMPONENTS_DIR) -name "component.yaml" -not -path "*/_template/*" 2>/dev/null | wc -l | tr -d ' '); \
		echo "✅ Found $$count components"; \
	fi
	@echo -e "$(GREEN)✅ Component loader tests passed$(NC)"

test-component-deps: ## Test component dependency resolution
	@echo -e "$(BLUE)Testing component dependencies...$(NC)"
	@if ! command -v yq >/dev/null 2>&1; then \
		echo -e "$(YELLOW)⚠️ yq not installed, skipping dependency tests$(NC)"; \
		exit 0; \
	fi
	@# Check for circular dependencies (basic check)
	@echo "Checking for dependency issues..."
	@for dir in $(COMPONENTS_DIR)/*/; do \
		name=$$(basename "$$dir"); \
		if [ "$$name" != "_template" ] && [ -f "$$dir/component.yaml" ]; then \
			deps=$$(yq -r '.requires.components // [] | .[]' "$$dir/component.yaml" 2>/dev/null); \
			for dep in $$deps; do \
				if [ ! -d "$(COMPONENTS_DIR)/$$dep" ]; then \
					echo "   ⚠️ $$name requires missing component: $$dep"; \
				fi; \
			done; \
		fi; \
	done
	@echo -e "$(GREEN)✅ Dependency check completed$(NC)"

# =============================================================================
# Dotfiles Management Targets
# =============================================================================

dotfiles-install: ## Run full dotfiles installation (install.sh)
	@echo -e "$(BLUE)Running dotfiles installation...$(NC)"
	@bash install.sh
	@echo -e "$(GREEN)✅ Dotfiles installation complete$(NC)"

dotfiles-inject: ## Create shell bootstrap files (~/.bashrc, ~/.zshrc)
	@echo -e "$(BLUE)Injecting shell bootstrap files...$(NC)"
	@# Create ~/.bashrc bootstrap
	@if [ ! -f ~/.bashrc ] || ! grep -q "DOTFILES_ROOT" ~/.bashrc 2>/dev/null; then \
		echo "Creating ~/.bashrc..."; \
		echo '# Dotfiles bootstrap - component-based architecture' > ~/.bashrc; \
		echo 'export DOTFILES_ROOT="$(DOTFILES_DIR)"' >> ~/.bashrc; \
		echo '[ -f "$$DOTFILES_ROOT/core/bootstrap.sh" ] && . "$$DOTFILES_ROOT/core/bootstrap.sh"' >> ~/.bashrc; \
		echo "✅ Created ~/.bashrc"; \
	else \
		echo "⚠️  ~/.bashrc already configured"; \
	fi
	@# Create ~/.zshrc bootstrap
	@if [ ! -f ~/.zshrc ] || ! grep -q "DOTFILES_ROOT" ~/.zshrc 2>/dev/null; then \
		echo "Creating ~/.zshrc..."; \
		echo '# Dotfiles bootstrap - component-based architecture' > ~/.zshrc; \
		echo 'export DOTFILES_ROOT="$(DOTFILES_DIR)"' >> ~/.zshrc; \
		echo '[ -f "$$DOTFILES_ROOT/core/bootstrap.sh" ] && . "$$DOTFILES_ROOT/core/bootstrap.sh"' >> ~/.zshrc; \
		echo "✅ Created ~/.zshrc"; \
	else \
		echo "⚠️  ~/.zshrc already configured"; \
	fi
	@# Create ~/.bash_profile to source ~/.bashrc
	@if [ ! -f ~/.bash_profile ] || ! grep -q "bashrc" ~/.bash_profile 2>/dev/null; then \
		echo "Creating ~/.bash_profile..."; \
		echo '# Source bashrc for login shells' > ~/.bash_profile; \
		echo '[ -f ~/.bashrc ] && . ~/.bashrc' >> ~/.bash_profile; \
		echo "✅ Created ~/.bash_profile"; \
	else \
		echo "⚠️  ~/.bash_profile already configured"; \
	fi
	@echo -e "$(GREEN)✅ Bootstrap files created$(NC)"

dotfiles-eject: ## Remove shell bootstrap files
	@echo -e "$(YELLOW)Removing shell bootstrap files...$(NC)"
	@echo "This will remove:"
	@echo "  - ~/.bashrc (if managed by dotfiles)"
	@echo "  - ~/.zshrc (if managed by dotfiles)"
	@echo "  - ~/.bash_profile (if managed by dotfiles)"
	@read -p "Continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		if grep -q "DOTFILES_ROOT" ~/.bashrc 2>/dev/null; then rm ~/.bashrc && echo "Removed ~/.bashrc"; fi; \
		if grep -q "DOTFILES_ROOT" ~/.zshrc 2>/dev/null; then rm ~/.zshrc && echo "Removed ~/.zshrc"; fi; \
		if grep -q "bashrc" ~/.bash_profile 2>/dev/null; then rm ~/.bash_profile && echo "Removed ~/.bash_profile"; fi; \
		echo -e "$(GREEN)✅ Bootstrap files removed$(NC)"; \
	else \
		echo "Cancelled"; \
	fi

dotfiles-link: ## Link app configurations (git, ghostty, vscode, claude)
	@echo -e "$(BLUE)Linking app configurations...$(NC)"
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source functions/core/inject.sh; inject_configs'
	@echo -e "$(GREEN)✅ App configurations linked$(NC)"

dotfiles-unlink: ## Remove app configuration links
	@echo -e "$(YELLOW)Removing app configuration links...$(NC)"
	@# Git
	@if [ -L ~/.gitconfig ]; then rm ~/.gitconfig && echo "Removed ~/.gitconfig"; fi
	@if [ -L ~/.gitignore ]; then rm ~/.gitignore && echo "Removed ~/.gitignore"; fi
	@# Ghostty
	@if [ -L ~/.config/ghostty/config ]; then rm ~/.config/ghostty/config && echo "Removed ghostty config"; fi
	@# VS Code
	@if [ -L ~/Library/Application\ Support/Code/User/settings.json ]; then rm ~/Library/Application\ Support/Code/User/settings.json && echo "Removed VS Code settings"; fi
	@# Claude
	@if [ -L ~/.config/claude/settings.json ]; then rm ~/.config/claude/settings.json && echo "Removed Claude settings"; fi
	@echo -e "$(GREEN)✅ App configuration links removed$(NC)"

dotfiles-status: ## Show dotfiles installation status
	@echo -e "$(BLUE)Dotfiles Installation Status$(NC)"
	@echo "============================="
	@echo ""
	@echo "📂 Repository: $(DOTFILES_DIR)"
	@echo ""
	@echo "🔗 Bootstrap Files:"
	@if grep -q "DOTFILES_ROOT" ~/.bashrc 2>/dev/null; then echo "   ✅ ~/.bashrc"; else echo "   ❌ ~/.bashrc"; fi
	@if grep -q "DOTFILES_ROOT" ~/.zshrc 2>/dev/null; then echo "   ✅ ~/.zshrc"; else echo "   ❌ ~/.zshrc"; fi
	@if [ -f ~/.bash_profile ]; then echo "   ✅ ~/.bash_profile"; else echo "   ❌ ~/.bash_profile"; fi
	@echo ""
	@echo "🔗 App Configs:"
	@if [ -L ~/.gitconfig ]; then echo "   ✅ Git (~/.gitconfig)"; else echo "   ❌ Git (~/.gitconfig)"; fi
	@if [ -L ~/.config/ghostty/config ]; then echo "   ✅ Ghostty"; else echo "   ❌ Ghostty"; fi
	@if [ -L ~/Library/Application\ Support/Code/User/settings.json ] 2>/dev/null; then echo "   ✅ VS Code"; else echo "   ❌ VS Code"; fi
	@if [ -L ~/.config/claude/settings.json ]; then echo "   ✅ Claude"; else echo "   ❌ Claude"; fi

dotfiles-reload: ## Reload shell configuration
	@echo -e "$(BLUE)Reloading shell configuration...$(NC)"
	@echo "Run this command to reload:"
	@echo ""
	@echo "  source ~/.bashrc   # for bash"
	@echo "  source ~/.zshrc    # for zsh"
	@echo ""
	@echo "Or start a new shell session."

dotfiles-update: ## Git pull and show reload instructions
	@echo -e "$(BLUE)Updating dotfiles...$(NC)"
	@git pull --rebase
	@echo ""
	@echo -e "$(GREEN)✅ Dotfiles updated$(NC)"
	@echo "Run 'source ~/.bashrc' or 'source ~/.zshrc' to reload."

# =============================================================================
# Python/UV Management Targets
# =============================================================================

uv-install: ## Install UV package manager
	@echo -e "$(BLUE)Installing UV...$(NC)"
	@if command -v uv >/dev/null 2>&1; then \
		echo -e "$(GREEN)✅ UV already installed: $$(uv --version)$(NC)"; \
	else \
		echo "Downloading UV installer..."; \
		curl -LsSf https://astral.sh/uv/install.sh | sh; \
		echo -e "$(GREEN)✅ UV installed$(NC)"; \
		echo ""; \
		echo "⚠️  Add UV to your PATH or restart your shell"; \
	fi

uv-setup: uv-install ## Complete UV setup with Python
	@echo -e "$(BLUE)Setting up UV environment...$(NC)"
	@# Ensure UV is in PATH for this session
	@export PATH="$$HOME/.local/bin:$$PATH"; \
	if command -v uv >/dev/null 2>&1; then \
		echo "Installing Python via UV..."; \
		uv python install 3.12 2>/dev/null || echo "Python 3.12 already installed or use system Python"; \
		echo ""; \
		echo -e "$(GREEN)✅ UV setup complete$(NC)"; \
	else \
		echo -e "$(RED)❌ UV not found in PATH$(NC)"; \
		exit 1; \
	fi

uv-status: ## Check UV and Python status
	@echo -e "$(BLUE)UV/Python Status$(NC)"
	@echo "================"
	@echo ""
	@echo "📦 UV:"
	@if command -v uv >/dev/null 2>&1; then \
		echo "   ✅ Installed: $$(uv --version)"; \
		echo "   Path: $$(which uv)"; \
	elif [ -f "$$HOME/.local/bin/uv" ]; then \
		echo "   ✅ Installed: $$($$HOME/.local/bin/uv --version)"; \
		echo "   Path: $$HOME/.local/bin/uv"; \
		echo "   ⚠️  Not in PATH"; \
	else \
		echo "   ❌ Not installed"; \
	fi
	@echo ""
	@echo "🐍 Python:"
	@if command -v python3 >/dev/null 2>&1; then \
		echo "   ✅ Version: $$(python3 --version)"; \
		echo "   Path: $$(which python3)"; \
	else \
		echo "   ❌ Not installed"; \
	fi
	@echo ""
	@echo "📦 pip:"
	@if command -v pip3 >/dev/null 2>&1; then \
		echo "   ✅ Version: $$(pip3 --version | awk '{print $$2}')"; \
	elif python3 -m pip --version >/dev/null 2>&1; then \
		echo "   ✅ Version: $$(python3 -m pip --version | awk '{print $$2}')"; \
	else \
		echo "   ❌ Not available"; \
	fi
	@echo ""
	@echo "🔧 Virtual Environment:"
	@if [ -n "$$VIRTUAL_ENV" ]; then \
		echo "   ✅ Active: $$VIRTUAL_ENV"; \
	else \
		echo "   ❌ None active"; \
	fi

python-status: uv-status ## Alias for uv-status

venv-create: ## Create a virtual environment using UV or venv
	@echo -e "$(BLUE)Creating virtual environment...$(NC)"
	@if command -v uv >/dev/null 2>&1; then \
		echo "Using UV to create .venv..."; \
		uv venv .venv; \
	else \
		echo "Using python3 -m venv..."; \
		python3 -m venv .venv; \
	fi
	@echo -e "$(GREEN)✅ Virtual environment created at .venv$(NC)"
	@echo ""
	@echo "Activate with: source .venv/bin/activate"

venv-status: ## Show active virtual environment info
	@echo -e "$(BLUE)Virtual Environment Status$(NC)"
	@echo "=========================="
	@if [ -n "$$VIRTUAL_ENV" ]; then \
		echo "✅ Active: $$VIRTUAL_ENV"; \
		echo "   Python: $$(python --version)"; \
		if command -v uv >/dev/null 2>&1; then \
			echo "   UV: $$(uv --version)"; \
		fi; \
		echo ""; \
		echo "📦 Installed packages:"; \
		pip list 2>/dev/null | head -15 || echo "   (none)"; \
	else \
		echo "❌ No virtual environment active"; \
		echo ""; \
		echo "To create: make venv-create"; \
		echo "To activate: source .venv/bin/activate"; \
	fi

# =============================================================================
# Go Management Targets
# =============================================================================

go-install: ## Install Go via Homebrew (macOS) or download (Linux)
	@echo -e "$(BLUE)Installing Go...$(NC)"
	@if command -v go >/dev/null 2>&1; then \
		echo -e "$(GREEN)✅ Go already installed: $$(go version)$(NC)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		echo "Installing Go via Homebrew..."; \
		brew install go; \
		echo -e "$(GREEN)✅ Go installed$(NC)"; \
	else \
		echo "Installing Go from official source..."; \
		curl -LO https://go.dev/dl/go1.22.0.linux-amd64.tar.gz; \
		sudo rm -rf /usr/local/go; \
		sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz; \
		rm go1.22.0.linux-amd64.tar.gz; \
		echo "Add to PATH: export PATH=\$$PATH:/usr/local/go/bin"; \
		echo -e "$(GREEN)✅ Go installed$(NC)"; \
	fi

go-setup: go-install ## Setup Go environment
	@echo -e "$(BLUE)Setting up Go environment...$(NC)"
	@# Create Go workspace directories
	@mkdir -p ~/go/{bin,src,pkg}
	@echo "✅ Go workspace created at ~/go"
	@# Show environment
	@if command -v go >/dev/null 2>&1; then \
		echo ""; \
		echo "Go Environment:"; \
		go env GOROOT GOPATH GOBIN; \
	fi
	@echo -e "$(GREEN)✅ Go setup complete$(NC)"

go-status: ## Check Go installation status
	@echo -e "$(BLUE)Go Status$(NC)"
	@echo "========="
	@echo ""
	@if command -v go >/dev/null 2>&1; then \
		echo "✅ Go installed"; \
		echo "   Version: $$(go version | awk '{print $$3}')"; \
		echo "   Path: $$(which go)"; \
		echo "   GOROOT: $$(go env GOROOT)"; \
		echo "   GOPATH: $$(go env GOPATH)"; \
		echo ""; \
		echo "📦 Installed tools:"; \
		ls $$(go env GOPATH)/bin 2>/dev/null | head -10 || echo "   (none)"; \
	else \
		echo "❌ Go not installed"; \
		echo "   Run: make go-install"; \
	fi

go-tools: ## Install common Go development tools
	@echo -e "$(BLUE)Installing Go tools...$(NC)"
	@if ! command -v go >/dev/null 2>&1; then \
		echo -e "$(RED)❌ Go not installed. Run: make go-install$(NC)"; \
		exit 1; \
	fi
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing gopls (language server)..."
	@go install golang.org/x/tools/gopls@latest
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Installing dlv (debugger)..."
	@go install github.com/go-delve/delve/cmd/dlv@latest
	@echo -e "$(GREEN)✅ Go tools installed$(NC)"

# =============================================================================
# Acorn CLI Build System
# =============================================================================

# Build configuration
BINARY_NAME := acorn
GO_MODULE := github.com/mistergrinvalds/acorn
VERSION_PKG := $(GO_MODULE)/internal/version
BUILD_DIR := bin
CMD_DIR := cmd/acorn

# Version information (extracted from git)
GIT_VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILT_BY := $(shell whoami)

# Build flags for version embedding
LDFLAGS := -s -w \
	-X '$(VERSION_PKG).Version=$(GIT_VERSION)' \
	-X '$(VERSION_PKG).Commit=$(GIT_COMMIT)' \
	-X '$(VERSION_PKG).Date=$(BUILD_DATE)' \
	-X '$(VERSION_PKG).BuiltBy=$(BUILT_BY)'

# Cross-compilation targets
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

.PHONY: acorn acorn-build acorn-install acorn-clean acorn-test acorn-lint acorn-fmt acorn-vet acorn-check acorn-release acorn-snapshot acorn-deps acorn-tidy acorn-verify acorn-coverage acorn-bench acorn-security

acorn: acorn-build ## Build acorn CLI (alias)

acorn-build: ## Build acorn CLI binary
	@echo -e "$(BLUE)Building acorn CLI...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo -e "$(GREEN)✅ Built $(BUILD_DIR)/$(BINARY_NAME)$(NC)"
	@echo "   Version: $(GIT_VERSION)"
	@echo "   Commit:  $(GIT_COMMIT)"

acorn-install: acorn-build ## Install acorn to GOPATH/bin
	@echo -e "$(BLUE)Installing acorn...$(NC)"
	@go install -ldflags "$(LDFLAGS)" ./$(CMD_DIR)
	@echo -e "$(GREEN)✅ Installed acorn to $$(go env GOPATH)/bin$(NC)"

acorn-clean: ## Clean build artifacts
	@echo -e "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@go clean -cache -testcache
	@echo -e "$(GREEN)✅ Clean complete$(NC)"

acorn-test: ## Run all tests with race detection
	@echo -e "$(BLUE)Running tests...$(NC)"
	@go test -race -v ./...
	@echo -e "$(GREEN)✅ Tests passed$(NC)"

acorn-coverage: ## Run tests with coverage report
	@echo -e "$(BLUE)Running tests with coverage...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go test -race -coverprofile=$(BUILD_DIR)/coverage.out -covermode=atomic ./...
	@go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@go tool cover -func=$(BUILD_DIR)/coverage.out | tail -1
	@echo -e "$(GREEN)✅ Coverage report: $(BUILD_DIR)/coverage.html$(NC)"

acorn-bench: ## Run benchmarks
	@echo -e "$(BLUE)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...
	@echo -e "$(GREEN)✅ Benchmarks complete$(NC)"

acorn-lint: ## Run golangci-lint
	@echo -e "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout 5m; \
		echo -e "$(GREEN)✅ Linting passed$(NC)"; \
	else \
		echo -e "$(YELLOW)⚠️  golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
		exit 1; \
	fi

acorn-fmt: ## Format Go code
	@echo -e "$(BLUE)Formatting code...$(NC)"
	@gofmt -s -w .
	@goimports -w . 2>/dev/null || true
	@echo -e "$(GREEN)✅ Code formatted$(NC)"

acorn-vet: ## Run go vet
	@echo -e "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo -e "$(GREEN)✅ Vet passed$(NC)"

acorn-security: ## Run security scanners (gosec, govulncheck)
	@echo -e "$(BLUE)Running security checks...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		echo "Running gosec..."; \
		gosec -quiet ./...; \
	else \
		echo -e "$(YELLOW)⚠️  gosec not installed. Run: go install github.com/securego/gosec/v2/cmd/gosec@latest$(NC)"; \
	fi
	@if command -v govulncheck >/dev/null 2>&1; then \
		echo "Running govulncheck..."; \
		govulncheck ./...; \
	else \
		echo -e "$(YELLOW)⚠️  govulncheck not installed. Run: go install golang.org/x/vuln/cmd/govulncheck@latest$(NC)"; \
	fi
	@echo -e "$(GREEN)✅ Security checks complete$(NC)"

acorn-check: acorn-fmt acorn-vet acorn-lint acorn-test ## Run all checks (fmt, vet, lint, test)
	@echo -e "$(GREEN)✅ All checks passed$(NC)"

acorn-deps: ## Install development dependencies
	@echo -e "$(BLUE)Installing development dependencies...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/goreleaser/goreleaser/v2@latest
	@echo -e "$(GREEN)✅ Dependencies installed$(NC)"

acorn-tidy: ## Tidy go modules
	@echo -e "$(BLUE)Tidying modules...$(NC)"
	@go mod tidy
	@go mod verify
	@echo -e "$(GREEN)✅ Modules tidied$(NC)"

acorn-verify: ## Verify module checksums
	@echo -e "$(BLUE)Verifying modules...$(NC)"
	@go mod verify
	@echo -e "$(GREEN)✅ Modules verified$(NC)"

acorn-cross: ## Build for all platforms
	@echo -e "$(BLUE)Cross-compiling for all platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		go build -ldflags "$(LDFLAGS)" \
			-o $(BUILD_DIR)/$(BINARY_NAME)-$${platform%/*}-$${platform#*/}$$([ "$${platform%/*}" = "windows" ] && echo ".exe") \
			./$(CMD_DIR); \
		echo "   Built: $(BINARY_NAME)-$${platform%/*}-$${platform#*/}"; \
	done
	@echo -e "$(GREEN)✅ Cross-compilation complete$(NC)"

acorn-snapshot: ## Build snapshot release with goreleaser
	@echo -e "$(BLUE)Building snapshot release...$(NC)"
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --clean; \
		echo -e "$(GREEN)✅ Snapshot built in dist/$(NC)"; \
	else \
		echo -e "$(YELLOW)⚠️  goreleaser not installed. Run: make acorn-deps$(NC)"; \
		exit 1; \
	fi

acorn-release: ## Create tagged release with goreleaser
	@echo -e "$(BLUE)Creating release...$(NC)"
	@if [ -z "$(shell git tag -l --points-at HEAD)" ]; then \
		echo -e "$(RED)❌ No tag on current commit. Tag with: git tag vX.Y.Z$(NC)"; \
		exit 1; \
	fi
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --clean; \
		echo -e "$(GREEN)✅ Release created$(NC)"; \
	else \
		echo -e "$(YELLOW)⚠️  goreleaser not installed. Run: make acorn-deps$(NC)"; \
		exit 1; \
	fi

acorn-version: ## Show version that would be embedded
	@echo "Version:  $(GIT_VERSION)"
	@echo "Commit:   $(GIT_COMMIT)"
	@echo "Date:     $(BUILD_DATE)"
	@echo "Built By: $(BUILT_BY)"

# =============================================================================
# Brew Package Management
# =============================================================================

# Lists of managed packages
BREW_DEVOPS := argocd awscli azure-cli cloudflared doctl helm k9s kind kubernetes-cli kustomize terraform vault
BREW_DEV := bash bash-completion ctags fd fzf gh jq lazygit neovim shellcheck tmux yq
BREW_DB := kafka mongocli mongosh mycli neo4j pgcli postgresql@14 redis iredis

brew-status: ## Show status of all managed brew packages
	@echo -e "$(BLUE)Brew Package Status$(NC)"
	@echo "==================="
	@echo ""
	@echo -e "$(BLUE)DevOps Tools:$(NC)"
	@for pkg in $(BREW_DEVOPS); do \
		if brew list --formula 2>/dev/null | grep -q "^$$pkg$$"; then \
			version=$$(brew info $$pkg 2>/dev/null | head -1 | awk '{print $$3}'); \
			printf "   ✅ $$pkg ($$version)\n"; \
		else \
			printf "   ❌ $$pkg\n"; \
		fi; \
	done
	@echo ""
	@echo -e "$(BLUE)Dev Tools:$(NC)"
	@for pkg in $(BREW_DEV); do \
		if brew list --formula 2>/dev/null | grep -q "^$$pkg$$"; then \
			version=$$(brew info $$pkg 2>/dev/null | head -1 | awk '{print $$3}'); \
			printf "   ✅ $$pkg ($$version)\n"; \
		else \
			printf "   ❌ $$pkg\n"; \
		fi; \
	done
	@echo ""
	@echo -e "$(BLUE)Database Tools:$(NC)"
	@for pkg in $(BREW_DB); do \
		if brew list --formula 2>/dev/null | grep -q "^$$pkg$$"; then \
			version=$$(brew info $$pkg 2>/dev/null | head -1 | awk '{print $$3}'); \
			printf "   ✅ $$pkg ($$version)\n"; \
		else \
			printf "   ❌ $$pkg\n"; \
		fi; \
	done

brew-update: ## Update all brew packages
	@echo -e "$(BLUE)Updating Homebrew and packages...$(NC)"
	@brew update
	@brew upgrade
	@brew cleanup
	@echo -e "$(GREEN)✅ Brew packages updated$(NC)"

brew-install-devops: ## Install all DevOps tools via brew
	@echo -e "$(BLUE)Installing DevOps tools...$(NC)"
	@brew install $(BREW_DEVOPS) || true
	@echo -e "$(GREEN)✅ DevOps tools installed$(NC)"

brew-install-dev: ## Install all dev tools via brew
	@echo -e "$(BLUE)Installing dev tools...$(NC)"
	@brew install $(BREW_DEV) || true
	@echo -e "$(GREEN)✅ Dev tools installed$(NC)"

brew-install-db: ## Install all database tools via brew
	@echo -e "$(BLUE)Installing database tools...$(NC)"
	@brew install $(BREW_DB) || true
	@echo -e "$(GREEN)✅ Database tools installed$(NC)"

brew-install-all: brew-install-devops brew-install-dev brew-install-db ## Install all managed brew packages
	@echo -e "$(GREEN)✅ All brew packages installed$(NC)"

# Individual database tool installs
db-install-mysql: ## Install MySQL client + mycli
	@echo -e "$(BLUE)Installing MySQL tools...$(NC)"
	@brew install mysql-client mycli
	@echo -e "$(GREEN)✅ MySQL tools installed$(NC)"

db-install-mongo: ## Install MongoDB shell + mongocli
	@echo -e "$(BLUE)Installing MongoDB tools...$(NC)"
	@brew install mongosh mongocli
	@echo -e "$(GREEN)✅ MongoDB tools installed$(NC)"

db-install-redis: ## Install Redis + iredis
	@echo -e "$(BLUE)Installing Redis tools...$(NC)"
	@brew install redis iredis
	@echo -e "$(GREEN)✅ Redis tools installed$(NC)"

db-install-neo4j: ## Install Neo4j
	@echo -e "$(BLUE)Installing Neo4j...$(NC)"
	@brew install neo4j
	@echo -e "$(GREEN)✅ Neo4j installed$(NC)"

db-install-kafka: ## Install Kafka
	@echo -e "$(BLUE)Installing Kafka...$(NC)"
	@brew install kafka
	@echo -e "$(GREEN)✅ Kafka installed$(NC)"

# =============================================================================
# Comprehensive Status Target
# =============================================================================

status: ## Show complete environment status
	@echo -e "$(BLUE)╔════════════════════════════════════════════════════════════════╗$(NC)"
	@echo -e "$(BLUE)║          Complete Environment Status                           ║$(NC)"
	@echo -e "$(BLUE)╚════════════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@# Shell Detection
	@echo -e "$(BLUE)🐚 Shell Environment$(NC)"
	@echo "==================="
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh 2>/dev/null; \
		echo "   Shell: $${CURRENT_SHELL:-unknown}"; \
		echo "   Platform: $${CURRENT_PLATFORM:-unknown}"; \
		echo "   Interactive: $${IS_INTERACTIVE:-unknown}"; \
		echo "   DOTFILES_ROOT: $${DOTFILES_ROOT:-not set}"'
	@echo ""
	@# XDG Directories
	@echo -e "$(BLUE)📁 XDG Directories$(NC)"
	@echo "=================="
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh 2>/dev/null; \
		echo "   XDG_CONFIG_HOME: $${XDG_CONFIG_HOME:-not set}"; \
		echo "   XDG_DATA_HOME: $${XDG_DATA_HOME:-not set}"; \
		echo "   XDG_CACHE_HOME: $${XDG_CACHE_HOME:-not set}"'
	@echo ""
	@# Development Tools
	@echo -e "$(BLUE)🛠️  Development Tools$(NC)"
	@echo "===================="
	@printf "   Git: "; command -v git >/dev/null && echo "✅ $$(git --version | awk '{print $$3}')" || echo "❌"
	@printf "   Node.js: "; command -v node >/dev/null && echo "✅ $$(node --version)" || echo "❌"
	@printf "   Python: "; command -v python3 >/dev/null && echo "✅ $$(python3 --version | awk '{print $$2}')" || echo "❌"
	@printf "   Go: "; command -v go >/dev/null && echo "✅ $$(go version | awk '{print $$3}')" || echo "❌"
	@printf "   UV: "; command -v uv >/dev/null && echo "✅ $$(uv --version 2>/dev/null)" || echo "❌"
	@printf "   pnpm: "; command -v pnpm >/dev/null && echo "✅ $$(pnpm --version)" || echo "❌"
	@printf "   fzf: "; command -v fzf >/dev/null && echo "✅ $$(fzf --version | cut -d' ' -f1)" || echo "❌"
	@printf "   fd: "; command -v fd >/dev/null && echo "✅ $$(fd --version | awk '{print $$2}')" || echo "❌"
	@printf "   jq: "; command -v jq >/dev/null && echo "✅ $$(jq --version)" || echo "❌"
	@printf "   lazygit: "; command -v lazygit >/dev/null && echo "✅ installed" || echo "❌"
	@echo ""
	@# DevOps Tools
	@echo -e "$(BLUE)☸️  DevOps Tools$(NC)"
	@echo "==============="
	@printf "   kubectl: "; command -v kubectl >/dev/null && echo "✅ $$(kubectl version --client -o yaml 2>/dev/null | grep gitVersion | awk '{print $$2}')" || echo "❌"
	@printf "   helm: "; command -v helm >/dev/null && echo "✅ $$(helm version --short 2>/dev/null)" || echo "❌"
	@printf "   terraform: "; command -v terraform >/dev/null && echo "✅ $$(terraform version -json 2>/dev/null | jq -r '.terraform_version' 2>/dev/null || terraform version | head -1 | awk '{print $$2}')" || echo "❌"
	@printf "   aws: "; command -v aws >/dev/null && echo "✅ $$(aws --version 2>/dev/null | awk '{print $$1}' | cut -d/ -f2)" || echo "❌"
	@printf "   az: "; command -v az >/dev/null && echo "✅ $$(az version 2>/dev/null | jq -r '.\"azure-cli\"' 2>/dev/null)" || echo "❌"
	@printf "   doctl: "; command -v doctl >/dev/null && echo "✅ $$(doctl version 2>/dev/null | head -1 | awk '{print $$3}')" || echo "❌"
	@printf "   vault: "; command -v vault >/dev/null && echo "✅ $$(vault version 2>/dev/null | awk '{print $$2}')" || echo "❌"
	@printf "   k9s: "; command -v k9s >/dev/null && echo "✅ installed" || echo "❌"
	@printf "   argocd: "; command -v argocd >/dev/null && echo "✅ installed" || echo "❌"
	@echo ""
	@# Database Tools
	@echo -e "$(BLUE)🗄️  Database Tools$(NC)"
	@echo "================="
	@printf "   pgcli: "; command -v pgcli >/dev/null && echo "✅ $$(pgcli --version 2>/dev/null | awk '{print $$2}')" || echo "❌"
	@printf "   psql: "; command -v psql >/dev/null && echo "✅ $$(psql --version 2>/dev/null | awk '{print $$3}')" || echo "❌"
	@printf "   mycli: "; command -v mycli >/dev/null && echo "✅ installed" || echo "❌"
	@printf "   mongosh: "; command -v mongosh >/dev/null && echo "✅ installed" || echo "❌"
	@printf "   redis-cli: "; command -v redis-cli >/dev/null && echo "✅ $$(redis-cli --version 2>/dev/null | awk '{print $$2}')" || echo "❌"
	@printf "   sqlite3: "; command -v sqlite3 >/dev/null && echo "✅ $$(sqlite3 --version 2>/dev/null | awk '{print $$1}')" || echo "❌"
	@echo ""
	@# AI/ML Tools
	@echo -e "$(BLUE)🤖 AI/ML Tools$(NC)"
	@echo "=============="
	@printf "   Ollama: "; command -v ollama >/dev/null && echo "✅ installed" || echo "❌"
	@printf "   HuggingFace: "; python3 -c "import transformers" 2>/dev/null && echo "✅ installed" || echo "❌"
	@echo ""
	@# Installation Status
	@echo -e "$(BLUE)📦 Installation Status$(NC)"
	@echo "====================="
	@printf "   ~/.bashrc: "; grep -q "DOTFILES_ROOT" ~/.bashrc 2>/dev/null && echo "✅" || echo "❌"
	@printf "   ~/.zshrc: "; grep -q "DOTFILES_ROOT" ~/.zshrc 2>/dev/null && echo "✅" || echo "❌"
	@printf "   Git config: "; [ -L ~/.gitconfig ] && echo "✅" || echo "❌"
	@printf "   Ghostty: "; [ -L ~/.config/ghostty/config ] && echo "✅" || echo "❌"
	@echo ""
	@echo -e "$(GREEN)Run 'make help' for available commands$(NC)"

# Default test for CI
.DEFAULT_GOAL := test-quick
# =============================================================================
# Claude Code Integration
# =============================================================================

.PHONY: claude claude-context claude-resume claude-status

CLAUDE_CONTEXT_FILE := $(DOTFILES_DIR)/.claude/SESSION_CONTEXT.md

claude: ## Launch Claude Code with session context (1M token context enabled)
	@echo "🤖 Launching Claude Code with session context..."
	@echo ""
	@if [ ! -f "$(CLAUDE_CONTEXT_FILE)" ]; then \
		echo "⚠️  No context file found at $(CLAUDE_CONTEXT_FILE)"; \
		echo "   Starting fresh session..."; \
		cd $(DOTFILES_DIR) && claude; \
	else \
		echo "📋 Loading context from: $(CLAUDE_CONTEXT_FILE)"; \
		echo ""; \
		cd $(DOTFILES_DIR) && claude "Read the file .claude/SESSION_CONTEXT.md to understand this repository's current state and architecture. Then ask me what I'd like to work on."; \
	fi

claude-context: ## Show current session context
	@if [ -f "$(CLAUDE_CONTEXT_FILE)" ]; then \
		cat "$(CLAUDE_CONTEXT_FILE)"; \
	else \
		echo "No context file found at $(CLAUDE_CONTEXT_FILE)"; \
	fi

claude-resume: ## Resume most recent Claude session
	@echo "🔄 Resuming most recent Claude session..."
	@cd $(DOTFILES_DIR) && claude --continue

claude-status: ## Show Claude Code configuration status
	@echo "🤖 Claude Code Status"
	@echo "===================="
	@echo ""
	@echo "Settings:"
	@if [ -f "$(DOTFILES_DIR)/.claude/settings.local.json" ]; then \
		echo "  Project settings: $(DOTFILES_DIR)/.claude/settings.local.json"; \
		echo "  Model: $$(jq -r '.model // "default"' $(DOTFILES_DIR)/.claude/settings.local.json)"; \
		echo "  Large Context: $$(jq -r '.largeContext // false' $(DOTFILES_DIR)/.claude/settings.local.json)"; \
	else \
		echo "  No project settings found"; \
	fi
	@echo ""
	@echo "Context File:"
	@if [ -f "$(CLAUDE_CONTEXT_FILE)" ]; then \
		echo "  ✅ $(CLAUDE_CONTEXT_FILE)"; \
		echo "  Size: $$(wc -l < $(CLAUDE_CONTEXT_FILE)) lines"; \
	else \
		echo "  ❌ No context file"; \
	fi
	@echo ""
	@echo "Claude Code Tooling:"
	@echo "  Commands: $$(ls -1 $(DOTFILES_DIR)/.claude/commands/*.md 2>/dev/null | wc -l | tr -d ' ') custom commands"
	@echo "  Agents: $$(ls -1 $(DOTFILES_DIR)/.claude/agents/*.md 2>/dev/null | wc -l | tr -d ' ') custom agents"
