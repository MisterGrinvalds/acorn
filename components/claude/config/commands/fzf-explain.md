---
description: Explain fzf concepts, configuration, and features
argument-hint: [topic]
allowed-tools: Read, Glob, Grep
---

## Task

Explain the requested topic about fzf. If no specific topic provided, give an overview of fzf and the dotfiles integration.

## Topics

If `$ARGUMENTS` is provided, explain that specific topic:

- **basics** - What fzf is and how fuzzy matching works
- **keybindings** - Shell keybindings (Ctrl+R, Ctrl+T, Alt+C)
- **preview** - Preview window configuration
- **options** - Common command-line options
- **environment** - FZF_DEFAULT_OPTS and other env vars
- **theme** - Catppuccin and color configuration
- **git** - Git integration functions
- **docker** - Docker integration functions
- **kubernetes** - Kubernetes integration functions
- **fd** - Using fd instead of find
- **bat** - Using bat for syntax-highlighted previews

## Context

Reference these files for accurate information:
@components/fzf/component.yaml
@components/fzf/functions.sh
@components/fzf/aliases.sh
@components/fzf/env.sh

## Response Format

1. **Concept overview** - What it is and why it matters
2. **Key options/commands** - Essential flags or syntax
3. **Dotfiles integration** - Available functions and aliases
4. **Examples** - Practical usage examples
5. **Tips** - Best practices and performance considerations
