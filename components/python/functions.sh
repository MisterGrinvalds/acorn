#!/bin/sh
# components/python/functions.sh - Python development functions
# UV: https://github.com/astral-sh/uv

# =============================================================================
# Virtual Environment Management
# =============================================================================

# Create a Python virtual environment using UV (with fallback)
mkvenv() {
    local name="${1:-.venv}"

    if ! command -v uv >/dev/null 2>&1; then
        echo "UV not found. Falling back to python3 -m venv"
        python3 -m venv "$name"
    else
        uv venv "$name"
    fi

    # Activate the environment
    if [ -f "$name/bin/activate" ]; then
        . "$name/bin/activate"
        echo "Created and activated: $name"
    else
        echo "Created: $name"
    fi
}

# Activate a Python virtual environment
venv() {
    local name="${1:-.venv}"

    if [ -f "$name/bin/activate" ]; then
        . "$name/bin/activate"
    elif [ -n "$ENVS_LOCATION" ] && [ -f "$ENVS_LOCATION/$name/bin/activate" ]; then
        . "$ENVS_LOCATION/$name/bin/activate"
    else
        echo "Virtual environment not found: $name"
        return 1
    fi
}

# Deactivate current virtual environment
dvenv() {
    if [ -n "$VIRTUAL_ENV" ]; then
        deactivate
    else
        echo "No active virtual environment"
    fi
}

# =============================================================================
# UV Convenience Functions
# =============================================================================

# Sync dependencies from pyproject.toml/uv.lock
uv_sync() {
    if ! command -v uv >/dev/null 2>&1; then
        echo "UV not installed. Install: curl -LsSf https://astral.sh/uv/install.sh | sh"
        return 1
    fi
    uv sync "$@"
}

# Add a package
uv_add() {
    if ! command -v uv >/dev/null 2>&1; then
        echo "UV not installed"
        return 1
    fi
    uv add "$@"
}

# Remove a package
uv_remove() {
    if ! command -v uv >/dev/null 2>&1; then
        echo "UV not installed"
        return 1
    fi
    uv remove "$@"
}

# Run a command in the project environment
uv_run() {
    if ! command -v uv >/dev/null 2>&1; then
        echo "UV not installed"
        return 1
    fi
    uv run "$@"
}

# Initialize a new Python project with UV
uv_init() {
    if ! command -v uv >/dev/null 2>&1; then
        echo "UV not installed"
        return 1
    fi
    uv init "$@"
}

# =============================================================================
# Development Environment Setup
# =============================================================================

# FastAPI development environment setup
fastapi_env() {
    local env_name="${1:-.venv}"

    # Create and activate environment
    mkvenv "$env_name"

    echo "Installing FastAPI development dependencies..."
    if command -v uv >/dev/null 2>&1; then
        uv pip install fastapi uvicorn python-multipart
        uv pip install pytest httpx pytest-asyncio
        uv pip install ruff
        uv pip install python-dotenv
    else
        pip install --upgrade pip
        pip install fastapi uvicorn python-multipart
        pip install pytest httpx pytest-asyncio
        pip install ruff
        pip install python-dotenv
    fi

    echo "FastAPI environment ready!"
    echo "Run: uvicorn main:app --reload"
}

# Lightweight IPython setup
setup_ipython() {
    if command -v uv >/dev/null 2>&1; then
        uv pip install ipython rich
    else
        pip install ipython rich
    fi
    echo "IPython installed with rich output"
}

# Install common development tools
setup_devtools() {
    local tools="ruff mypy pytest pytest-cov pre-commit"

    if command -v uv >/dev/null 2>&1; then
        echo "Installing dev tools with UV..."
        uv pip install $tools
    else
        echo "Installing dev tools with pip..."
        pip install $tools
    fi
    echo "Development tools installed: $tools"
}
