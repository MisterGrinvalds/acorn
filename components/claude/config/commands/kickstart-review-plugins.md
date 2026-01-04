# Review Plugins

List all plugins currently in use in this Neovim configuration.

## Instructions

Scan all plugin definition locations and create a comprehensive inventory.

### Step 1: Scan Plugin Locations

Read and parse plugins from:
1. `init.lua` - Extract all plugin definitions from `require('lazy').setup({...})`
2. `lua/kickstart/plugins/*.lua` - Optional plugins that are `require`d
3. `lua/custom/plugins/*.lua` - User-added plugins

### Step 2: Extract Plugin Information

For each plugin, identify:
- **Plugin name**: `owner/repo` format or local `dir` path
- **Forked?**: Check if it uses `require('lib.forks').fork()` wrapper
- **Source file**: Where the plugin is defined
- **Key options**: lazy loading, dependencies, etc.

### Step 3: Present Results

Create a table grouped by source:

```
## Core Plugins (init.lua)
| Plugin | Owner/Repo | Forked? |
|--------|------------|---------|
| telescope.nvim | nvim-telescope/telescope.nvim | No |
| ... | ... | ... |

## Kickstart Optional (lua/kickstart/plugins/)
| Plugin | Owner/Repo | Forked? |
|--------|------------|---------|
| neo-tree | nvim-neo-tree/neo-tree.nvim | No |
| ... | ... | ... |

## Custom Plugins (lua/custom/plugins/)
| Plugin | Owner/Repo | Forked? |
|--------|------------|---------|
| harpoon | ThePrimeagen/harpoon | No |
| ... | ... | ... |

## Forked Plugins
| Plugin | Original | Fork URL |
|--------|----------|----------|
| (none yet) | | |
```

### Step 4: Summary Stats

Provide counts:
- Total plugins
- Plugins using forks
- Plugins that could be forked (external GitHub repos)

## Notes

- Local `dir` plugins (like claude-code.nvim) cannot be forked
- Plugins with `url` overrides should note the custom URL
- Check for `fork('...')` pattern to detect existing forks
