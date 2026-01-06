---
description: Set up LSP (Language Server Protocol) for a language
argument-hint: <language>
allowed-tools: Read, Write, Edit, Glob, Grep, Bash
---

## Task

Help the user set up LSP support for a specific programming language in Neovim.

## Process

1. **Identify language** from `$ARGUMENTS` or ask
2. **Check Mason** - See what's available
3. **Install server** - Via Mason or manual
4. **Configure lspconfig** - Add server configuration
5. **Verify** - Test with a file of that language

## Common Language Servers

| Language | Server | Mason Name |
|----------|--------|------------|
| Python | pyright | `pyright` |
| TypeScript/JS | typescript-language-server | `typescript-language-server` |
| Go | gopls | `gopls` |
| Rust | rust-analyzer | `rust-analyzer` |
| Lua | lua-language-server | `lua-language-server` |
| C/C++ | clangd | `clangd` |
| JSON | json-lsp | `json-lsp` |
| YAML | yaml-language-server | `yaml-language-server` |
| Bash | bash-language-server | `bash-language-server` |
| HTML/CSS | cssls, html | `cssls`, `html-lsp` |

## Mason Installation

```vim
:Mason
" Navigate to server, press 'i' to install
```

Or in config:
```lua
require("mason-lspconfig").setup({
  ensure_installed = { "pyright", "gopls", "lua_ls" },
})
```

## LSP Configuration

### Basic Setup
```lua
local lspconfig = require("lspconfig")

lspconfig.pyright.setup({})
lspconfig.gopls.setup({})
lspconfig.lua_ls.setup({
  settings = {
    Lua = {
      diagnostics = { globals = { "vim" } },
    },
  },
})
```

### With Capabilities (for completion)
```lua
local capabilities = require("cmp_nvim_lsp").default_capabilities()

lspconfig.pyright.setup({
  capabilities = capabilities,
})
```

### With on_attach (for keymaps)
```lua
local on_attach = function(client, bufnr)
  local opts = { buffer = bufnr }
  vim.keymap.set("n", "gd", vim.lsp.buf.definition, opts)
  vim.keymap.set("n", "K", vim.lsp.buf.hover, opts)
  vim.keymap.set("n", "<leader>rn", vim.lsp.buf.rename, opts)
  vim.keymap.set("n", "<leader>ca", vim.lsp.buf.code_action, opts)
end

lspconfig.pyright.setup({
  on_attach = on_attach,
  capabilities = capabilities,
})
```

## Common LSP Keymaps

| Key | Action |
|-----|--------|
| `gd` | Go to definition |
| `gr` | Find references |
| `K` | Hover documentation |
| `<leader>rn` | Rename symbol |
| `<leader>ca` | Code actions |
| `<leader>f` | Format buffer |
| `[d` / `]d` | Previous/next diagnostic |

## Verification

1. Open a file of the language
2. Run `:LspInfo` - Should show attached server
3. Try `gd` on a function call
4. Check `:LspLog` if issues

## Troubleshooting

1. **Server not starting**: Check `:Mason` installed correctly
2. **No diagnostics**: Check server configuration
3. **Slow**: Try limiting workspace scope
4. **Missing features**: Check server capabilities

## Context

@~/.config/nvim/lua/plugins/lsp.lua
