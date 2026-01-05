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

### Shell Integration (NEW)
- [x] `acorn shell status` - Show integration status
- [x] `acorn shell generate [component...]` - Generate shell scripts (per-component)
- [x] `acorn shell inject` - Add to shell rc
- [x] `acorn shell install` - Full setup (generate + inject)
- [x] `acorn shell eject` - Remove from shell rc
- [x] JSON/YAML output for all shell commands (`-o json`)
- [x] Shell components in `internal/shell/components.go`

### Go Development Commands
- [x] `acorn go new` - Initialize Go project
- [x] `acorn go test` - Run tests
- [x] `acorn go cover` - Coverage report
- [x] `acorn go bench` - Benchmarks
- [x] `acorn go build-all` - Cross-compile
- [x] `acorn go clean` - Clean artifacts
- [x] `acorn go env` - Show Go environment
- [x] `acorn go cobra new/add` - Cobra scaffolding

### VS Code Commands
- [x] `acorn vscode workspaces` - List workspaces
- [x] `acorn vscode workspace` - Open workspace
- [x] `acorn vscode project new` - Create project
- [x] `acorn vscode ext list/install/export/essentials` - Extensions
- [x] `acorn vscode config sync/path` - Config management

### Python Commands
- [x] `acorn python venv new/list` - Virtual environment management
- [x] `acorn python init` - Initialize UV project
- [x] `acorn python sync` - Sync dependencies
- [x] `acorn python add/remove` - Package management
- [x] `acorn python run` - Run in project env
- [x] `acorn python env` - Show environment info
- [x] `acorn python fastapi` - FastAPI setup
- [x] `acorn python setup ipython/devtools` - Dev tools setup

### Tmux Commands
- [x] `acorn tmux info` - Show tmux information
- [x] `acorn tmux session list` - List active sessions
- [x] `acorn tmux tpm install/update` - TPM management
- [x] `acorn tmux tpm plugins-install/plugins-update` - Plugin management
- [x] `acorn tmux config reload` - Reload config
- [x] `acorn tmux smug list/new` - Smug session management
- [x] `acorn tmux smug install/link` - Install and link smug
- [x] `acorn tmux smug repo-init/status/pull/push/sync` - Git sync

### Components Migrated & Removed
- [x] `components/go/` - DELETED (shell in acorn)
- [x] `components/tools/` - DELETED (shell in acorn)
- [x] `components/vscode/*.sh` - DELETED (config/ kept for settings.json)
- [x] `components/python/*.sh` - DELETED (config/ kept for startup.py)
- [x] `components/tmux/*.sh` - DELETED (config/ kept for tmux.conf, smug configs)
- [x] `components/claude/*.sh` - DELETED (config/ kept for settings.json, agents/, commands/, subagents/)

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
| 1 | tools | 10 | ✅ DONE | Tool version management - shell files removed |
| 2 | go | 8 | ✅ DONE | Go dev helpers - shell files removed |
| 3 | vscode | 12 | ✅ DONE | VS Code helpers - shell files removed, config kept |
| 4 | secrets | 11 | ⬜ TODO | Secret loading/validation |
| 5 | python | 11 | ✅ DONE | Python/UV environment - shell files removed, config kept |
| 6 | ghostty | 7 | ⬜ TODO | Terminal config |
| 7 | ollama | 11 | ⬜ TODO | Local AI management |
| 8 | huggingface | 5 | ⬜ TODO | HF model management |
| 9 | neovim | 5 | ⬜ TODO | Neovim helpers |
| 10 | claude | 18 | ✅ DONE | Claude Code integration - shell files removed, config kept |
| 11 | iterm2 | 10 | ⬜ TODO | iTerm2 config |

### Tier 2: High Priority (90%+ score)
| # | Component | Functions | Status | Notes |
|---|-----------|-----------|--------|-------|
| 12 | github | 11 | ⬜ TODO | GitHub CLI helpers |
| 13 | kubernetes | 9 | ⬜ TODO | kubectl/helm helpers |
| 14 | tmux | 27 | ✅ DONE | Session management - shell files removed, config kept |
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
- [ ] Identify action functions (migrate to Go)
- [ ] Identify shell-only functions (keep in shell integration)
- [ ] Map to CLI command structure

### 2. Implementation
- [ ] Create `internal/<component>/` package
- [ ] Create `internal/cmd/<component>.go`
- [ ] Implement each subcommand
- [ ] Add output format support (-o table|json|yaml)
- [ ] Add shell completion

### 3. Shell Integration
- [ ] Add component to `internal/shell/components.go`
- [ ] Include env vars, aliases, wrapper functions
- [ ] Test with `acorn shell generate <component> --dry-run`

### 4. Cleanup
- [ ] If no config files: `rm -rf components/<name>/`
- [ ] If has config files: keep only `components/<name>/config/`
  - Remove `aliases.sh`, `env.sh`, `functions.sh`, `completions.sh`
  - Remove `component.yaml`

### 5. Testing
- [ ] Run `make acorn-check`
- [ ] Test CLI commands manually
- [ ] Test shell integration: `acorn shell generate <component>`

### 6. Documentation
- [ ] Update this migration doc (mark as DONE)

---

## Migrated Component Structure

After migration, a component directory is **removed entirely** unless it has config files.

**No config files** → Delete entire directory:
```bash
rm -rf components/go
rm -rf components/tools
```

**Has config files** → Keep only config directory:
```
components/vscode/
└── config/
    ├── settings.json
    ├── keybindings.json
    └── extensions.txt
```

### Currently Removed
- `components/go/` - fully removed
- `components/tools/` - fully removed
- `components/vscode/*.sh` - removed, `config/` kept

---

## File Structure Goal

```
internal/
├── cmd/
│   ├── root.go           ✅ Done
│   ├── component.go      ✅ Done
│   ├── migrate.go        ✅ Done
│   ├── shell.go          ✅ Done (NEW)
│   ├── tools.go          ✅ Done
│   ├── golang.go         ✅ Done
│   ├── vscode.go         ✅ Done
│   ├── python.go         ✅ Done
│   ├── tmux.go           ✅ Done
│   ├── claude.go         ✅ Done
│   ├── dotfiles.go       ⬜ TODO
│   └── ...
├── component/            ✅ Done
├── config/               ✅ Done
├── migrate/              ✅ Done
├── output/               ✅ Done
├── shell/                ✅ Done (NEW - shell script generation)
├── version/              ✅ Done
├── tools/                ✅ Done
├── golang/               ✅ Done
├── vscode/               ✅ Done
├── python/               ✅ Done
├── tmux/                 ✅ Done
├── claude/               ✅ Done
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

# Shell integration
acorn shell status          # Show integration status
acorn shell generate        # Generate shell scripts
acorn shell inject          # Add to shell rc
acorn shell install         # Generate + inject (full setup)
acorn shell eject           # Remove from shell rc
acorn shell uninstall       # Full removal
acorn shell list            # List available components

# Go development
acorn go new <module>       # Create new Go project
acorn go test [pattern]     # Run tests
acorn go cover              # Run tests with coverage
acorn go bench [pattern]    # Run benchmarks
acorn go build-all [name]   # Cross-compile
acorn go clean              # Clean artifacts
acorn go env                # Show Go environment
acorn go cobra new <app>    # Create Cobra CLI
acorn go cobra add <cmd>    # Add Cobra command

# VS Code
acorn vscode workspaces     # List workspaces
acorn vscode workspace <n>  # Open workspace
acorn vscode project new    # Create project
acorn vscode ext list       # List extensions
acorn vscode ext install    # Install from file
acorn vscode ext export     # Export to file
acorn vscode ext essentials # Install essentials
acorn vscode config sync    # Sync from dotfiles
acorn vscode config path    # Show config paths

# Python
acorn python venv new [n]   # Create virtual environment
acorn python venv list      # List venvs
acorn python init [name]    # Initialize UV project
acorn python sync           # Sync dependencies
acorn python add <pkg...>   # Add packages
acorn python remove <pkg>   # Remove packages
acorn python run <cmd>      # Run in project env
acorn python env            # Show Python environment
acorn python fastapi [n]    # Setup FastAPI env
acorn python setup ipython  # Install IPython
acorn python setup devtools # Install dev tools

# Tmux
acorn tmux info             # Show tmux info
acorn tmux session list     # List active sessions
acorn tmux tpm install      # Install TPM
acorn tmux tpm update       # Update TPM
acorn tmux tpm plugins-install  # Install plugins
acorn tmux tpm plugins-update   # Update plugins
acorn tmux config reload    # Reload config
acorn tmux smug list        # List smug sessions
acorn tmux smug new <name>  # Create new smug config
acorn tmux smug install     # Install smug
acorn tmux smug link        # Link smug configs
acorn tmux smug repo-init   # Init smug git repo
acorn tmux smug status      # Show repo status
acorn tmux smug pull        # Pull from remote
acorn tmux smug push [msg]  # Push to remote
acorn tmux smug sync        # Full sync

# Claude Code
acorn claude info           # Show Claude Code info
acorn claude stats          # View usage statistics
acorn claude stats tokens   # View token usage by model
acorn claude stats daily [n]# View daily token usage
acorn claude permissions    # View permissions
acorn claude permissions add <rule> [allow|deny]  # Add permission
acorn claude permissions remove <rule> [allow|deny] # Remove permission
acorn claude settings [global|local|config]  # View settings
acorn claude settings edit [type]  # Edit settings
acorn claude projects       # List projects
acorn claude mcp            # List MCP servers
acorn claude mcp add <name> <url> [type]  # Add MCP server
acorn claude commands       # List custom commands
acorn claude aggregate [dir]  # Aggregate agents/commands from repos
acorn claude aggregate list # List all agents/commands
acorn claude clear [cache|stats]  # Clear cache or stats
acorn claude help           # Show help

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

# 2. See what's done
./build/acorn shell status
./build/acorn shell list

# 3. Next component to migrate: claude
./build/acorn migrate analyze claude

# 4. Migration pattern:
#    a) Create internal/<component>/ package
#    b) Create internal/cmd/<component>.go
#    c) Add <Component>Component() to internal/shell/components.go
#    d) Remove components/<component>/*.sh files (keep config/ if present)
#    e) Test: acorn shell generate <component> -o json
```

## Completed: Claude Component

The `claude` component migration is complete:
- `acorn claude info` - Show Claude Code info summary
- `acorn claude stats` - View usage statistics
- `acorn claude stats tokens` - View token usage by model
- `acorn claude stats daily` - View daily token usage
- `acorn claude permissions` - View/manage permissions
- `acorn claude permissions add/remove` - Modify permission rules
- `acorn claude settings` - View settings (global/local/config)
- `acorn claude projects` - List projects with costs
- `acorn claude mcp` - List MCP servers
- `acorn claude mcp add` - Add MCP server
- `acorn claude commands` - List custom commands
- `acorn claude aggregate` - Aggregate agents/commands from repos
- `acorn claude aggregate list` - List all agents/commands
- `acorn claude clear` - Clear cache/stats
- `acorn claude help` - Show all available functions

Files created:
- `internal/cmd/claude.go`
- `internal/claude/claude.go`
- `internal/claude/stats.go`
- `internal/claude/settings.go`
- `internal/claude/config.go`
- `internal/claude/aggregate.go`

Shell integration added to `internal/shell/components.go`.

## Next Up: Remaining Components
- Tier 1: secrets (11), ollama (11), ghostty (7), huggingface (5), neovim (5), iterm2 (10)
- Tier 2: node (6), github (11), kubernetes (9)

---

## Shell Integration Architecture

Acorn generates shell scripts and injects them into the user's shell configuration.

### Directory Structure

```
~/.config/acorn/           # XDG_CONFIG_HOME/acorn
├── shell.sh               # Main entrypoint (sources all component scripts)
├── go.sh                  # Go: env, aliases, function wrappers
├── vscode.sh              # VSCode: aliases, function wrappers
└── tools.sh               # Tools: function wrappers
```

### How It Works

1. **Generate**: `acorn shell generate` creates shell scripts for each component
2. **Inject**: `acorn shell inject` adds a source line to `~/.bashrc` or `~/.zshrc`
3. **Install**: `acorn shell install` does both (full setup)

### Shell Script Contents

Each component script contains:
- **Environment variables** (e.g., GOPATH, PATH additions)
- **Aliases** (e.g., `gob='go build'`, `c='code .'`)
- **Function wrappers** that call acorn commands

### Example: Go Component

```bash
# ~/.config/acorn/go.sh

# Environment
export GOPATH="${GOPATH:-$HOME/go}"
export GOBIN="$GOPATH/bin"
export PATH="$GOBIN:$PATH"

# Aliases
alias gob='go build'
alias got='go test'
alias gomt='go mod tidy'

# Function wrappers (call acorn)
gonew() { acorn go new "$1" && cd "$1"; }
gobuildall() { acorn go build-all "${1:-app}"; }
```

### Commands

```bash
# Status and listing
acorn shell status              # Show integration status
acorn shell list                # List available components

# Per-component generation
acorn shell generate            # Generate all components + entrypoint
acorn shell generate go         # Generate only go.sh
acorn shell generate go vscode  # Generate specific components

# Injection
acorn shell inject              # Add to shell rc
acorn shell install             # Generate + inject (full setup)
acorn shell eject               # Remove from shell rc
acorn shell uninstall           # Full removal

# Structured output (JSON/YAML)
acorn shell generate go -o json # Returns JSON with file content, target path
acorn shell status -o json      # Returns JSON with status info
acorn shell install -o yaml     # Returns YAML with generate + inject results
```

### Benefits

1. **Consistent aliases** across machines
2. **Functions wrap acorn** for rich CLI features
3. **Easy updates**: regenerate scripts after acorn updates
4. **Clean removal**: eject removes cleanly

---

## Notes

- Aliases stay in shell for tab completion - acorn generates them
- Environment variables stay in shell - acorn generates them
- Simple wrappers call acorn for the heavy lifting
- Action functions (create, modify, status) → Go commands
- Use `acorn-migration` Claude agent for code generation
