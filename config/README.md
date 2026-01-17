# Acorn Configuration System

This package provides runtime configuration loading from the `.sapling/config` directory with support for template rendering.

## Overview

Acorn's configuration system has been redesigned to source component configurations from the `.sapling/config` directory. This allows for:

1. **Runtime Configuration**: Configs are loaded at runtime from `.sapling/config`, not embedded at build time
2. **Template Rendering**: Config files can use Go template syntax for dynamic values
3. **Git-Tracked Configs**: Your dotfiles can be tracked in a separate git repository within `.sapling`

## Directory Structure

```
.sapling/
└── config/
    ├── git/
    │   └── config.yaml
    ├── neovim/
    │   └── config.yaml
    ├── tmux/
    │   └── config.yaml
    └── ... (other components)
```

## Basic Usage

### Loading a Component Config

```go
import "github.com/mistergrinvalds/acorn/config"

// Load raw config
data, err := config.GetConfig("git")
if err != nil {
    log.Fatal(err)
}

// Parse into your struct
var gitConfig GitConfig
if err := yaml.Unmarshal(data, &gitConfig); err != nil {
    log.Fatal(err)
}
```

### List All Components

```go
components, err := config.ListComponents()
if err != nil {
    log.Fatal(err)
}

for _, component := range components {
    fmt.Println(component)
}
```

### Check if Component Exists

```go
if config.HasConfig("git") {
    fmt.Println("Git component exists")
}
```

## Template Rendering

Config files can use Go template syntax for dynamic values. This is useful for:
- User-specific paths
- Environment-specific values
- Computed values

### Example Template Config

```yaml
# .sapling/config/git/config.yaml
name: git
description: Git configuration for {{ .Username }}
env:
  GIT_AUTHOR_NAME: {{ .Username }}
  GIT_AUTHOR_EMAIL: {{ .Email }}
paths:
  - {{ .HomeDir }}/bin
files:
  - target: ~/.gitconfig
    format: gitconfig
    values:
      user.name: {{ .Username }}
      user.email: {{ .Email }}
```

### Loading with Template Data

```go
templateData := map[string]any{
    "Username": "john",
    "Email":    "john@example.com",
    "HomeDir":  "/home/john",
}

data, err := config.GetConfigWithTemplate("git", templateData)
if err != nil {
    log.Fatal(err)
}

// Now data contains the rendered template
```

## Environment Variables

### SAPLING_DIR

Override the location of the `.sapling` directory:

```bash
export SAPLING_DIR=/path/to/my/dotfiles/.sapling
acorn component list
```

If not set, Acorn searches upward from the current directory to find `.sapling/config`.

## Walking All Configs

Process all component configs:

```go
err := config.WalkConfigs(func(component string, path string) error {
    fmt.Printf("Processing %s at %s\n", component, path)
    return nil
})
```

## Migration from Embedded Configs

The previous system embedded configs at build time using `//go:embed`. The new system:

1. **Loads at runtime**: More flexible, no rebuild needed for config changes
2. **Supports templates**: Dynamic values based on environment
3. **Git-trackable**: Store your configs in a git repository within `.sapling`

### Migration Steps

If you're migrating from the old embedded system:

1. Move your `config/` directory to `.sapling/config`
2. Update imports to use the new `config` package
3. Add `.sapling/` to `.gitignore` if you want to track it separately
4. (Optional) Convert static configs to use template syntax

## Best Practices

1. **Use Templates Sparingly**: Only use template syntax where dynamic values are actually needed
2. **Validate Rendered Output**: When using templates, validate the rendered YAML is valid
3. **Document Template Variables**: Comment what template variables your config expects
4. **Test Without Templates First**: Ensure static configs work before adding template logic

## Example: Complete Component Config

```yaml
# .sapling/config/neovim/config.yaml
name: neovim
description: Neovim text editor configuration
version: "0.10.0"

env:
  EDITOR: nvim
  VISUAL: nvim

paths:
  - {{ .ConfigHome }}/nvim/bin

files:
  - target: ~/.config/nvim/init.lua
    format: raw
    values:
      content: |
        -- Neovim configuration for {{ .Username }}
        vim.opt.number = true
        vim.opt.relativenumber = true

install:
  tools:
    - name: neovim
      check: command -v nvim
      methods:
        - platform: darwin
          type: homebrew
          package: neovim
        - platform: linux
          type: apt
          package: neovim
```

## Testing

Run the config package tests:

```bash
go test ./config -v
```

The test suite includes:
- Loading component configs
- Listing all components
- Template rendering
- Config existence checks
- Custom SAPLING_DIR handling
