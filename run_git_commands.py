#!/usr/bin/env python3
import os
import subprocess
import sys

# Set the working directory
repo_path = '/Users/mistergrinvalds/Repos/personal/bash-profile'
os.chdir(repo_path)

def execute_command(cmd):
    print(f"\n>>> Executing: {cmd}")
    try:
        result = subprocess.run(cmd, shell=True, capture_output=True, text=True, cwd=repo_path)
        if result.stdout:
            print(result.stdout)
        if result.stderr:
            print(f"STDERR: {result.stderr}")
        print(f"Return code: {result.returncode}")
        return result.returncode == 0
    except Exception as e:
        print(f"Error executing command: {e}")
        return False

# Execute the git operations
print("Starting git operations for bash-profile repository")
print(f"Working directory: {os.getcwd()}")

# Check current status
execute_command("git status")

# Create and switch to new branch
execute_command("git checkout -b fix/automation-framework-compatibility")

# Show the current branch
execute_command("git branch")

print("\nGit operations completed!")
print("Next steps:")
print("1. Add your files with: git add <files>")
print("2. Make commits with: git commit -m 'your message'")
print("3. Push with: git push -u origin fix/automation-framework-compatibility")