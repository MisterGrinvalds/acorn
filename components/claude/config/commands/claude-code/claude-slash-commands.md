---
description: Create custom Claude Code slash commands
argument-hint: [command-name] [description]
allowed-tools: Read, Write, Edit, Glob
---

# Custom Slash Commands

Help the user create custom slash commands.

## Command Locations

| Type | Location | Shown As |
|------|----------|----------|
| Project | `.claude/commands/` | (project) |
| User | `~/.claude/commands/` | (user) |

Project commands override user commands with same name.

## Your Task

Based on the user's request: $ARGUMENTS

1. Create the command markdown file
2. Add appropriate frontmatter
3. Include argument handling if needed

## Command Format

```markdown
---
description: Brief description (shown in /help)
argument-hint: [arg1] [arg2]
allowed-tools: Tool1, Tool2
model: claude-sonnet-4-5-20250929
---

Your prompt goes here.

Use $ARGUMENTS for all arguments.
Use $1, $2, etc. for positional arguments.
```

## Frontmatter Fields

| Field | Purpose |
|-------|---------|
| `description` | Shown in help, enables SlashCommand tool |
| `argument-hint` | Shows expected arguments |
| `allowed-tools` | Restrict available tools |
| `model` | Override model for this command |
| `disable-model-invocation` | Prevent Claude from calling this |

## Arguments

### All Arguments
```markdown
Create a commit with message: $ARGUMENTS
```
Usage: `/commit fix: resolve login bug`

### Positional Arguments
```markdown
Review PR #$1 with priority $2
```
Usage: `/review-pr 123 high`

## Special Features

### Bash Execution
Use `!` prefix to run bash commands:
```markdown
## Context
- Git status: !`git status --short`
- Current branch: !`git branch --show-current`
- Recent commits: !`git log --oneline -5`

## Task
Based on the above, create a commit.
```

### File References
Use `@` to include file contents:
```markdown
Review the implementation in @src/main.ts
Compare @old.js with @new.js
```

## Examples

### Quick Commit
```markdown
---
description: Create a quick git commit
argument-hint: [message]
allowed-tools: Bash(git:*)
---

## Context
!`git status --short`
!`git diff --staged`

Create a git commit with message: $ARGUMENTS
```

### Code Review
```markdown
---
description: Review code changes
allowed-tools: Read, Grep, Glob, Bash
---

## Current Changes
!`git diff HEAD`

Review these changes for:
- Code quality
- Security issues
- Test coverage
- Performance
```

### PR Description
```markdown
---
description: Generate PR description
allowed-tools: Bash(git:*), Read
---

## Branch Info
!`git log main..HEAD --oneline`
!`git diff main...HEAD --stat`

Generate a pull request description summarizing these changes.
```

## Namespacing

Use subdirectories to organize:
```
.claude/commands/
├── git/
│   ├── commit.md      → /commit (project:git)
│   └── pr.md          → /pr (project:git)
└── review/
    └── security.md    → /security (project:review)
```

## Built-in Commands Reference

Essential: `/init`, `/memory`, `/config`, `/agents`, `/hooks`
Session: `/clear`, `/compact`, `/rewind`, `/resume`
Info: `/help`, `/cost`, `/context`, `/status`
