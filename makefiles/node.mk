# Node.js and NVM Management
# NVM installation, Node versions, pnpm

.PHONY: nvm-install nvm-setup nvm-status node-install node-lts node-status node-update
.PHONY: pnpm-install pnpm-setup

# NVM installation
nvm-install: ## Install NVM (Node Version Manager)
	@if [ -d "$$HOME/.nvm" ] || [ -d "$$HOME/nvm" ]; then \
		echo "NVM already installed"; \
	else \
		curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash; \
		echo "Restart your shell or run: source ~/.bashrc"; \
	fi

nvm-setup: nvm-install node-lts pnpm-install ## Complete NVM setup with LTS Node and pnpm

nvm-status: ## Check NVM installation status
	@echo "NVM Status"
	@echo "=========="
	@if [ -d "$$HOME/.nvm" ]; then \
		echo "NVM_DIR: $$HOME/.nvm"; \
		export NVM_DIR="$$HOME/.nvm"; \
		[ -s "$$NVM_DIR/nvm.sh" ] && . "$$NVM_DIR/nvm.sh" && nvm --version && echo "Installed versions:" && nvm ls; \
	elif [ -d "$$HOME/nvm" ]; then \
		echo "NVM_DIR: $$HOME/nvm"; \
	else \
		echo "NVM not installed. Run: make nvm-install"; \
	fi

# Node.js installation
node-install: ## Install latest Node.js via NVM
	@export NVM_DIR="$${NVM_DIR:-$$HOME/.nvm}"; \
	[ -s "$$NVM_DIR/nvm.sh" ] && . "$$NVM_DIR/nvm.sh" && nvm install node && nvm use node

node-lts: ## Install latest LTS Node.js via NVM
	@export NVM_DIR="$${NVM_DIR:-$$HOME/.nvm}"; \
	[ -s "$$NVM_DIR/nvm.sh" ] && . "$$NVM_DIR/nvm.sh" && nvm install --lts && nvm use --lts && nvm alias default lts/*

node-status: ## Check complete Node.js ecosystem status
	@echo "Node.js Ecosystem Status"
	@echo "========================"
	@echo ""
	@echo "NVM:"
	@if [ -d "$$HOME/.nvm" ] || [ -d "$$HOME/nvm" ]; then \
		export NVM_DIR="$${NVM_DIR:-$$HOME/.nvm}"; \
		[ -s "$$NVM_DIR/nvm.sh" ] && . "$$NVM_DIR/nvm.sh" && echo "  Version: $$(nvm --version)" && echo "  Active: $$(nvm current)"; \
	else \
		echo "  Not installed"; \
	fi
	@echo ""
	@echo "Node.js:"
	@command -v node >/dev/null && echo "  Version: $$(node --version)" && echo "  Path: $$(which node)" || echo "  Not installed"
	@echo ""
	@echo "npm:"
	@command -v npm >/dev/null && echo "  Version: $$(npm --version)" || echo "  Not installed"
	@echo ""
	@echo "pnpm:"
	@command -v pnpm >/dev/null && echo "  Version: $$(pnpm --version)" || echo "  Not installed"
	@echo ""
	@echo "corepack:"
	@command -v corepack >/dev/null && echo "  Available" || echo "  Not available"
	@echo ""
	@echo "Global packages:"
	@command -v npm >/dev/null && npm list -g --depth=0 2>/dev/null | tail -10 || true

node-update: ## Update Node.js to latest LTS and reinstall globals
	@export NVM_DIR="$${NVM_DIR:-$$HOME/.nvm}"; \
	[ -s "$$NVM_DIR/nvm.sh" ] && . "$$NVM_DIR/nvm.sh" && \
	CURRENT=$$(nvm current) && \
	nvm install --lts --reinstall-packages-from=$$CURRENT && \
	nvm use --lts && \
	nvm alias default lts/*

# pnpm installation
pnpm-install: ## Install pnpm globally via npm
	@if command -v pnpm >/dev/null 2>&1; then \
		echo "pnpm already installed: $$(pnpm --version)"; \
	elif command -v npm >/dev/null 2>&1; then \
		npm install -g pnpm; \
		echo "pnpm installed: $$(pnpm --version)"; \
	else \
		echo "npm not found. Install Node.js first: make node-lts"; exit 1; \
	fi

pnpm-setup: ## Setup pnpm with corepack (Node 16.13+)
	@if command -v corepack >/dev/null 2>&1; then \
		corepack enable && corepack prepare pnpm@latest --activate; \
		echo "pnpm enabled via corepack: $$(pnpm --version)"; \
	else \
		echo "corepack not available. Using npm install instead."; \
		$(MAKE) pnpm-install; \
	fi
