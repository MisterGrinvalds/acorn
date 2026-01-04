# Integrate Plugin

Set up keymaps, commands, and documentation for a plugin following consistent conventions.

## Arguments

$ARGUMENTS - Plugin name (e.g., `telescope`, `harpoon`, `gitsigns`)

## Mission

Ensure a plugin is properly integrated into the configuration with:
- Consistent keymap conventions that match the plugin's purpose
- Useful Neovim commands that expose functionality
- Integration with the help system for in-editor documentation
- Proper categorization for discoverability

## Keymap Convention System

### Leader Key Groups

| Prefix | Category | Purpose | Examples |
|--------|----------|---------|----------|
| `<leader>f` | Find | File/content finding | `ff` files, `fg` grep, `fb` buffers |
| `<leader>s` | Search | Search/explore features | `sf` files, `sg` grep, `sc` commands |
| `<leader>g` | Git | Git operations | `gs` status, `gc` commits, `gb` branches |
| `<leader>c` | Code | Code actions/LSP | `ca` actions, `cr` rename, `cd` definition |
| `<leader>d` | Debug | Debugging | `dc` continue, `db` breakpoint |
| `<leader>t` | Terminal | Terminal management | `tt` toggle, `tn` new |
| `<leader>b` | Buffer | Buffer management | `bd` delete, `bp` pick |
| `<leader>w` | Window | Window management | `wv` vsplit, `wh` hsplit |
| `<leader>x` | Diagnostics | Errors/warnings | `xx` list, `xw` workspace |
| `<leader>h` | Hunk/Harpoon | Quick marks/git hunks | `ha` add, `h1-4` jump |
| `<leader>l` | LSP | Language server | `lr` restart, `li` info |
| `<leader>n` | Notes/Neo | Note-taking, neo-tree | `nt` toggle tree |
| `<leader>p` | Palette/Project | Command palette, projects | `pp` projects |
| `<leader>q` | Quit/Session | Session management | `qs` save, `qr` restore |
| `<leader>u` | UI/Toggle | Toggle UI elements | `uw` wrap, `un` numbers |
| `<leader>r` | Run/Replace | Run commands, search-replace | `rr` run, `rs` replace |

### Keymap Best Practices

1. **Mnemonic Keys**: First letter matches action (`f`ind, `g`it, `c`ode)
2. **Consistent Depth**: Most actions are 2 keys after leader
3. **Related Grouping**: Related actions share prefix
4. **Descriptive desc**: Always include `desc` for which-key
5. **Mode Awareness**: Use appropriate modes (n, v, x, i)

### Keymap Template

```lua
-- Pattern: <leader>{category}{action}
vim.keymap.set('n', '<leader>XX', function()
  require('{plugin}').{action}()
end, { desc = '[X] {Description}' })

-- For visual mode actions
vim.keymap.set('v', '<leader>XX', function()
  require('{plugin}').{action}()
end, { desc = '[X] {Description}' })
```

## Command Convention System

### User Commands

Every plugin should expose commonly-used functions as commands:

```lua
-- Pattern: {Plugin}{Action}
vim.api.nvim_create_user_command('{Plugin}{Action}', function(opts)
  require('{plugin}').{action}(opts.args)
end, {
  desc = '{Description}',
  nargs = '*',  -- or '?', '+', '1' as needed
  complete = function() return {} end,  -- optional completion
})
```

### Standard Commands to Create

| Pattern | Purpose | Example |
|---------|---------|---------|
| `{Plugin}Toggle` | Toggle visibility/state | `TelescopeToggle` |
| `{Plugin}Open` | Open main interface | `TroubleOpen` |
| `{Plugin}Close` | Close interface | `TroubleClose` |
| `{Plugin}Refresh` | Refresh/reload | `GitsignsRefresh` |
| `{Plugin}Help` | Show plugin help | `HarpoonHelp` |

## Integration Workflow

### Step 1: Analyze Plugin

Read the plugin configuration:
```
lua/custom/plugins/{plugin}.lua
```

Identify:
- Current keymaps
- Exposed functions
- Category/purpose

### Step 2: Categorize

Assign plugin to appropriate category:

| Category | Plugins | Leader Group |
|----------|---------|--------------|
| Navigation | telescope, harpoon, oil | `<leader>f`, `<leader>h` |
| Git | gitsigns, diffview, lazygit | `<leader>g` |
| LSP/Code | trouble, lspconfig | `<leader>c`, `<leader>x` |
| Terminal | toggleterm, terminal-manager | `<leader>t` |
| UI | bufferline, lualine | `<leader>u` |
| Editing | comment, surround | `gc`, `ys` (vim-style) |
| Debug | nvim-dap | `<leader>d` |

### Step 3: Audit Keymaps

Check for:
1. **Conflicts**: Does this keymap override something?
2. **Consistency**: Does it follow conventions?
3. **Discoverability**: Is `desc` set for which-key?
4. **Completeness**: Are all main functions mapped?

### Step 4: Create Help Integration

Register with the help system:

```lua
-- Add to lua/lib/plugin-registry.lua
registry.register({
  name = '{plugin}',
  category = '{category}',
  description = '{one-line description}',
  keymaps = {
    { key = '<leader>xx', desc = 'Action description', mode = 'n' },
    -- ...
  },
  commands = {
    { cmd = '{Plugin}Action', desc = 'What it does' },
    -- ...
  },
  help_topics = {
    'Basic Usage',
    'Advanced Features',
    'Configuration',
  },
})
```

### Step 5: Generate Documentation

Create quick-reference card for the plugin:

```markdown
## {Plugin Name}

**Category:** {category}
**Purpose:** {one-line description}

### Keymaps
| Key | Action |
|-----|--------|
| `<leader>xx` | Action |

### Commands
| Command | Description |
|---------|-------------|
| `:PluginAction` | What it does |

### Quick Tips
- Tip 1
- Tip 2
```

## Output

After integration, provide:

### 1. Keymap Summary

```
Keymaps for {plugin}:
  <leader>xx → {action} (category: {cat})
  <leader>xy → {action} (category: {cat})
```

### 2. Commands Created

```
Commands:
  :{Plugin}Action  - {description}
  :{Plugin}Help    - Show help for {plugin}
```

### 3. Conflicts Detected

```
Potential conflicts:
  <leader>xx was previously mapped to {other}
```

### 4. Suggestions

```
Recommended additions:
  - Add <leader>xx for {commonly-used-action}
  - Consider remapping {key} for consistency
```

## Example

**Input:** `/integrate-plugin harpoon`

**Analysis:**
- Category: Navigation (quick file marks)
- Current keymaps: `<leader>ha`, `<leader>1-4`
- Missing: Help command, menu toggle

**Keymap Audit:**
```
✓ <leader>ha (add) - follows convention
✓ <leader>1-4 (jump) - intuitive
✗ No toggle for menu list
✗ No help command
```

**Additions:**
```lua
-- Toggle menu
vim.keymap.set('n', '<leader>hm', function()
  require('harpoon.ui').toggle_quick_menu()
end, { desc = '[H]arpoon [M]enu' })

-- Help command
vim.api.nvim_create_user_command('HarpoonHelp', function()
  require('lib.plugin-help').show('harpoon')
end, { desc = 'Show Harpoon help' })
```

**Registration:**
```lua
-- Added to plugin-registry
{
  name = 'harpoon',
  category = 'navigation',
  description = 'Quick file marks for fast navigation',
  keymaps = {
    { key = '<leader>ha', desc = 'Add file to harpoon', mode = 'n' },
    { key = '<leader>hm', desc = 'Toggle harpoon menu', mode = 'n' },
    { key = '<leader>1-4', desc = 'Jump to mark 1-4', mode = 'n' },
  },
}
```

## Categories Reference

| Category | Description | Typical Plugins |
|----------|-------------|-----------------|
| navigation | File finding, jumping, marks | telescope, harpoon, oil |
| git | Version control | gitsigns, diffview, lazygit |
| lsp | Language intelligence | lspconfig, trouble, none-ls |
| completion | Autocompletion | blink.cmp, nvim-cmp |
| editing | Text manipulation | comment, surround, autopairs |
| terminal | Terminal management | toggleterm, terminal-manager |
| ui | Visual enhancements | bufferline, lualine, noice |
| treesitter | Syntax/text objects | treesitter, ts-context |
| debug | Debugging | nvim-dap, neotest |
| formatting | Code formatting | conform, null-ls |
| utility | General utilities | which-key, nvim-notify |
