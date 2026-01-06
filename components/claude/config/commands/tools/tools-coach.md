---
description: Interactive coaching session to learn system tools management
argument-hint: [skill-level: beginner|intermediate|advanced]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Guide the user through learning system tools management interactively.

## Approach

1. **Assess level** - Ask about experience with package managers
2. **Detect platform** - Identify macOS or Linux distribution
3. **Progressive exercises** - Build from basics to advanced
4. **Real-time practice** - Have them run actual commands
5. **Reinforce** - Summarize and suggest next steps

## Skill Levels

### Beginner
- What is a package manager?
- Installing your first tool
- Checking tool versions
- Updating packages
- Finding installed tools

### Intermediate
- Managing multiple versions
- Cleaning up old packages
- Using dotfiles functions
- Troubleshooting path issues
- Cross-platform considerations

### Advanced
- Creating automation scripts
- Building tool manifests
- Managing team tool versions
- Security considerations
- Performance optimization

## Interactive Exercises

### Beginner Exercises
```bash
# Exercise 1: Check what's available
versions

# Exercise 2: Find where a tool is
whichx git
whichx python3

# Exercise 3: Check package manager
which brew    # macOS
which apt-get # Linux

# Exercise 4: Update package lists
# macOS:
brew update

# Linux (Debian/Ubuntu):
sudo apt update
```

### Intermediate Exercises
```bash
# Exercise 5: Check all tools
tools-list
tools-check

# Exercise 6: Find missing tools
tools-missing

# Exercise 7: Update everything
system-update

# Exercise 8: Install a new tool
# macOS:
brew install fzf

# Linux:
sudo apt install fzf

# Exercise 9: Clean up
# macOS:
brew cleanup
brew autoremove
```

### Advanced Exercises
```bash
# Exercise 10: Check specific tool versions
tools-check git
tools-check node

# Exercise 11: Upgrade bash on macOS
upgrade_bash

# Exercise 12: Create version audit script
quick_versions > ~/tool-versions-$(date +%Y%m%d).txt

# Exercise 13: Compare tool locations
type -a python3  # bash
whence -a python3  # zsh
```

## Platform-Specific Tracks

### macOS Track
1. Homebrew basics
2. Cask for GUI apps
3. Services management
4. Path configuration
5. bash vs zsh

### Linux (Debian/Ubuntu) Track
1. apt/apt-get basics
2. PPA repositories
3. Snap packages
4. Manual installations
5. PATH management

## Context

@components/tools/functions.sh
@components/tools/aliases.sh

## Coaching Style

- Start with version checks (`versions`)
- Use platform-appropriate commands
- Show dotfiles shortcuts for efficiency
- Emphasize keeping tools updated
- Build toward automation
- Address common path issues
