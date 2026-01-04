---
description: Explain Go concepts, tools, and workflows
argument-hint: [topic]
allowed-tools: Read, Glob, Grep
---

## Task

Explain the requested topic about Go development. If no specific topic provided, give an overview of Go tooling and the dotfiles integration.

## Topics

If `$ARGUMENTS` is provided, explain that specific topic:

- **modules** - Go modules and dependency management
- **packages** - Package organization and imports
- **testing** - Testing patterns and coverage
- **errors** - Error handling patterns
- **interfaces** - Interface design and usage
- **concurrency** - Goroutines and channels
- **cobra** - Cobra CLI framework
- **building** - Build process and cross-compilation
- **embedding** - Embedding files with go:embed
- **generics** - Generics in Go
- **context** - Context usage patterns
- **json** - JSON encoding/decoding
- **http** - HTTP server and client patterns

## Context

Reference these files for accurate information:
@components/go/component.yaml
@components/go/functions.sh
@components/go/aliases.sh

## Response Format

1. **Concept overview** - What it is and why it matters
2. **Key commands** - Essential Go commands
3. **Dotfiles integration** - Available functions and aliases
4. **Examples** - Practical code examples
5. **Best practices** - Idiomatic Go patterns
