# Python and UV Management
# UV package manager, virtual environments

.PHONY: uv-install uv-setup uv-status python-status venv-create venv-status

# UV installation
uv-install: ## Install UV package manager
	@if command -v uv >/dev/null 2>&1; then \
		echo "UV already installed: $$(uv --version)"; \
	else \
		curl -LsSf https://astral.sh/uv/install.sh | sh; \
		echo "UV installed. Add to PATH or restart shell."; \
	fi

uv-setup: uv-install ## Complete UV setup with Python
	@export PATH="$$HOME/.local/bin:$$PATH"; \
	if command -v uv >/dev/null 2>&1; then \
		uv python install 3.12 2>/dev/null || echo "Python 3.12 installed or using system"; \
	else \
		echo "UV not found in PATH"; exit 1; \
	fi

# Status
uv-status: ## Check UV and Python status
	@echo "UV/Python Status"
	@echo "================"
	@echo ""
	@echo "UV:"
	@if command -v uv >/dev/null 2>&1; then \
		echo "  Version: $$(uv --version)"; \
		echo "  Path: $$(which uv)"; \
	elif [ -f "$$HOME/.local/bin/uv" ]; then \
		echo "  Version: $$($$HOME/.local/bin/uv --version)"; \
		echo "  Path: $$HOME/.local/bin/uv (not in PATH)"; \
	else \
		echo "  Not installed"; \
	fi
	@echo ""
	@echo "Python:"
	@if command -v python3 >/dev/null 2>&1; then \
		echo "  Version: $$(python3 --version)"; \
		echo "  Path: $$(which python3)"; \
	else \
		echo "  Not installed"; \
	fi
	@echo ""
	@echo "pip:"
	@if command -v pip3 >/dev/null 2>&1; then \
		echo "  Version: $$(pip3 --version | awk '{print $$2}')"; \
	elif python3 -m pip --version >/dev/null 2>&1; then \
		echo "  Version: $$(python3 -m pip --version | awk '{print $$2}')"; \
	else \
		echo "  Not available"; \
	fi
	@echo ""
	@echo "Virtual Environment:"
	@if [ -n "$$VIRTUAL_ENV" ]; then \
		echo "  Active: $$VIRTUAL_ENV"; \
	else \
		echo "  None active"; \
	fi

python-status: uv-status ## Alias for uv-status

# Virtual environments
venv-create: ## Create a virtual environment using UV or venv
	@if command -v uv >/dev/null 2>&1; then \
		uv venv .venv; \
	else \
		python3 -m venv .venv; \
	fi
	@echo "Created .venv - activate with: source .venv/bin/activate"

venv-status: ## Show active virtual environment info
	@echo "Virtual Environment Status"
	@echo "=========================="
	@if [ -n "$$VIRTUAL_ENV" ]; then \
		echo "Active: $$VIRTUAL_ENV"; \
		echo "Python: $$(python --version)"; \
		command -v uv >/dev/null && echo "UV: $$(uv --version)" || true; \
		echo ""; \
		echo "Installed packages:"; \
		pip list 2>/dev/null | head -15 || echo "  (none)"; \
	else \
		echo "No virtual environment active"; \
		echo ""; \
		echo "Create: make venv-create"; \
		echo "Activate: source .venv/bin/activate"; \
	fi
