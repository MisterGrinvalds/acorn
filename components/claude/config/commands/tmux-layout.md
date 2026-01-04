---
description: Configure and manage tmux window and pane layouts
argument-hint: [layout-type]
allowed-tools: Read, Bash
---

## Task

Help the user understand and configure tmux layouts for windows and panes.

## Built-in Layouts

Tmux has 5 preset layouts (cycle with `prefix + space`):

1. **even-horizontal** - Panes spread horizontally
2. **even-vertical** - Panes spread vertically
3. **main-horizontal** - Main pane on top, others below
4. **main-vertical** - Main pane on left, others right
5. **tiled** - Equal-sized grid

## Layout Commands

### Pane Splitting
- `prefix + %` - Split horizontally (left/right)
- `prefix + "` - Split vertically (top/bottom)
- `tmux split-window -h` - Horizontal split
- `tmux split-window -v` - Vertical split

### Pane Resizing
- `prefix + Ctrl+arrow` - Resize in direction
- `prefix + z` - Zoom/unzoom current pane
- `tmux resize-pane -D 5` - Resize down 5 cells
- `tmux resize-pane -R 10` - Resize right 10 cells

### Layout Selection
- `prefix + space` - Cycle through layouts
- `tmux select-layout main-vertical`
- `tmux select-layout tiled`

### Custom Layouts
Capture current layout:
```bash
tmux list-windows -F "#{window_layout}"
```

Apply saved layout:
```bash
tmux select-layout "layout-string-here"
```

## Common Workflow Layouts

### IDE Layout (main-vertical)
```
+----------------+--------+
|                |        |
|     Editor     | Tests  |
|                |        |
|                +--------+
|                | Logs   |
+----------------+--------+
```

### Monitoring Layout (tiled)
```
+--------+--------+
|  App1  |  App2  |
+--------+--------+
|  Logs  |  K9s   |
+--------+--------+
```

### Presentation Layout
```
+------------------------+
|                        |
|      Main Content      |
|                        |
+------------------------+
|   Notes (small pane)   |
+------------------------+
```

## Dotfiles Functions

Use `dev_session` for a quick 3-pane dev layout:
```bash
dev_session myproject
# Creates: Editor | Main + Logs layout
```

## Tips

1. Use `prefix + z` to temporarily zoom a pane
2. Save complex layouts in smug configs for repeatability
3. Consider `tmux-resurrect` plugin to persist layouts across restarts
