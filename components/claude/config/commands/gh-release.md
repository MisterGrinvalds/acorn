---
description: Create and manage GitHub releases
argument-hint: <version> [--draft] [--prerelease]
allowed-tools: Read, Bash
---

## Task

Help the user create and manage GitHub releases.

## Create Release

### Basic Release
```bash
gh release create v1.0.0
# Opens editor for release notes
```

### With Options
```bash
gh release create v1.0.0 \
  --title "Version 1.0.0" \
  --notes "Release notes here"
```

### Generate Notes Automatically
```bash
gh release create v1.0.0 --generate-notes
# Uses commit messages since last release
```

### Draft Release
```bash
gh release create v1.0.0 --draft
# Review before publishing
```

### Pre-release
```bash
gh release create v1.0.0-beta.1 --prerelease
```

### With Assets
```bash
gh release create v1.0.0 \
  ./dist/app-linux-amd64 \
  ./dist/app-darwin-arm64 \
  ./dist/app-windows-amd64.exe
```

## Release Notes Template

```markdown
## What's New

### Features
- Added X feature (#123)
- Implemented Y functionality (#456)

### Bug Fixes
- Fixed Z issue (#789)

### Breaking Changes
- Changed API endpoint format

### Dependencies
- Updated dependency A to v2.0

## Contributors
@user1, @user2
```

## Manage Releases

### List Releases
```bash
gh release list
```

### View Release
```bash
gh release view v1.0.0
gh release view v1.0.0 --web
```

### Edit Release
```bash
gh release edit v1.0.0 --title "New Title"
gh release edit v1.0.0 --notes "Updated notes"
gh release edit v1.0.0 --draft=false  # Publish draft
```

### Delete Release
```bash
gh release delete v1.0.0
gh release delete v1.0.0 --yes  # Skip confirmation
```

### Download Assets
```bash
gh release download v1.0.0
gh release download v1.0.0 --pattern "*.tar.gz"
gh release download v1.0.0 --dir ./downloads
```

## Release Workflow

### Standard Workflow
```bash
# 1. Ensure main is up to date
git checkout main
git pull

# 2. Create version tag
git tag v1.0.0
git push origin v1.0.0

# 3. Create release
gh release create v1.0.0 --generate-notes

# 4. Upload binaries (if applicable)
gh release upload v1.0.0 ./dist/*
```

### With Build Artifacts
```bash
# Build for multiple platforms
gobuildall myapp  # dotfiles function

# Create release with all artifacts
gh release create v1.0.0 \
  --generate-notes \
  ./dist/myapp-linux-amd64 \
  ./dist/myapp-darwin-arm64 \
  ./dist/myapp-windows-amd64.exe
```

## Semantic Versioning

| Version | Meaning |
|---------|---------|
| v1.0.0 | Major.Minor.Patch |
| v1.0.0-alpha | Pre-release alpha |
| v1.0.0-beta.1 | Pre-release beta |
| v1.0.0-rc.1 | Release candidate |

### When to Bump

- **Major (1.x.x)**: Breaking changes
- **Minor (x.1.x)**: New features, backward compatible
- **Patch (x.x.1)**: Bug fixes, backward compatible

## CI Integration

### GitHub Actions Release
```yaml
name: Release
on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build
        run: make build-all
      - name: Create Release
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          gh release create ${{ github.ref_name }} \
            --generate-notes \
            ./dist/*
```

## Tips

1. Use `--generate-notes` for automatic changelog
2. Create draft releases for review before publishing
3. Include checksums for security
4. Use pre-releases for testing
5. Tag before releasing for proper linking
