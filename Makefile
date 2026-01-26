# Dotfiles Build System
# =====================
# Modular Makefile for dotfiles, development tools, and testing
#
# Usage:
#   make help             Show all available targets
#   make status           Show environment status
#   make test             Run quick tests
#   make acorn-build      Build the acorn CLI
#   make docker-status    Show Docker status
#   make kubernetes-status Show Kubernetes status
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
#
# Component makefiles (per-tool):
#   makefiles/docker.mk, kubernetes.mk, helm.mk, k9s.mk, argocd.mk,
#   git.mk, github.mk, fzf.mk, tmux.mk, jq.mk, yq.mk, neovim.mk,
#   terraform.mk, aws.mk, cloudflare.mk

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

# Component-specific makefiles
-include makefiles/docker.mk
-include makefiles/kubernetes.mk
-include makefiles/helm.mk
-include makefiles/k9s.mk
-include makefiles/argocd.mk
-include makefiles/git.mk
-include makefiles/github.mk
-include makefiles/fzf.mk
-include makefiles/tmux.mk
-include makefiles/jq.mk
-include makefiles/yq.mk
-include makefiles/neovim.mk
-include makefiles/terraform.mk
-include makefiles/aws.mk
-include makefiles/cloudflare.mk

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
