---
description: Clean Node.js project caches and reinstall dependencies
argument-hint: [scope: project|all|cache]
allowed-tools: Read, Bash
---

## Task

Help the user clean up Node.js project artifacts and reinstall dependencies.

## Scopes

Based on `$ARGUMENTS`:

### project
Clean current project node_modules:

```bash
# Using dotfiles function
nclean

# Manual (npm)
rm -rf node_modules
npm install

# Manual (pnpm)
rm -rf node_modules
pnpm install

# Manual (yarn)
rm -rf node_modules
yarn install
```

The `nclean` function auto-detects package manager from lock file.

### all
Clean all node_modules in directory tree:

```bash
# Find all node_modules with sizes
nfind

# Remove all (with confirmation)
ncleanall
```

This is useful for:
- Freeing disk space
- Cleaning monorepos
- Starting fresh after switching Node versions

### cache
Clean package manager caches:

```bash
# npm
npm cache clean --force
npm cache verify

# pnpm
pnpm store prune

# yarn
yarn cache clean
```

## Full Clean Sequence

For a complete project reset:

```bash
# 1. Remove node_modules
rm -rf node_modules

# 2. Remove build artifacts
rm -rf dist build .next .nuxt

# 3. Remove cache directories
rm -rf .cache .parcel-cache

# 4. Remove lock file (optional - if want fresh resolution)
# rm package-lock.json  # or pnpm-lock.yaml

# 5. Clean npm cache
npm cache clean --force

# 6. Reinstall
npm install  # or pnpm install
```

## When to Clean

### Clean node_modules when:
- Switching Node versions
- Strange import/module errors
- After modifying dependencies manually
- Build issues that seem cache-related

### Clean all node_modules when:
- Low disk space
- Switching machines/restoring from backup
- Major Node version upgrade

### Clean cache when:
- Corrupted cache errors
- Very low disk space
- Suspected package corruption

## Disk Space Check

```bash
# Check node_modules size in current project
du -sh node_modules

# Find largest directories
du -sh node_modules/* | sort -hr | head -20

# Find all node_modules recursively
nfind
```

## Monorepo Cleaning

```bash
# pnpm workspace
pnpm recursive clean
rm -rf node_modules **/node_modules

# npm workspace
npm run clean --workspaces

# Manual
ncleanall  # dotfiles function
```

## Dotfiles Functions

- `nclean` - Clean and reinstall current project
- `nfind` - Find all node_modules with sizes
- `ncleanall` - Remove all node_modules (with confirmation)
- `npm_detect` - Detect package manager for proper reinstall
