# Enhanced Python virtual environment management
mkvenv() { 
    if [ -z "$1" ]; then
        # Create local .venv in current directory
        python3 -m venv .venv
        echo "Created .venv in current directory"
    else
        # Create named environment in ENVS_LOCATION
        export VENV="$1"
        python3 -m venv "$ENVS_LOCATION/$VENV"
        echo "Created virtual environment: $VENV"
    fi
}

# FastAPI development environment setup
fastapi_env() {
    local env_name="${1:-.venv}"
    
    if [ "$env_name" = ".venv" ]; then
        python3 -m venv .venv
        source .venv/bin/activate
    else
        python3 -m venv "$ENVS_LOCATION/$env_name"
        source "$ENVS_LOCATION/$env_name/bin/activate"
    fi
    
    echo "Installing FastAPI development dependencies..."
    pip install --upgrade pip
    pip install fastapi uvicorn python-multipart
    pip install pytest httpx pytest-asyncio  # Testing
    pip install black isort flake8          # Code formatting/linting
    pip install python-dotenv               # Environment variables
    
    echo "FastAPI environment ready!"
    echo "Run: uvicorn main:app --reload"
}

# Lightweight IPython setup
setup_ipython() {
    pip install ipython rich
    echo "Lightweight IPython installed with rich output"
}
