# dotfiles.nvim

A minimal Neovim plugin to install and manage the bash-profile dotfiles system.
Perfect for setting up new or foreign environments directly from Neovim.

## Installation

### Using lazy.nvim

```lua
{
  dir = vim.env.DOTFILES_ROOT and (vim.env.DOTFILES_ROOT .. "/components/neovim/plugin") or nil,
  name = "dotfiles",
  config = function()
    require("dotfiles").setup({
      -- Optional: auto-prompt if dotfiles not installed
      auto_setup = { enabled = true },
    })
  end,
  cond = function()
    -- Only load if DOTFILES_ROOT is set or plugin directory exists
    return vim.env.DOTFILES_ROOT ~= nil
  end,
}
```

### Using packer.nvim

```lua
use {
  "~/Repos/personal/bash-profile/components/neovim/plugin",
  as = "dotfiles",
  config = function()
    require("dotfiles").setup()
  end,
}
```

### Manual Installation

Add to your runtimepath in `init.lua`:

```lua
vim.opt.runtimepath:append(vim.env.DOTFILES_ROOT .. "/components/neovim/plugin")
require("dotfiles").setup()
```

## Quick Start (No Config Needed)

If you just want to run the installer without any setup:

```vim
:DotfilesSetup
```

This will auto-configure and launch the interactive installer.

## Commands

| Command | Description |
|---------|-------------|
| `:Dotfiles` | Run interactive installer |
| `:Dotfiles minimal` | Quick install (dotfiles + configs only) |
| `:Dotfiles auto` | Non-interactive full install |
| `:DotfilesMinimal` | Same as `:Dotfiles minimal` |
| `:DotfilesComponents` | Component-based installer menu |
| `:DotfilesComponent <name>` | Install specific component (e.g., `python`, `go`) |
| `:DotfilesListComponents` | Show all components and their status |
| `:DotfilesUpdate` | Git pull dotfiles repo |
| `:DotfilesStatus` | Check installation status |
| `:DotfilesHelp` | Show installer help |

## Configuration

```lua
require("dotfiles").setup({
  -- Path to dotfiles repository (auto-detected if nil)
  dotfiles_root = nil,

  -- Installation options (maps to install.sh flags)
  install = {
    dotfiles = true,          -- Install shell bootstrap files
    app_configs = true,       -- Link application configs (git, ssh)
    package_manager = false,  -- Install package manager tools
    dev_tools = false,        -- Install development tools
    cloud_tools = false,      -- Install cloud/k8s tools
    skip_gui = true,          -- Skip GUI apps (recommended for remote)
  },

  -- Component-based installation
  components = {
    enabled = false,          -- Use component-based installation
    categories = {},          -- Categories: "core", "dev", "cloud", "ai", "database"
    list = {},                -- Specific components to install
  },

  -- Auto-setup on Neovim start
  auto_setup = {
    enabled = false,          -- Prompt to install if not detected
    check_installed = true,   -- Only prompt if not already installed
    silent = false,           -- Suppress notifications
    minimal = true,           -- Use minimal install by default
  },

  -- Notifications
  notify = {
    enabled = true,
    level = vim.log.levels.INFO,
  },
})
```

## Lua API

You can also call functions directly from Lua:

```lua
local dotfiles = require("dotfiles")

-- Check if installed
if not dotfiles.is_installed() then
  dotfiles.install_minimal()
end

-- Install specific component
dotfiles.installer.install_component("python")

-- Update dotfiles
dotfiles.update()
```

## Use Cases

### Setting up a new remote machine

1. SSH into the machine
2. Install Neovim (if not present): `brew install neovim` or `apt install neovim`
3. Clone dotfiles: `git clone <repo> ~/.config/dotfiles`
4. Open Neovim and run: `:DotfilesSetup`
5. Done! Your shell environment is ready.

### Adding tools incrementally

After initial setup, add components as needed:

```vim
:DotfilesComponent python
:DotfilesComponent go
:DotfilesComponent kubernetes
```

### Updating dotfiles

Keep your config in sync:

```vim
:DotfilesUpdate
```

## Integration with Existing Config

If you're using a Neovim config like kickstart.nvim, add the plugin
to your plugin list. The plugin auto-detects DOTFILES_ROOT from:

1. `$DOTFILES_ROOT` environment variable
2. `~/.config/dotfiles`
3. `~/.dotfiles`
4. `~/dotfiles`
5. `~/Repos/personal/bash-profile`
