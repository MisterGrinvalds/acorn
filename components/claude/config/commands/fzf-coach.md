---
description: Interactive coaching session to learn fzf step by step
argument-hint: [skill-level: beginner|intermediate|advanced]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Guide the user through learning fzf interactively. Assess their current level and provide hands-on exercises.

## Approach

1. **Assess level** - Ask about current fzf experience if not provided
2. **Set goals** - Identify what they want to find/filter
3. **Progressive exercises** - Start simple, build complexity
4. **Real-time feedback** - Have them run commands and report results
5. **Reinforce** - Summarize key learnings

## Skill Levels

### Beginner
- What is fuzzy finding?
- Basic usage: `ls | fzf`
- Search syntax and matching
- Selecting items
- Using Ctrl+R for history

### Intermediate
- Multi-select with Tab
- Preview windows
- Shell keybindings (Ctrl+T, Alt+C)
- Using dotfiles functions (ff, fcd, fgb)
- FZF_DEFAULT_OPTS

### Advanced
- Custom finders with preview
- Integration with fd and bat
- Writing reusable functions
- Action bindings (--bind)
- Performance optimization

## Interactive Exercises

### Beginner Exercises
```bash
# Exercise 1: Basic filtering
ls | fzf

# Exercise 2: History search
# Press Ctrl+R and type part of a command

# Exercise 3: File finding
# Press Ctrl+T to find files
```

### Intermediate Exercises
```bash
# Exercise 4: Multi-select
ls | fzf -m
# Press Tab to select multiple, Enter to confirm

# Exercise 5: Preview
fzf --preview 'head -20 {}'

# Exercise 6: Git branches
fgb  # or fzf_git_branch
```

### Advanced Exercises
```bash
# Exercise 7: Custom finder
fd --type f | fzf --preview 'bat --color=always {}'

# Exercise 8: Action binding
fzf --bind 'ctrl-y:execute(echo {} | pbcopy)'
```

## Context

@components/fzf/functions.sh
@components/fzf/aliases.sh

## Coaching Style

- Start with simple pipes: `command | fzf`
- Build up to complex finders
- Use dotfiles functions as examples
- Encourage experimentation
- Celebrate when they "get it"
