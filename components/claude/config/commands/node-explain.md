---
description: Explain Node.js concepts, tools, and workflows
argument-hint: [topic]
allowed-tools: Read, Glob, Grep
---

## Task

Explain the requested topic about Node.js development. If no specific topic provided, give an overview of Node.js tooling and the dotfiles integration.

## Topics

If `$ARGUMENTS` is provided, explain that specific topic:

- **nvm** - Node Version Manager usage
- **npm** - npm package management
- **pnpm** - pnpm package manager benefits
- **yarn** - Yarn package manager
- **typescript** - TypeScript setup and configuration
- **modules** - ES modules vs CommonJS
- **package.json** - Package configuration
- **scripts** - npm scripts and task running
- **monorepo** - Monorepo with pnpm workspaces
- **testing** - Jest, Vitest, Node test runner
- **building** - Build tools (esbuild, tsup, tsc)
- **eslint** - ESLint configuration
- **prettier** - Code formatting

## Context

Reference these files for accurate information:
@components/node/component.yaml
@components/node/functions.sh
@components/node/aliases.sh

## Response Format

1. **Concept overview** - What it is and why it matters
2. **Key commands** - Essential commands
3. **Dotfiles integration** - Available functions and aliases
4. **Examples** - Practical usage examples
5. **Best practices** - Modern recommendations
