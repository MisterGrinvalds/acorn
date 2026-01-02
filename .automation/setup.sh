#!/usr/bin/env bash
# Automation Framework Setup Script

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🚀 Setting up Automation Framework..."

# Make all scripts executable
echo "📝 Making scripts executable..."
chmod +x "$SCRIPT_DIR/auto"
chmod +x "$SCRIPT_DIR/framework"/*.sh
chmod +x "$SCRIPT_DIR/modules"/*.sh

# Create required directories
echo "📁 Creating directories..."
mkdir -p "$SCRIPT_DIR/config"
mkdir -p "$SCRIPT_DIR/logs"
mkdir -p "$SCRIPT_DIR/cache"
mkdir -p "$SCRIPT_DIR/k8s/manifests"
mkdir -p "$SCRIPT_DIR/github/templates"
mkdir -p "$SCRIPT_DIR/config/profiles"
mkdir -p "$SCRIPT_DIR/config/templates"
mkdir -p "$SCRIPT_DIR/config/environments"

# Initialize configuration
echo "⚙️  Initializing configuration..."
if [ ! -f "$SCRIPT_DIR/config/automation.conf" ]; then
    cat > "$SCRIPT_DIR/config/automation.conf" << 'EOF'
# Automation Framework Configuration
AUTO_LOG_LEVEL=INFO
AUTO_PARALLEL_JOBS=4
AUTO_TIMEOUT=300
AUTO_RETRY_COUNT=3
AUTO_INIT_ON_STARTUP=true

# Default project directory
DEV_PROJECTS_DIR=$HOME/projects

# GitHub settings
GITHUB_DEFAULT_VISIBILITY=public

# Kubernetes settings
K8S_DEFAULT_NAMESPACE=default

# System settings
SYSTEM_BACKUP_RETENTION_DAYS=30
EOF
fi

# Test the CLI
echo "🧪 Testing automation CLI..."
if "$SCRIPT_DIR/auto" --version >/dev/null 2>&1; then
    echo "✅ Automation CLI is working"
else
    echo "❌ Automation CLI test failed"
    exit 1
fi

# Add to shell profile if not already present
SHELL_PROFILE="$HOME/.bash_profile"
if [ -f "$HOME/.zshrc" ] && [ ! -f "$SHELL_PROFILE" ]; then
    SHELL_PROFILE="$HOME/.zshrc"
fi

AUTO_PATH_EXPORT="export PATH=\"$SCRIPT_DIR:\$PATH\""
if ! grep -q "$SCRIPT_DIR" "$SHELL_PROFILE" 2>/dev/null; then
    echo "🔗 Adding automation to PATH in $SHELL_PROFILE"
    echo "" >> "$SHELL_PROFILE"
    echo "# Automation Framework" >> "$SHELL_PROFILE"
    echo "$AUTO_PATH_EXPORT" >> "$SHELL_PROFILE"
fi

# Create initial project directory
mkdir -p "$HOME/projects"

echo ""
echo "✅ Automation Framework setup complete!"
echo ""
echo "🎯 Quick start:"
echo "   auto --help                    # Show all commands"
echo "   auto dev init python my-api    # Create new Python project"  
echo "   auto k8s cluster info          # Check Kubernetes cluster"
echo "   auto github repo create my-app # Create GitHub repository"
echo "   auto system setup              # Setup development environment"
echo ""
echo "☁️  Cloud automation:"
echo "   auto cloud status              # Check all cloud providers"
echo "   auto aws ec2 list              # List AWS EC2 instances"
echo "   auto azure vm list             # List Azure VMs"
echo "   auto digitalocean droplets list # List DigitalOcean droplets"
echo ""
echo "🔄 Reload your shell or run: source $SHELL_PROFILE"
echo ""
echo "📚 For detailed help on any module:"
echo "   auto <module> --help"
echo ""