---
description: Configure Neovim keymaps and shortcuts
argument-hint: [action: list|add|which-key]
allowed-tools: Read, Write, Edit, Glob, Grep, Bash
---

## Task

Help the user configure and manage Neovim keymaps.

## Actions

Based on `$ARGUMENTS`:

### list
Show how to view current keymaps:
```vim
:map              " All mappings
:nmap             " Normal mode
:imap             " Insert mode
:vmap             " Visual mode
:Telescope keymaps  " Search all (with Telescope)
```

### add
Help create a new keymap with proper format.

### which-key
Configure which-key for keymap discovery.

## Keymap API

### Modern Lua API (Recommended)
```lua
vim.keymap.set(mode, lhs, rhs, opts)
```

Parameters:
- `mode`: "n", "i", "v", "x", "t", or table {"n", "v"}
- `lhs`: Left-hand side (the key)
- `rhs`: Right-hand side (the action)
- `opts`: Options table

### Examples
```lua
-- Normal mode mapping
vim.keymap.set("n", "<leader>w", ":w<CR>", { desc = "Save file" })

-- Multiple modes
vim.keymap.set({"n", "v"}, "<leader>y", '"+y', { desc = "Copy to clipboard" })

-- Lua function as action
vim.keymap.set("n", "<leader>ff", function()
  require("telescope.builtin").find_files()
end, { desc = "Find files" })

-- Buffer-local
vim.keymap.set("n", "gd", vim.lsp.buf.definition, { buffer = bufnr })
```

### Options
```lua
{
  desc = "Description",     -- Shows in which-key
  silent = true,            -- Don't show command
  noremap = true,           -- Default true in Lua API
  expr = false,             -- Expression mapping
  buffer = bufnr,           -- Buffer-local
}
```

## Common Keymap Patterns

### Leader Mappings
```lua
vim.g.mapleader = " "  -- Space as leader

-- File operations
vim.keymap.set("n", "<leader>w", ":w<CR>", { desc = "Save" })
vim.keymap.set("n", "<leader>q", ":q<CR>", { desc = "Quit" })

-- Buffer navigation
vim.keymap.set("n", "<leader>bn", ":bnext<CR>", { desc = "Next buffer" })
vim.keymap.set("n", "<leader>bp", ":bprev<CR>", { desc = "Prev buffer" })

-- Window navigation
vim.keymap.set("n", "<C-h>", "<C-w>h", { desc = "Window left" })
vim.keymap.set("n", "<C-j>", "<C-w>j", { desc = "Window down" })
vim.keymap.set("n", "<C-k>", "<C-w>k", { desc = "Window up" })
vim.keymap.set("n", "<C-l>", "<C-w>l", { desc = "Window right" })
```

### Plugin Keymaps
```lua
-- Telescope
vim.keymap.set("n", "<leader>ff", "<cmd>Telescope find_files<cr>")
vim.keymap.set("n", "<leader>fg", "<cmd>Telescope live_grep<cr>")
vim.keymap.set("n", "<leader>fb", "<cmd>Telescope buffers<cr>")

-- File tree
vim.keymap.set("n", "<leader>e", "<cmd>NvimTreeToggle<cr>")

-- Git
vim.keymap.set("n", "<leader>gs", "<cmd>Git<cr>")
vim.keymap.set("n", "<leader>gb", "<cmd>Git blame<cr>")
```

## Which-Key Integration

```lua
{
  "folke/which-key.nvim",
  event = "VeryLazy",
  config = function()
    local wk = require("which-key")
    wk.setup()

    -- Register groups
    wk.register({
      ["<leader>f"] = { name = "+find" },
      ["<leader>g"] = { name = "+git" },
      ["<leader>l"] = { name = "+lsp" },
      ["<leader>b"] = { name = "+buffer" },
    })
  end,
}
```

## Recommended Leader Structure

| Prefix | Category |
|--------|----------|
| `<leader>f` | Find (Telescope) |
| `<leader>g` | Git |
| `<leader>l` | LSP |
| `<leader>b` | Buffers |
| `<leader>w` | Windows |
| `<leader>t` | Toggle |
| `<leader>x` | Trouble/Diagnostics |

## Context

@~/.config/nvim/lua/config/keymaps.lua

## Tips

1. Always add `desc` for which-key visibility
2. Use `<leader>` for user commands
3. Keep related keymaps grouped
4. Check for conflicts with `:map <key>`
