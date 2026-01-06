---
description: Configure tmux settings and customization
argument-hint: [setting-type]
allowed-tools: Read, Write, Edit, Bash
---

## Task

Help the user configure and customize tmux settings.

## Configuration File

Location: `~/.config/tmux/tmux.conf` (XDG) or `~/.tmux.conf` (legacy)

Edit with: `tmux_config`
Reload with: `tmux_reload` or `prefix + r`

## Common Configuration Areas

### Prefix Key
```bash
# Change prefix from Ctrl-b to Ctrl-a
unbind C-b
set -g prefix C-a
bind C-a send-prefix
```

### Mouse Support
```bash
set -g mouse on
```

### Window/Pane Indexing
```bash
set -g base-index 1        # Start windows at 1
setw -g pane-base-index 1  # Start panes at 1
set -g renumber-windows on # Renumber when window closed
```

### Vi Mode
```bash
setw -g mode-keys vi
bind -T copy-mode-vi v send -X begin-selection
bind -T copy-mode-vi y send -X copy-selection-and-cancel
```

### Pane Navigation (Vim-style)
```bash
bind h select-pane -L
bind j select-pane -D
bind k select-pane -U
bind l select-pane -R
```

### Pane Splitting (intuitive keys)
```bash
bind | split-window -h -c "#{pane_current_path}"
bind - split-window -v -c "#{pane_current_path}"
```

### History
```bash
set -g history-limit 50000
```

### Colors/Terminal
```bash
set -g default-terminal "tmux-256color"
set -ag terminal-overrides ",xterm-256color:RGB"
```

### Status Bar
```bash
set -g status-position top
set -g status-interval 5
set -g status-left-length 50
set -g status-right-length 50
```

## Catppuccin Theme Configuration

```bash
set -g @catppuccin_flavour 'mocha'
set -g @catppuccin_window_status_style "rounded"
set -g @catppuccin_status_modules_right "session date_time"
```

## Dotfiles Functions

- `tmux_config` - Open config in editor
- `tmux_reload` - Reload configuration
- `tmux_info` - Show current tmux info

## Config Structure Template

```bash
# =============================================================================
# General Settings
# =============================================================================
set -g prefix C-a
set -g mouse on
set -g base-index 1

# =============================================================================
# Key Bindings
# =============================================================================
bind r source-file ~/.config/tmux/tmux.conf \; display "Reloaded!"
bind | split-window -h -c "#{pane_current_path}"
bind - split-window -v -c "#{pane_current_path}"

# =============================================================================
# Appearance
# =============================================================================
set -g status-position top
set -g default-terminal "tmux-256color"

# =============================================================================
# Plugins
# =============================================================================
set -g @plugin 'tmux-plugins/tpm'
set -g @plugin 'catppuccin/tmux'

# Initialize TPM (keep at bottom)
run '~/.config/tmux/plugins/tpm/tpm'
```

## Context

@components/tmux/config.yaml
@components/tmux/config/tmux.conf

## Tips

1. Use `tmux show-options -g` to see all global options
2. Use `tmux list-keys` to see all key bindings
3. Test changes with `tmux source-file ~/.config/tmux/tmux.conf`
4. Keep TPM initialization at the very end of config
