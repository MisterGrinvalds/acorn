---
description: Configure fzf defaults, theme, and keybindings
argument-hint: [aspect: theme|keybindings|defaults]
allowed-tools: Read, Write, Edit, Bash
---

## Task

Help the user configure fzf settings, theme, and shell integration.

## Configuration Aspects

Based on `$ARGUMENTS`:

### theme
Configure Catppuccin or custom colors:

```bash
# Catppuccin Mocha theme
export FZF_DEFAULT_OPTS="
  --color=bg+:#313244,bg:#1e1e2e,spinner:#f5e0dc,hl:#f38ba8
  --color=fg:#cdd6f4,header:#f38ba8,info:#cba6f7,pointer:#f5e0dc
  --color=marker:#f5e0dc,fg+:#cdd6f4,prompt:#cba6f7,hl+:#f38ba8
  --color=border:#6c7086
"
```

Color elements:
- `bg` - Background
- `fg` - Foreground text
- `hl` - Highlighted matches
- `pointer` - Selection pointer
- `marker` - Multi-select marker
- `prompt` - Prompt text
- `border` - Border color

### keybindings
Configure shell integration:

```bash
# In .bashrc or .zshrc after fzf init
# Ctrl+R - History
# Ctrl+T - File finder
# Alt+C - Directory changer

# Custom keybinding example
export FZF_CTRL_T_OPTS="--preview 'bat --color=always {}'"
export FZF_ALT_C_OPTS="--preview 'ls -la {}'"
```

### defaults
Configure FZF_DEFAULT_OPTS:

```bash
export FZF_DEFAULT_OPTS="
  --height 40%
  --layout=reverse
  --border rounded
  --preview-window=right:50%
  --bind 'ctrl-/:toggle-preview'
  --bind 'ctrl-a:select-all'
  --bind 'ctrl-y:execute-silent(echo {} | pbcopy)'
"
```

Common options:
- `--height` - Percentage of terminal height
- `--layout` - `default`, `reverse`, `reverse-list`
- `--border` - `rounded`, `sharp`, `bold`, `none`
- `--preview` - Preview command
- `--preview-window` - Position and size
- `--bind` - Custom keybindings

## File Generation Commands

FZF_DEFAULT_COMMAND:
```bash
# Using fd (recommended)
export FZF_DEFAULT_COMMAND='fd --type f --hidden --follow --exclude .git'

# Using find (fallback)
export FZF_DEFAULT_COMMAND='find . -type f -not -path "*/\.git/*"'
```

## Context

@components/fzf/env.sh

## Configuration Location

Add to `components/fzf/env.sh` or `~/.config/shell/local.sh`

## Tips

1. Use `--height 40%` to avoid fullscreen
2. Use `--layout=reverse` for top-down display
3. Preview with bat for syntax highlighting
4. Test options with: `fzf --option1 --option2`
