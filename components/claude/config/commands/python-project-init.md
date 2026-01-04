---
description: Initialize a new Python project with modern tooling
argument-hint: <project-name>
allowed-tools: Read, Write, Bash
---

## Task

Help the user initialize a new Python project with UV, modern tooling, and best practices.

## Quick Start

```bash
# Using UV (recommended)
uv init <project-name>
cd <project-name>

# Using dotfiles alias
uvi <project-name>
```

## Full Project Setup

### Step 1: Initialize
```bash
mkdir <project-name>
cd <project-name>
uv init
```

### Step 2: Add Dependencies
```bash
# Main dependencies
uv add <your-dependencies>

# Development dependencies
uv add --dev pytest ruff mypy pre-commit
```

### Step 3: Create Structure
```bash
mkdir -p src/<project_name> tests
touch src/<project_name>/__init__.py
touch src/<project_name>/main.py
touch tests/__init__.py
touch tests/test_main.py
```

## Recommended pyproject.toml

```toml
[project]
name = "project-name"
version = "0.1.0"
description = "Project description"
requires-python = ">=3.11"
dependencies = []

[project.optional-dependencies]
dev = [
    "pytest>=7.0",
    "pytest-cov>=4.0",
    "ruff>=0.1.0",
    "mypy>=1.0",
    "pre-commit>=3.0",
]

[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"

[tool.ruff]
line-length = 88
select = ["E", "F", "I", "UP", "B", "SIM"]
ignore = ["E501"]

[tool.ruff.isort]
known-first-party = ["project_name"]

[tool.mypy]
python_version = "3.11"
strict = true
warn_return_any = true

[tool.pytest.ini_options]
testpaths = ["tests"]
addopts = "-v --cov=src"
```

## Project Structure

```
project-name/
├── pyproject.toml
├── uv.lock
├── README.md
├── .gitignore
├── .pre-commit-config.yaml
├── src/
│   └── project_name/
│       ├── __init__.py
│       ├── main.py
│       └── py.typed          # For type hints
└── tests/
    ├── __init__.py
    ├── conftest.py
    └── test_main.py
```

## Pre-commit Configuration

Create `.pre-commit-config.yaml`:
```yaml
repos:
  - repo: https://github.com/astral-sh/ruff-pre-commit
    rev: v0.1.6
    hooks:
      - id: ruff
        args: [--fix]
      - id: ruff-format

  - repo: https://github.com/pre-commit/mirrors-mypy
    rev: v1.7.1
    hooks:
      - id: mypy
        additional_dependencies: []
```

Install hooks:
```bash
uv run pre-commit install
```

## .gitignore Essentials

```gitignore
# Virtual environment
.venv/
venv/

# Python
__pycache__/
*.py[cod]
*.egg-info/
dist/
build/

# IDE
.idea/
.vscode/
*.swp

# Testing
.coverage
htmlcov/
.pytest_cache/

# Type checking
.mypy_cache/

# UV
uv.lock  # Include or exclude based on preference
```

## Verification

```bash
# Install dependencies
uv sync

# Run tests
uv run pytest

# Check linting
uv run ruff check .

# Check types
uv run mypy src/
```

## Project Types

### CLI Application
```bash
uv add typer rich
```

### Web API (FastAPI)
```bash
fastapi_env  # dotfiles function
# or: uv add fastapi uvicorn
```

### Data Science
```bash
uv add pandas numpy matplotlib jupyter
```
