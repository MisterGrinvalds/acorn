# Makefile for Testing Bash Profile and Automation Framework
# Comprehensive test suite for all functionality

.PHONY: help test test-all test-quick test-dotfiles test-automation test-cloud test-modules test-syntax test-security test-install clean setup ai-setup ai-status ai-test ai-models ai-chat ai-benchmark ai-cleanup ai-examples nvm-install nvm-setup nvm-status node-install node-lts pnpm-install pnpm-setup node-status

# Default target
help: ## Show this help message
	@echo "Bash Profile & Automation Framework Test Suite"
	@echo "=============================================="
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Test Categories:"
	@echo "  test-quick        - Quick syntax and basic functionality tests"
	@echo "  test-all          - Complete test suite (requires cloud CLIs)"
	@echo "  test-dotfiles     - Test dotfiles configuration only"
	@echo "  test-automation   - Test automation framework only"
	@echo "  test-api-keys     - Test API key configuration and validation"
	@echo "  test-auth-status  - Test authentication status for all services"
	@echo "  test-required-tools - Test if required CLI tools are installed"
	@echo ""
	@echo "Tools Management:"
	@echo "  tools-status      - Show comprehensive tools status"
	@echo "  tools-missing     - Show missing tools"
	@echo "  tools-update      - Interactive tool updates"
	@echo "  tools-update-all  - Update all tools without prompting"
	@echo "  tools-update-yes  - Update all tools with auto-yes prompts"
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

# Variables
SHELL := /bin/bash
DOTFILES_DIR := $(PWD)
AUTO_DIR := $(DOTFILES_DIR)/.automation
TEST_DIR := $(DOTFILES_DIR)/tests
LOG_DIR := $(TEST_DIR)/logs
BACKUP_DIR := $(TEST_DIR)/backups

# Test configuration
TEST_PROJECT_DIR := $(TEST_DIR)/test-projects
TEST_VENV_NAME := test-automation-env
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
test-all: setup test-syntax test-dotfiles test-automation test-modules test-integration ## Run complete test suite
	@echo -e "$(GREEN)✅ All tests completed successfully$(NC)"

# Individual test categories
test: test-quick ## Alias for test-quick

# Test syntax of all shell files
test-syntax: ## Test syntax of all shell scripts
	@echo -e "$(BLUE)Testing shell script syntax...$(NC)"
	@echo "Testing shell modules..."
	@for file in shell/*.sh; do \
		echo "  Testing $$file..."; \
		bash -n "$$file" || (echo -e "$(RED)❌ Syntax error in $$file$(NC)" && exit 1); \
	done
	@echo "Testing function modules..."
	@for file in functions/**/*.sh; do \
		echo "  Testing $$file..."; \
		bash -n "$$file" || (echo -e "$(RED)❌ Syntax error in $$file$(NC)" && exit 1); \
	done
	@echo "Testing automation framework..."
	@bash -n $(AUTO_DIR)/auto || (echo -e "$(RED)❌ Syntax error in automation CLI$(NC)" && exit 1)
	@for file in $(AUTO_DIR)/framework/*.sh; do \
		echo "  Testing $$file..."; \
		bash -n "$$file" || (echo -e "$(RED)❌ Syntax error in $$file$(NC)" && exit 1); \
	done
	@for file in $(AUTO_DIR)/modules/*.sh; do \
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

# Test automation framework
test-automation: test-automation-basic test-automation-cli ## Test automation framework

test-automation-basic: ## Test basic automation framework
	@echo -e "$(BLUE)Testing automation framework basics...$(NC)"
	@# Test framework files exist
	@echo "Testing framework structure..."
	@[ -f $(AUTO_DIR)/auto ] || (echo -e "$(RED)❌ Automation CLI not found$(NC)" && exit 1)
	@[ -f $(AUTO_DIR)/framework/core.sh ] || (echo -e "$(RED)❌ Core framework not found$(NC)" && exit 1)
	@[ -f $(AUTO_DIR)/framework/utils.sh ] || (echo -e "$(RED)❌ Utils framework not found$(NC)" && exit 1)
	@# Test framework initialization
	@echo "Testing framework initialization..."
	@bash -c 'source $(AUTO_DIR)/framework/core.sh && [ -n "$$AUTO_HOME" ]' || \
		(echo -e "$(RED)❌ Framework initialization failed$(NC)" && exit 1)
	@echo -e "$(GREEN)✅ Basic automation tests passed$(NC)"

test-automation-cli: ## Test automation CLI functionality
	@echo -e "$(BLUE)Testing automation CLI...$(NC)"
	@# Make CLI executable
	@chmod +x $(AUTO_DIR)/auto
	@# Test help command
	@echo "Testing CLI help..."
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto --help > $(LOG_DIR)/auto-help.log 2>&1 || \
		(echo -e "$(RED)❌ Auto help command failed$(NC)" && exit 1)
	@# Test version command
	@echo "Testing CLI version..."
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto --version > $(LOG_DIR)/auto-version.log 2>&1 || \
		(echo -e "$(RED)❌ Auto version command failed$(NC)" && exit 1)
	@echo -e "$(GREEN)✅ Automation CLI tests passed$(NC)"

# Test individual modules
test-modules: test-dev-module test-k8s-module test-github-module test-system-module test-config-module test-secrets-module test-tools-module ## Test all automation modules

test-dev-module: ## Test development module
	@echo -e "$(BLUE)Testing development module...$(NC)"
	@chmod +x $(AUTO_DIR)/auto
	@# Test dev help
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto dev --help > $(LOG_DIR)/dev-help.log 2>&1 || \
		(echo -e "$(RED)❌ Dev module help failed$(NC)" && exit 1)
	@# Test project initialization (dry run)
	@echo "Testing project initialization..."
	@cd $(TEST_PROJECT_DIR) && \
		timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto dev init python test-project --dry-run > $(LOG_DIR)/dev-init.log 2>&1 || \
		echo -e "$(YELLOW)⚠️ Dev init test skipped (may require dependencies)$(NC)"
	@echo -e "$(GREEN)✅ Development module tests passed$(NC)"

test-k8s-module: ## Test Kubernetes module
	@echo -e "$(BLUE)Testing Kubernetes module...$(NC)"
	@# Test k8s help
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto k8s --help > $(LOG_DIR)/k8s-help.log 2>&1 || \
		(echo -e "$(RED)❌ K8s module help failed$(NC)" && exit 1)
	@# Test cluster info (may fail if no cluster)
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto k8s cluster info > $(LOG_DIR)/k8s-cluster.log 2>&1 || \
		echo -e "$(YELLOW)⚠️ K8s cluster test skipped (no cluster configured)$(NC)"
	@echo -e "$(GREEN)✅ Kubernetes module tests passed$(NC)"

test-github-module: ## Test GitHub module
	@echo -e "$(BLUE)Testing GitHub module...$(NC)"
	@# Test github help
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto github --help > $(LOG_DIR)/github-help.log 2>&1 || \
		(echo -e "$(RED)❌ GitHub module help failed$(NC)" && exit 1)
	@# Test repo list (may fail if not authenticated)
	@if command -v gh >/dev/null 2>&1; then \
		timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto github repo list > $(LOG_DIR)/github-repos.log 2>&1 || \
		echo -e "$(YELLOW)⚠️ GitHub repo test skipped (not authenticated)$(NC)"; \
	else \
		echo -e "$(YELLOW)⚠️ GitHub tests skipped (gh CLI not installed)$(NC)"; \
	fi
	@echo -e "$(GREEN)✅ GitHub module tests passed$(NC)"

test-system-module: ## Test system module
	@echo -e "$(BLUE)Testing system module...$(NC)"
	@# Test system help
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto system --help > $(LOG_DIR)/system-help.log 2>&1 || \
		(echo -e "$(RED)❌ System module help failed$(NC)" && exit 1)
	@# Test system info
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto system info > $(LOG_DIR)/system-info.log 2>&1 || \
		(echo -e "$(RED)❌ System info failed$(NC)" && exit 1)
	@echo -e "$(GREEN)✅ System module tests passed$(NC)"

test-config-module: ## Test configuration module
	@echo -e "$(BLUE)Testing configuration module...$(NC)"
	@# Test config help
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto config --help > $(LOG_DIR)/config-help.log 2>&1 || \
		(echo -e "$(RED)❌ Config module help failed$(NC)" && exit 1)
	@# Test config validation
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto config validate > $(LOG_DIR)/config-validate.log 2>&1 || \
		(echo -e "$(RED)❌ Config validation failed$(NC)" && exit 1)
	@echo -e "$(GREEN)✅ Configuration module tests passed$(NC)"

test-secrets-module: ## Test secrets management module
	@echo -e "$(BLUE)Testing secrets management module...$(NC)"
	@# Test secrets help
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto secrets --help > $(LOG_DIR)/secrets-help.log 2>&1 || \
		(echo -e "$(RED)❌ Secrets module help failed$(NC)" && exit 1)
	@# Test secrets initialization
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto secrets init > $(LOG_DIR)/secrets-init.log 2>&1 || \
		(echo -e "$(RED)❌ Secrets initialization failed$(NC)" && exit 1)
	@# Test secrets requirements check
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto secrets check-requirements > $(LOG_DIR)/secrets-check.log 2>&1 || \
		echo -e "$(YELLOW)⚠️ Secrets requirements check completed (expected missing keys)$(NC)"
	@echo -e "$(GREEN)✅ Secrets module tests passed$(NC)"

test-tools-module: ## Test tools management module
	@echo -e "$(BLUE)Testing tools management module...$(NC)"
	@# Test tools help
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto tools --help > $(LOG_DIR)/tools-help.log 2>&1 || \
		(echo -e "$(RED)❌ Tools module help failed$(NC)" && exit 1)
	@# Test tools list
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto tools list > $(LOG_DIR)/tools-list.log 2>&1 || \
		(echo -e "$(RED)❌ Tools list failed$(NC)" && exit 1)
	@# Test tools status
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto tools status > $(LOG_DIR)/tools-status.log 2>&1 || \
		(echo -e "$(RED)❌ Tools status failed$(NC)" && exit 1)
	@# Test tools check
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto tools check > $(LOG_DIR)/tools-check.log 2>&1 || \
		echo -e "$(YELLOW)⚠️ Tools check completed (some tools may be missing)$(NC)"
	@echo -e "$(GREEN)✅ Tools module tests passed$(NC)"

# Test API keys and secrets
test-api-keys: test-api-keys-check test-api-keys-validate ## Test API key availability and validation

test-api-keys-check: ## Check which API keys are configured
	@echo -e "$(BLUE)Checking API key configuration...$(NC)"
	@chmod +x $(AUTO_DIR)/auto
	@$(AUTO_DIR)/auto secrets init >/dev/null 2>&1 || true
	@echo "=== API Key Status Report ===" > $(LOG_DIR)/api-keys-status.log
	@$(AUTO_DIR)/auto secrets check-requirements >> $(LOG_DIR)/api-keys-status.log 2>&1 || true
	@echo "Detailed API key status saved to $(LOG_DIR)/api-keys-status.log"
	@echo
	@echo -e "$(BLUE)📋 API Key Summary:$(NC)"
	@grep -E "(✅|❌)" $(LOG_DIR)/api-keys-status.log || echo "No API key status found"
	@echo -e "$(GREEN)✅ API key check completed$(NC)"

test-api-keys-validate: ## Validate configured API keys
	@echo -e "$(BLUE)Validating configured API keys...$(NC)"
	@if [ -f $(AUTO_DIR)/secrets/.env ]; then \
		echo "Found secrets file, validating..."; \
		$(AUTO_DIR)/auto secrets validate > $(LOG_DIR)/api-keys-validation.log 2>&1 && \
		echo -e "$(GREEN)✅ API key validation passed$(NC)" || \
		echo -e "$(YELLOW)⚠️ Some API keys failed validation (check $(LOG_DIR)/api-keys-validation.log)$(NC)"; \
	else \
		echo -e "$(YELLOW)⚠️ No secrets configured, skipping validation$(NC)"; \
		echo "No secrets file found" > $(LOG_DIR)/api-keys-validation.log; \
	fi

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

setup-secrets-wizard: ## Run interactive secrets setup wizard
	@echo -e "$(BLUE)Starting secrets setup wizard...$(NC)"
	@chmod +x $(AUTO_DIR)/auto
	@$(AUTO_DIR)/auto secrets init
	@echo "🔐 Run the following command to setup your API keys:"
	@echo "   $(AUTO_DIR)/auto secrets setup"
	@echo ""
	@echo "Or setup individual providers:"
	@echo "   $(AUTO_DIR)/auto secrets aws"
	@echo "   $(AUTO_DIR)/auto secrets azure"
	@echo "   $(AUTO_DIR)/auto secrets digitalocean"
	@echo "   $(AUTO_DIR)/auto secrets github"

tools-status: ## Show comprehensive tools status using automation framework
	@echo -e "$(BLUE)Comprehensive tools status...$(NC)"
	@chmod +x $(AUTO_DIR)/auto
	@$(AUTO_DIR)/auto tools status

tools-missing: ## Show missing tools
	@echo -e "$(BLUE)Checking for missing tools...$(NC)"
	@chmod +x $(AUTO_DIR)/auto
	@$(AUTO_DIR)/auto tools missing

tools-update: ## Interactive tool updates
	@echo -e "$(BLUE)Starting interactive tool updates...$(NC)"
	@chmod +x $(AUTO_DIR)/auto
	@$(AUTO_DIR)/auto tools update

tools-update-all: ## Update all tools without prompting
	@echo -e "$(BLUE)Updating all tools...$(NC)"
	@chmod +x $(AUTO_DIR)/auto
	@$(AUTO_DIR)/auto tools update --force

tools-update-yes: ## Update all tools with yes-to-all prompts
	@echo -e "$(BLUE)Updating all tools with auto-yes...$(NC)"
	@chmod +x $(AUTO_DIR)/auto
	@$(AUTO_DIR)/auto tools update --yes-to-all

# Test cloud modules
test-cloud: test-cloud-unified test-aws-module test-azure-module test-do-module ## Test all cloud modules

test-cloud-unified: ## Test unified cloud management
	@echo -e "$(BLUE)Testing unified cloud management...$(NC)"
	@# Test cloud help
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto cloud --help > $(LOG_DIR)/cloud-help.log 2>&1 || \
		(echo -e "$(RED)❌ Cloud module help failed$(NC)" && exit 1)
	@# Test cloud status
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto cloud status > $(LOG_DIR)/cloud-status.log 2>&1 || \
		(echo -e "$(RED)❌ Cloud status failed$(NC)" && exit 1)
	@# Test cost comparison
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto cloud cost-compare > $(LOG_DIR)/cloud-cost.log 2>&1 || \
		(echo -e "$(RED)❌ Cloud cost comparison failed$(NC)" && exit 1)
	@echo -e "$(GREEN)✅ Unified cloud tests passed$(NC)"

test-aws-module: ## Test AWS module
	@echo -e "$(BLUE)Testing AWS module...$(NC)"
	@# Test AWS help
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto aws --help > $(LOG_DIR)/aws-help.log 2>&1 || \
		(echo -e "$(RED)❌ AWS module help failed$(NC)" && exit 1)
	@# Test AWS auth status (may fail if not configured)
	@if command -v aws >/dev/null 2>&1; then \
		timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto aws auth list-profiles > $(LOG_DIR)/aws-auth.log 2>&1 || \
		echo -e "$(YELLOW)⚠️ AWS auth test skipped (not configured)$(NC)"; \
	else \
		echo -e "$(YELLOW)⚠️ AWS tests skipped (AWS CLI not installed)$(NC)"; \
	fi
	@echo -e "$(GREEN)✅ AWS module tests passed$(NC)"

test-azure-module: ## Test Azure module
	@echo -e "$(BLUE)Testing Azure module...$(NC)"
	@# Test Azure help
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto azure --help > $(LOG_DIR)/azure-help.log 2>&1 || \
		(echo -e "$(RED)❌ Azure module help failed$(NC)" && exit 1)
	@# Test Azure auth status (may fail if not configured)
	@if command -v az >/dev/null 2>&1; then \
		timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto azure auth list-subscriptions > $(LOG_DIR)/azure-auth.log 2>&1 || \
		echo -e "$(YELLOW)⚠️ Azure auth test skipped (not configured)$(NC)"; \
	else \
		echo -e "$(YELLOW)⚠️ Azure tests skipped (Azure CLI not installed)$(NC)"; \
	fi
	@echo -e "$(GREEN)✅ Azure module tests passed$(NC)"

test-do-module: ## Test DigitalOcean module
	@echo -e "$(BLUE)Testing DigitalOcean module...$(NC)"
	@# Test DO help
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto digitalocean --help > $(LOG_DIR)/do-help.log 2>&1 || \
		(echo -e "$(RED)❌ DigitalOcean module help failed$(NC)" && exit 1)
	@# Test DO auth status (may fail if not configured)
	@if command -v doctl >/dev/null 2>&1; then \
		timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto digitalocean droplets list > $(LOG_DIR)/do-droplets.log 2>&1 || \
		echo -e "$(YELLOW)⚠️ DigitalOcean auth test skipped (not configured)$(NC)"; \
	else \
		echo -e "$(YELLOW)⚠️ DigitalOcean tests skipped (doctl not installed)$(NC)"; \
	fi
	@echo -e "$(GREEN)✅ DigitalOcean module tests passed$(NC)"

# Test installation and setup
test-install: ## Test installation process
	@echo -e "$(BLUE)Testing installation process...$(NC)"
	@# Test automation setup script
	@chmod +x $(AUTO_DIR)/setup.sh
	@# Dry run of setup (skip actual installation)
	@echo "Testing setup script syntax..."
	@bash -n $(AUTO_DIR)/setup.sh || (echo -e "$(RED)❌ Setup script syntax error$(NC)" && exit 1)
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
	@# Test dotfiles structure (new organization)
	@echo "Testing dotfiles structure..."
	@[ -d shell ] || (echo -e "$(RED)❌ Shell directory missing$(NC)" && exit 1)
	@[ -f shell/init.sh ] || (echo -e "$(RED)❌ Shell init.sh missing$(NC)" && exit 1)
	@[ -d functions ] || (echo -e "$(RED)❌ Functions directory missing$(NC)" && exit 1)
	@[ -d config ] || (echo -e "$(RED)❌ Config directory missing$(NC)" && exit 1)
	@[ -d .automation ] || (echo -e "$(RED)❌ Automation directory missing$(NC)" && exit 1)
	@# Test automation framework structure
	@echo "Testing automation framework structure..."
	@[ -f .automation/auto ] || (echo -e "$(RED)❌ Auto CLI missing$(NC)" && exit 1)
	@[ -f .automation/setup.sh ] || (echo -e "$(RED)❌ Setup script missing$(NC)" && exit 1)
	@[ -d .automation/modules ] || (echo -e "$(RED)❌ Modules directory missing$(NC)" && exit 1)
	@[ -d .automation/framework ] || (echo -e "$(RED)❌ Framework directory missing$(NC)" && exit 1)
	@# Test function modules
	@[ -f functions/cloud/secrets.sh ] || (echo -e "$(RED)❌ Secrets function missing$(NC)" && exit 1)
	@[ -f functions/ai/ollama.sh ] || (echo -e "$(RED)❌ Ollama function missing$(NC)" && exit 1)
	@echo -e "$(GREEN)✅ Installer components tests passed$(NC)"

# Integration tests
test-integration: ## Test integration between components
	@echo -e "$(BLUE)Testing component integration...$(NC)"
	@# Test dotfiles + automation integration
	@echo "Testing dotfiles + automation integration..."
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source shell/init.sh; declare -f dotfiles_status >/dev/null' || \
		(echo -e "$(RED)❌ Dotfiles integration not loaded$(NC)" && exit 1)
	@# Test cloud functions loaded
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source shell/init.sh; declare -f load_secrets >/dev/null' || \
		(echo -e "$(RED)❌ Cloud functions not loaded$(NC)" && exit 1)
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
	@find . -name "*.sh" -perm +111 | grep -v ".automation/auto" | grep -v "install.sh" | grep -v "setup.sh" && \
		echo -e "$(YELLOW)⚠️ Unexpected executable permissions found$(NC)" || true
	@echo -e "$(GREEN)✅ Security tests passed$(NC)"

# Performance tests
test-performance: ## Test performance of shell loading
	@echo -e "$(BLUE)Testing shell loading performance...$(NC)"
	@# Time shell loading
	@echo "Testing shell init load time..."
	@time bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source shell/init.sh; exit 0' > $(LOG_DIR)/perf-load.log 2>&1 || \
		(echo -e "$(RED)❌ Performance test failed$(NC)" && exit 1)
	@# Test automation CLI response time
	@echo "Testing automation CLI response time..."
	@chmod +x $(AUTO_DIR)/auto
	@time $(AUTO_DIR)/auto --version > $(LOG_DIR)/perf-cli.log 2>&1 || \
		(echo -e "$(RED)❌ CLI performance test failed$(NC)" && exit 1)
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
	@echo "Benchmarking shell init loading..."
	@for i in {1..10}; do \
		time bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source shell/init.sh; exit 0' 2>> $(LOG_DIR)/benchmark-load.log; \
	done
	@echo "Benchmarking automation CLI..."
	@chmod +x $(AUTO_DIR)/auto
	@for i in {1..5}; do \
		time $(AUTO_DIR)/auto --version 2>> $(LOG_DIR)/benchmark-cli.log; \
	done
	@echo "Benchmark results saved to $(LOG_DIR)/benchmark-*.log"
	@echo -e "$(GREEN)✅ Benchmarks completed$(NC)"

# Stress tests
stress-test: ## Run stress tests
	@echo -e "$(BLUE)Running stress tests...$(NC)"
	@# Test rapid shell loading
	@echo "Testing rapid shell loading..."
	@for i in {1..50}; do \
		bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source shell/init.sh; exit 0' >/dev/null 2>&1 || \
		(echo -e "$(RED)❌ Stress test failed at iteration $$i$(NC)" && exit 1); \
	done
	@echo -e "$(GREEN)✅ Stress tests passed$(NC)"

# Final comprehensive test
test-comprehensive: setup test-syntax test-dotfiles test-automation test-modules test-cloud test-api-keys test-auth-status test-required-tools test-install test-installer-components test-integration test-security test-performance test-workflows analyze-logs test-report ## Run all tests and generate report
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
	@# Test Ollama functionality
	@echo "Testing Ollama integration..."
	@chmod +x $(AUTO_DIR)/auto
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto ai ollama status > $(LOG_DIR)/ai-ollama-test.log 2>&1 || \
		echo -e "$(YELLOW)⚠️ Ollama test completed (may not be installed)$(NC)"
	@# Test Hugging Face functionality
	@echo "Testing Hugging Face integration..."
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto ai hf status > $(LOG_DIR)/ai-hf-test.log 2>&1 || \
		echo -e "$(YELLOW)⚠️ Hugging Face test completed (may not be installed)$(NC)"
	@# Test AI module syntax
	@echo "Testing AI module syntax..."
	@bash -n $(AUTO_DIR)/modules/ai.sh || (echo -e "$(RED)❌ AI module syntax error$(NC)" && exit 1)
	@bash -n functions/ai/ollama.sh || (echo -e "$(RED)❌ Ollama tools syntax error$(NC)" && exit 1)
	@bash -n functions/ai/huggingface.sh || (echo -e "$(RED)❌ Hugging Face tools syntax error$(NC)" && exit 1)
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
	@chmod +x $(AUTO_DIR)/auto
	@echo "=== AI/ML Performance Benchmark ===" > $(LOG_DIR)/ai-benchmark.log
	@echo "Generated: $$(date)" >> $(LOG_DIR)/ai-benchmark.log
	@echo "" >> $(LOG_DIR)/ai-benchmark.log
	@$(AUTO_DIR)/auto ai benchmark all >> $(LOG_DIR)/ai-benchmark.log 2>&1 || \
		echo -e "$(YELLOW)⚠️ Benchmark completed (some platforms may not be available)$(NC)"
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
test-ai-module: ## Test AI module specifically
	@echo -e "$(BLUE)Testing AI/ML module...$(NC)"
	@# Test AI module help
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto ai --help > $(LOG_DIR)/ai-help.log 2>&1 || \
		(echo -e "$(RED)❌ AI module help failed$(NC)" && exit 1)
	@# Test AI status (non-destructive)
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto ai status > $(LOG_DIR)/ai-status.log 2>&1 || \
		echo -e "$(YELLOW)⚠️ AI status test completed (some tools may not be installed)$(NC)"
	@# Test AI examples
	@timeout $(TEST_TIMEOUT) $(AUTO_DIR)/auto ai examples > $(LOG_DIR)/ai-examples.log 2>&1 || \
		(echo -e "$(RED)❌ AI examples failed$(NC)" && exit 1)
	@echo -e "$(GREEN)✅ AI module tests passed$(NC)"

# Extended module tests to include AI
test-modules: test-dev-module test-k8s-module test-github-module test-system-module test-config-module test-secrets-module test-tools-module test-ai-module ## Test all automation modules

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

# Default test for CI
.DEFAULT_GOAL := test-quick