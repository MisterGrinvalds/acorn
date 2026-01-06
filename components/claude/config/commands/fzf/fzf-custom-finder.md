---
description: Create a custom fzf finder function
argument-hint: <finder-name> [data-source]
allowed-tools: Read, Write, Edit, Bash
---

## Task

Help the user create a custom fzf finder function for their specific workflow.

## Process

1. **Understand the goal**: What do they want to find/select?
2. **Identify data source**: Files, processes, git objects, API output, etc.
3. **Design the finder**: Command pipeline, preview, actions
4. **Implement**: Write the function
5. **Add alias**: Create short alias for the function

## Finder Template

```bash
fzf_<name>() {
    local selection

    # Generate list | filter with fzf | capture selection
    selection=$(<data-source-command> | fzf \
        --prompt="<Prompt>: " \
        --preview '<preview-command> {}' \
        --preview-window=right:50% \
        --height=80% \
        --border rounded)

    # Act on selection
    [ -n "$selection" ] && <action-command> "$selection"
}

# Alias
alias f<short>='fzf_<name>'
```

## Example Finders

### NPM Scripts Finder
```bash
fzf_npm_script() {
    local script
    script=$(jq -r '.scripts | keys[]' package.json 2>/dev/null | fzf \
        --prompt="Run script: " \
        --preview 'jq -r ".scripts.{}" package.json')
    [ -n "$script" ] && npm run "$script"
}
alias fnpm='fzf_npm_script'
```

### SSH Host Finder
```bash
fzf_ssh() {
    local host
    host=$(grep "^Host " ~/.ssh/config | awk '{print $2}' | fzf \
        --prompt="SSH to: " \
        --preview 'grep -A5 "^Host {}" ~/.ssh/config')
    [ -n "$host" ] && ssh "$host"
}
alias fssh='fzf_ssh'
```

### Brew Package Finder
```bash
fzf_brew_install() {
    local pkg
    pkg=$(brew search --formulae "" | fzf \
        --prompt="Install: " \
        --preview 'brew info {}')
    [ -n "$pkg" ] && brew install "$pkg"
}
alias fbrew='fzf_brew_install'
```

### Project Directory Finder
```bash
fzf_project() {
    local dir
    dir=$(fd --type d --max-depth 2 . ~/projects ~/work | fzf \
        --prompt="Project: " \
        --preview 'ls -la {}')
    [ -n "$dir" ] && cd "$dir"
}
alias fp='fzf_project'
```

## Key Components

### Data Sources
- `fd` / `find` - Files and directories
- `git` - Branches, commits, stashes
- `ps aux` - Processes
- `kubectl get` - Kubernetes resources
- `jq` - JSON data
- API calls - External services

### Preview Options
```bash
--preview 'cat {}'           # File content
--preview 'head -100 {}'     # First 100 lines
--preview 'bat --color=always {}'  # Syntax highlighted
--preview 'ls -la {}'        # Directory listing
--preview 'git show {}'      # Git object
```

### Actions
```bash
# Simple action
[ -n "$selection" ] && command "$selection"

# Multiple selections
echo "$selection" | xargs command

# With confirmation
[ -n "$selection" ] && read -p "Confirm? " && command "$selection"
```

## Context

@components/fzf/functions.sh

## Output

After designing, add the function to `components/fzf/functions.sh` and alias to `components/fzf/aliases.sh`.
