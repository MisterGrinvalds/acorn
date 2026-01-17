# Dotfiles Build System
# =====================
# Modular Makefile for dotfiles, development tools, and testing
#
# Usage:
#   make help          Show all available targets
#   make status        Show environment status
#   make test          Run quick tests
#   make acorn-build   Build the acorn CLI
#
# Structure:
#   makefiles/core.mk      - Variables and setup
#   makefiles/test.mk      - Testing targets
#   makefiles/ai.mk        - AI/ML (Ollama, HuggingFace)
#   makefiles/node.mk      - Node.js/NVM/pnpm
#   makefiles/shell.mk     - Shell module testing
#   makefiles/component.mk - Component management
#   makefiles/dotfiles.mk  - Dotfiles install/inject
#   makefiles/python.mk    - Python/UV
#   makefiles/go.mk        - Go development
#   makefiles/acorn.mk     - Acorn CLI build
#   makefiles/brew.mk      - Homebrew packages
#   makefiles/claude.mk    - Claude Code integration
#   makefiles/status.mk    - Status overview

# Include all component makefiles
include makefiles/core.mk
include makefiles/test.mk
include makefiles/ai.mk
include makefiles/node.mk
include makefiles/shell.mk
include makefiles/component.mk
include makefiles/dotfiles.mk
include makefiles/python.mk
include makefiles/go.mk
include makefiles/acorn.mk
include makefiles/brew.mk
include makefiles/claude.mk
include makefiles/status.mk

# Default target
.DEFAULT_GOAL := help

# Help target - scans all included makefiles for ## comments
help: ## Show this help message
	@echo "Dotfiles Build System"
	@echo "====================="
	@echo ""
	@echo "Targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' | \
		sort
	@echo ""
	@echo "Categories:"
	@echo "  test-*            Testing and validation"
	@echo "  shell-*           Shell module testing"
	@echo "  component-*       Component management"
	@echo "  dotfiles-*        Dotfiles installation"
	@echo "  ai-*              AI/ML tools (Ollama, HuggingFace)"
	@echo "  nvm-*, node-*     Node.js ecosystem"
	@echo "  uv-*, venv-*      Python/UV management"
	@echo "  go-*              Go development"
	@echo "  acorn-*           Acorn CLI build system"
	@echo "  brew-*            Homebrew packages"
	@echo "  db-*              Database tools"
	@echo "  claude-*          Claude Code integration"
