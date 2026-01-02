# Template Component

This is a template for creating new components. Copy this directory to create a new component.

## Quick Start

```bash
# Copy template to new component
cp -r components/_template components/mycomponent

# Edit component.yaml with your metadata
$EDITOR components/mycomponent/component.yaml

# Add your aliases, functions, etc.
$EDITOR components/mycomponent/aliases.sh
$EDITOR components/mycomponent/functions.sh
```

## Files

| File | Purpose | When Loaded |
|------|---------|-------------|
| `component.yaml` | Component metadata | Discovery phase |
| `env.sh` | Environment variables | All shells |
| `aliases.sh` | Shell aliases | Interactive only |
| `functions.sh` | Shell functions | Interactive only |
| `completions.sh` | Tab completions | Interactive only |
| `setup.sh` | Installation script | On demand |
| `README.md` | Documentation | Reference only |

## Component Metadata (component.yaml)

Required fields:
- `name`: Component identifier (lowercase, no spaces)
- `version`: Semantic version (e.g., 1.0.0)
- `description`: Brief description
- `category`: One of: core, dev, cloud, ai, database

Optional fields:
- `requires.tools`: List of CLI tools that must be installed
- `requires.components`: List of components that must load first
- `provides`: What aliases/functions/completions this provides
- `xdg`: XDG subdirectories to create
- `platforms`: Limit to specific platforms (darwin, linux)
- `shells`: Limit to specific shells (bash, zsh)
- `setup`: Installation configuration

## XDG Integration

Use the XDG helper functions to get paths:

```bash
# In your scripts
config_dir=$(xdg_config_dir mycomponent)  # ~/.config/mycomponent
data_dir=$(xdg_data_dir mycomponent)      # ~/.local/share/mycomponent
cache_dir=$(xdg_cache_dir mycomponent)    # ~/.cache/mycomponent
state_dir=$(xdg_state_dir mycomponent)    # ~/.local/state/mycomponent

# Create all directories at once
xdg_ensure_dirs mycomponent
```

## Setup Script

The setup.sh script should:
1. Check prerequisites
2. Install dependencies
3. Configure the component
4. Validate the setup

Run with: `make setup-mycomponent`
