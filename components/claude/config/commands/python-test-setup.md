---
description: Set up pytest testing for a Python project
allowed-tools: Read, Write, Edit, Bash, Glob
---

## Task

Help the user set up pytest testing infrastructure for their Python project.

## Quick Setup

```bash
# Install pytest with coverage
uv add --dev pytest pytest-cov

# Or using dotfiles
setup_devtools  # Installs pytest, ruff, mypy, pre-commit
```

## Directory Structure

```
project/
├── src/
│   └── mypackage/
│       ├── __init__.py
│       └── calculator.py
└── tests/
    ├── __init__.py
    ├── conftest.py        # Shared fixtures
    ├── test_calculator.py
    └── integration/       # Integration tests
        └── test_api.py
```

## pytest Configuration

Add to `pyproject.toml`:
```toml
[tool.pytest.ini_options]
testpaths = ["tests"]
python_files = ["test_*.py"]
python_functions = ["test_*"]
addopts = [
    "-v",
    "--cov=src",
    "--cov-report=term-missing",
    "--cov-report=html",
]
markers = [
    "slow: marks tests as slow",
    "integration: marks integration tests",
]
```

## Basic Test Example

`tests/test_calculator.py`:
```python
import pytest
from mypackage.calculator import add, divide


def test_add():
    assert add(2, 3) == 5


def test_add_negative():
    assert add(-1, 1) == 0


def test_divide():
    assert divide(10, 2) == 5


def test_divide_by_zero():
    with pytest.raises(ZeroDivisionError):
        divide(10, 0)
```

## Fixtures

`tests/conftest.py`:
```python
import pytest


@pytest.fixture
def sample_data():
    return {"name": "test", "value": 42}


@pytest.fixture
def temp_file(tmp_path):
    file = tmp_path / "test.txt"
    file.write_text("test content")
    return file


@pytest.fixture(scope="session")
def db_connection():
    # Setup
    conn = create_connection()
    yield conn
    # Teardown
    conn.close()
```

## Parametrized Tests

```python
import pytest


@pytest.mark.parametrize("input,expected", [
    (1, 2),
    (2, 4),
    (3, 6),
    (0, 0),
])
def test_double(input, expected):
    assert double(input) == expected


@pytest.mark.parametrize("a,b,expected", [
    (1, 1, 2),
    (0, 0, 0),
    (-1, 1, 0),
])
def test_add_parametrized(a, b, expected):
    assert add(a, b) == expected
```

## Async Tests

```bash
uv add --dev pytest-asyncio
```

```python
import pytest


@pytest.mark.asyncio
async def test_async_function():
    result = await async_fetch_data()
    assert result is not None
```

## Running Tests

```bash
# Run all tests
uv run pytest

# Run with verbose output
uv run pytest -v

# Run specific file
uv run pytest tests/test_calculator.py

# Run specific test
uv run pytest tests/test_calculator.py::test_add

# Run with coverage
uv run pytest --cov=src --cov-report=html

# Run marked tests
uv run pytest -m "not slow"

# Run in parallel (requires pytest-xdist)
uv run pytest -n auto
```

## Coverage Report

```bash
# Terminal report
uv run pytest --cov=src --cov-report=term-missing

# HTML report
uv run pytest --cov=src --cov-report=html
open htmlcov/index.html
```

## CI Integration

GitHub Actions example:
```yaml
- name: Run tests
  run: |
    uv sync
    uv run pytest --cov=src --cov-report=xml

- name: Upload coverage
  uses: codecov/codecov-action@v3
```

## Common Patterns

### Testing Exceptions
```python
def test_raises():
    with pytest.raises(ValueError, match="invalid"):
        validate("bad input")
```

### Testing Logs
```python
def test_logging(caplog):
    my_function()
    assert "expected message" in caplog.text
```

### Mocking
```python
from unittest.mock import patch, MagicMock

def test_with_mock():
    with patch("mypackage.api.requests.get") as mock_get:
        mock_get.return_value.json.return_value = {"data": "test"}
        result = fetch_data()
        assert result == {"data": "test"}
```
