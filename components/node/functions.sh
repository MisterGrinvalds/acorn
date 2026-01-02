#!/bin/sh
# components/node/functions.sh - Node.js development functions

# =============================================================================
# NVM Management
# =============================================================================

# Install and use latest LTS Node version
nvm_latest() {
    if ! command -v nvm >/dev/null 2>&1; then
        echo "NVM not installed"
        return 1
    fi
    nvm install --lts
    nvm use --lts
}

# Install NVM if not present
nvm_setup() {
    if command -v nvm >/dev/null 2>&1; then
        echo "NVM already installed"
        nvm --version
        return 0
    fi

    echo "Installing NVM..."
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

    # Reload NVM
    export NVM_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/nvm"
    [ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"

    echo "NVM installed. Installing latest LTS Node..."
    nvm install --lts
}

# =============================================================================
# Project Initialization
# =============================================================================

# Create new Node.js project with TypeScript
node_init() {
    local name="${1:-.}"

    if [ "$name" != "." ]; then
        mkdir -p "$name"
        cd "$name" || return 1
    fi

    echo "Initializing Node.js project..."

    # Initialize package.json
    if command -v pnpm >/dev/null 2>&1; then
        pnpm init
        pnpm add -D typescript @types/node tsx
    else
        npm init -y
        npm install --save-dev typescript @types/node tsx
    fi

    # Create tsconfig.json
    cat > tsconfig.json << 'EOF'
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "outDir": "./dist"
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules"]
}
EOF

    mkdir -p src
    echo 'console.log("Hello, TypeScript!");' > src/index.ts

    echo "Node.js TypeScript project initialized!"
}

# =============================================================================
# Utility Functions
# =============================================================================

# Clean node_modules and reinstall
nclean() {
    echo "Removing node_modules..."
    rm -rf node_modules

    echo "Reinstalling dependencies..."
    if [ -f "pnpm-lock.yaml" ]; then
        pnpm install
    elif [ -f "yarn.lock" ]; then
        yarn install
    else
        npm install
    fi
}

# Find and list all node_modules in current directory tree
nfind() {
    find . -name "node_modules" -type d -prune | while read -r dir; do
        du -sh "$dir"
    done
}

# Remove all node_modules in current directory tree
ncleanall() {
    echo "Finding all node_modules directories..."
    find . -name "node_modules" -type d -prune -print

    printf "Remove all? [y/N] "
    read -r confirm
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        find . -name "node_modules" -type d -prune -exec rm -rf {} +
        echo "All node_modules removed"
    fi
}

# Check which package manager to use based on lock files
npm_detect() {
    if [ -f "pnpm-lock.yaml" ]; then
        echo "pnpm"
    elif [ -f "yarn.lock" ]; then
        echo "yarn"
    elif [ -f "package-lock.json" ]; then
        echo "npm"
    else
        echo "npm"
    fi
}
