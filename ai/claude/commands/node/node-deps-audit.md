---
description: Audit and manage Node.js dependencies
argument-hint: [action: audit|outdated|update|duplicates]
allowed-tools: Read, Bash
---

## Task

Help the user audit and manage their Node.js project dependencies.

## Actions

Based on `$ARGUMENTS`:

### audit
Check for security vulnerabilities:

```bash
# npm
npm audit
npm audit fix  # Auto-fix

# pnpm
pnpm audit
pnpm audit --fix

# Detailed report
npm audit --json > audit-report.json
```

### outdated
Check for outdated packages:

```bash
# npm
npm outdated

# pnpm
pnpm outdated

# Output shows:
# Package  Current  Wanted  Latest
```

### update
Update dependencies:

```bash
# Update within semver ranges
npm update
pnpm update

# Update to latest (may break)
npm update --latest
pnpm update --latest

# Update specific package
npm update package-name
pnpm update package-name

# Interactive update (npm)
npx npm-check -u
```

### duplicates
Find duplicate dependencies:

```bash
# npm
npm dedupe
npm ls --all | grep "deduped"

# pnpm (automatically handles this)
pnpm why <package>
```

## Dependency Analysis

### Check package size
```bash
# Install bundle analyzer
npx bundlephobia-cli <package-name>

# Or use online tool
# https://bundlephobia.com
```

### Find unused dependencies
```bash
npx depcheck

# Output shows:
# - Unused dependencies
# - Missing dependencies
# - Unused devDependencies
```

### Visualize dependency tree
```bash
# npm
npm ls
npm ls --depth=1

# pnpm
pnpm list
pnpm list --depth=1

# Why is package installed?
npm explain <package>
pnpm why <package>
```

## Best Practices

### Security
```bash
# Run audit in CI
npm audit --audit-level=high

# Ignore specific advisories (if needed)
npm audit --ignore <advisory-id>
```

### Version Pinning
```json
// package.json
{
  "dependencies": {
    "express": "4.18.2"      // Exact version
  },
  "devDependencies": {
    "typescript": "^5.0.0"   // Allow minor updates
  }
}
```

### Lock File Maintenance
```bash
# Regenerate lock file
rm package-lock.json  # or pnpm-lock.yaml
npm install           # or pnpm install
```

## Common Issues

### Peer dependency warnings
```bash
# Check peer deps
npm ls --depth=0

# Install missing peer deps
npm install <peer-dep>
```

### Conflicting versions
```bash
# Check why conflict exists
pnpm why <package>

# Override version
# Add to package.json:
{
  "pnpm": {
    "overrides": {
      "package-name": "^2.0.0"
    }
  }
}
```

## Dotfiles Integration

- `nclean` - Remove node_modules and reinstall
- `nfind` - Find all node_modules with sizes
- `npm_detect` - Detect package manager
