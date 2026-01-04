---
description: Delete a smug session configuration
argument-hint: <session-name>
allowed-tools: Bash
---

## Task

Delete a smug session configuration file.

## Execution

### Safety Check

Always confirm before deletion:
```bash
echo "Delete session '$ARGUMENTS'? This will remove:"
echo "  ${SMUG_CONFIG_DIR:-$HOME/.config/smug}/$ARGUMENTS.yml"
```

### Delete Command

```bash
rm "${SMUG_CONFIG_DIR:-$HOME/.config/smug}/$ARGUMENTS.yml"
```

### Stop Running Session First

If the session is running, stop it:
```bash
smug stop $ARGUMENTS 2>/dev/null || true
rm "${SMUG_CONFIG_DIR:-$HOME/.config/smug}/$ARGUMENTS.yml"
```

## Post-Deletion

Sync the deletion to git:
```bash
smug_push "Remove $ARGUMENTS session"
```

## Recovery

If accidentally deleted and using git sync:
```bash
# Check git status
smug_status

# Restore from git
cd "${SMUG_REPO_DIR:-$HOME/.local/share/smug-sessions}"
git checkout -- "$ARGUMENTS.yml"
```

## Context

@components/tmux/functions.sh
