---
description: Guide on using Claude Code checkpointing and rewind features
allowed-tools: Read, Glob, Bash
---

# Claude Code Checkpointing Guide

Help the user understand and use Claude Code's checkpointing and rewind features.

## What is Checkpointing?

Claude Code automatically tracks the state of your code before each edit, creating checkpoints that let you:
- Quickly undo changes
- Rewind to previous states
- Experiment fearlessly

## How It Works

- **Every user prompt creates a checkpoint**
- **Persistent across sessions** - available when resuming
- **Auto-cleanup** - removed after 30 days (configurable via `cleanupPeriodDays`)

## Using Rewind

### Access Methods
1. **Keyboard**: Press `Esc` twice (`Esc + Esc`)
2. **Command**: `/rewind`

### Rewind Options

| Option | What It Does |
|--------|--------------|
| **Conversation only** | Rewind to a message, keep code changes |
| **Code only** | Revert files, keep conversation history |
| **Both** | Restore both to a prior point |

## Use Cases

### Exploring Alternatives
Try different implementations:
```
> Refactor this using strategy pattern
[Examine result]
Esc + Esc → Rewind code only
> Now try using factory pattern instead
```

### Recovering from Mistakes
```
> Refactor the entire authentication system
[Something breaks]
Esc + Esc → Rewind both code and conversation
> Let's take a more incremental approach...
```

### Iterating on Features
```
> Add dark mode support
[Review implementation]
Esc + Esc → Rewind code only
> Add dark mode but use CSS variables instead
```

## Important Limitations

### Not Tracked by Checkpoints

1. **Bash command changes**:
   ```bash
   rm file.txt        # NOT tracked
   mv old.txt new.txt # NOT tracked
   cp a.txt b.txt     # NOT tracked
   ```

2. **External changes**:
   - Manual edits outside Claude Code
   - Changes from other concurrent sessions
   - Git operations

### Not a Replacement for Git

| Checkpoints | Git |
|-------------|-----|
| Session-level safety | Permanent history |
| Quick local undo | Collaboration |
| 30-day retention | Forever |
| Automatic | Explicit commits |

## Best Practices

1. **Use for experimentation** - Try bold refactorings knowing you can rewind
2. **Commit milestones to Git** - Don't rely solely on checkpoints
3. **Understand the scope** - Bash and external changes aren't tracked
4. **Combine with version control** - Checkpoints for session safety, Git for history

## Configuration

Set cleanup period in settings.json:
```json
{
  "cleanupPeriodDays": 30
}
```

## Commands

- `/rewind` - Open rewind menu
- `Esc + Esc` - Quick access to rewind
