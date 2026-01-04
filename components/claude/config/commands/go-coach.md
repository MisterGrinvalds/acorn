---
description: Interactive coaching session to learn Go development
argument-hint: [skill-level: beginner|intermediate|advanced]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Guide the user through learning Go development workflows interactively.

## Approach

1. **Assess level** - Ask about Go experience
2. **Set goals** - Identify what they want to build
3. **Progressive exercises** - Start simple, build complexity
4. **Real-time practice** - Have them run commands
5. **Reinforce** - Summarize and suggest next steps

## Skill Levels

### Beginner
- Go installation and setup
- Hello World and basic syntax
- Variables, functions, types
- Go modules basics
- Running and building
- Basic testing

### Intermediate
- Package organization
- Error handling patterns
- Interfaces and composition
- Testing with table-driven tests
- Using dotfiles functions
- Building for multiple platforms

### Advanced
- Concurrency (goroutines, channels)
- Context and cancellation
- Generics
- Performance optimization
- Cobra CLI development
- Profiling and benchmarking

## Interactive Exercises

### Beginner Exercises
```bash
# Exercise 1: Check Go installation
gover  # go version (alias)

# Exercise 2: Create first project
gonew hello-go
cd hello-go

# Exercise 3: Run it
gor main.go  # go run main.go (alias)

# Exercise 4: Build it
gob  # go build (alias)
./hello-go
```

### Intermediate Exercises
```bash
# Exercise 5: Add a test
# Create main_test.go
got  # go test (alias)

# Exercise 6: Test with coverage
gotestcover  # dotfiles function

# Exercise 7: Format and vet
gof  # go fmt ./...
gov  # go vet ./...

# Exercise 8: Add dependency
gog github.com/spf13/cobra
gomt  # go mod tidy
```

### Advanced Exercises
```bash
# Exercise 9: Create Cobra CLI
cobranew mycli
cobradd serve

# Exercise 10: Cross-compile
gobuildall mycli

# Exercise 11: Run benchmarks
gobench
```

## Context

@components/go/functions.sh
@components/go/aliases.sh

## Coaching Style

- Start with basic Go syntax
- Emphasize error handling early
- Show idiomatic patterns
- Use dotfiles functions and aliases
- Build toward a working CLI app
