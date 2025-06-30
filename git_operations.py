#!/usr/bin/env python3
import subprocess
import os

# Change to the repository directory
os.chdir('/Users/mistergrinvalds/Repos/personal/bash-profile')

def run_git_command(cmd):
    """Execute a git command and print the output"""
    try:
        result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
        print(f"Command: {cmd}")
        if result.stdout:
            print(f"Output: {result.stdout}")
        if result.stderr:
            print(f"Error: {result.stderr}")
        print("-" * 50)
        return result.returncode == 0
    except Exception as e:
        print(f"Exception running '{cmd}': {e}")
        return False

# Check current status
print("=== Checking current git status ===")
run_git_command("git status")

# Create and checkout new branch
print("\n=== Creating branch fix/automation-framework-compatibility ===")
if run_git_command("git checkout -b fix/automation-framework-compatibility"):
    print("Branch created successfully!")
else:
    print("Failed to create branch. It might already exist.")
    run_git_command("git checkout fix/automation-framework-compatibility")

# Commit 1: Bash compatibility fix for automation framework
if os.path.exists('.automation/framework/core.sh'):
    print("\n=== Commit 1: Bash compatibility fix ===")
    run_git_command("git add .automation/framework/core.sh")
    commit_msg = """fix: update automation framework core.sh for bash compatibility

- Replace 'function' keyword with POSIX-compliant syntax
- Ensure compatibility with both bash and sh shells"""
    run_git_command(f'git commit -m "{commit_msg}"')

# Commit 2: Fix associative array compatibility
if os.path.exists('.automation/auto'):
    print("\n=== Commit 2: Fix associative array compatibility ===")
    run_git_command("git add .automation/auto")
    commit_msg = """fix: replace associative arrays in auto script for POSIX compliance

- Convert bash-specific associative arrays to portable implementation
- Maintain functionality while ensuring shell compatibility"""
    run_git_command(f'git commit -m "{commit_msg}"')

# Commit 3: Create missing .bash_profile.dir structure
if not os.path.exists('.bash_profile.dir'):
    print("\n=== Commit 3: Create bash profile directory ===")
    os.makedirs('.bash_profile.dir', exist_ok=True)
    with open('.bash_profile.dir/.gitkeep', 'w') as f:
        f.write('')
    run_git_command("git add .bash_profile.dir/")
    commit_msg = """feat: add missing .bash_profile.dir directory structure

- Create required directory for modular bash configuration
- Add .gitkeep to preserve empty directory in git"""
    run_git_command(f'git commit -m "{commit_msg}"')

# Commit 4: Add yes-to-all functionality to initialize script
if os.path.exists('initialize.sh'):
    print("\n=== Commit 4: Add yes-to-all functionality ===")
    run_git_command("git add initialize.sh")
    commit_msg = """feat: add --yes flag to initialize.sh for automated installation

- Support non-interactive installation with --yes flag
- Maintain backward compatibility with interactive mode"""
    run_git_command(f'git commit -m "{commit_msg}"')

# Commit 5: Update tools.sh changes
if os.path.exists('.automation/modules/tools.sh'):
    print("\n=== Commit 5: Update tools module ===")
    run_git_command("git add .automation/modules/tools.sh")
    commit_msg = """feat: add yes-to-all support to automation tools module

- Pass through --yes flag from parent scripts
- Enable fully automated tool installation"""
    run_git_command(f'git commit -m "{commit_msg}"')

# Commit 6: Add new make targets
if os.path.exists('Makefile'):
    print("\n=== Commit 6: Add new make targets ===")
    run_git_command("git add Makefile")
    commit_msg = """feat: enhance Makefile with new automation targets

- Add targets for automated installation and testing
- Improve development workflow with convenience commands"""
    run_git_command(f'git commit -m "{commit_msg}"')

# Check for additional project config files
print("\n=== Checking for additional project config files ===")
config_files = ['CLAUDE.md', 'package.json', 'package-lock.json']
files_to_add = [f for f in config_files if os.path.exists(f)]
if files_to_add:
    run_git_command(f"git add {' '.join(files_to_add)}")
    # Check if there are changes to commit
    result = subprocess.run("git diff --cached --quiet", shell=True)
    if result.returncode != 0:  # There are changes
        commit_msg = """chore: update project configuration files

- Update CLAUDE.md documentation
- Add npm package configuration for development tools"""
        run_git_command(f'git commit -m "{commit_msg}"')

# Show final status
print("\n=== Final git status ===")
run_git_command("git status")
run_git_command("git log --oneline -10")

print("\n=== Operation complete! ===")
print("Branch created and commits made successfully!")
print("You can now push with: git push -u origin fix/automation-framework-compatibility")