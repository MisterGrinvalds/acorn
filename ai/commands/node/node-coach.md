---
description: Interactive coaching session to learn Node.js development
argument-hint: [skill-level: beginner|intermediate|advanced]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Guide the user through learning Node.js development workflows interactively.

## Approach

1. **Assess level** - Ask about Node.js/JavaScript experience
2. **Set goals** - Identify what they want to build
3. **Progressive exercises** - Start simple, build complexity
4. **Real-time practice** - Have them run commands
5. **Reinforce** - Summarize and suggest next steps

## Skill Levels

### Beginner
- What is Node.js?
- Installing Node with NVM
- Running JavaScript files
- npm basics (install, scripts)
- package.json structure

### Intermediate
- pnpm vs npm vs yarn
- TypeScript setup
- Project structure
- npm scripts
- Development workflow (dev, build, test)
- Using dotfiles functions

### Advanced
- Monorepo setup
- Build optimization
- Testing strategies
- CI/CD integration
- Performance profiling

## Interactive Exercises

### Beginner Exercises
```bash
# Exercise 1: Check Node installation
node --version
npm --version

# Exercise 2: Setup NVM
nvm_setup  # dotfiles function

# Exercise 3: Create simple project
mkdir my-first-node
cd my-first-node
npm init -y

# Exercise 4: Create and run a file
echo 'console.log("Hello Node!")' > index.js
node index.js
```

### Intermediate Exercises
```bash
# Exercise 5: Initialize TypeScript project
node_init my-ts-project

# Exercise 6: Use pnpm
pi  # pnpm install (alias)
pa express  # pnpm add
pad typescript  # pnpm add -D

# Exercise 7: Run scripts
prd  # pnpm run dev
prb  # pnpm run build
```

### Advanced Exercises
```bash
# Exercise 8: Setup monorepo
pnpm init
# Add pnpm-workspace.yaml
# Create packages structure

# Exercise 9: Configure build tool
pa -D tsup
# Setup build scripts

# Exercise 10: Setup testing
pa -D vitest
# Write first test
```

## Context

@components/node/functions.sh
@components/node/aliases.sh

## Coaching Style

- Start with NVM for version management
- Prefer pnpm for new projects
- Use dotfiles aliases for common tasks
- Show TypeScript as default for new projects
- Explain package.json structure
