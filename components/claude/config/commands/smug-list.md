---
description: List all smug session configurations
allowed-tools: Bash, Read
---

## Task

List all available smug session configurations with their descriptions.

## Execution

```bash
smug_list
```

## Output Format

Display sessions in a table format:
```
Session Name    | Description                  | Windows
----------------|------------------------------|--------
dotfiles        | Dotfiles development         | 3
my-project      | Main project workspace       | 4
```

## Alternative: Direct Listing

If `smug_list` is not available:
```bash
ls -1 "${SMUG_CONFIG_DIR:-$HOME/.config/smug}"/*.yml 2>/dev/null | xargs -I {} basename {} .yml
```

## Context

@components/tmux/functions.sh
