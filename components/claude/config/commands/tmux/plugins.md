---
description: Manage tmux plugins with TPM
argument-hint: [action: install|update|list|add]
allowed-tools: Read, Write, Edit, Bash
---

## Task

Help the user manage tmux plugins using TPM (Tmux Plugin Manager).

## Actions

Based on `$ARGUMENTS`:

### install
Install TPM and all plugins:
```bash
tmux_install_tpm      # Install TPM
tmux_install_plugins  # Install plugins (run outside tmux)
# OR inside tmux: prefix + I
```

### update
Update TPM and plugins:
```bash
tmux_update_tpm       # Update TPM itself
tmux_update_plugins   # Update all plugins
# OR inside tmux: prefix + U
```

### list
Show installed plugins. Check `~/.config/tmux/tmux.conf` for plugin declarations.

### add
Add a new plugin to tmux.conf:
1. Add to plugins section: `set -g @plugin 'author/plugin-name'`
2. Reload config: `tmux_reload`
3. Install: `prefix + I` (inside tmux) or `tmux_install_plugins`

## TPM Plugin Format

In `~/.config/tmux/tmux.conf`:
```bash
# Plugin manager
set -g @plugin 'tmux-plugins/tpm'

# Plugins
set -g @plugin 'tmux-plugins/tmux-sensible'
set -g @plugin 'catppuccin/tmux'
set -g @plugin 'tmux-plugins/tmux-resurrect'

# Initialize TPM (keep at bottom)
run '~/.config/tmux/plugins/tpm/tpm'
```

## Recommended Plugins

### Essential
| Plugin | Purpose |
|--------|---------|
| `tmux-plugins/tpm` | Plugin manager (required) |
| `tmux-plugins/tmux-sensible` | Sensible defaults |

### Productivity
| Plugin | Purpose |
|--------|---------|
| `tmux-plugins/tmux-resurrect` | Save/restore sessions |
| `tmux-plugins/tmux-continuum` | Auto-save sessions |
| `christoomey/vim-tmux-navigator` | Seamless vim navigation |

### Session Management
| Plugin | Purpose |
|--------|---------|
| `omerxx/tmux-sessionx` | Enhanced session picker |
| `27medkamal/tmux-session-wizard` | fzf session management |

### Appearance
| Plugin | Purpose |
|--------|---------|
| `catppuccin/tmux` | Catppuccin theme |
| `tmux-plugins/tmux-prefix-highlight` | Show when prefix is active |

### Utilities
| Plugin | Purpose |
|--------|---------|
| `tmux-plugins/tmux-yank` | Better copy/paste |
| `tmux-plugins/tmux-open` | Open files/URLs from terminal |

## TPM Key Bindings

- `prefix + I` - Install plugins
- `prefix + U` - Update plugins
- `prefix + alt + u` - Uninstall plugins not in config

## Dotfiles Functions

```bash
tmux_install_tpm      # Install TPM
tmux_update_tpm       # Update TPM
tmux_install_plugins  # Install all plugins
tmux_update_plugins   # Update all plugins
```

## Context

@components/tmux/config.yaml
@components/tmux/config/tmux.conf

## Troubleshooting

1. **Plugins not loading**: Ensure TPM init is at bottom of tmux.conf
2. **Install fails**: Check `~/.config/tmux/plugins/` permissions
3. **Theme not working**: Some themes need specific terminal capabilities
