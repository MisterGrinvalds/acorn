---
description: Initialize a new Makefile for Go/Cobra CLI project
model: claude-sonnet-4-5-20250929
allowed-tools: Write, Read, Bash
---

# Initialize Makefile

Create a new Makefile for this Go/Cobra CLI project following best practices.

## Task

1. Check if a Makefile already exists in the project root
2. If it exists, ask the user if they want to overwrite it
3. If not exists or user confirms overwrite, create a new Makefile with:
   - Project metadata variables (APP_NAME, VERSION, BUILD_TIME)
   - Go command variables (GOCMD, GOBUILD, GOTEST, etc.)
   - Build directory and binary name configuration
   - Essential target groups:
     - **General**: help
     - **Development**: run, fmt, lint, vet
     - **Building**: build, install
     - **Testing**: test, test-coverage
     - **Cobra**: cobra-add
     - **Utilities**: clean, deps, version
   - Self-documenting help system using awk
   - .PHONY declarations for all targets

4. Use the project name from go.mod if available, otherwise default to directory name
5. Follow the style guide from the makefile-expert agent
6. Display the help output after creation

## Requirements

- All targets must have ## comments for the help system
- Use @ prefix for silent commands
- Use := for variable assignment (immediate expansion)
- Group targets with ##@ section headers
- Follow kebab-case naming for multi-word targets
