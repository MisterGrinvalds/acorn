---
description: Explain Git concepts, commands, and workflows
argument-hint: [topic]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Explain the requested topic about Git. If no specific topic provided, give an overview of Git and the dotfiles integration.

## Topics

If `$ARGUMENTS` is provided, explain that specific topic:

- **basics** - Git fundamentals (init, add, commit, push, pull)
- **branches** - Branch creation, switching, merging
- **rebase** - Rebasing vs merging, interactive rebase
- **merge** - Merge strategies and fast-forward
- **conflicts** - Understanding and resolving conflicts
- **stash** - Stashing changes for later
- **reset** - Reset modes (soft, mixed, hard)
- **reflog** - Reference log for recovery
- **cherry-pick** - Applying specific commits
- **bisect** - Binary search for bugs
- **hooks** - Git hooks and automation
- **config** - Git configuration levels
- **submodules** - Working with submodules
- **worktrees** - Multiple working directories
- **internals** - Objects, refs, and packfiles

## Context

Reference these files for accurate information:
@components/git/component.yaml
@components/git/functions.sh
@components/git/aliases.sh
@components/git/config

## Response Format

1. **Concept overview** - What it is and why it matters
2. **Key commands** - Essential Git commands for this topic
3. **Dotfiles integration** - Available aliases and functions
4. **Examples** - Practical usage examples
5. **Common pitfalls** - Mistakes to avoid
