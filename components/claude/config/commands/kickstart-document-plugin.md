# Document Plugin

Research and document a Neovim plugin comprehensively.

## Arguments

$ARGUMENTS - Plugin name or `owner/repo` format (e.g., `telescope.nvim` or `nvim-telescope/telescope.nvim`)

## Mission

Create thorough documentation for a plugin that helps the user understand:
- What the plugin does and why it's useful
- How it works under the hood
- All available features and commands
- Configuration options and their effects
- Keymaps and how to use them effectively

## Research Workflow

### Step 1: Gather Information

1. **Fetch GitHub README**:
   ```
   WebFetch https://github.com/{owner}/{repo}
   ```

2. **Check for documentation site**:
   - Look for `/doc` folder, wiki, or external docs
   - Search for `{plugin} documentation site`

3. **Find usage examples**:
   - Search GitHub for config examples
   - Look at popular dotfiles using this plugin

### Step 2: Analyze Local Configuration

Read how the plugin is configured in this repo:
- `init.lua` - check if it's a core plugin
- `lua/custom/plugins/` - check for custom config
- `lua/kickstart/plugins/` - check kickstart config

### Step 3: Document Structure (Docusaurus-Compatible)

Create documentation in `.claude/docs/plugins/{plugin-name}.mdx`:

```mdx
---
title: {Plugin Name}
description: {One-line description for SEO and previews}
sidebar_position: {number based on category}
sidebar_label: {Short name for sidebar}
tags:
  - {category}
  - {sub-category if applicable}
keywords:
  - neovim
  - {plugin-name}
  - {related keywords}
---

# {Plugin Name}

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

> {One-line description}

| | |
|---|---|
| **Repository** | [{owner}/{repo}](https://github.com/{owner}/{repo}) |
| **Category** | {category from list below} |
| **Our Config** | `{path to config file}` |

## What It Does

{2-3 paragraph explanation of the plugin's purpose and benefits}

:::info Why We Use This
{Brief explanation of why this plugin is in our config}
:::

## Key Concepts

### {Concept 1}
{Explanation}

### {Concept 2}
{Explanation}

## Features

| Feature | Description | Command/Keymap |
|---------|-------------|----------------|
| ... | ... | ... |

## Our Configuration

```lua title="lua/custom/plugins/{plugin}.lua" showLineNumbers
{Current config with inline comments explaining each option}
```

<details>
<summary>Why These Settings</summary>

- **{setting}**: {why we chose this value}
- **{setting}**: {why we chose this value}

</details>

## Keymaps

:::tip Quick Reference
The most important keymaps to remember are highlighted below.
:::

| Keymap | Mode | Action | When to Use |
|--------|------|--------|-------------|
| ... | ... | ... | ... |

## Commands

| Command | Description | Example |
|---------|-------------|---------|
| `:Command` | {description} | `:Command arg` |

## Common Workflows

<Tabs>
<TabItem value="workflow1" label="{Workflow 1 Name}" default>

1. Step one
2. Step two
3. Step three

</TabItem>
<TabItem value="workflow2" label="{Workflow 2 Name}">

1. Step one
2. Step two

</TabItem>
</Tabs>

## Integration Points

```mermaid
graph LR
    A[{Plugin}] --> B[telescope]
    A --> C[which-key]
    A --> D[other-plugin]
```

- **Works with:** {other plugins it integrates with}
- **Requires:** {dependencies}
- **Conflicts with:** {incompatible plugins, if any}

## Tips & Tricks

:::tip {Tip 1 Title}
{Explanation of the tip}
:::

:::tip {Tip 2 Title}
{Explanation of the tip}
:::

## Troubleshooting

<details>
<summary>{Common Issue 1}</summary>

**Problem:** {description}

**Solution:**
```lua
{fix code if applicable}
```

</details>

<details>
<summary>{Common Issue 2}</summary>

**Problem:** {description}

**Solution:** {fix}

</details>

## Advanced Configuration

:::warning Advanced Users Only
These settings are for power users. The defaults work well for most cases.
:::

```lua title="Advanced options"
{Advanced config options with explanations}
```

## Resources

- [GitHub Repository](https://github.com/{owner}/{repo})
- [Documentation]({docs link if available})
- [Wiki]({wiki link if available})

---

*Last updated: {date}*
```

### Docusaurus Features Reference

| Feature | Syntax | Use Case |
|---------|--------|----------|
| Frontmatter | `---\ntitle: ...\n---` | SEO, sidebar config |
| Admonitions | `:::tip`, `:::info`, `:::warning`, `:::danger`, `:::note` | Callouts |
| Tabs | `<Tabs>/<TabItem>` | Alternative workflows |
| Details | `<details><summary>` | Collapsible sections |
| Mermaid | ` ```mermaid ` | Diagrams |
| Code titles | ` ```lua title="file.lua" ` | Show filename |
| Line numbers | ` ```lua showLineNumbers ` | Reference specific lines |
| Line highlight | ` ```lua {2-4} ` | Highlight lines 2-4 |

### Sidebar Position by Category

| Category | Position Range |
|----------|---------------|
| Navigation | 100-199 |
| Editor | 200-299 |
| UI | 300-399 |
| Git | 400-499 |
| LSP | 500-599 |
| Completion | 600-699 |
| Treesitter | 700-799 |
| Terminal | 800-899 |
| Debugging | 900-999 |
| Formatting | 1000-1099 |
| Utilities | 1100-1199 |

## Categories Reference

| Category | Description |
|----------|-------------|
| Navigation | File finding, jumping, marks |
| Editor | Editing enhancements, surround, comments |
| UI | Statusline, bufferline, notifications |
| Git | Version control integration |
| LSP | Language server features |
| Completion | Autocompletion, snippets |
| Treesitter | Syntax, text objects |
| Terminal | Terminal integration |
| Debugging | DAP, testing |
| Formatting | Code formatting, linting |
| Utilities | General utilities |

## Output

After researching, create:

1. **Documentation file** at `.claude/docs/plugins/{plugin-name}.mdx`
2. **Registry entry** in `lua/lib/plugin-registry.lua` (if not exists)
3. **Summary** for the user with key takeaways
4. **Suggestions** for improving our current config (if applicable)

### Registry Entry Format

Add to `lua/lib/plugin-registry.lua`:

```lua
help.register {
  name = '{plugin}',
  category = '{category}',
  description = '{one-line description}',
  keymaps = {
    { key = '<leader>xx', desc = 'Action description', mode = 'n' },
    -- ... all keymaps
  },
  commands = {
    { cmd = '{Plugin}Action', desc = 'What it does' },
    -- ... all commands
  },
  tips = {
    'Quick tip 1',
    'Quick tip 2',
  },
}
```

This enables:
- `:PluginHelp {plugin}` - Show in-editor help
- `<leader>?` - Pick plugin from list
- `:PluginKeymaps` - Browse all keymaps

## Example

**Input:** `/document-plugin telescope.nvim`

**Creates:** `.claude/docs/plugins/telescope.mdx`

**Summary:**
> Telescope is a fuzzy finder that can search files, buffers, git commits, LSP symbols, and more. Our config uses fzf-native for speed and ui-select for consistent UI. Key keymaps: `<leader>ff` (files), `<leader>fg` (grep), `<leader>sb` (buffers).

**Suggestions:**
- Consider adding `<leader>ft` for Telescope resume (continue last search)
- The `hidden = true` option would include dotfiles in searches
