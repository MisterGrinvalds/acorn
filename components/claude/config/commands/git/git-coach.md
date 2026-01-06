---
description: Interactive coaching session to learn Git step by step
argument-hint: [skill-level: beginner|intermediate|advanced]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Guide the user through learning Git interactively. Assess their current level and provide hands-on exercises.

## Approach

1. **Assess level** - Ask about current Git experience
2. **Set goals** - Identify what workflows they need
3. **Progressive exercises** - Start simple, build complexity
4. **Real-time practice** - Have them run commands and report
5. **Reinforce** - Summarize and suggest next steps

## Skill Levels

### Beginner
- What is version control?
- git init, add, commit basics
- Viewing history with git log
- Working with remotes (clone, push, pull)
- Basic branching

### Intermediate
- Branch management strategies
- Merging and conflict resolution
- Rebasing basics
- Stashing changes
- Using aliases and shortcuts
- Interactive staging (git add -p)

### Advanced
- Interactive rebase for history cleanup
- Cherry-picking and bisect
- Reflog for recovery
- Git internals understanding
- Hooks and automation
- Large repo optimization

## Interactive Exercises

### Beginner Exercises
```bash
# Exercise 1: Initialize and first commit
mkdir git-practice && cd git-practice
git init
echo "Hello" > file.txt
git add file.txt
git commit -m "Initial commit"

# Exercise 2: View history
git log --oneline

# Exercise 3: Make changes and commit
echo "World" >> file.txt
git diff
git add file.txt
git commit -m "Add world"
```

### Intermediate Exercises
```bash
# Exercise 4: Branching
gcob feature/test  # or: git checkout -b feature/test
# make changes, commit
git checkout main
git merge feature/test

# Exercise 5: Resolve a conflict
# Create conflicting changes on two branches

# Exercise 6: Interactive add
gadd  # or: git add -p
```

### Advanced Exercises
```bash
# Exercise 7: Interactive rebase
gri HEAD~3  # Squash/reorder last 3 commits

# Exercise 8: Recovery with reflog
git reflog
git checkout <lost-commit>

# Exercise 9: Bisect
git bisect start
git bisect bad
git bisect good <known-good>
```

## Context

@components/git/functions.sh
@components/git/aliases.sh

## Coaching Style

- Use dotfiles aliases when teaching shortcuts
- Always explain what commands do before running
- Emphasize safety (don't rebase shared branches)
- Show recovery options for mistakes
- Build confidence with small wins
