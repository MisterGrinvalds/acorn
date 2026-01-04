---
description: Create and configure a Python virtual environment
argument-hint: [name]
allowed-tools: Read, Bash
---

## Task

Help the user create a Python virtual environment using UV (with pip fallback).

## Process

1. **Choose name** - From `$ARGUMENTS` or default to `.venv`
2. **Create environment** - Using UV or python3
3. **Activate** - Source the activate script
4. **Verify** - Check Python path and version

## Quick Creation

Using dotfiles function:
```bash
# Create and activate .venv
mkvenv

# Create with custom name
mkvenv myenv

# Using alias
mkv myenv
```

## Manual Creation

### With UV (Recommended - Faster)
```bash
# Create
uv venv .venv

# With specific Python version
uv venv .venv --python 3.11

# Activate
source .venv/bin/activate
```

### With python3 (Fallback)
```bash
# Create
python3 -m venv .venv

# Activate
source .venv/bin/activate
```

## Verification

```bash
# Check Python location
which python
# Should show: /path/to/.venv/bin/python

# Check version
python --version

# Check environment variable
echo $VIRTUAL_ENV

# List installed packages
pip list
```

## Deactivation

```bash
# Using dotfiles function
dvenv

# Using alias
dv

# Standard command
deactivate
```

## Best Practices

1. **Name**: Use `.venv` (standard, gitignored)
2. **Location**: Project root directory
3. **Git**: Add `.venv/` to `.gitignore`
4. **Python version**: Match project requirements

## Environment Locations

```bash
# Per-project (recommended)
project/
├── .venv/
├── src/
└── pyproject.toml

# Centralized (alternative)
export ENVS_LOCATION=~/.virtualenvs
mkvenv myproject  # Creates ~/.virtualenvs/myproject
venv myproject    # Activates from ENVS_LOCATION
```

## Troubleshooting

### "command not found: python"
```bash
# Check if activated
echo $VIRTUAL_ENV

# Reactivate
source .venv/bin/activate
```

### Wrong Python version
```bash
# Recreate with specific version
rm -rf .venv
uv venv .venv --python 3.11
```

## Dotfiles Integration

- `mkvenv [name]` / `mkv` - Create and activate
- `venv [name]` - Activate existing
- `dvenv` / `dv` - Deactivate
