#!/bin/bash

# Script to create branch and commits for bash-profile repository

# Change to the repository directory
cd /Users/mistergrinvalds/Repos/personal/bash-profile

# Check current status
echo "Current git status:"
git status

# Create and checkout new branch
echo -e "\nCreating branch fix/automation-framework-compatibility..."
git checkout -b fix/automation-framework-compatibility

# Stage and commit changes for automation framework core.sh
if [ -f .automation/framework/core.sh ]; then
    echo -e "\nCommit 1: Bash compatibility fix for automation framework"
    git add .automation/framework/core.sh
    git commit -m "fix: update automation framework core.sh for bash compatibility

- Replace 'function' keyword with POSIX-compliant syntax
- Ensure compatibility with both bash and sh shells"
fi

# Stage and commit changes for automation auto script
if [ -f .automation/auto ]; then
    echo -e "\nCommit 2: Fix associative array compatibility"
    git add .automation/auto
    git commit -m "fix: replace associative arrays in auto script for POSIX compliance

- Convert bash-specific associative arrays to portable implementation
- Maintain functionality while ensuring shell compatibility"
fi

# Create missing .bash_profile.dir structure if needed
if [ ! -d .bash_profile.dir ]; then
    echo -e "\nCommit 3: Create bash profile directory structure"
    mkdir -p .bash_profile.dir
    touch .bash_profile.dir/.gitkeep
    git add .bash_profile.dir/
    git commit -m "feat: add missing .bash_profile.dir directory structure

- Create required directory for modular bash configuration
- Add .gitkeep to preserve empty directory in git"
fi

# Stage and commit initialize.sh changes
if [ -f initialize.sh ]; then
    echo -e "\nCommit 4: Add yes-to-all functionality to initialize script"
    git add initialize.sh
    git commit -m "feat: add --yes flag to initialize.sh for automated installation

- Support non-interactive installation with --yes flag
- Maintain backward compatibility with interactive mode"
fi

# Stage and commit tools.sh changes
if [ -f .automation/modules/tools.sh ]; then
    echo -e "\nCommit 5: Update tools module to support yes-to-all mode"
    git add .automation/modules/tools.sh
    git commit -m "feat: add yes-to-all support to automation tools module

- Pass through --yes flag from parent scripts
- Enable fully automated tool installation"
fi

# Stage and commit Makefile changes
if [ -f Makefile ]; then
    echo -e "\nCommit 6: Add new make targets for automation"
    git add Makefile
    git commit -m "feat: enhance Makefile with new automation targets

- Add targets for automated installation and testing
- Improve development workflow with convenience commands"
fi

# Stage any remaining project config files
echo -e "\nChecking for additional project config files..."
if [ -f CLAUDE.md ] || [ -f package.json ] || [ -f package-lock.json ]; then
    git add CLAUDE.md package.json package-lock.json 2>/dev/null
    if git diff --cached --quiet; then
        echo "No additional changes to commit"
    else
        git commit -m "chore: update project configuration files

- Update CLAUDE.md documentation
- Add npm package configuration for development tools"
    fi
fi

# Show final status
echo -e "\nFinal status:"
git status
echo -e "\nBranch created and commits made successfully!"
echo "You can now push with: git push -u origin fix/automation-framework-compatibility"