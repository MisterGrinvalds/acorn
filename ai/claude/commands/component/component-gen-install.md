---
description: Generate installation config for a component
argument_hints:
  - tmux
  - go
  - python
  - node
  - kubernetes
---

Generate `install:` section in config.yaml for: $ARGUMENTS

## Instructions

Add declarative installation configuration to a component's config.yaml file.

### 1. Find Component Config

Read `internal/componentconfig/config/$ARGUMENTS/config.yaml`

### 2. Identify Required Tools

Determine what tools this component needs:

| Component | Tools |
|-----------|-------|
| tmux | tmux, tpm (plugin manager) |
| go | go |
| python | python3, uv (package manager) |
| node | node, nvm, pnpm |
| kubernetes | kubectl, helm, k9s |
| ollama | ollama |
| neovim | nvim |
| fzf | fzf |

### 3. Research Package Names by Platform

For each tool, determine installation methods:

| Install Type | Use For | Example |
|--------------|---------|---------|
| brew | macOS packages | `type: brew` |
| apt | Debian/Ubuntu | `type: apt` |
| npm | Node packages | `type: npm`, `global: true` |
| pip | Python packages | `type: pip`, `global: true` |
| go | Go binaries | `type: go`, `package: github.com/...` |
| curl | Script installers | `type: curl`, `url: https://...` |

### 4. Generate install: Section

Add to the component's config.yaml:

```yaml
# Installation configuration
install:
  tools:
    - name: <tool-name>
      description: <what it does>
      check: "command -v <tool-name>"
      methods:
        darwin:
          type: brew
          package: <brew-package-name>
        linux:
          type: apt
          package: <apt-package-name>
      requires:
        - <prerequisite-component:tool>  # e.g., node:npm
      post_install:
        message: "Run '<command>' to complete setup"
```

### 5. Schema Reference

**ToolInstall fields:**
- `name` (required): Tool/binary name
- `description`: Human-readable description
- `check` (required): Command to verify installation (e.g., `command -v go`)
- `methods` (required): Platform-specific install methods
- `requires`: Prerequisites (format: `component:tool` or just `tool`)
- `post_install.message`: Message to show after install
- `post_install.commands`: Commands to run after install

**InstallMethod fields:**
- `type` (required): brew, apt, npm, pip, go, curl, binary
- `package`: Package name if different from tool name
- `global`: For npm/pip, install globally (default: false)
- `url`: For curl type, the script URL
- `args`: Additional install arguments

### 6. Example: Node Component

```yaml
install:
  tools:
    - name: node
      description: Node.js runtime (includes npm)
      check: "command -v node"
      methods:
        darwin:
          type: brew
          package: node
        linux:
          type: apt
          package: nodejs

    - name: pnpm
      description: Fast, disk space efficient package manager
      check: "command -v pnpm"
      methods:
        darwin:
          type: npm
          package: pnpm
          global: true
        linux:
          type: npm
          package: pnpm
          global: true
      requires:
        - node:npm
```

### 7. Example: Go Component

```yaml
install:
  tools:
    - name: go
      description: Go programming language
      check: "command -v go"
      methods:
        darwin:
          type: brew
          package: go
        linux:
          type: apt
          package: golang-go
      post_install:
        message: "Ensure GOPATH is set in your shell config"
```

### 8. Test Installation

After adding the config, test with:

```bash
acorn $ARGUMENTS install --dry-run
```

This shows the installation plan without executing anything.

### 9. Report

```
Generated Installation Config: $ARGUMENTS
==========================================

Config file: internal/componentconfig/config/$ARGUMENTS/config.yaml

Tools configured:
  - <tool1>: <methods>
  - <tool2>: <methods>

Prerequisites:
  - <any requires entries>

Test command:
  acorn $ARGUMENTS install --dry-run

Full install:
  acorn $ARGUMENTS install
```

## CLI Usage

After configuration, users can install with:

```bash
# Preview what will be installed
acorn <component> install --dry-run

# Install all tools
acorn <component> install

# Verbose output
acorn <component> install --verbose
```

The installer automatically:
- Detects platform (darwin/linux)
- Resolves prerequisites recursively
- Skips already-installed tools
- Shows post-install messages
