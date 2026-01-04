---
description: Cross-compile Go application for multiple platforms
argument-hint: [app-name]
allowed-tools: Read, Bash
---

## Task

Help the user build their Go application for multiple platforms.

## Quick Build

Using dotfiles function:
```bash
gobuildall myapp
# Builds for: linux, darwin, windows (amd64, arm64)
# Output in dist/
```

## Manual Cross-Compilation

### Single Platform
```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o dist/myapp-linux-amd64 .

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o dist/myapp-darwin-arm64 .

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o dist/myapp-windows-amd64.exe .
```

### Common Targets

| GOOS | GOARCH | Description |
|------|--------|-------------|
| linux | amd64 | Linux x86-64 |
| linux | arm64 | Linux ARM64 |
| linux | arm | Linux ARM (32-bit) |
| darwin | amd64 | macOS Intel |
| darwin | arm64 | macOS Apple Silicon |
| windows | amd64 | Windows x86-64 |
| windows | arm64 | Windows ARM64 |

### List Supported Platforms
```bash
go tool dist list
```

## Build with Version Info

```bash
VERSION=$(git describe --tags --always)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT=$(git rev-parse --short HEAD)

go build -ldflags "\
  -X main.version=$VERSION \
  -X main.commit=$COMMIT \
  -X main.buildTime=$BUILD_TIME" \
  -o dist/myapp .
```

In your Go code:
```go
var (
    version   = "dev"
    commit    = "none"
    buildTime = "unknown"
)

func main() {
    fmt.Printf("Version: %s, Commit: %s, Built: %s\n",
        version, commit, buildTime)
}
```

## Build Optimizations

### Smaller Binary
```bash
# Strip debug info
go build -ldflags "-s -w" -o dist/myapp .

# Further compress with upx (if installed)
upx --best dist/myapp
```

### Reproducible Builds
```bash
go build -trimpath -ldflags "-s -w" -o dist/myapp .
```

### Static Binary (Linux)
```bash
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o dist/myapp .
```

## Makefile Integration

```makefile
APP_NAME := myapp
VERSION := $(shell git describe --tags --always)

.PHONY: build-all

build-all:
	GOOS=linux GOARCH=amd64 go build -o dist/$(APP_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o dist/$(APP_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o dist/$(APP_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o dist/$(APP_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o dist/$(APP_NAME)-windows-amd64.exe .
```

## CGO Considerations

```bash
# Disable CGO for pure Go builds
CGO_ENABLED=0 go build -o dist/myapp .

# CGO requires platform-specific toolchain for cross-compile
# Usually easier to build on target platform or use Docker
```

## Docker Multi-Platform Build

```dockerfile
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /myapp .

FROM scratch
COPY --from=builder /myapp /myapp
ENTRYPOINT ["/myapp"]
```

```bash
docker buildx build --platform linux/amd64,linux/arm64 -t myapp .
```

## Dotfiles Functions

- `gobuildall [name]` - Build for all major platforms
- `goclean` - Clean dist/ and artifacts
- `gob` - go build (alias)
