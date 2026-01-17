# Core Makefile variables and configuration
# Shared across all makefiles

# Shell and directories
SHELL := /bin/bash
DOTFILES_DIR := $(PWD)
TEST_DIR := $(DOTFILES_DIR)/tests
LOG_DIR := $(TEST_DIR)/logs
BACKUP_DIR := $(TEST_DIR)/backups
TEST_PROJECT_DIR := $(TEST_DIR)/test-projects
COMPONENTS_DIR := $(DOTFILES_DIR)/components

# Test configuration
TEST_TIMEOUT := 30

# Terminal colors
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m

# Helper function to check command exists
define check_cmd
	@command -v $(1) >/dev/null 2>&1
endef

# Setup test environment
setup: ## Create test directories and backup configs
	@mkdir -p $(LOG_DIR) $(BACKUP_DIR) $(TEST_PROJECT_DIR)
	@if [ -f ~/.bash_profile ]; then cp ~/.bash_profile $(BACKUP_DIR)/bash_profile.backup 2>/dev/null || true; fi
	@if [ -d ~/.config/shell ]; then cp -r ~/.config/shell $(BACKUP_DIR)/ 2>/dev/null || true; fi

# Clean test artifacts
clean: ## Remove test artifacts and restore backups
	@rm -rf $(TEST_DIR)/test-projects $(LOG_DIR)
	@if [ -f $(BACKUP_DIR)/bash_profile.backup ]; then cp $(BACKUP_DIR)/bash_profile.backup ~/.bash_profile; fi
	@if [ -d $(BACKUP_DIR)/shell ]; then cp -r $(BACKUP_DIR)/shell ~/.config/; fi
