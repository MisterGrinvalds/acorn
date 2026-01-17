# Acorn Output Management System

**Date:** 2026-01-17
**Status:** ✅ Implemented and tested

## Overview

Acorn's output management system has been redesigned to use the `.sapling` directory as a unified location for both configuration sources and generated outputs. This enables:

1. **Git-tracked dotfiles** - Version control your dotfiles in a separate repository
2. **Template rendering** - Dynamic config generation with Go templates
3. **Runtime configuration** - No rebuild needed when configs change
4. **Organized outputs** - Clear separation between source and generated files

## Architecture

### Directory Structure

```
.sapling/
├── .git/                    # Separate git repository for dotfiles
├── .gitignore              # Git ignore patterns
├── README.md               # Quick start guide
├── STRUCTURE.md            # Detailed documentation
│
├── config/                 # Source configurations (YAML)
│   ├── git/
│   │   └── config.yaml    # Defines git dotfiles to generate
│   ├── neovim/
│   │   └── config.yaml    # Defines neovim dotfiles to generate
│   ├── tmux/
│   │   └── config.yaml
│   └── ...                # One directory per component
│
└── generated/              # Generated outputs (created by acorn)
    ├── git/
    │   └── .gitconfig     # Generated from config/git/config.yaml
    ├── neovim/
    │   └── init.lua       # Generated from config/neovim/config.yaml
    ├── shell/             # Shell integration scripts
    │   ├── git.sh         # Git aliases and functions
    │   ├── neovim.sh      # Neovim aliases and functions
    │   └── entrypoint.sh  # Main acorn shell integration
    └── ...
```

### How It Works

```
┌──────────────────────────────────────────────────────────────┐
│                      Configuration Flow                       │
└──────────────────────────────────────────────────────────────┘

1. SOURCE CONFIGS
   .sapling/config/git/config.yaml
   ├── Defines what to generate
   ├── Specifies target paths
   ├── Contains values/templates
   └── Declares formats

2. ACORN READS & RENDERS
   acorn component generate git
   ├── Reads .sapling/config/git/config.yaml
   ├── Renders templates (if using)
   ├── Applies format writers (gitconfig, json, yaml, etc.)
   └── Generates output files

3. OUTPUTS WRITTEN
   .sapling/generated/git/.gitconfig
   ├── Written to .sapling/generated/git/
   ├── Metadata tracked (symlink target, format, etc.)
   └── Ready for symlinking

4. SYMLINKS CREATED
   ~/.gitconfig -> .sapling/generated/git/.gitconfig
   ├── XDG path points to generated file
   ├── System uses .gitconfig normally
   └── Changes tracked in .sapling git repo
```

## Implementation Details

### Config Package (config/config.go)

New package at repository root for managing `.sapling` paths:

```go
package config

// Get the .sapling root directory
root, err := config.SaplingRoot()
// Returns: /path/to/project/.sapling

// Get the generated directory
genDir, err := config.GeneratedDir()
// Returns: /path/to/project/.sapling/generated

// Load a component config
data, err := config.GetConfig("git")
// Reads: .sapling/config/git/config.yaml

// Load with template rendering
data, err := config.GetConfigWithTemplate("git", map[string]any{
    "UserName": "john",
    "UserEmail": "john@example.com",
})

// List all components
components, err := config.ListComponents()
// Returns: ["git", "neovim", "tmux", ...]
```

**Features:**
- ✅ Auto-discovery of `.sapling` directory (walks up from CWD)
- ✅ `SAPLING_DIR` environment variable support
- ✅ Template rendering with Go's text/template
- ✅ Component enumeration and validation
- ✅ Comprehensive test suite (7/7 tests passing)

### ConfigFile Manager (internal/utils/configfile/writer.go)

Updated to use `.sapling/generated` by default:

```go
// Creates manager that writes to .sapling/generated
manager := configfile.NewManager(dryRun)

// Generate a file for a component
result, err := manager.GenerateFileForComponent("git", fileConfig)
// Writes to: .sapling/generated/git/{filename}
// Tracks:    Symlink target (e.g., ~/.gitconfig)
```

**Behavior:**
- If `.sapling` exists → writes to `.sapling/generated/{component}/`
- If `.sapling` doesn't exist → falls back to direct writes (legacy mode)
- DryRun mode → previews without writing

### Shell Manager (internal/components/terminal/shell/shell.go)

Updated to write shell scripts to `.sapling/generated/shell/`:

```go
manager := shell.NewManager(config)
result, err := manager.GenerateAll()
// Writes to: .sapling/generated/shell/
// Creates:
//   - .sapling/generated/shell/git.sh
//   - .sapling/generated/shell/neovim.sh
//   - .sapling/generated/shell/entrypoint.sh
```

### Setup Command (internal/cmd/setup.go)

Updated to use `.sapling/generated`:

```bash
acorn setup
```

**Steps:**
1. Build acorn binary
2. Generate shell scripts → `.sapling/generated/shell/`
3. Inject shell integration → `~/.bashrc` or `~/.zshrc`
4. Create symlinks from XDG paths → `.sapling/generated/`
5. Sync component configurations

## Usage Examples

### Example 1: Generate Git Config

```bash
# 1. Edit source config
vim .sapling/config/git/config.yaml

# 2. Generate outputs
acorn component generate git

# 3. Files created:
#    .sapling/generated/git/.gitconfig
#
# 4. Symlink created:
#    ~/.gitconfig -> .sapling/generated/git/.gitconfig
```

### Example 2: Template Rendering

```yaml
# .sapling/config/git/config.yaml
name: git
files:
  - target: "${HOME}/.gitconfig"
    format: gitconfig
    values:
      user:
        name: {{ .GitUserName }}
        email: {{ .GitUserEmail }}
```

```go
// In code
templateData := map[string]any{
    "GitUserName":  "John Doe",
    "GitUserEmail": "john@example.com",
}

manager := configfile.NewManager(false)
result, err := manager.GenerateFileForComponent("git", fileConfig)
```

### Example 3: New Machine Setup

```bash
# 1. Clone dotfiles repo
git clone https://github.com/yourusername/dotfiles.git ~/.sapling

# 2. Clone acorn (or use system installation)
cd ~/path/to/acorn

# 3. Run setup
acorn setup

# Everything is now:
# - Generated from .sapling/config/
# - Written to .sapling/generated/
# - Symlinked to XDG paths
```

## Environment Variables

### SAPLING_DIR

Override the `.sapling` location:

```bash
export SAPLING_DIR=/path/to/my/dotfiles/.sapling
acorn setup
```

Default behavior: Searches upward from CWD to find `.sapling/`

### DOTFILES_ROOT (Legacy)

Still supported for backward compatibility, but `.sapling` takes precedence.

## Git Repository Setup

The `.sapling` directory is designed to be its own git repository:

```bash
cd .sapling
git init
git remote add origin https://github.com/yourusername/dotfiles.git

# Decide what to track
git add config/           # Always track source configs
git add generated/        # Optional: track generated outputs
git add README.md STRUCTURE.md .gitignore

git commit -m "Initial dotfiles commit"
git push -u origin main
```

**Recommendation:** Track `generated/` for backup and visibility.

## Testing

All tests pass:

```bash
$ go test ./...
ok      github.com/mistergrinvalds/acorn/config                       0.402s
ok      github.com/mistergrinvalds/acorn/internal/components/terminal/shell    0.207s
ok      github.com/mistergrinvalds/acorn/internal/utils/config                 0.383s
ok      github.com/mistergrinvalds/acorn/internal/utils/configfile             0.375s
ok      github.com/mistergrinvalds/acorn/internal/utils/installer              0.593s
```

**Test coverage:**
- ✅ Config loading from `.sapling/config/`
- ✅ Template rendering
- ✅ Path discovery and fallback
- ✅ Generated directory creation
- ✅ File generation with multiple formats
- ✅ Legacy mode fallback

## Migration Path

From the old embedded config system to `.sapling`:

1. **Move configs** ✅ Completed
   - Moved `config/` → `.sapling/config/`
   - Created new `config/` package at repository root

2. **Update code** ✅ Completed
   - `configfile.Manager` uses `.sapling/generated`
   - `shell.Manager` uses `.sapling/generated/shell/`
   - `setup` command updated

3. **Initialize git** (Manual)
   ```bash
   cd .sapling
   git init
   git add config/ generated/
   git commit -m "Initialize dotfiles"
   ```

## Benefits

### ✅ Version Control
- Your dotfiles are in git
- Track changes over time
- Sync across machines

### ✅ Template Support
- Dynamic values per machine
- Shared configs with personal customization
- Environment-specific settings

### ✅ Runtime Configuration
- No rebuild needed for config changes
- Edit, regenerate, test, commit workflow
- Faster iteration

### ✅ Clear Organization
- Source configs in `config/`
- Generated outputs in `generated/`
- Symlinks to XDG paths
- Everything in one `.sapling` directory

### ✅ Backward Compatible
- Falls back to legacy mode if `.sapling` doesn't exist
- Existing code still works
- Gradual migration possible

## Files Modified

### New Files
- `config/config.go` - New config package
- `config/config_test.go` - Tests
- `config/README.md` - Documentation
- `.sapling/STRUCTURE.md` - Detailed docs
- `.sapling/README.md` - Quick start guide
- `.sapling/.gitignore` - Git ignore patterns
- `ACORN_OUTPUT_MANAGEMENT.md` - This file

### Modified Files
- `internal/utils/configfile/writer.go` - Use `.sapling/generated`
- `internal/utils/configfile/writer_test.go` - Fix test for new behavior
- `internal/components/terminal/shell/shell.go` - Use `.sapling/generated/shell/`
- `internal/cmd/setup.go` - Use `.sapling/generated`
- `.gitignore` - Ignore `.sapling/`

### Deleted Files (Moved to .sapling/)
- `config/embed.go` - Old embedded config system
- `config/*/config.yaml` - All component configs (moved to `.sapling/config/`)
- `ai/` - Claude Code configuration (moved to `.sapling/ai/`)

## Managing the .sapling Repository

Acorn provides built-in commands to manage the `.sapling` git repository:

```bash
# Check git status
acorn sapling status

# Commit changes
acorn sapling commit -a -m "Update dotfiles"

# Push to remote
acorn sapling push

# Pull from remote
acorn sapling pull

# Full sync (pull + commit + push)
acorn sapling sync
acorn sapling sync -m "Custom message"
```

**Sync Workflow:**
1. Pulls latest changes from remote
2. Commits any local changes (optional auto-commit message)
3. Pushes to remote

This ensures your dotfiles are always in sync across machines.

## Next Steps

1. **Set up remote repository** (first time only)
   ```bash
   cd .sapling
   git remote add origin https://github.com/MisterGrinvalds/sapling.git
   git branch -M main
   cd ..
   ```

2. **Initial sync to remote**
   ```bash
   acorn sapling sync -m "Initial dotfiles configuration"
   ```

3. **Test the workflow**
   ```bash
   acorn component generate git
   ls -la .sapling/generated/git/
   ls -la ~/.gitconfig  # Should be symlink
   ```

4. **Daily usage**
   ```bash
   # Make changes to configs
   vim .sapling/config/git/config.yaml

   # Regenerate
   acorn component generate git

   # Sync to remote
   acorn sapling sync -m "Update git aliases"
   ```

## References

- `.sapling/STRUCTURE.md` - Detailed structure documentation
- `.sapling/README.md` - User-facing quick start
- `config/README.md` - Config package API documentation

---

**Implementation Complete:** 2026-01-17
**All Tests Passing:** ✅
**Build Status:** ✅
**Ready for Use:** ✅
