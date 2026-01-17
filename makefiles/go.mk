# Go Development
# Installation, environment setup, and common tools

.PHONY: go-install go-setup go-status go-tools

# Installation
go-install: ## Install Go via Homebrew (macOS) or download (Linux)
	@if command -v go >/dev/null 2>&1; then \
		echo "Go already installed: $$(go version)"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		brew install go; \
	else \
		curl -LO https://go.dev/dl/go1.22.0.linux-amd64.tar.gz; \
		sudo rm -rf /usr/local/go; \
		sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz; \
		rm go1.22.0.linux-amd64.tar.gz; \
		echo "Add to PATH: export PATH=\$$PATH:/usr/local/go/bin"; \
	fi

go-setup: go-install ## Setup Go environment
	@mkdir -p ~/go/{bin,src,pkg}
	@echo "Go workspace created at ~/go"
	@command -v go >/dev/null && go env GOROOT GOPATH GOBIN || true

# Status
go-status: ## Check Go installation status
	@echo "Go Status"
	@echo "========="
	@echo ""
	@if command -v go >/dev/null 2>&1; then \
		echo "Version: $$(go version | awk '{print $$3}')"; \
		echo "Path: $$(which go)"; \
		echo "GOROOT: $$(go env GOROOT)"; \
		echo "GOPATH: $$(go env GOPATH)"; \
		echo ""; \
		echo "Installed tools:"; \
		ls $$(go env GOPATH)/bin 2>/dev/null | head -10 || echo "  (none)"; \
	else \
		echo "Go not installed. Run: make go-install"; \
	fi

# Development tools
go-tools: ## Install common Go development tools
	@if ! command -v go >/dev/null 2>&1; then \
		echo "Go not installed. Run: make go-install"; exit 1; \
	fi
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing gopls (language server)..."
	@go install golang.org/x/tools/gopls@latest
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Installing dlv (debugger)..."
	@go install github.com/go-delve/delve/cmd/dlv@latest
	@echo "Go tools installed"
