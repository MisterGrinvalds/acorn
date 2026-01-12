# Select Command

Help the user select and run the most appropriate command for their task.

## Mission

Act as an interactive command selector that:
- Understands what the user wants to accomplish
- Presents relevant commands from the available options
- Helps them choose the right one
- Runs the selected command with proper arguments

## Available Commands

| Command | Purpose | When to Use |
|---------|---------|-------------|
| `/nvim-config` | Interactive configuration wizard | Modifying settings, options, or preferences |
| `/nvim-research` | Research Neovim ecosystem | Exploring new features, trends, or alternatives |
| `/audit-config` | Audit configuration health | Checking for issues, orphans, conflicts |
| `/add-plugin` | Add a new plugin | Installing something new to the config |
| `/fork-plugin` | Fork a plugin to your account | Taking ownership of a plugin for customization |
| `/replace-with-fork` | Convert plugin to use fork | Switching an existing plugin to your fork |
| `/update-forks` | Sync forks with upstream | Keeping forked plugins up to date |
| `/research-plugin` | Deep-dive plugin analysis | Understanding plugin internals and API |
| `/document-plugin` | Create plugin documentation | Writing comprehensive docs for a plugin |
| `/integrate-plugin` | Set up keymaps and commands | Configuring keybindings and help integration |
| `/review-plugins` | List all plugins and status | Getting an overview of installed plugins |

## Workflow

### Step 1: Understand the Request

Ask clarifying questions if the user's intent isn't clear:

```
What would you like to do with your Neovim config?

1. **Configure** - Change settings, keymaps, or options
2. **Add/Install** - Add a new plugin
3. **Document** - Create or update documentation
4. **Maintain** - Audit, update forks, fix issues
5. **Research** - Explore ecosystem, learn about plugins
6. **Review** - See what's installed and how it's organized
```

### Step 2: Present Relevant Options

Based on their response, show only the relevant commands:

**For Configuration:**
- `/nvim-config` - Interactive wizard for settings
- `/integrate-plugin` - Set up keymaps for a plugin

**For Adding Plugins:**
- `/add-plugin {plugin}` - Research and add a new plugin
- `/fork-plugin {owner/repo}` - Fork first, then add

**For Documentation:**
- `/document-plugin {plugin}` - Full MDX documentation
- Location: `.claude/docs/plugins/{plugin}.mdx`

**For Maintenance:**
- `/audit-config` - Check for issues
- `/update-forks` - Sync with upstream
- `/replace-with-fork {plugin}` - Convert to fork

**For Research:**
- `/nvim-research` - Ecosystem trends and recommendations
- `/research-plugin {plugin}` - Deep-dive analysis

**For Review:**
- `/review-plugins` - List all plugins and their status

### Step 3: Execute the Command

Once the user selects a command:

1. Confirm the selection and any required arguments
2. Use the Skill tool to invoke the selected command
3. Pass any arguments the user provided

## Quick Selection Mode

If the user provides context, skip to the relevant command:

| User Says | Suggested Command |
|-----------|-------------------|
| "I want to add X" | `/add-plugin X` |
| "Document X" | `/document-plugin X` |
| "Research X" | `/research-plugin X` |
| "What plugins do I have?" | `/review-plugins` |
| "Check for issues" | `/audit-config` |
| "Update my forks" | `/update-forks` |
| "Set up keymaps for X" | `/integrate-plugin X` |
| "Fork X" | `/fork-plugin X` |

## Documentation Output Locations

All documentation commands write to Docusaurus-compatible locations:

| Type | Location | Format |
|------|----------|--------|
| Plugin docs | `.claude/docs/plugins/{name}.mdx` | MDX with frontmatter |
| General docs | `.claude/docs/` | Markdown |
| Quick reference | `.claude/docs/quick-reference.md` | Markdown |
| User guide | `.claude/docs/user-guide.md` | Markdown |

## Example Interaction

**User:** I want to document the telescope plugin

**Assistant:**
I'll help you document telescope. The `/document-plugin` command will:

1. Research the plugin on GitHub
2. Analyze your current configuration
3. Create comprehensive MDX documentation at `.claude/docs/plugins/telescope.mdx`
4. Add a registry entry for in-editor help

Run `/document-plugin telescope.nvim`?

**User:** yes

**Assistant:** *Invokes `/document-plugin telescope.nvim`*
