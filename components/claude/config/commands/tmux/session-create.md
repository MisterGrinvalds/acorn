---
description: Create a custom tmux session with windows and panes
argument-hint: <session-name> [project-path]
allowed-tools: Read, Write, Edit, Bash
---

## Task

Help the user create a custom tmux session. Can be:
1. **Interactive** - Create a one-time session with specific layout
2. **Template** - Create a reusable shell function
3. **Smug config** - Create a persistent YAML session definition

## Process

1. **Gather requirements**:
   - Session name from `$1` or ask
   - Number of windows needed
   - Pane layout for each window
   - Commands to run in each pane
   - Working directory

2. **Choose approach**:
   - Quick one-off: Direct tmux commands
   - Reusable: Shell function in dotfiles
   - Persistent: Smug YAML config

3. **Create the session** with appropriate method

## Session Layouts

Common layouts to offer:
- **Dev**: Editor | Terminal + Logs (3 panes)
- **Full-stack**: Frontend | Backend | Database | Logs (4 panes)
- **K8s**: kubectl | logs (split) | k9s
- **Simple**: Single window, single pane

## Example Outputs

### Direct Command
```bash
tmux new-session -d -s myproject -c ~/projects/myproject
tmux split-window -h -t myproject
tmux send-keys -t myproject:0.0 'nvim .' Enter
tmux attach -t myproject
```

### Shell Function
Add to `components/tmux/functions.sh`:
```bash
myproject_session() {
    tmux new-session -d -s myproject -c ~/projects/myproject
    # ... layout commands
    tmux attach-session -t myproject
}
```

### Smug Config
Create `~/.config/smug/myproject.yml`:
```yaml
session: myproject
root: ~/projects/myproject
windows:
  - name: code
    commands: [nvim .]
  - name: terminal
```

## Context

@components/tmux/config.yaml

## Reference

Use `project_session` as a model - it auto-detects project type and creates appropriate windows.
