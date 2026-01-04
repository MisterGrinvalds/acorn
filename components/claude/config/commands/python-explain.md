---
description: Explain Python concepts, tools, and workflows
argument-hint: [topic]
allowed-tools: Read, Glob, Grep
---

## Task

Explain the requested topic about Python development. If no specific topic provided, give an overview of Python tooling and the dotfiles integration.

## Topics

If `$ARGUMENTS` is provided, explain that specific topic:

- **venv** - Virtual environments and isolation
- **uv** - UV package manager basics
- **pip** - pip package management
- **pyproject** - pyproject.toml configuration
- **ruff** - Linting and formatting with ruff
- **mypy** - Type checking with mypy
- **pytest** - Testing with pytest
- **typing** - Type hints and annotations
- **packaging** - Building and distributing packages
- **asyncio** - Async programming basics
- **fastapi** - FastAPI web framework
- **ipython** - Interactive Python shell
- **poetry** - Poetry vs UV comparison
- **conda** - Conda vs venv/UV

## Context

Reference these files for accurate information:
@components/python/component.yaml
@components/python/functions.sh
@components/python/aliases.sh

## Response Format

1. **Concept overview** - What it is and why it matters
2. **Key commands** - Essential commands for this topic
3. **Dotfiles integration** - Available functions and aliases
4. **Examples** - Practical usage examples
5. **Best practices** - Modern recommendations
