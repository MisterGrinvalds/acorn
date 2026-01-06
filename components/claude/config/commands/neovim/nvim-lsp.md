# LSP Management

Comprehensive Language Server Protocol management for your Neovim configuration.

## Arguments

$ARGUMENTS - Action or language (e.g., `add python`, `diagnose`, `status`, `optimize lua_ls`)

## Mission

Provide complete LSP management including:
- Adding and configuring new language servers
- Diagnosing and troubleshooting LSP issues
- Optimizing server settings for performance
- Managing Mason packages and tools
- Checking LSP health and connectivity
- Researching best configurations for languages

## Quick Actions Reference

| Command | Action |
|---------|--------|
| `/lsp status` | Show all active servers and their status |
| `/lsp add {language}` | Add a new language server |
| `/lsp diagnose` | Troubleshoot current LSP issues |
| `/lsp diagnose {server}` | Diagnose specific server |
| `/lsp optimize {server}` | Optimize server configuration |
| `/lsp mason` | Review Mason packages |
| `/lsp health` | Run comprehensive health check |
| `/lsp research {language}` | Research best LSP setup for language |
| `/lsp remove {server}` | Remove a language server |

## Workflow

### Mode Detection

Parse $ARGUMENTS to determine the action:

1. **No arguments or `status`** - Show LSP overview
2. **`add {language}`** - Add new server workflow
3. **`diagnose [server]`** - Troubleshooting workflow
4. **`optimize {server}`** - Optimization workflow
5. **`mason`** - Mason management workflow
6. **`health`** - Health check workflow
7. **`research {language}`** - Research workflow
8. **`remove {server}`** - Removal workflow

---

## Workflow: Status Overview

### Step 1: Read Current Configuration

```
Read lua/plugins/lsp.lua
```

Extract:
- Configured servers from `servers = { ... }`
- Mason ensure_installed list
- Custom server settings

### Step 2: Check Active Servers

Provide commands for user to run:
```vim
:LspInfo          " Shows attached clients
:Mason            " Opens Mason UI
:checkhealth lsp  " LSP health check
```

### Step 3: Present Status Report

```
## LSP Status Report

### Configured Servers
| Server | Language | Settings | Status |
|--------|----------|----------|--------|
| lua_ls | Lua | callSnippet=Replace | Configured |
| pyright | Python | (defaults) | Configured |
| ts_ls | TypeScript | (defaults) | Configured |
| gopls | Go | staticcheck, gofumpt | Configured |
| bashls | Bash | globPattern | Configured |
| yamlls | YAML | SchemaStore | Configured |
| jsonls | JSON | SchemaStore | Configured |

### Mason Tools
| Tool | Type | Purpose |
|------|------|---------|
| stylua | Formatter | Lua formatting |
| black | Formatter | Python formatting |
| prettier | Formatter | JS/TS/JSON/HTML/CSS/YAML |
| goimports | Formatter | Go imports |
| isort | Formatter | Python import sorting |
| shfmt | Formatter | Shell script formatting |

### Quick Commands
- `:LspInfo` - View attached clients
- `:LspRestart` - Restart LSP for current buffer
- `:Mason` - Open Mason package manager
```

---

## Workflow: Add Language Server

### Step 1: Identify the Language Server

Use the reference table below to map language to recommended server:

| Language | Server | Mason Name | Dependencies |
|----------|--------|------------|--------------|
| Python | pyright | pyright | python3 |
| Python (alt) | pylsp | python-lsp-server | python3, pip |
| TypeScript | ts_ls | typescript-language-server | node, npm |
| JavaScript | ts_ls | typescript-language-server | node, npm |
| Go | gopls | gopls | go |
| Rust | rust_analyzer | rust-analyzer | rustup |
| Lua | lua_ls | lua-language-server | - |
| C/C++ | clangd | clangd | llvm |
| Java | jdtls | jdtls | java |
| Ruby | solargraph | solargraph | ruby, gem |
| PHP | intelephense | intelephense | node, npm |
| HTML | html | html-lsp | node |
| CSS | cssls | css-lsp | node |
| JSON | jsonls | json-lsp | node |
| YAML | yamlls | yaml-language-server | node |
| Bash | bashls | bash-language-server | node |
| Docker | dockerls | dockerfile-language-server | node |
| Terraform | terraformls | terraform-ls | terraform |
| SQL | sqlls | sqlls | node |
| GraphQL | graphql | graphql-language-service-cli | node |
| Markdown | marksman | marksman | - |
| Tailwind | tailwindcss | tailwindcss-language-server | node |
| Vue | volar | vue-language-server | node |
| Svelte | svelte | svelte-language-server | node |
| Astro | astro | astro-language-server | node |
| Elixir | elixirls | elixir-ls | elixir |
| Haskell | hls | haskell-language-server | ghcup |
| Zig | zls | zls | zig |
| Kotlin | kotlin_language_server | kotlin-language-server | java |
| Scala | metals | metals | java, scala |
| Clojure | clojure_lsp | clojure-lsp | java |
| Nix | nil_ls | nil | nix |
| Prisma | prismals | prisma-language-server | node |

### Step 2: Research Server Configuration

```
WebSearch "{server_name} neovim lspconfig configuration 2025"
WebFetch https://github.com/neovim/nvim-lspconfig/blob/master/doc/configs.md
```

Check for:
- Recommended settings
- Required capabilities
- Common options to enable
- Performance considerations

### Step 3: Check for Formatters/Linters

| Language | Formatters | Linters |
|----------|-----------|---------|
| Python | black, ruff, autopep8 | ruff, pylint, mypy |
| JavaScript/TypeScript | prettier, biome | eslint, biome |
| Go | gofmt, goimports, gofumpt | golangci-lint |
| Rust | rustfmt | clippy |
| Lua | stylua | luacheck, selene |
| Ruby | rubocop | rubocop |
| PHP | php-cs-fixer | phpstan |
| C/C++ | clang-format | clang-tidy |
| Shell | shfmt | shellcheck |
| JSON | prettier | jsonlint |
| YAML | prettier | yamllint |
| Markdown | prettier | markdownlint |
| SQL | sql-formatter | sqlfluff |
| HTML | prettier | htmlhint |
| CSS | prettier, stylelint | stylelint |

### Step 4: Update lsp.lua

Add to the servers table in `lua/plugins/lsp.lua`:

```lua
-- In the servers = { ... } table:
{server_name} = {
  settings = {
    [{ServerCapitalName}] = {
      -- configuration options
    },
  },
},
```

### Step 5: Update Mason Ensure Installed

Add formatter/linter to ensure_installed:

```lua
vim.list_extend(ensure_installed, {
  '{formatter}',
})
```

### Step 6: Update Conform Formatters

If adding formatters, update `lua/plugins/formatting.lua`:

```lua
formatters_by_ft = {
  {filetype} = { '{formatter}' },
},
```

### Step 7: Provide Summary

```
## Added: {Language} LSP Support

### Server Configuration
- **Server:** {server_name}
- **Mason package:** {mason_name}
- **Settings:** {key settings}

### Tools Added
| Tool | Type | Purpose |
|------|------|---------|
| {formatter} | Formatter | {description} |

### Files Modified
1. `lua/plugins/lsp.lua` - Added server configuration
2. `lua/plugins/formatting.lua` - Added formatter (if applicable)

### Next Steps
1. Run `:Lazy sync` to install plugins
2. Run `:Mason` and install any pending tools
3. Open a {language} file to verify LSP attaches
4. Check `:LspInfo` to confirm connection
```

---

## Workflow: Diagnose LSP Issues

### Step 1: Gather Diagnostic Information

Provide user commands to run:

```vim
:LspInfo              " Check attached clients
:LspLog               " View LSP logs
:checkhealth lsp      " Health check
:checkhealth mason    " Mason health
:messages             " Check for errors
```

### Step 2: Common Issues Checklist

```
## LSP Diagnostic Checklist

### Server Attachment
- [ ] Is the server configured in `lua/plugins/lsp.lua`?
- [ ] Is the server installed via Mason? (`:Mason`)
- [ ] Does `:LspInfo` show the server attached?
- [ ] Is the filetype correct? (`:set ft?`)

### Mason Installation
- [ ] Is Mason installed and loaded? (`:Mason`)
- [ ] Is the language server installed? (check Mason UI)
- [ ] Are dependencies met? (node, python, etc.)

### File Issues
- [ ] Is the file in a project root? (has .git, package.json, etc.)
- [ ] Does the project have required config files? (tsconfig.json, pyproject.toml)
- [ ] Are there syntax errors preventing LSP attachment?

### Server-Specific
- [ ] Check `:LspLog` for server errors
- [ ] Is the server executable in PATH?
- [ ] Are language-specific requirements met?
```

### Step 3: Read Configuration

```
Read lua/plugins/lsp.lua
```

Verify:
- Server is in `servers` table
- Mason-lspconfig handler is correct
- No syntax errors in config

### Step 4: Provide Fixes

Based on findings, suggest specific fixes:

**Server Not Attaching:**
```lua
-- Add to servers table in lua/plugins/lsp.lua
{server} = {},
```

**Mason Package Missing:**
```lua
-- Add to ensure_installed
vim.list_extend(ensure_installed, { '{package}' })
```

**Root Detection Issue:**
```lua
-- Add custom root_dir function
{server} = {
  root_dir = function(fname)
    return require('lspconfig.util').find_git_ancestor(fname)
      or vim.fn.getcwd()
  end,
},
```

**Filetype Not Recognized:**
```lua
-- Add filetype detection
vim.filetype.add({
  extension = { {ext} = '{filetype}' },
  filename = { ['{filename}'] = '{filetype}' },
})
```

### Step 5: Verification Commands

```
## Verification Steps

1. Save and source config or restart Neovim
2. Open a {language} file
3. Run `:LspInfo` - should show server attached
4. Test functionality:
   - `K` - Hover documentation
   - `grd` - Go to definition
   - `grr` - Find references
   - `gra` - Code actions
```

---

## Workflow: Optimize Server

### Step 1: Analyze Current Configuration

```
Read lua/plugins/lsp.lua
```

Find current settings for the server.

### Step 2: Research Optimal Settings

```
WebSearch "{server} performance optimization neovim 2025"
```

### Step 3: Server-Specific Optimizations

**lua_ls Optimizations:**
```lua
lua_ls = {
  settings = {
    Lua = {
      completion = { callSnippet = 'Replace' },
      diagnostics = {
        disable = { 'missing-fields' },
      },
      workspace = {
        checkThirdParty = false,
        maxPreload = 1000,
        preloadFileSize = 150,
      },
      telemetry = { enable = false },
    },
  },
},
```

**pyright Optimizations:**
```lua
pyright = {
  settings = {
    python = {
      analysis = {
        autoSearchPaths = true,
        useLibraryCodeForTypes = true,
        diagnosticMode = 'openFilesOnly',
        typeCheckingMode = 'basic',
      },
    },
  },
},
```

**ts_ls Optimizations:**
```lua
ts_ls = {
  settings = {
    typescript = {
      inlayHints = {
        includeInlayParameterNameHints = 'all',
        includeInlayFunctionParameterTypeHints = true,
        includeInlayVariableTypeHints = true,
      },
    },
  },
  on_attach = function(client)
    -- Disable semantic tokens if causing slowness
    client.server_capabilities.semanticTokensProvider = nil
  end,
},
```

**gopls Optimizations:**
```lua
gopls = {
  settings = {
    gopls = {
      analyses = {
        unusedparams = true,
        shadow = true,
        nilness = true,
        unusedwrite = true,
      },
      staticcheck = true,
      gofumpt = true,
      usePlaceholders = true,
      completeUnimported = true,
      hints = {
        assignVariableTypes = true,
        compositeLiteralFields = true,
        constantValues = true,
        functionTypeParameters = true,
        parameterNames = true,
        rangeVariableTypes = true,
      },
    },
  },
},
```

**rust_analyzer Optimizations:**
```lua
rust_analyzer = {
  settings = {
    ['rust-analyzer'] = {
      checkOnSave = { command = 'clippy' },
      cargo = { allFeatures = true },
      procMacro = { enable = true },
      diagnostics = {
        enable = true,
        experimental = { enable = true },
      },
      inlayHints = {
        enable = true,
        chainingHints = { enable = true },
        typeHints = { enable = true },
      },
    },
  },
},
```

### Step 4: Apply and Verify

Update `lua/plugins/lsp.lua` with recommended settings.

```
## Performance Verification

1. Restart Neovim
2. Open a large project file
3. Monitor with `:LspInfo` for status
4. Test response times:
   - Hover (`K`) should be < 200ms
   - Completion should be responsive
   - Diagnostics should update in < 500ms

### Commands
:lua print(vim.inspect(vim.lsp.get_clients()))
:LspRestart
```

---

## Workflow: Mason Management

### Step 1: List Current Packages

```
Read lua/plugins/lsp.lua
```

Extract `ensure_installed` list.

### Step 2: Audit Package Usage

Cross-reference with:
- Servers in `servers` table
- Formatters in `lua/plugins/formatting.lua`

### Step 3: Present Mason Overview

```
## Mason Package Audit

### Installed via ensure_installed
| Package | Type | Used By | Status |
|---------|------|---------|--------|
| stylua | Formatter | lua (conform) | Active |
| black | Formatter | python (conform) | Active |
| prettier | Formatter | js/ts/json/yaml (conform) | Active |
| goimports | Formatter | go (conform) | Active |
| isort | Formatter | python (conform) | Active |
| shfmt | Formatter | sh/bash (conform) | Active |

### Language Servers (via servers table)
| Server | Mason Package | Status |
|--------|--------------|--------|
| lua_ls | lua-language-server | Active |
| pyright | pyright | Active |
| ts_ls | typescript-language-server | Active |
| gopls | gopls | Active |
| bashls | bash-language-server | Active |
| yamlls | yaml-language-server | Active |
| jsonls | json-lsp | Active |
```

### Step 4: Mason Commands Reference

```
## Mason Commands

| Command | Action |
|---------|--------|
| `:Mason` | Open Mason UI |
| `:MasonInstall {pkg}` | Install package |
| `:MasonUninstall {pkg}` | Remove package |
| `:MasonUpdate` | Update all packages |
| `:MasonLog` | View Mason logs |

### Mason UI Navigation
- `i` - Install package
- `X` - Uninstall package
- `u` - Update package
- `g?` - Show help
- `<CR>` - Expand package info
```

---

## Workflow: Health Check

### Step 1: Run Comprehensive Checks

Provide all health check commands:

```vim
:checkhealth lsp
:checkhealth mason
:checkhealth vim.lsp
```

### Step 2: Analyze Configuration

Read and validate:
- `lua/plugins/lsp.lua`
- `lua/plugins/formatting.lua`

### Step 3: Present Health Report

```
## LSP Health Report

### Core Status
| Component | Status | Notes |
|-----------|--------|-------|
| nvim-lspconfig | OK | Loaded |
| mason.nvim | OK | Loaded |
| mason-lspconfig | OK | Loaded |
| mason-tool-installer | OK | Loaded |
| blink.cmp | OK | Providing capabilities |
| schemastore.nvim | OK | Providing schemas |

### Server Health
| Server | Configured | Installed | Binary |
|--------|------------|-----------|--------|
| lua_ls | Yes | Yes | ~/.local/share/nvim/mason/bin/lua-language-server |
| pyright | Yes | Yes | ~/.local/share/nvim/mason/bin/pyright |
| ts_ls | Yes | Yes | ~/.local/share/nvim/mason/bin/typescript-language-server |
| gopls | Yes | Yes | ~/.local/share/nvim/mason/bin/gopls |
| bashls | Yes | Yes | ~/.local/share/nvim/mason/bin/bash-language-server |
| yamlls | Yes | Yes | ~/.local/share/nvim/mason/bin/yaml-language-server |
| jsonls | Yes | Yes | ~/.local/share/nvim/mason/bin/vscode-json-language-server |

### Formatter Health
| Formatter | Configured | Installed | Works |
|-----------|------------|-----------|-------|
| stylua | Yes | Yes | Yes |
| black | Yes | Yes | Yes |
| prettier | Yes | Yes | Yes |
| goimports | Yes | Yes | Yes |
| shfmt | Yes | Yes | Yes |
```

---

## Workflow: Research Language

### Step 1: Research Current Best Practices

```
WebSearch "{language} neovim LSP setup best practices 2025"
WebSearch "{language} neovim formatter linter recommendations"
```

### Step 2: Present Research Report

```
## {Language} LSP Research

### Recommended Server
**Primary:** {server_name}
- Mason package: {mason_name}
- Maturity: Stable/Beta
- Community adoption: High/Medium

**Alternatives:**
- {alt_server}: {pros/cons}

### Recommended Tools
| Type | Tool | Why |
|------|------|-----|
| Formatter | {name} | {reason} |
| Linter | {name} | {reason} |

### Optimal Configuration
```lua
{server} = {
  settings = {
    -- Recommended settings
  },
},
```

### Formatter Configuration
```lua
formatters_by_ft = {
  {filetype} = { '{formatter}' },
},
```

### Resources
- [nvim-lspconfig docs]({url})
- [Server Repository]({url})
```

---

## Workflow: Remove Server

### Step 1: Identify Removal Targets

```
Read lua/plugins/lsp.lua
Read lua/plugins/formatting.lua
```

### Step 2: List Dependencies

Show what will be affected:
- Server entry in `servers` table
- Mason packages
- Formatter configuration

### Step 3: Perform Removal

Remove from `lua/plugins/lsp.lua`:
```lua
-- Remove this entry from servers = { ... }
{server} = { ... },
```

Remove from `ensure_installed` if no longer needed:
```lua
'{formatter}',
```

Remove from `lua/plugins/formatting.lua`:
```lua
{filetype} = { '{formatter}' },
```

### Step 4: Cleanup Mason

```
## Cleanup Steps

1. Remove from config (done above)
2. Run `:Lazy sync`
3. Open `:Mason`
4. Uninstall unused packages
```

---

## Error Recovery

### Common Errors and Fixes

| Error | Cause | Fix |
|-------|-------|-----|
| "Client X quit with exit code 1" | Server crashed | Check `:LspLog`, reinstall server |
| "No matching language server" | Not configured | Add to `servers` table |
| "Spawning server failed" | Binary not found | Install via Mason or system |
| "Client X is not running" | Not started | Check filetype, root directory |
| "Capability not supported" | Server limitation | Use different server |

### Recovery Commands

```vim
:LspRestart          " Restart all LSP clients
:LspStop             " Stop all clients
:LspStart            " Start clients for buffer
:e                   " Reload buffer (triggers LSP)
:MasonInstall {pkg}  " Reinstall package
```

---

## LSP Keymaps Reference

Current keymaps defined in `lua/plugins/lsp.lua`:

| Keymap | Mode | Action | Description |
|--------|------|--------|-------------|
| `grn` | n | `vim.lsp.buf.rename` | Rename symbol |
| `gra` | n, x | `vim.lsp.buf.code_action` | Code actions |
| `grr` | n | Telescope lsp_references | Find references |
| `gri` | n | Telescope lsp_implementations | Find implementations |
| `grd` | n | Telescope lsp_definitions | Go to definition |
| `grD` | n | `vim.lsp.buf.declaration` | Go to declaration |
| `gO` | n | Telescope lsp_document_symbols | Document symbols |
| `gW` | n | Telescope lsp_dynamic_workspace_symbols | Workspace symbols |
| `grt` | n | Telescope lsp_type_definitions | Type definition |
| `<leader>th` | n | Toggle inlay hints | Toggle inlay hints |
| `<leader>lf` | n | Conform format | Format buffer |

---

## Examples

### Example 1: Add Rust Support
**Input:** `/lsp add rust`

**Output:**
```
## Adding Rust LSP Support

### Configuration Added

**Server:** rust_analyzer
**Mason package:** rust-analyzer

Added to lua/plugins/lsp.lua:
```lua
rust_analyzer = {
  settings = {
    ['rust-analyzer'] = {
      checkOnSave = { command = 'clippy' },
      cargo = { allFeatures = true },
      procMacro = { enable = true },
    },
  },
},
```

**Formatters:** rustfmt (via rust-analyzer)

**Next Steps:**
1. `:Lazy sync`
2. `:Mason` - Install rust-analyzer
3. Open a .rs file
4. `:LspInfo` to verify
```

### Example 2: Diagnose TypeScript Issues
**Input:** `/lsp diagnose ts_ls`

**Output:**
```
## Diagnosing ts_ls

### Run These Commands
:LspInfo
:LspLog
:set ft?

### Common ts_ls Issues

1. **No tsconfig.json**
   - ts_ls requires tsconfig.json or jsconfig.json in project root
   - Create: `npx tsc --init`

2. **Wrong Node Version**
   - Needs Node 14+
   - Check: `node --version`

3. **Large Project Slowness**
   - Add: `diagnosticMode = 'openFilesOnly'`

### Fix
If not attaching, ensure project has:
- package.json or tsconfig.json at root
- Valid TypeScript/JavaScript files
```

### Example 3: Check Status
**Input:** `/lsp status`

**Output:**
```
## LSP Status Overview

### Configured Servers (7)
| Server | Language | Custom Settings |
|--------|----------|-----------------|
| lua_ls | Lua | callSnippet=Replace |
| pyright | Python | (defaults) |
| ts_ls | TypeScript | (defaults) |
| gopls | Go | staticcheck, gofumpt |
| bashls | Bash | globPattern |
| yamlls | YAML | SchemaStore |
| jsonls | JSON | SchemaStore |

### Mason Tools (6 formatters)
stylua, black, isort, prettier, goimports, shfmt

### Quick Commands
:LspInfo - View active clients
:Mason - Package manager
:LspRestart - Restart servers
```

---

## Notes

- Always backup `lua/plugins/lsp.lua` before major changes
- Use `:Lazy sync` after configuration changes
- Mason packages stored in `~/.local/share/nvim/mason`
- LSP logs at `:LspLog` or `~/.local/state/nvim/lsp.log`
- For debugging: `vim.lsp.set_log_level('debug')`
