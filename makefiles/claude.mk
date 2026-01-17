# Claude Code Integration
# Session management and configuration

.PHONY: claude claude-context claude-resume claude-status

CLAUDE_CONTEXT_FILE := $(DOTFILES_DIR)/.claude/SESSION_CONTEXT.md

# Launch Claude
claude: ## Launch Claude Code with session context
	@if [ ! -f "$(CLAUDE_CONTEXT_FILE)" ]; then \
		echo "No context file found. Starting fresh session..."; \
		cd $(DOTFILES_DIR) && claude; \
	else \
		echo "Loading context from: $(CLAUDE_CONTEXT_FILE)"; \
		cd $(DOTFILES_DIR) && claude "Read the file .claude/SESSION_CONTEXT.md to understand this repository's current state and architecture. Then ask me what I'd like to work on."; \
	fi

# Context management
claude-context: ## Show current session context
	@if [ -f "$(CLAUDE_CONTEXT_FILE)" ]; then \
		cat "$(CLAUDE_CONTEXT_FILE)"; \
	else \
		echo "No context file found at $(CLAUDE_CONTEXT_FILE)"; \
	fi

claude-resume: ## Resume most recent Claude session
	@cd $(DOTFILES_DIR) && claude --continue

# Status
claude-status: ## Show Claude Code configuration status
	@echo "Claude Code Status"
	@echo "=================="
	@echo ""
	@echo "Settings:"
	@if [ -f "$(DOTFILES_DIR)/.claude/settings.local.json" ]; then \
		echo "  Project: $(DOTFILES_DIR)/.claude/settings.local.json"; \
		echo "  Model: $$(jq -r '.model // "default"' $(DOTFILES_DIR)/.claude/settings.local.json 2>/dev/null)"; \
	else \
		echo "  No project settings found"; \
	fi
	@echo ""
	@echo "Context File:"
	@if [ -f "$(CLAUDE_CONTEXT_FILE)" ]; then \
		echo "  $(CLAUDE_CONTEXT_FILE)"; \
		echo "  Size: $$(wc -l < $(CLAUDE_CONTEXT_FILE)) lines"; \
	else \
		echo "  Not found"; \
	fi
	@echo ""
	@echo "Custom Assets:"
	@echo "  Commands: $$(find $(DOTFILES_DIR)/.claude/commands -name "*.md" 2>/dev/null | wc -l | tr -d ' ')"
	@echo "  Agents: $$(find $(DOTFILES_DIR)/.claude/agents -name "*.md" 2>/dev/null | wc -l | tr -d ' ')"
