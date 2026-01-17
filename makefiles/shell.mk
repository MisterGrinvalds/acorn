# Shell Layer Testing
# Module-by-module testing for shell configuration

.PHONY: shell-status shell-test-all
.PHONY: shell-test-discovery shell-test-xdg shell-test-theme shell-test-env
.PHONY: shell-test-options shell-test-aliases shell-test-functions shell-test-prompt

# Status overview
shell-status: ## Show status of all shell modules
	@echo "Shell Module Status"
	@echo "==================="
	@echo ""
	@echo "Core Modules:"
	@for mod in discovery xdg theme environment options aliases functions prompt; do \
		if [ -f "shell/$$mod.sh" ]; then \
			if bash -n "shell/$$mod.sh" 2>/dev/null; then \
				echo "  shell/$$mod.sh: ok"; \
			else \
				echo "  shell/$$mod.sh: syntax error"; \
			fi; \
		else \
			echo "  shell/$$mod.sh: missing"; \
		fi; \
	done
	@echo ""
	@echo "Function Modules:"
	@for mod in functions/*.sh; do \
		if [ -f "$$mod" ]; then \
			if bash -n "$$mod" 2>/dev/null; then \
				echo "  $$mod: ok"; \
			else \
				echo "  $$mod: syntax error"; \
			fi; \
		fi; \
	done | head -15

# Individual module tests
shell-test-discovery: ## Test shell/discovery.sh module
	@bash -c '\
		export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source shell/discovery.sh; \
		[ -n "$$CURRENT_SHELL" ] && echo "CURRENT_SHELL: $$CURRENT_SHELL" || exit 1; \
		[ -n "$$CURRENT_PLATFORM" ] && echo "CURRENT_PLATFORM: $$CURRENT_PLATFORM" || exit 1; \
	'

shell-test-xdg: ## Test shell/xdg.sh module
	@bash -c '\
		export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source shell/init.sh; \
		[ -n "$$XDG_CONFIG_HOME" ] && echo "XDG_CONFIG_HOME: $$XDG_CONFIG_HOME" || exit 1; \
		[ -n "$$XDG_DATA_HOME" ] && echo "XDG_DATA_HOME: $$XDG_DATA_HOME" || exit 1; \
	'

shell-test-theme: ## Test shell/theme.sh module
	@bash -c '\
		export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source shell/init.sh; \
		[ -n "$$CURRENT_THEME" ] && echo "CURRENT_THEME: $$CURRENT_THEME" || echo "CURRENT_THEME: not set"; \
	'

shell-test-env: ## Test shell/environment.sh module
	@bash -c '\
		export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source shell/init.sh; \
		[ -n "$$EDITOR" ] && echo "EDITOR: $$EDITOR" || echo "EDITOR: not set"; \
		[ -n "$$VISUAL" ] && echo "VISUAL: $$VISUAL" || echo "VISUAL: not set"; \
	'

shell-test-options: ## Test shell/options.sh module
	@bash -c '\
		export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source shell/init.sh; \
		shopt -q histappend && echo "histappend: enabled" || echo "histappend: disabled"; \
	'

shell-test-aliases: ## Test shell/aliases.sh module
	@bash -c '\
		export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source shell/init.sh; \
		alias | grep -q "ll=" && echo "ll alias: defined" || exit 1; \
		alias | grep -q "la=" && echo "la alias: defined" || exit 1; \
	'

shell-test-functions: ## Test all function modules
	@echo "Testing function modules..."
	@for func_dir in functions/*/; do \
		for file in "$$func_dir"*.sh; do \
			[ -f "$$file" ] || continue; \
			if bash -n "$$file" 2>/dev/null; then \
				echo "  $$file: ok"; \
			else \
				echo "  $$file: syntax error"; \
			fi; \
		done; \
	done | head -20
	@bash -c '\
		export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source shell/init.sh; \
		declare -f h >/dev/null && echo "h(): defined" || echo "h(): not found"; \
		declare -f mkvenv >/dev/null && echo "mkvenv(): defined" || echo "mkvenv(): not found"; \
	'

shell-test-prompt: ## Test shell/prompt.sh module
	@bash -c '\
		export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source shell/init.sh; \
		declare -f git_branch >/dev/null && echo "git_branch(): defined" || exit 1; \
		declare -f git_color >/dev/null && echo "git_color(): defined" || exit 1; \
		[ -n "$$PS1" ] && echo "PS1: set" || echo "PS1: not set"; \
	'

# Complete shell test
shell-test-all: shell-test-discovery shell-test-xdg shell-test-theme shell-test-env shell-test-options shell-test-aliases shell-test-functions shell-test-prompt ## Test complete shell loading sequence
