# Audit Neovim Configuration

You are a Neovim configuration auditor. Analyze the current config structure and report issues.

## Expected Structure

```
init.lua                    (minimal - just requires core)

lua/
  core/
    init.lua                (loads options, keymaps, autocmds, lazy)
    options.lua             (vim.o/vim.g settings)
    keymaps.lua             (non-plugin keymaps)
    autocmds.lua            (autocommands)
    lazy.lua                (lazy.nvim bootstrap)

  plugins/
    init.lua                (returns all plugin specs)
    telescope.lua           (telescope + extensions)
    lsp.lua                 (lspconfig + mason)
    treesitter.lua          (treesitter + context)
    completion.lua          (blink.cmp + snippets)
    formatting.lua          (conform)
    colorscheme.lua         (catppuccin)
    ui.lua                  (which-key, scrollbar, indent)
    bufferline.lua          (bufferline + scope)
    statusline.lua          (lualine)
    git.lua                 (gitsigns + lazygit + diffview)
    navigation.lua          (harpoon + neo-tree + smart-splits)
    editing.lua             (mini + comment + autopairs)
    diagnostics.lua         (trouble)
    terminal.lua            (tmux integration)
    claude-code.lua         (Claude Code)
    docusaurus.lua          (Docusaurus editing)
    plugin-help.lua         (help system)

  lib/                      (shared utilities)
    forks.lua
    plugin-help.lua
    plugin-registry.lua

  custom/                   (standalone modules)
    tab-picker.lua
    tmux/
```

## Audit Steps

1. **Structure Check**: Verify all expected files exist
2. **Orphan Detection**: Find files not in expected locations
3. **Plugin Format Validation**: Each plugins/*.lua should return a table
4. **Keymap Audit**: Check for duplicate leader keymaps
5. **Load Test**: Verify Neovim starts without errors
6. **Import Check**: Verify plugins/init.lua imports all plugin files

## Run the Audit

Execute these checks and produce a report:

### 1. Structure Check
```bash
# Check core files exist
for f in init options keymaps autocmds lazy; do
  test -f lua/core/$f.lua && echo "OK: core/$f.lua" || echo "MISSING: core/$f.lua"
done

# Check plugin files exist
for f in init telescope lsp treesitter completion formatting colorscheme ui bufferline statusline git navigation editing diagnostics terminal claude-code docusaurus plugin-help; do
  test -f lua/plugins/$f.lua && echo "OK: plugins/$f.lua" || echo "MISSING: plugins/$f.lua"
done
```

### 2. Orphan Detection
```bash
# Find lua files not in expected directories
find lua -name "*.lua" | grep -v -E "(core|plugins|lib|custom)/" | head -20
```

### 3. Keymap Conflicts
```lua
-- In Neovim, check for duplicate leader keymaps
:lua for _, m in ipairs(vim.api.nvim_get_keymap('n')) do if m.lhs:match('^<[Ll]eader>') then print(m.lhs, m.desc or '') end end
```

### 4. Plugin Return Check
Each plugin file should return a table. Grep for patterns like:
- `return {` at start of return statement
- No stray code after the return

### 5. Load Test
```bash
nvim --headless -c 'lua print("OK")' -c 'qall' 2>&1
```

## Report Format

```
## Neovim Config Audit Report

### Structure: [PASS/FAIL]
- Missing files: [list]
- Orphan files: [list]

### Plugin Format: [PASS/FAIL]
- Issues: [list]

### Keymaps: [PASS/WARN]
- Conflicts: [list]
- Unused leader keys: [suggestions]

### Load Test: [PASS/FAIL]
- Errors: [list]

### Recommendations
1. [actionable improvements]
```

## Execute Now

Run all audit checks and produce a complete report. Use Read/Glob/Grep tools to analyze files. Use Bash for Neovim tests.
