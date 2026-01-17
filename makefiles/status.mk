# Comprehensive Status
# Full environment overview

.PHONY: status

status: ## Show complete environment status
	@echo "Environment Status"
	@echo "=================="
	@echo ""
	@# Shell
	@echo "Shell:"
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh 2>/dev/null; \
		echo "  Shell: $${CURRENT_SHELL:-unknown}"; \
		echo "  Platform: $${CURRENT_PLATFORM:-unknown}"; \
		echo "  DOTFILES_ROOT: $${DOTFILES_ROOT:-not set}"'
	@echo ""
	@# XDG
	@echo "XDG Directories:"
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; source shell/init.sh 2>/dev/null; \
		echo "  XDG_CONFIG_HOME: $${XDG_CONFIG_HOME:-not set}"; \
		echo "  XDG_DATA_HOME: $${XDG_DATA_HOME:-not set}"; \
		echo "  XDG_CACHE_HOME: $${XDG_CACHE_HOME:-not set}"'
	@echo ""
	@# Dev tools
	@echo "Development Tools:"
	@printf "  %-12s" "git:"; command -v git >/dev/null && echo "$$(git --version | awk '{print $$3}')" || echo "missing"
	@printf "  %-12s" "node:"; command -v node >/dev/null && echo "$$(node --version)" || echo "missing"
	@printf "  %-12s" "python:"; command -v python3 >/dev/null && echo "$$(python3 --version | awk '{print $$2}')" || echo "missing"
	@printf "  %-12s" "go:"; command -v go >/dev/null && echo "$$(go version | awk '{print $$3}')" || echo "missing"
	@printf "  %-12s" "uv:"; command -v uv >/dev/null && echo "$$(uv --version 2>/dev/null)" || echo "missing"
	@printf "  %-12s" "pnpm:"; command -v pnpm >/dev/null && echo "$$(pnpm --version)" || echo "missing"
	@printf "  %-12s" "fzf:"; command -v fzf >/dev/null && echo "$$(fzf --version | cut -d' ' -f1)" || echo "missing"
	@printf "  %-12s" "jq:"; command -v jq >/dev/null && echo "$$(jq --version)" || echo "missing"
	@echo ""
	@# DevOps
	@echo "DevOps Tools:"
	@printf "  %-12s" "kubectl:"; command -v kubectl >/dev/null && echo "installed" || echo "missing"
	@printf "  %-12s" "helm:"; command -v helm >/dev/null && echo "$$(helm version --short 2>/dev/null)" || echo "missing"
	@printf "  %-12s" "terraform:"; command -v terraform >/dev/null && echo "installed" || echo "missing"
	@printf "  %-12s" "aws:"; command -v aws >/dev/null && echo "installed" || echo "missing"
	@printf "  %-12s" "docker:"; command -v docker >/dev/null && echo "installed" || echo "missing"
	@echo ""
	@# Database
	@echo "Database Tools:"
	@printf "  %-12s" "pgcli:"; command -v pgcli >/dev/null && echo "installed" || echo "missing"
	@printf "  %-12s" "mycli:"; command -v mycli >/dev/null && echo "installed" || echo "missing"
	@printf "  %-12s" "mongosh:"; command -v mongosh >/dev/null && echo "installed" || echo "missing"
	@printf "  %-12s" "redis-cli:"; command -v redis-cli >/dev/null && echo "installed" || echo "missing"
	@echo ""
	@# AI
	@echo "AI/ML Tools:"
	@printf "  %-12s" "ollama:"; command -v ollama >/dev/null && echo "installed" || echo "missing"
	@printf "  %-12s" "transformers:"; python3 -c "import transformers" 2>/dev/null && echo "installed" || echo "missing"
	@echo ""
	@# Installation
	@echo "Installation Status:"
	@printf "  %-12s" "~/.bashrc:"; grep -q "DOTFILES_ROOT" ~/.bashrc 2>/dev/null && echo "configured" || echo "not configured"
	@printf "  %-12s" "~/.zshrc:"; grep -q "DOTFILES_ROOT" ~/.zshrc 2>/dev/null && echo "configured" || echo "not configured"
	@printf "  %-12s" "git:"; [ -L ~/.gitconfig ] && echo "linked" || echo "not linked"
	@printf "  %-12s" "ghostty:"; [ -L ~/.config/ghostty/config ] && echo "linked" || echo "not linked"
	@echo ""
	@echo "Run 'make help' for available commands"
