---
description: Explain tmux concepts, configuration, and features
argument-hint: [topic]
allowed-tools: Read, Glob, Grep
---

## Task

Explain the requested topic about tmux. If no specific topic provided, give an overview of tmux and the dotfiles integration.

## Topics

If `$ARGUMENTS` is provided, explain that specific topic. Common topics include:

- **sessions** - Session management, naming, attaching/detaching
- **windows** - Window creation, navigation, naming
- **panes** - Pane splitting, resizing, navigation
- **prefix** - Prefix key and key bindings
- **config** - Configuration file structure and options
- **tpm** - Tmux Plugin Manager setup and usage
- **plugins** - Popular plugins and their purpose
- **smug** - Smug session manager and YAML configs
- **copy-mode** - Copy/paste and scrollback buffer
- **status-bar** - Status line customization
- **themes** - Catppuccin and other themes

## Context

Reference these files for accurate information:
@components/tmux/component.yaml
@components/tmux/config.yaml
@components/tmux/config/tmux.conf

## Response Format

1. **Concept overview** - What it is and why it matters
2. **Key commands/keys** - Essential shortcuts or commands
3. **Dotfiles integration** - Available functions and aliases
4. **Examples** - Practical usage examples
5. **Tips** - Best practices and common pitfalls
