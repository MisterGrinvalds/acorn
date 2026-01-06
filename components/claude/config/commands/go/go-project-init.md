---
description: Initialize a new Go project with proper structure
argument-hint: <project-name> [type: basic|cli|api]
allowed-tools: Read, Write, Bash
---

## Task

Help the user initialize a new Go project with proper structure.

## Project Types

Based on second argument in `$ARGUMENTS`:

### basic (default)
Simple Go project:

```bash
# Using dotfiles function
gonew myproject

# Creates:
# myproject/
# ├── go.mod
# └── main.go
```

### cli
Cobra CLI application:

```bash
# Using dotfiles function
cobranew mycli
cd mycli

# Add commands
cobradd serve
cobradd version
```

Structure:
```
mycli/
├── go.mod
├── main.go
└── cmd/
    ├── root.go
    ├── serve.go
    └── version.go
```

### api
HTTP API server:

```bash
gonew myapi
# Then add server code
```

Structure:
```
myapi/
├── go.mod
├── main.go
├── cmd/
│   └── server.go
├── internal/
│   ├── handlers/
│   ├── middleware/
│   └── models/
└── pkg/
    └── utils/
```

## Recommended Structure

```
project/
├── go.mod
├── go.sum
├── main.go              # Entry point
├── Makefile             # Build automation
├── README.md
├── cmd/                 # Command implementations
│   └── root.go
├── internal/            # Private packages
│   ├── config/
│   ├── handlers/
│   └── services/
├── pkg/                 # Public packages
│   └── utils/
├── api/                 # API definitions (OpenAPI, etc.)
├── scripts/             # Build/deploy scripts
├── tests/               # Integration tests
└── dist/                # Build outputs
```

## Essential Files

### go.mod
```go
module github.com/username/myproject

go 1.22

require (
    github.com/spf13/cobra v1.8.0
)
```

### Makefile
```makefile
.PHONY: build test clean

APP_NAME := myproject
VERSION := $(shell git describe --tags --always)

build:
	go build -ldflags "-X main.version=$(VERSION)" -o dist/$(APP_NAME) .

test:
	go test -v ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf dist/ coverage.out

lint:
	golangci-lint run

fmt:
	go fmt ./...
	go vet ./...
```

### .gitignore
```gitignore
# Binaries
*.exe
dist/

# Test
coverage.out
coverage.html

# IDE
.idea/
.vscode/
*.swp

# OS
.DS_Store
```

## Verification

```bash
# Build
gob  # go build

# Run
gor main.go  # go run

# Test
got  # go test

# Format
gof  # go fmt ./...
gov  # go vet ./...
```

## Dotfiles Functions

- `gonew <name>` - Create basic Go project
- `cobranew <name>` - Create Cobra CLI project
- `cobradd <cmd>` - Add command to Cobra project
