# Homebrew Package Management
# Install and manage brew packages by category

.PHONY: brew-status brew-update brew-install-devops brew-install-dev brew-install-db brew-install-all
.PHONY: db-install-mysql db-install-mongo db-install-redis db-install-neo4j db-install-kafka

# Package lists by category
BREW_DEVOPS := argocd awscli azure-cli cloudflared doctl helm k9s kind kubernetes-cli kustomize terraform vault
BREW_DEV := bash bash-completion ctags fd fzf gh jq lazygit neovim shellcheck tmux yq
BREW_DB := kafka mongocli mongosh mycli neo4j pgcli postgresql@14 redis iredis

# Status
brew-status: ## Show status of all managed brew packages
	@echo "Homebrew Package Status"
	@echo "======================="
	@echo ""
	@echo "DevOps Tools:"
	@for pkg in $(BREW_DEVOPS); do \
		if brew list --formula 2>/dev/null | grep -q "^$$pkg$$"; then \
			version=$$(brew info $$pkg 2>/dev/null | head -1 | awk '{print $$3}'); \
			printf "  %-15s %s\n" "$$pkg" "$$version"; \
		else \
			printf "  %-15s %s\n" "$$pkg" "(not installed)"; \
		fi; \
	done
	@echo ""
	@echo "Dev Tools:"
	@for pkg in $(BREW_DEV); do \
		if brew list --formula 2>/dev/null | grep -q "^$$pkg$$"; then \
			version=$$(brew info $$pkg 2>/dev/null | head -1 | awk '{print $$3}'); \
			printf "  %-15s %s\n" "$$pkg" "$$version"; \
		else \
			printf "  %-15s %s\n" "$$pkg" "(not installed)"; \
		fi; \
	done
	@echo ""
	@echo "Database Tools:"
	@for pkg in $(BREW_DB); do \
		if brew list --formula 2>/dev/null | grep -q "^$$pkg$$"; then \
			version=$$(brew info $$pkg 2>/dev/null | head -1 | awk '{print $$3}'); \
			printf "  %-15s %s\n" "$$pkg" "$$version"; \
		else \
			printf "  %-15s %s\n" "$$pkg" "(not installed)"; \
		fi; \
	done

# Update
brew-update: ## Update all brew packages
	@brew update
	@brew upgrade
	@brew cleanup

# Install by category
brew-install-devops: ## Install all DevOps tools via brew
	@brew install $(BREW_DEVOPS) || true

brew-install-dev: ## Install all dev tools via brew
	@brew install $(BREW_DEV) || true

brew-install-db: ## Install all database tools via brew
	@brew install $(BREW_DB) || true

brew-install-all: brew-install-devops brew-install-dev brew-install-db ## Install all managed brew packages

# Individual database tools
db-install-mysql: ## Install MySQL client + mycli
	@brew install mysql-client mycli

db-install-mongo: ## Install MongoDB shell + mongocli
	@brew install mongosh mongocli

db-install-redis: ## Install Redis + iredis
	@brew install redis iredis

db-install-neo4j: ## Install Neo4j
	@brew install neo4j

db-install-kafka: ## Install Kafka
	@brew install kafka
