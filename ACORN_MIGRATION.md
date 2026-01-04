# Acorn CLI Migration Plan

This document tracks the migration of shell components to the acorn Go CLI.

## Current State

### Completed
- [x] Initial Go CLI structure (Cobra + Viper)
- [x] `acorn component list` - List all components
- [x] `acorn component status` - Health checks
- [x] `acorn component validate` - Validate YAML
- [x] `acorn component info` - Detailed info
- [x] `acorn migrate analyze` - Analyze shell functions
- [x] `acorn migrate plan` - Prioritized migration plan
- [x] `acorn migrate report` - Comprehensive report
- [x] Build system (Makefile targets)
- [x] golangci-lint configuration
- [x] goreleaser configuration
- [x] Migration analysis agent (`acorn-migration.md`)

### Build Commands
```bash
make acorn-build      # Build binary
make acorn-install    # Install to GOPATH/bin
make acorn-check      # Run all checks (fmt, vet, lint, test)
make acorn-test       # Run tests
```

---

## Migration Priority Queue

Based on `acorn migrate plan` output (sorted by migration score):

### Tier 1: High Priority (100% score)
These have ALL action functions - maximum benefit from Go migration.

| # | Component | Functions | Status | Notes |
|---|-----------|-----------|--------|-------|
| 1 | tools | 10 | ⬜ TODO | Tool version management - good starter |
| 2 | go | 8 | ⬜ TODO | Go dev helpers |
| 3 | vscode | 12 | ⬜ TODO | VS Code project helpers |
| 4 | secrets | 11 | ⬜ TODO | Secret loading/validation |
| 5 | python | 11 | ⬜ TODO | Python/UV environment |
| 6 | ghostty | 7 | ⬜ TODO | Terminal config |
| 7 | ollama | 11 | ⬜ TODO | Local AI management |
| 8 | huggingface | 5 | ⬜ TODO | HF model management |
| 9 | neovim | 5 | ⬜ TODO | Neovim helpers |
| 10 | claude | 18 | ⬜ TODO | Claude Code integration |
| 11 | iterm2 | 10 | ⬜ TODO | iTerm2 config |

### Tier 2: High Priority (90%+ score)
| # | Component | Functions | Status | Notes |
|---|-----------|-----------|--------|-------|
| 12 | github | 11 | ⬜ TODO | GitHub CLI helpers |
| 13 | kubernetes | 9 | ⬜ TODO | kubectl/helm helpers |
| 14 | tmux | 27 | ⬜ TODO | Session management (largest) |
| 15 | node | 6 | ⬜ TODO | Node/NVM helpers |

### Tier 3: Medium Priority (50-90% score)
| # | Component | Functions | Status | Notes |
|---|-----------|-----------|--------|-------|
| 16 | cloudflare | 14 | ⬜ TODO | Workers/Pages |
| 17 | database | 4 | ⬜ TODO | DB service management |
| 18 | shell | 8 | ⬜ TODO | Core shell functions |
| 19 | git | 6 | ⬜ TODO | Git helpers |

### Tier 4: Low Priority (<50% score)
| # | Component | Functions | Status | Notes |
|---|-----------|-----------|--------|-------|
| 20 | fzf | 5 | ⬜ SKIP | Mostly shell integration |

---

## Recommended First Migration: `tools`

The `tools` component is the best starting point:
- 100% migration score
- Self-contained (no complex dependencies)
- Clear command mapping
- Useful for dogfooding

### Proposed Command Structure
```
acorn tools
├── status      # Show all tool versions and status
├── check       # Check if specific tools are installed
├── outdated    # Show outdated tools
├── install     # Install missing tools
└── update      # Update tools via brew/go/npm
```

### Functions to Migrate
```bash
acorn migrate analyze tools
```

| Shell Function | Go Command | Complexity |
|---------------|------------|------------|
| tools_status() | acorn tools status | medium |
| check_tool() | acorn tools check | low |
| check_version() | acorn tools version | low |
| install_go_tools() | acorn tools install go | medium |
| update_brew() | acorn tools update brew | low |
| ... | ... | ... |

---

## Migration Checklist per Component

For each component migration:

### 1. Analysis
- [ ] Run `acorn migrate analyze <component>`
- [ ] Identify action functions (migrate)
- [ ] Identify wrapper functions (keep as shell)
- [ ] Map to CLI command structure

### 2. Implementation
- [ ] Create `internal/<component>/` package
- [ ] Create `internal/cmd/<component>.go`
- [ ] Implement each subcommand
- [ ] Add output format support (-o table|json|yaml)
- [ ] Add shell completion

### 3. Testing
- [ ] Write unit tests for domain logic
- [ ] Test CLI commands manually
- [ ] Run `make acorn-check`

### 4. Shell Integration (Optional)
- [ ] Update `functions.sh` to call acorn
- [ ] Keep aliases unchanged
- [ ] Keep env vars unchanged

### 5. Documentation
- [ ] Update component README
- [ ] Add command examples
- [ ] Update this migration doc

---

## File Structure Goal

```
internal/
├── cmd/
│   ├── root.go           ✅ Done
│   ├── component.go      ✅ Done
│   ├── migrate.go        ✅ Done
│   ├── tools.go          ⬜ TODO
│   ├── python.go         ⬜ TODO
│   ├── claude.go         ⬜ TODO
│   ├── tmux.go           ⬜ TODO
│   ├── dotfiles.go       ⬜ TODO
│   └── ...
├── component/            ✅ Done
├── config/               ✅ Done
├── migrate/              ✅ Done
├── output/               ✅ Done
├── version/              ✅ Done
├── tools/                ⬜ TODO
├── python/               ⬜ TODO
├── claude/               ⬜ TODO
├── tmux/                 ⬜ TODO
├── dotfiles/             ⬜ TODO
└── ...
```

---

## Core Functionality to Add

Beyond component migration, these core features should be added:

### Dotfiles Management
```
acorn dotfiles
├── status      # Git status + bootstrap check
├── audit       # Full drift report
├── drift       # Per-component drift view
├── link        # Link app configs
├── unlink      # Remove config links
├── inject      # Install bootstrap files
├── eject       # Remove bootstrap files
├── pull        # Git pull
├── push        # Git commit + push
└── sync        # Full sync (pull + push)
```

### Shell Integration
```
acorn shell
├── generate    # Generate shell functions that call acorn
├── completions # Generate shell completions
└── install     # Add acorn to shell config
```

---

## Commands Reference

```bash
# Current working commands
acorn version
acorn component list
acorn component status [name]
acorn component validate [name]
acorn component info <name>
acorn migrate analyze [component]
acorn migrate plan
acorn migrate report

# Output formats
acorn component list -o json
acorn component list -o yaml
acorn component list -o table  # default
```

---

## Quick Start for Next Session

```bash
# 1. Build acorn
make acorn-build

# 2. Check current migration status
./build/acorn migrate report

# 3. Analyze the next component
./build/acorn migrate analyze tools

# 4. Start migration (use acorn-migration agent)
# The agent will help generate Go code

# 5. Test
make acorn-check
./build/acorn tools status
```

---

## Notes

- Keep shell aliases - they must stay in shell for tab completion
- Keep environment variables - they set runtime state
- Simple wrappers can stay as shell functions
- Action functions (create, modify, status) → Go commands
- Use `acorn-migration` Claude agent for code generation
