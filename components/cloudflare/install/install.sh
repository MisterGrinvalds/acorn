#!/usr/bin/env bash
# CloudFlare Wrangler CLI installation script

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

info() { echo -e "${GREEN}[INFO]${NC} $*"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*" >&2; }

# Check if wrangler is already installed
if command -v wrangler &>/dev/null; then
    info "Wrangler is already installed: $(wrangler --version)"
    exit 0
fi

# Check for npm
if ! command -v npm &>/dev/null; then
    error "npm is required to install wrangler"
    error "Please install Node.js first: https://nodejs.org/"
    exit 1
fi

info "Installing wrangler globally via npm..."
npm install -g wrangler

if command -v wrangler &>/dev/null; then
    info "Wrangler installed successfully: $(wrangler --version)"
    info "Run 'wrangler login' to authenticate with CloudFlare"
else
    error "Wrangler installation failed"
    exit 1
fi
