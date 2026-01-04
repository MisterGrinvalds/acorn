---
description: Run Go tests with coverage reporting
argument-hint: [package]
allowed-tools: Read, Bash
---

## Task

Help the user run Go tests with coverage analysis.

## Quick Coverage

Using dotfiles function:
```bash
gotestcover
# Generates coverage.html and shows summary
```

## Manual Coverage

### Basic Coverage
```bash
# Run with coverage
go test -cover ./...

# Using alias
gotc  # go test -cover
```

### Coverage Report
```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View summary
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Open report
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

### Package-Specific Coverage
```bash
# Single package
go test -cover ./internal/handlers

# Multiple packages
go test -cover ./internal/... ./pkg/...
```

## Coverage Modes

```bash
# set: Statement coverage (default)
go test -covermode=set -coverprofile=coverage.out ./...

# count: How many times each statement runs
go test -covermode=count -coverprofile=coverage.out ./...

# atomic: Like count, but safe for concurrent tests
go test -covermode=atomic -coverprofile=coverage.out ./...
```

## Table-Driven Test Pattern

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, 1, 0},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.a, tt.b, got, tt.expected)
            }
        })
    }
}
```

## Coverage Thresholds

### Check minimum coverage
```bash
# Parse coverage percentage
COVERAGE=$(go test -cover ./... | grep -o '[0-9]\+\.[0-9]\+%' | head -1)
echo "Coverage: $COVERAGE"
```

### CI Integration
```yaml
# GitHub Actions example
- name: Test with coverage
  run: |
    go test -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out | grep total | awk '{print $3}'

- name: Upload coverage
  uses: codecov/codecov-action@v3
  with:
    files: coverage.out
```

## Excluding Files

```go
//go:build !test
// +build !test

// This file is excluded from test builds
```

Or by naming:
- `*_test.go` - Only in tests
- Don't test `main.go` directly

## Coverage Tips

1. **Aim for meaningful coverage**, not 100%
2. **Test edge cases**, not just happy paths
3. **Focus on critical paths** first
4. **Use subtests** for organized coverage
5. **Review uncovered lines** in HTML report

## Dotfiles Integration

- `gotestcover` - Run tests with HTML coverage report
- `gotest [pattern]` - Run tests with optional filter
- `got` - go test (alias)
- `gotv` - go test -v (alias)
- `gotc` - go test -cover (alias)
