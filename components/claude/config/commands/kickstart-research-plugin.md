# Research Plugin

Deep-dive research into how a Neovim plugin works internally.

## Arguments

$ARGUMENTS - Plugin name or `owner/repo` format (e.g., `telescope.nvim` or `nvim-telescope/telescope.nvim`)

## Mission

Conduct comprehensive research on a plugin's internals, architecture, and API to understand:
- How the plugin works under the hood
- Its Lua API and extension points
- Configuration options and their effects
- How to extend or customize it
- Common patterns and advanced usage

## Research Workflow

### Step 1: Gather Source Information

1. **Fetch GitHub Repository**:
   ```
   WebFetch https://github.com/{owner}/{repo}
   ```

2. **Explore Documentation**:
   - Check `/doc/` folder for Vim help files
   - Look for wiki or external documentation
   - Search: `site:github.com/{owner}/{repo} documentation`

3. **Find Source Code Structure**:
   - `lua/{plugin}/init.lua` - Main entry point
   - `lua/{plugin}/config.lua` - Default configuration
   - `lua/{plugin}/actions.lua` - Actions/commands
   - `plugin/{plugin}.vim` or `plugin/{plugin}.lua` - Auto-loaded setup

### Step 2: Analyze Architecture

1. **Module Structure**:
   - What modules does the plugin expose?
   - What's the public API vs internal?
   - How does lazy loading work?

2. **Configuration System**:
   - How does `setup({})` process options?
   - What are ALL available options?
   - What are the defaults?

3. **Event System**:
   - What autocommands does it create?
   - What user events does it emit?
   - What hooks are available?

4. **Integration Points**:
   - Does it expose a Telescope extension?
   - Does it provide which-key integration?
   - Does it emit events other plugins can hook?

### Step 3: Document API

Create detailed API documentation:

```markdown
## Public API

### require('{plugin}').setup(opts)
Main configuration function.

**Options:**
| Option | Type | Default | Description |
|--------|------|---------|-------------|
| ... | ... | ... | ... |

### require('{plugin}').{function}()
{Description of what this function does}

**Parameters:**
- `param1` (type): description
- `param2` (type): description

**Returns:** {return type and description}

**Example:**
```lua
require('{plugin}').{function}({ ... })
```
```

### Step 4: Discover Advanced Features

Research advanced usage patterns:

1. **Custom Actions/Extensions**:
   - How to add custom commands
   - How to create custom pickers/sources
   - How to add custom keymaps dynamically

2. **Performance Considerations**:
   - Lazy loading recommendations
   - Memory/CPU impact
   - Async operations

3. **Debugging**:
   - How to enable debug mode
   - Common error messages
   - Health checks (`:checkhealth {plugin}`)

### Step 5: Find Community Examples

Search for advanced usage:
- `site:github.com "require('{plugin}')" setup`
- `site:reddit.com/r/neovim {plugin} config`
- `site:dotfyle.com {plugin}`

## Output Format

### Research Summary

Create `.claude/docs/plugins/{plugin-name}-research.md`:

```markdown
# {Plugin Name} - Deep Dive Research

## Overview
{Brief description of what the plugin does}

## Architecture

### Module Structure
```
{plugin}/
├── init.lua        # Main entry, exports public API
├── config.lua      # Default configuration
├── {module}.lua    # Feature module
└── ...
```

### Data Flow
{Explain how data flows through the plugin}

## Configuration Reference

### All Options
| Option | Type | Default | Description |
|--------|------|---------|-------------|
| ... | ... | ... | ... |

### Option Deep Dive

#### `{important_option}`
{Detailed explanation of what this option does, when to change it, edge cases}

## Public API

### Core Functions
| Function | Description |
|----------|-------------|
| `setup(opts)` | Initialize the plugin |
| ... | ... |

### Detailed API

#### `require('{plugin}').{function}()`
{Full documentation}

## Extension Points

### Adding Custom {Feature}
```lua
-- Example of extending the plugin
```

### Event Hooks
| Event | When Fired | Payload |
|-------|------------|---------|
| ... | ... | ... |

## Integration Guide

### With Telescope
```lua
-- How to integrate with telescope
```

### With Which-Key
```lua
-- How to add to which-key groups
```

### With Other Plugins
{List of compatible plugins and how to integrate}

## Advanced Patterns

### Pattern 1: {Name}
{Description and example}

### Pattern 2: {Name}
{Description and example}

## Performance

- **Lazy Loading:** {recommended strategy}
- **Memory Usage:** {notes}
- **Startup Impact:** {ms estimate}

## Debugging

### Enable Debug Mode
```lua
{how to enable verbose logging}
```

### Common Issues
| Issue | Cause | Solution |
|-------|-------|----------|
| ... | ... | ... |

### Health Check
```vim
:checkhealth {plugin}
```

## Source Code References

| File | Purpose |
|------|---------|
| `lua/{plugin}/init.lua` | Main entry |
| ... | ... |

## Resources

- [GitHub]({repo_url})
- [Documentation]({docs_url})
- [Reddit Discussions]({reddit_search})
```

## Example

**Input:** `/research-plugin telescope.nvim`

**Creates:** `.claude/docs/plugins/telescope-research.md`

**Key Findings:**
- Public API has 47 functions across 12 modules
- Supports custom pickers via `pickers.new()`
- Emits `TelescopePreviewerLoaded` event
- Extension system via `telescope.load_extension()`
- Debug mode: `require('telescope').setup({ debug = true })`
