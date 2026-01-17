# Dotfiles Management
# Install, inject, eject, link, and status

.PHONY: dotfiles-install dotfiles-inject dotfiles-eject dotfiles-link dotfiles-unlink
.PHONY: dotfiles-status dotfiles-reload dotfiles-update

# Installation
dotfiles-install: ## Run full dotfiles installation (install.sh)
	@bash install.sh

# Shell bootstrap injection
dotfiles-inject: ## Create shell bootstrap files (~/.bashrc, ~/.zshrc)
	@# Create ~/.bashrc
	@if [ ! -f ~/.bashrc ] || ! grep -q "DOTFILES_ROOT" ~/.bashrc 2>/dev/null; then \
		echo '# Dotfiles bootstrap' > ~/.bashrc; \
		echo 'export DOTFILES_ROOT="$(DOTFILES_DIR)"' >> ~/.bashrc; \
		echo '[ -f "$$DOTFILES_ROOT/core/bootstrap.sh" ] && . "$$DOTFILES_ROOT/core/bootstrap.sh"' >> ~/.bashrc; \
		echo "Created ~/.bashrc"; \
	else \
		echo "~/.bashrc already configured"; \
	fi
	@# Create ~/.zshrc
	@if [ ! -f ~/.zshrc ] || ! grep -q "DOTFILES_ROOT" ~/.zshrc 2>/dev/null; then \
		echo '# Dotfiles bootstrap' > ~/.zshrc; \
		echo 'export DOTFILES_ROOT="$(DOTFILES_DIR)"' >> ~/.zshrc; \
		echo '[ -f "$$DOTFILES_ROOT/core/bootstrap.sh" ] && . "$$DOTFILES_ROOT/core/bootstrap.sh"' >> ~/.zshrc; \
		echo "Created ~/.zshrc"; \
	else \
		echo "~/.zshrc already configured"; \
	fi
	@# Create ~/.bash_profile
	@if [ ! -f ~/.bash_profile ] || ! grep -q "bashrc" ~/.bash_profile 2>/dev/null; then \
		echo '# Source bashrc for login shells' > ~/.bash_profile; \
		echo '[ -f ~/.bashrc ] && . ~/.bashrc' >> ~/.bash_profile; \
		echo "Created ~/.bash_profile"; \
	else \
		echo "~/.bash_profile already configured"; \
	fi

# Shell bootstrap removal
dotfiles-eject: ## Remove shell bootstrap files
	@echo "This will remove managed bootstrap files."
	@read -p "Continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		grep -q "DOTFILES_ROOT" ~/.bashrc 2>/dev/null && rm ~/.bashrc && echo "Removed ~/.bashrc" || true; \
		grep -q "DOTFILES_ROOT" ~/.zshrc 2>/dev/null && rm ~/.zshrc && echo "Removed ~/.zshrc" || true; \
		grep -q "bashrc" ~/.bash_profile 2>/dev/null && rm ~/.bash_profile && echo "Removed ~/.bash_profile" || true; \
	else \
		echo "Cancelled"; \
	fi

# App config linking
dotfiles-link: ## Link app configurations (git, ghostty, vscode, claude)
	@bash -c 'export DOTFILES_ROOT="$(DOTFILES_DIR)"; source functions/core/inject.sh; inject_configs'

dotfiles-unlink: ## Remove app configuration links
	@[ -L ~/.gitconfig ] && rm ~/.gitconfig && echo "Removed ~/.gitconfig" || true
	@[ -L ~/.gitignore ] && rm ~/.gitignore && echo "Removed ~/.gitignore" || true
	@[ -L ~/.config/ghostty/config ] && rm ~/.config/ghostty/config && echo "Removed ghostty config" || true
	@[ -L ~/Library/Application\ Support/Code/User/settings.json ] && rm ~/Library/Application\ Support/Code/User/settings.json && echo "Removed VS Code settings" || true
	@[ -L ~/.config/claude/settings.json ] && rm ~/.config/claude/settings.json && echo "Removed Claude settings" || true

# Status
dotfiles-status: ## Show dotfiles installation status
	@echo "Dotfiles Status"
	@echo "==============="
	@echo ""
	@echo "Repository: $(DOTFILES_DIR)"
	@echo ""
	@echo "Bootstrap Files:"
	@grep -q "DOTFILES_ROOT" ~/.bashrc 2>/dev/null && echo "  ~/.bashrc: ok" || echo "  ~/.bashrc: not configured"
	@grep -q "DOTFILES_ROOT" ~/.zshrc 2>/dev/null && echo "  ~/.zshrc: ok" || echo "  ~/.zshrc: not configured"
	@[ -f ~/.bash_profile ] && echo "  ~/.bash_profile: ok" || echo "  ~/.bash_profile: missing"
	@echo ""
	@echo "App Configs:"
	@[ -L ~/.gitconfig ] && echo "  Git: linked" || echo "  Git: not linked"
	@[ -L ~/.config/ghostty/config ] && echo "  Ghostty: linked" || echo "  Ghostty: not linked"
	@[ -L ~/Library/Application\ Support/Code/User/settings.json ] 2>/dev/null && echo "  VS Code: linked" || echo "  VS Code: not linked"
	@[ -L ~/.config/claude/settings.json ] && echo "  Claude: linked" || echo "  Claude: not linked"

# Reload and update
dotfiles-reload: ## Reload shell configuration
	@echo "To reload, run:"
	@echo "  source ~/.bashrc   # for bash"
	@echo "  source ~/.zshrc    # for zsh"

dotfiles-update: ## Git pull and show reload instructions
	@git pull --rebase
	@echo ""
	@echo "Updated. Run 'source ~/.bashrc' or 'source ~/.zshrc' to reload."
