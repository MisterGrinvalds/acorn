# Acorn CLI Build System
# Build, test, lint, and release

.PHONY: acorn acorn-build acorn-install acorn-clean acorn-test acorn-coverage
.PHONY: acorn-bench acorn-lint acorn-fmt acorn-vet acorn-security acorn-check
.PHONY: acorn-deps acorn-tidy acorn-verify acorn-cross acorn-snapshot acorn-release acorn-version

# Build configuration
BINARY_NAME := acorn
GO_MODULE := github.com/mistergrinvalds/acorn
VERSION_PKG := $(GO_MODULE)/internal/version
BUILD_DIR := bin
CMD_DIR := cmd/acorn

# Version info from git
GIT_VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILT_BY := $(shell whoami)

# Linker flags for version embedding
LDFLAGS := -s -w \
	-X '$(VERSION_PKG).Version=$(GIT_VERSION)' \
	-X '$(VERSION_PKG).Commit=$(GIT_COMMIT)' \
	-X '$(VERSION_PKG).Date=$(BUILD_DATE)' \
	-X '$(VERSION_PKG).BuiltBy=$(BUILT_BY)'

# Cross-compilation platforms
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

# Build targets
acorn: acorn-build ## Build acorn CLI (alias)

acorn-build: ## Build acorn CLI binary
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME) ($(GIT_VERSION))"

acorn-install: acorn-build ## Install acorn to GOPATH/bin
	@go install -ldflags "$(LDFLAGS)" ./$(CMD_DIR)
	@echo "Installed to $$(go env GOPATH)/bin"

acorn-clean: ## Clean build artifacts
	@rm -rf $(BUILD_DIR)
	@go clean -cache -testcache

# Testing
acorn-test: ## Run all tests with race detection
	@go test -race -v ./...

acorn-coverage: ## Run tests with coverage report
	@mkdir -p $(BUILD_DIR)
	@go test -race -coverprofile=$(BUILD_DIR)/coverage.out -covermode=atomic ./...
	@go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@go tool cover -func=$(BUILD_DIR)/coverage.out | tail -1
	@echo "Coverage report: $(BUILD_DIR)/coverage.html"

acorn-bench: ## Run benchmarks
	@go test -bench=. -benchmem ./...

# Code quality
acorn-lint: ## Run golangci-lint
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout 5m; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

acorn-fmt: ## Format Go code
	@gofmt -s -w .
	@goimports -w . 2>/dev/null || true

acorn-vet: ## Run go vet
	@go vet ./...

acorn-security: ## Run security scanners (gosec, govulncheck)
	@if command -v gosec >/dev/null 2>&1; then \
		gosec -quiet ./...; \
	else \
		echo "gosec not installed"; \
	fi
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "govulncheck not installed"; \
	fi

acorn-check: acorn-fmt acorn-vet acorn-lint acorn-test ## Run all checks (fmt, vet, lint, test)

# Dependencies
acorn-deps: ## Install development dependencies
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/goreleaser/goreleaser/v2@latest

acorn-tidy: ## Tidy go modules
	@go mod tidy
	@go mod verify

acorn-verify: ## Verify module checksums
	@go mod verify

# Cross-compilation
acorn-cross: ## Build for all platforms
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		go build -ldflags "$(LDFLAGS)" \
			-o $(BUILD_DIR)/$(BINARY_NAME)-$${platform%/*}-$${platform#*/} \
			./$(CMD_DIR); \
		echo "Built: $(BINARY_NAME)-$${platform%/*}-$${platform#*/}"; \
	done

# Release
acorn-snapshot: ## Build snapshot release with goreleaser
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --clean; \
	else \
		echo "goreleaser not installed. Run: make acorn-deps"; exit 1; \
	fi

acorn-release: ## Create tagged release with goreleaser
	@if [ -z "$$(git tag -l --points-at HEAD)" ]; then \
		echo "No tag on current commit. Tag with: git tag vX.Y.Z"; exit 1; \
	fi
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --clean; \
	else \
		echo "goreleaser not installed. Run: make acorn-deps"; exit 1; \
	fi

acorn-version: ## Show version that would be embedded
	@echo "Version:  $(GIT_VERSION)"
	@echo "Commit:   $(GIT_COMMIT)"
	@echo "Date:     $(BUILD_DATE)"
	@echo "Built By: $(BUILT_BY)"
