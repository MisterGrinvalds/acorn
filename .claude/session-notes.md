# Component Config File Generation System

Last updated: 2026-01-06

## Current Focus

Adding generic config file generation support to acorn components via a new `files:` array in config.yaml.

## Session Summary

See: `SESSION-2026-01-06-1854.md` for full details.

**Accomplished:**
- Created `/component:config-files-add` Claude command
- Designed `files:` array schema for multi-format config generation
- Analyzed Ghostty component structure
- Explored existing component config system

**Key Decisions:**
- No comment preservation (values only)
- config.yaml wins on merge
- One format at a time (start with Ghostty)

## Next Session Priorities

### 1. Schema Extension (High Priority)
- Add FileConfig and FieldSchema types to `internal/componentconfig/schema.go`
- Add Files field to BaseConfig
- Update loader MergeConfigs to handle Files

### 2. Ghostty Format Writer (High Priority)
- Create `internal/configfile/` package
- Implement Writer interface
- Create GhosttyWriter with multi-value key support

### 3. Integration (High Priority)
- Update `internal/shell/shell.go` to call file generation
- Create `components/ghostty/config.yaml` with files section
- Test end-to-end

## Design Specification

### files: Array Schema

```yaml
files:
  - target: "${XDG_CONFIG_HOME}/ghostty/config"
    format: ghostty
    schema:
      theme: { type: string, default: "Catppuccin Mocha" }
      font-size: { type: int, default: 14 }
      keybind: { type: list, items: string }
    values:
      theme: "Catppuccin Mocha"
      font-size: 14
      keybind:
        - "super+d=new_split:right"
```

### Supported Formats

| Format | Output Style | Implementation |
|--------|-------------|----------------|
| ghostty | `key = value` | Priority 1 |
| json | `{"key": "value"}` | Future |
| yaml | `key: value` | Future |
| toml | `key = "value"` | Future |
| ini | `[section]\nkey=value` | Future |

## Files to Create

**New:**
- `internal/configfile/writer.go` - Writer interface
- `internal/configfile/ghostty.go` - Ghostty writer
- `internal/configfile/ghostty_reader.go` - Ghostty reader
- `internal/configfile/manager.go` - File processing

**Modified:**
- `internal/componentconfig/schema.go` - Add types
- `internal/componentconfig/loader.go` - Handle Files
- `internal/shell/shell.go` - Call generation
- `components/ghostty/config.yaml` - Add files section

## Resources

- Command: `components/claude/config/commands/component/component-config-files-add.md`
- Plan: `.claude/plans/reactive-mixing-hare.md`
- Ghostty docs: https://ghostty.org/docs/config

## Notes

- Use `/component:config-files-add ghostty` to start implementation
- Ghostty will be first component to use new system
- Framework is generic and reusable for all future components
