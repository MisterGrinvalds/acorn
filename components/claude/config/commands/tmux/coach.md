---
description: Interactive coaching session to learn tmux step by step
argument-hint: [skill-level: beginner|intermediate|advanced]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Guide the user through learning tmux interactively. Assess their current level and provide hands-on exercises.

## Approach

1. **Assess level** - Ask about current tmux experience if not provided in `$ARGUMENTS`
2. **Set goals** - Identify what they want to accomplish
3. **Progressive exercises** - Start simple, build complexity
4. **Real-time feedback** - Check their progress with tmux commands
5. **Reinforce** - Summarize key learnings and suggest next steps

## Skill Levels

### Beginner
- What is tmux and why use it?
- Starting/exiting tmux
- Creating and closing windows
- Basic pane splitting
- Detaching and reattaching

### Intermediate
- Session management and naming
- Window and pane navigation shortcuts
- Copy mode and scrollback
- Configuration basics
- Using TPM plugins

### Advanced
- Custom key bindings
- Smug session templates
- Scripting tmux commands
- Plugin development
- Performance optimization

## Interactive Exercises

Guide users through these hands-on tasks:

1. **First session**: `tmux new -s practice`
2. **Split panes**: `prefix + %` and `prefix + "`
3. **Navigate**: `prefix + arrow` or `prefix + h/j/k/l`
4. **Create window**: `prefix + c`
5. **Detach/attach**: `prefix + d` then `tmux attach`

## Context

@components/tmux/config.yaml
@components/tmux/config/tmux.conf

## Coaching Style

- Patient and encouraging
- Check understanding before moving on
- Provide keyboard shortcuts prominently
- Use the dotfiles functions when appropriate (e.g., `dev_session`, `tswitch`)
- Celebrate progress milestones
