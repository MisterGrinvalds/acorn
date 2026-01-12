# Replace Plugin with Fork

Update an existing plugin spec to use the `fork()` wrapper function.

## Arguments

$ARGUMENTS - Plugin name in `owner/repo` format (e.g., `ThePrimeagen/harpoon`)

## Prerequisites

The fork must already exist at `MisterGrinvalds/{owner}.{repo}` on GitHub.

If the fork doesn't exist, tell the user to run `/fork-plugin $ARGUMENTS` first.

## Workflow

### Step 1: Verify Fork Exists

```bash
gh repo view MisterGrinvalds/{owner}.{repo} --json name 2>/dev/null
```

If this fails, the fork doesn't exist. Stop and suggest running `/fork-plugin` first.

### Step 2: Find Plugin in Config

Search for `$ARGUMENTS` pattern in:
1. `init.lua`
2. `lua/kickstart/plugins/*.lua`
3. `lua/custom/plugins/*.lua`

### Step 3: Identify Current Spec Pattern

**Pattern A: String only**
```lua
'owner/repo',
```

**Pattern B: Table with string first element**
```lua
{ 'owner/repo', opts = {...} }
```

### Step 4: Add Fork Import

If the file doesn't already have the fork import, add it near the top:

```lua
local fork = require('lib.forks').fork
```

### Step 5: Apply Fork Wrapper

**For Pattern A:**
```lua
-- Before:
'owner/repo',

-- After:
fork('owner/repo'),
```

**For Pattern B:**
```lua
-- Before:
{ 'owner/repo', opts = {...} }

-- After:
{ fork('owner/repo'), opts = {...} }
```

### Step 6: Verify

Suggest running `:Lazy sync` to verify the fork resolves correctly.

## Example

**Input:** `ThePrimeagen/harpoon`

**Before** (`lua/custom/plugins/harpoon.lua`):
```lua
return {
  'ThePrimeagen/harpoon',
  branch = 'harpoon2',
  dependencies = { 'nvim-lua/plenary.nvim' },
  ...
}
```

**After:**
```lua
local fork = require('lib.forks').fork

return {
  fork('ThePrimeagen/harpoon'),
  branch = 'harpoon2',
  dependencies = { 'nvim-lua/plenary.nvim' },
  ...
}
```

## Notes

- The `fork()` function transforms `owner/repo` to `MisterGrinvalds/owner.repo`
- Dependencies are NOT automatically forked - only the main plugin
- If you want to fork dependencies too, run this command for each one
