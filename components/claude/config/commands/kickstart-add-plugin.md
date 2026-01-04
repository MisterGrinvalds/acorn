# Add Plugin

Add a new Neovim plugin to the configuration with fork integration.

## Arguments

$ARGUMENTS - Plugin in `owner/repo` format (e.g., `dmtrKovalenko/fff.nvim`)

## Workflow

### Step 1: Research the Plugin

Fetch the GitHub README to understand:
- What the plugin does
- Dependencies required
- Default configuration
- Key features and keymaps
- Installation instructions

```
WebFetch https://github.com/$ARGUMENTS
```

### Step 2: Check for Redundancy

Scan existing plugins for overlapping functionality:

1. Read current plugin inventory:
   - `init.lua` - core plugins
   - `lua/kickstart/plugins/*.lua` - optional plugins
   - `lua/custom/plugins/*.lua` - custom plugins

2. Identify potential overlaps:
   - File pickers: telescope.nvim, fff.nvim, fzf-lua
   - File navigation: harpoon, arrow.nvim
   - Git: gitsigns, diffview, lazygit, fugitive
   - LSP: nvim-lspconfig, none-ls, trouble
   - Completion: blink.cmp, nvim-cmp, coq
   - Comments: Comment.nvim, ts-comments
   - etc.

3. **Ask the user** if overlap is found:
   - "This plugin overlaps with {existing}. Do you want to:"
     - A) Replace {existing} with this plugin
     - B) Use both (complementary features)
     - C) Cancel

### Step 3: Fork the Plugin

Use the fork workflow:

```bash
# Check if fork exists
gh repo view MisterGrinvalds/{owner}.{repo} --json name 2>/dev/null

# If not, create and rename
gh repo fork $ARGUMENTS --clone=false
sleep 2
gh repo rename {owner}.{repo} --repo MisterGrinvalds/{repo} --yes
```

### Step 4: Create Plugin Config

Create `lua/custom/plugins/{plugin-name}.lua`:

```lua
-- {plugin-name} - {brief description}
-- {link to original repo}
local fork = require('lib.forks').fork

return {
  fork('$ARGUMENTS'),
  -- Add dependencies (also forked if external)
  dependencies = {
    fork('dep/repo'),
  },
  -- Lazy loading strategy
  event = 'VeryLazy', -- or cmd, keys, ft as appropriate
  -- Configuration
  config = function()
    require('{plugin}').setup({
      -- baseline config from README
    })

    -- Keymaps following our convention
    vim.keymap.set('n', '<leader>xx', function()
      -- action
    end, { desc = '[X] Description' })
  end,
}
```

### Step 5: Configuration Guidance

Based on the README, provide:

1. **Recommended keymaps** following our convention:
   - `<leader>` prefix for most actions
   - Descriptive `desc` for which-key
   - Avoid conflicts with existing bindings

2. **Lazy loading strategy**:
   - `cmd = {...}` for command-triggered plugins
   - `keys = {...}` for keymap-triggered
   - `event = 'VeryLazy'` for background plugins
   - `ft = {...}` for filetype-specific

3. **Integration points**:
   - Does it integrate with telescope?
   - Does it need LSP?
   - Does it conflict with other plugins?

### Step 6: Verify and Report

1. Confirm fork exists: `gh repo view MisterGrinvalds/{owner}.{repo}`
2. Show the created config file
3. List any keymaps added
4. Note any manual steps (`:Lazy sync`, etc.)

## Example

**Input:** `/add-plugin dmtrKovalenko/fff.nvim`

**Research Output:**
- fff.nvim is a smart fuzzy file picker with Rust backend
- Requires Neovim 0.10+, Rustup
- Overlaps with telescope.nvim for file finding

**User Choice:** Use both (fff for files, telescope for everything else)

**Created:** `lua/custom/plugins/fff.lua`
```lua
local fork = require('lib.forks').fork

return {
  fork('dmtrKovalenko/fff.nvim'),
  build = function()
    require("fff.download").download_or_build_binary()
  end,
  lazy = false,
  keys = {
    { '<leader>ff', function() require('fff').find_files() end, desc = '[F]ind [F]iles (fff)' },
  },
}
```

**Keymaps:**
- `<leader>ff` - Find files with fff (fast fuzzy finder)

**Note:** Run `:Lazy sync` to install, then restart Neovim.

## Category Reference

When checking for redundancy, consider these categories:

| Category | Common Plugins |
|----------|---------------|
| File Picker | telescope, fff, fzf-lua, snap |
| File Navigation | harpoon, arrow, grapple |
| File Explorer | neo-tree, nvim-tree, oil |
| Git | gitsigns, diffview, lazygit, fugitive, neogit |
| LSP | nvim-lspconfig, none-ls, lsp-zero |
| Completion | blink.cmp, nvim-cmp, coq_nvim |
| Snippets | LuaSnip, UltiSnips, snippy |
| Comments | Comment.nvim, ts-context-commentstring |
| Surround | mini.surround, nvim-surround |
| Statusline | lualine, mini.statusline, heirline |
| Bufferline | bufferline, barbar, tabby |
| Terminal | toggleterm, terminal-manager, FTerm |
| Debugging | nvim-dap, vimspector |
| Testing | neotest, vim-test |
| Formatting | conform, null-ls, formatter |
| Diagnostics | trouble, lsp_lines |
