# Docker Management
# Container runtime and CLI tooling

.PHONY: docker-install docker-setup docker-status docker-test

# Installation
docker-install: ## Install Docker
	@if command -v docker >/dev/null 2>&1; then \
		echo "Docker already installed: $$(docker --version)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install --cask docker; \
		echo "Docker Desktop installed. Please start Docker Desktop from Applications."; \
	else \
		curl -fsSL https://get.docker.com -o get-docker.sh; \
		sudo sh get-docker.sh; \
		rm get-docker.sh; \
		sudo usermod -aG docker $$USER; \
		echo "Docker installed. Log out and back in for group changes to take effect."; \
	fi

docker-setup: docker-install ## Setup Docker environment
	@if command -v docker >/dev/null 2>&1; then \
		mkdir -p "$${XDG_CONFIG_HOME:-$$HOME/.config}/docker"; \
		echo "Docker config directory created"; \
		docker info >/dev/null 2>&1 && echo "Docker daemon is running" || echo "Docker daemon is not running. Start Docker Desktop or enable the daemon."; \
	fi

# Status
docker-status: ## Check Docker installation status
	@echo "Docker Status"
	@echo "============="
	@echo ""
	@if command -v docker >/dev/null 2>&1; then \
		echo "Version: $$(docker --version | awk '{print $$3}' | tr -d ',')"; \
		echo "Path: $$(which docker)"; \
		echo ""; \
		if docker info >/dev/null 2>&1; then \
			echo "Status: Running"; \
			echo "Server Version: $$(docker version --format '{{.Server.Version}}' 2>/dev/null || echo 'N/A')"; \
			echo "Containers: $$(docker ps -q 2>/dev/null | wc -l | tr -d ' ') running, $$(docker ps -aq 2>/dev/null | wc -l | tr -d ' ') total"; \
			echo "Images: $$(docker images -q 2>/dev/null | wc -l | tr -d ' ')"; \
		else \
			echo "Status: Not running"; \
		fi; \
	else \
		echo "Docker not installed. Run: make docker-install"; \
	fi

# Test
docker-test: ## Test Docker functionality
	@echo "Testing Docker..."
	@command -v docker >/dev/null || (echo "Docker not installed"; exit 1)
	@docker info >/dev/null 2>&1 || (echo "Docker daemon not running"; exit 1)
	@docker run --rm hello-world >/dev/null 2>&1 && echo "Docker test passed" || (echo "Docker test failed"; exit 1)
