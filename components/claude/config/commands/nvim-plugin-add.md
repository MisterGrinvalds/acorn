---
description: Add a new plugin to Neovim configuration
argument-hint: <plugin-name>
allowed-tools: Read, Write, Edit, Glob, Grep, Bash
---

## Task

Help the user add a new plugin to their Neovim configuration using lazy.nvim.

## Process

1. **Identify the plugin** from `$ARGUMENTS` or ask
2. **Find config location** - Check for lazy.nvim setup
3. **Create plugin spec** - Write the Lua configuration
4. **Add to config** - Insert in appropriate location
5. **Install** - Run `:Lazy sync`

## Plugin Spec Format (lazy.nvim)

### Basic
```lua
{ "author/plugin-name" }
```

### With Configuration
```lua
{
  "author/plugin-name",
  dependencies = { "nvim-lua/plenary.nvim" },
  event = "VeryLazy",  -- Lazy load
  config = function()
    require("plugin-name").setup({
      -- options
    })
  end,
}
```

### With Keys
```lua
{
  "author/plugin-name",
  keys = {
    { "<leader>p", "<cmd>PluginCommand<cr>", desc = "Run Plugin" },
  },
  config = function()
    require("plugin-name").setup()
  end,
}
```

## Common Plugin Patterns

### UI Plugin
```lua
{
  "author/ui-plugin",
  event = "VimEnter",
  config = function()
    require("ui-plugin").setup()
  end,
}
```

### LSP Enhancement
```lua
{
  "author/lsp-plugin",
  dependencies = { "neovim/nvim-lspconfig" },
  event = "LspAttach",
  config = function()
    require("lsp-plugin").setup()
  end,
}
```

### Filetype Specific
```lua
{
  "author/language-plugin",
  ft = { "python", "lua" },
  config = function()
    require("language-plugin").setup()
  end,
}
```

## Lazy Loading Events

| Event | When |
|-------|------|
| `VeryLazy` | After UI loads |
| `BufReadPre` | Before reading buffer |
| `InsertEnter` | Entering insert mode |
| `LspAttach` | LSP attaches to buffer |
| `VimEnter` | Neovim startup complete |

## Context

Look for existing plugin configuration:
@~/.config/nvim/lua/plugins

## After Adding

1. Save the file
2. Run `:Lazy` to open plugin manager
3. Press `I` to install or `S` to sync
4. Check `:Lazy health` if issues

## Tips

1. Check plugin README for recommended config
2. Use `lazy = true` if plugin has specific triggers
3. Add dependencies for plugins that require them
4. Group related plugins in the same spec file
