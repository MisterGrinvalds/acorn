---
description: Migrate a Python project from pip/poetry to UV
argument-hint: [from: pip|poetry|requirements]
allowed-tools: Read, Write, Edit, Bash, Glob
---

## Task

Help the user migrate their Python project to use UV package manager.

## Migration Paths

Based on `$ARGUMENTS`:

### requirements (from requirements.txt)

```bash
# Step 1: Initialize UV project
uv init

# Step 2: Add dependencies from requirements.txt
# UV can read requirements.txt directly
uv add $(cat requirements.txt | grep -v '^#' | grep -v '^$' | tr '\n' ' ')

# Or use UV's pip interface
uv pip install -r requirements.txt

# Step 3: Generate lock file
uv lock

# Step 4: (Optional) Remove old files
# Keep requirements.txt for reference or delete
```

### pip (from setup.py/setup.cfg)

```bash
# Step 1: Create pyproject.toml
uv init

# Step 2: Copy dependencies from setup.py
# Extract install_requires and extras_require

# Step 3: Add to pyproject.toml
uv add <dependencies>

# Step 4: Add dev dependencies
uv add --dev pytest ruff mypy
```

### poetry (from Poetry)

```bash
# Poetry uses pyproject.toml, but different format

# Step 1: Export requirements
poetry export -f requirements.txt --output requirements.txt

# Step 2: Initialize UV
uv init

# Step 3: Add dependencies
uv add $(cat requirements.txt | cut -d';' -f1 | tr '\n' ' ')

# Step 4: Remove Poetry files
rm poetry.lock
# Edit pyproject.toml to remove [tool.poetry] section
```

## pyproject.toml Conversion

### From Poetry format:
```toml
# Poetry style
[tool.poetry.dependencies]
python = "^3.11"
fastapi = "^0.100.0"
```

### To UV/PEP 621 format:
```toml
# Standard format (UV compatible)
[project]
name = "myproject"
version = "0.1.0"
requires-python = ">=3.11"
dependencies = [
    "fastapi>=0.100.0",
]

[project.optional-dependencies]
dev = [
    "pytest>=7.0",
    "ruff>=0.1.0",
]
```

## Post-Migration Checklist

```bash
# 1. Verify dependencies install
uv sync

# 2. Run tests
uv run pytest

# 3. Check imports work
uv run python -c "import mypackage"

# 4. Update CI/CD scripts
# Replace: pip install -r requirements.txt
# With: uv sync

# 5. Update documentation
```

## Common Issues

### Version conflicts
```bash
# UV will show conflicts during sync
uv sync
# Adjust versions in pyproject.toml as needed
```

### Missing optional dependencies
```bash
# Add to optional-dependencies
uv add --dev <package>
```

### Scripts/entry points
```toml
[project.scripts]
mycommand = "mypackage.cli:main"
```

## Benefits After Migration

1. **Speed**: UV is 10-100x faster than pip
2. **Locking**: uv.lock ensures reproducible builds
3. **Standards**: PEP 621 compliant pyproject.toml
4. **Simplicity**: One tool for venv + packages
