# Update Component Commands/Agents for Config Generation Strategy

## Summary

Update component-related Claude commands and agents to reflect the new strategy: **generate config files from `files:` section in config.yaml** instead of storing static config files in `config/` directories.

## Current State

The system has TWO mechanisms:
- **`sync_files:`** - Static files copied/symlinked from dotfiles repo (OLD pattern)
- **`files:`** - Config generated from declarative values in config.yaml (NEW pattern)

Commands and agents still reference the old static file pattern.

## Files to Update

### 1. `ai/claude/agents/component-refactor-expert.md`
Update to emphasize:
- `files:` section is the primary way to manage tool configs
- Static `config/` directories are deprecated for new components
- Show examples of `files:` with format writers (tmux, ghostty, iterm2, json, yaml)
- Clarify when `sync_files:` is still appropriate (SSH configs, credentials)

### 2. `ai/claude/commands/component/component-new.md`
Update creation workflow to:
- Generate `files:` section in config.yaml by default
- Remove references to `config/` directory for storing static files
- Include guidance on choosing format (json, yaml, toml, ghostty, tmux, iterm2)

### 3. `ai/claude/commands/component/component-gen-template.md`
Update template structure:
- Remove `config/.gitkeep` placeholder
- Add example `files:` section in generated config.yaml
- Update workflow to include config file generation setup

### 4. `ai/claude/commands/component/component-config-files-add.md`
This is already correct but needs enhancement:
- Add more format examples (iterm2, tmux patterns)
- Clarify this is the PRIMARY way to add config files
- Add migration guidance from static files to generated

### 5. `ai/claude/commands/component/component-validate.md`
Add validation for:
- `files:` section structure
- Format writer existence
- Target path validity
- Values match expected schema

## Key Message Changes

### OLD Pattern (Deprecate)
```yaml
# Store static files in config/ directory
sync_files:
  - source: "config/tool/settings.json"
    target: "${HOME}/.config/tool/settings.json"
    mode: "symlink"
```

### NEW Pattern (Promote)
```yaml
# Generate from declarative values
files:
  - target: "${HOME}/.config/tool/settings.json"
    format: json
    values:
      setting1: "value1"
      setting2: true
```

## When to Use Each Pattern

| Pattern | Use Case |
|---------|----------|
| `files:` (generate) | Tool configs (tmux, ghostty, vscode, iterm2) |
| `sync_files: symlink` | Credentials, SSH keys, git configs with includes |
| `sync_files: copy` | Files needing strict permissions (600) |
| `sync_files: merge` | User customizations overlay base config |

## Implementation Steps

1. Read and update `component-refactor-expert.md` agent
2. Read and update `component-new.md` command
3. Read and update `component-gen-template.md` command
4. Read and update `component-config-files-add.md` command
5. Read and update `component-validate.md` command
6. Commit changes

## Verification

1. Run `/component-new test-component` - should generate config.yaml with `files:` example
2. Review generated template structure - no `config/` directory
3. Validate updated agents provide correct guidance
