---
description: Validate a component against template standards
---

Validate component: $ARGUMENTS

## Instructions

Thoroughly validate the specified component against all standards.

### 1. Check Component Location

Verify component exists at:
- `internal/componentconfig/config/$ARGUMENTS/config.yaml`

### 2. Validate config.yaml Structure

**Required fields:**
- [ ] `name` - matches directory name
- [ ] `version` - valid semver (X.Y.Z)
- [ ] `description` - non-empty string

**Shell integration (at least one):**
- [ ] `env` - environment variables
- [ ] `aliases` - shell aliases
- [ ] `shell_functions` - functions requiring shell state

### 3. Validate files: Section (Config Generation)

If `files:` section exists:
- [ ] Each entry has `target` path (with valid env vars)
- [ ] Each entry has `format` (json, yaml, toml, ghostty, tmux, iterm2)
- [ ] Each entry has `values` object
- [ ] Format writer exists in `internal/configfile/`

**Test generation:**
```bash
acorn shell generate
ls generated/$ARGUMENTS/
```

### 4. Validate sync_files: Section (Static Files)

If `sync_files:` section exists:
- [ ] Each entry has `source` path (relative to $DOTFILES_ROOT)
- [ ] Each entry has `target` path
- [ ] Each entry has `mode` (symlink, copy, merge)
- [ ] Source files exist in `config/$ARGUMENTS/`
- [ ] Appropriate use case (credentials, SSH, permissions)

### 5. Validate install: Section

If `install:` section exists:
- [ ] Each tool has `name`
- [ ] Each tool has `check` command
- [ ] Each tool has `methods` for darwin/linux
- [ ] Method types are valid (brew, apt, npm, pip, go, curl)

### 6. Build Check

```bash
go build ./...
```

Component should not cause build failures.

## Output Format

```
Component Validation: $ARGUMENTS
================================

Location:            [PASS/FAIL]
  - config.yaml      [OK/MISSING]

Metadata:            [PASS/FAIL]
  - name             [OK/ERROR: <reason>]
  - version          [OK/ERROR: <reason>]
  - description      [OK/ERROR: <reason>]

Shell Integration:   [PASS/FAIL]
  - env              [N vars defined]
  - aliases          [N aliases defined]
  - shell_functions  [N functions defined]

Config Generation:   [PASS/SKIP]
  - files section    [N files configured]
  - format writers   [OK/ERROR: missing <format>]
  - generation test  [OK/ERROR: <reason>]

Static Files:        [PASS/SKIP]
  - sync_files       [N files configured]
  - source files     [OK/MISSING: <file>]

Installation:        [PASS/SKIP]
  - tools defined    [N tools]
  - check commands   [OK/ERROR]
  - methods          [OK/ERROR]

Build Check:         [PASS/FAIL]
  - go build         [OK/ERROR]

Overall:             [VALID/INVALID]

Issues Found:
  1. <issue description>
  2. <issue description>

Recommendations:
  - Consider using files: section instead of sync_files: for tool configs
  - <other suggestions>
```

## Config Strategy Recommendations

After validation, suggest improvements:

| Current | Recommendation |
|---------|---------------|
| `sync_files:` with JSON | Migrate to `files:` with `format: json` |
| `sync_files:` with YAML | Migrate to `files:` with `format: yaml` |
| Static tool configs | Use `files:` for declarative generation |
| `sync_files:` for SSH | Keep as-is (needs 600 permissions) |
| `sync_files:` for credentials | Keep as-is (sensitive data) |
