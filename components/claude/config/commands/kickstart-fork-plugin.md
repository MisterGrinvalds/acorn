# Fork Plugin

Fork a Neovim plugin to MisterGrinvalds GitHub account and update the config to use it.

## Arguments

$ARGUMENTS - Plugin name in `owner/repo` format (e.g., `nvim-telescope/telescope.nvim`)

## Workflow

### Step 1: Parse and Validate

Extract owner and repo from `$ARGUMENTS`:
- owner = first part before `/`
- repo = second part after `/`

Validate it's a valid GitHub repository format.

### Step 2: Check if Fork Already Exists

```bash
gh repo view MisterGrinvalds/{owner}.{repo} --json name 2>/dev/null
```

If fork exists, skip to Step 5 (config update).

### Step 3: Create GitHub Fork

```bash
gh repo fork $ARGUMENTS --clone=false
```

This creates `MisterGrinvalds/{repo}` (GitHub's default naming).

### Step 4: Rename Fork to Convention

Rename to match our naming convention `{owner}.{repo}`:

```bash
gh repo rename {owner}.{repo} --repo MisterGrinvalds/{repo}
```

**Example:**
- Original: `ThePrimeagen/harpoon`
- After fork: `MisterGrinvalds/harpoon`
- After rename: `MisterGrinvalds/ThePrimeagen.harpoon`

### Step 5: Find Plugin in Config

Search for `$ARGUMENTS` in:
1. `init.lua` - inline plugin specs
2. `lua/kickstart/plugins/*.lua` - optional plugins
3. `lua/custom/plugins/*.lua` - custom plugins

### Step 6: Update Config to Use Fork

Add the fork import if not present:
```lua
local fork = require('lib.forks').fork
```

Update the plugin spec:
```lua
-- Before:
{ 'owner/repo', opts = {...} }

-- After:
{ fork('owner/repo'), opts = {...} }
```

### Step 7: Verify

1. Run `gh repo view MisterGrinvalds/{owner}.{repo}` to confirm fork exists
2. Suggest running `:Lazy sync` in Neovim to verify resolution

### Step 8: Report

Summarize:
- Fork URL: `https://github.com/MisterGrinvalds/{owner}.{repo}`
- Config file updated: `path/to/file.lua`
- Original upstream: `https://github.com/{owner}/{repo}`

## Example

**Input:** `/fork-plugin ThePrimeagen/harpoon`

**Actions:**
1. `gh repo fork ThePrimeagen/harpoon --clone=false`
2. `gh repo rename ThePrimeagen.harpoon --repo MisterGrinvalds/harpoon`
3. Update `lua/custom/plugins/harpoon.lua`:

```lua
local fork = require('lib.forks').fork

return {
  fork('ThePrimeagen/harpoon'),
  branch = 'harpoon2',
  ...
}
```

**Result:**
- Fork: `https://github.com/MisterGrinvalds/ThePrimeagen.harpoon`
- Config updated to use fork

## Error Handling

- **Fork already exists with correct name**: Skip fork/rename, proceed to config update
- **Fork exists but wrong name**: Rename it to convention
- **Permission denied**: User may need to authenticate with `gh auth login`
- **Plugin not found in config**: Create new file in `lua/custom/plugins/` or ask user where to add it

## Notes

- Only the main plugin is forked, not its dependencies
- The fork maintains upstream as a remote for easy syncing
- Use `/update-forks` to sync with upstream later
