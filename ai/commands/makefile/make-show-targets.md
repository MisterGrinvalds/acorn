---
description: Show all available Makefile targets
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Bash
---

# Show Makefile Targets

Display all available Makefile targets with their descriptions in an organized format.

## Task

1. Check if a Makefile exists
2. Parse and display all targets organized by section:
   - Extract section headers (##@)
   - Extract target names and descriptions (##)
   - Show dependencies for each target
   - Identify phony vs file targets

3. Display in organized format:
   ```
   General
     help                 Display this help message

   Development
     run                  Run the application
     fmt                  Format code
     lint                 Run linters

   Building
     build                Build binary [clean]
     install              Install binary to GOPATH/bin [build]
   ```

4. Show additional information:
   - Total number of targets
   - Missing .PHONY declarations (warnings)
   - Undocumented targets (warnings)
   - Dependencies in brackets

5. Optionally filter by section or search pattern

## Usage Examples

- Show all targets: `/make-show-targets`
- Show specific section: `/make-show-targets section=Development`
- Search targets: `/make-show-targets search=test`

## Output Format

```
Makefile Targets (15 total)
===========================

##@ Development (5)
  run                    Run the application
  fmt                    Format code
  lint                   Run linters
  vet                    Run go vet
  watch                  Watch for changes [requires: air]

##@ Building (3)
  build                  Build binary [deps: clean]
  build-all              Build for all platforms [deps: clean]
  install                Install to GOPATH/bin [deps: build]

⚠ Warnings:
  • 2 targets without help comments
  • 1 target missing .PHONY
```
