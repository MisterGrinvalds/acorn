---
description: Explain GitHub CLI concepts, commands, and workflows
argument-hint: [topic]
allowed-tools: Read, Glob, Grep
---

## Task

Explain the requested topic about GitHub CLI and workflows. If no specific topic provided, give an overview of gh CLI and the dotfiles integration.

## Topics

If `$ARGUMENTS` is provided, explain that specific topic:

- **auth** - Authentication and tokens
- **pr** - Pull request workflows
- **issues** - Issue management
- **actions** - GitHub Actions and workflows
- **releases** - Release management
- **repo** - Repository operations
- **api** - Using gh api for custom queries
- **gist** - Gist creation and management
- **codespaces** - GitHub Codespaces
- **secrets** - Managing repository secrets
- **labels** - Label management
- **milestones** - Milestone tracking
- **projects** - GitHub Projects

## Context

Reference these files for accurate information:
@components/github/component.yaml
@components/github/functions.sh
@components/github/aliases.sh

## Response Format

1. **Concept overview** - What it is and why it matters
2. **Key commands** - Essential gh commands
3. **Dotfiles integration** - Available functions and aliases
4. **Examples** - Practical usage examples
5. **Best practices** - Workflow recommendations
