---
description: Configure Claude Code memory with CLAUDE.md files
argument-hint: [scope: project|user|local]
allowed-tools: Read, Write, Edit, Glob
---

# Claude Code Memory Configuration

Help the user configure Claude Code memory files.

## Memory Hierarchy (Highest to Lowest)

| Type | Location | Shared | Purpose |
|------|----------|--------|---------|
| Enterprise | System-level | Yes | Organization policies |
| Project | `./CLAUDE.md` | Yes | Team instructions |
| Project Rules | `./.claude/rules/*.md` | Yes | Modular rules |
| User | `~/.claude/CLAUDE.md` | No | Personal preferences |
| Local | `./CLAUDE.local.md` | No | Personal project overrides |

## Your Task

Based on the user's request: $ARGUMENTS

1. Determine appropriate scope
2. Create or update the memory file
3. Use proper structure and formatting

## CLAUDE.md Best Practices

### Structure
```markdown
# Project Name

## Build & Test
- `npm run build` - Build the project
- `npm test` - Run tests

## Code Style
- Use 2-space indentation
- Prefer const over let
- Use TypeScript strict mode

## Architecture
- Components in src/components/
- Utils in src/utils/
- Tests alongside source files
```

### Imports
Reference other files with `@`:
```markdown
See @README.md for overview
Build commands: @package.json
Git workflow: @docs/git-guide.md
Personal settings: @~/.claude/my-preferences.md
```

## Modular Rules (.claude/rules/)

Create focused rule files:
```
.claude/rules/
├── code-style.md
├── testing.md
├── security.md
└── frontend/
    └── react.md
```

### Path-Specific Rules
```markdown
---
paths: src/api/**/*.ts
---

# API Development Rules
- All endpoints must validate input
- Use proper error handling
```

## Commands

- `/init` - Bootstrap CLAUDE.md
- `/memory` - Open memory files in editor

## Tips

- Be specific: "Use 2-space indentation" > "Format code properly"
- Use bullet points under descriptive headings
- Review and update as project evolves
- Use symlinks to share common rules across projects
