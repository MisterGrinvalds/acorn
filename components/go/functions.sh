#!/bin/sh
# components/go/functions.sh - Go development functions

# =============================================================================
# Project Initialization
# =============================================================================

# Initialize new Go project
gonew() {
    if [ -z "$1" ]; then
        echo "Usage: gonew <module-name>"
        return 1
    fi

    mkdir -p "$1"
    cd "$1" || return 1
    go mod init "$1"

    cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
EOF

    echo "Go project '$1' initialized!"
}

# =============================================================================
# Cobra CLI Helpers (for building CLI tools)
# =============================================================================

# Create new Cobra CLI project
cobranew() {
    if [ -z "$1" ]; then
        echo "Usage: cobranew <app-name>"
        return 1
    fi

    if ! command -v cobra-cli >/dev/null 2>&1; then
        echo "cobra-cli not installed. Installing..."
        go install github.com/spf13/cobra-cli@latest
    fi

    cobra-cli init "$1"
    cd "$1" || return 1
    go mod tidy
}

# Add command to Cobra project
cobradd() {
    if [ -z "$1" ]; then
        echo "Usage: cobradd <command-name>"
        return 1
    fi

    if ! command -v cobra-cli >/dev/null 2>&1; then
        echo "cobra-cli not installed"
        return 1
    fi

    cobra-cli add "$1"
}

# =============================================================================
# Testing Helpers
# =============================================================================

# Run tests with optional pattern filter
gotest() {
    if [ -z "$1" ]; then
        go test ./...
    else
        go test ./... -run "$1"
    fi
}

# Run tests with coverage report
gotestcover() {
    go test ./... -coverprofile=coverage.out
    go tool cover -html=coverage.out -o coverage.html
    echo "Coverage report generated: coverage.html"

    # Show summary
    go tool cover -func=coverage.out | tail -1
}

# Run benchmarks
gobench() {
    if [ -z "$1" ]; then
        go test -bench=. ./...
    else
        go test -bench="$1" ./...
    fi
}

# =============================================================================
# Build Helpers
# =============================================================================

# Build for multiple platforms
gobuildall() {
    local name="${1:-app}"

    echo "Building for multiple platforms..."

    GOOS=linux GOARCH=amd64 go build -o "dist/${name}-linux-amd64" .
    GOOS=linux GOARCH=arm64 go build -o "dist/${name}-linux-arm64" .
    GOOS=darwin GOARCH=amd64 go build -o "dist/${name}-darwin-amd64" .
    GOOS=darwin GOARCH=arm64 go build -o "dist/${name}-darwin-arm64" .
    GOOS=windows GOARCH=amd64 go build -o "dist/${name}-windows-amd64.exe" .

    echo "Builds complete in dist/"
    ls -la dist/
}

# Clean build artifacts
goclean() {
    go clean
    rm -rf dist/ coverage.out coverage.html
    echo "Build artifacts cleaned"
}
