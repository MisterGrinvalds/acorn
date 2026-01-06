---
description: Interactive coaching session to learn Python development workflows
argument-hint: [skill-level: beginner|intermediate|advanced]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Guide the user through learning Python development workflows interactively.

## Approach

1. **Assess level** - Ask about Python experience
2. **Set goals** - Identify what they want to build
3. **Progressive exercises** - Start simple, build complexity
4. **Real-time practice** - Have them run commands
5. **Reinforce** - Summarize and suggest next steps

## Skill Levels

### Beginner
- Python basics (variables, functions, classes)
- Virtual environments - why and how
- Installing packages
- Running scripts
- Basic project structure

### Intermediate
- UV package manager
- pyproject.toml configuration
- Testing with pytest
- Linting with ruff
- Type hints basics
- Using the dotfiles functions

### Advanced
- Type checking with mypy
- Advanced pytest (fixtures, parametrize)
- Async/await patterns
- Package publishing
- CI/CD integration
- Performance optimization

## Interactive Exercises

### Beginner Exercises
```bash
# Exercise 1: Create virtual environment
mkvenv my-first-project
# or: mkv my-first-project

# Exercise 2: Check activation
which python
echo $VIRTUAL_ENV

# Exercise 3: Install a package
pip install requests
# or with UV: uv pip install requests

# Exercise 4: Deactivate
dvenv
# or: dv
```

### Intermediate Exercises
```bash
# Exercise 5: Initialize UV project
uv init my-project
cd my-project

# Exercise 6: Add dependencies
uv add fastapi uvicorn

# Exercise 7: Add dev dependencies
uv add --dev pytest ruff

# Exercise 8: Run with UV
uv run python -c "import fastapi; print(fastapi.__version__)"
```

### Advanced Exercises
```bash
# Exercise 9: Setup full dev environment
fastapi_env

# Exercise 10: Configure ruff in pyproject.toml
# Add [tool.ruff] section

# Exercise 11: Add type hints and run mypy
# Exercise 12: Write pytest tests
```

## Context

@components/python/functions.sh
@components/python/aliases.sh

## Coaching Style

- Emphasize UV for modern workflows
- Use dotfiles functions and aliases
- Explain the "why" behind virtual environments
- Build toward real project structure
- Celebrate working code
