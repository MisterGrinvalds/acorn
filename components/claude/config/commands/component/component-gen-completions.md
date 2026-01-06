---
description: Generate bash/zsh completion scripts for a component
argument_hints:
  - tmux
  - go
  - python
  - node
  - kubernetes
---

Generate completions for: $ARGUMENTS

## Instructions

Create shell completion scripts for the component's functions and aliases.

### 1. Identify Completable Functions

Read the component config and identify functions that would benefit from completions:
- Functions with file/directory arguments
- Functions with predefined options
- Functions that wrap other commands

### 2. Update shell/completions.sh

Create `components/$ARGUMENTS/shell/completions.sh`:

```bash
#!/bin/sh
# $ARGUMENTS completions
# Supports: bash, zsh

# =============================================================================
# Bash Completions
# =============================================================================

if [ -n "$BASH_VERSION" ]; then
    # Load bash-completion if available
    if [ -f /opt/homebrew/etc/profile.d/bash_completion.sh ]; then
        . /opt/homebrew/etc/profile.d/bash_completion.sh
    elif [ -f /usr/share/bash-completion/bash_completion ]; then
        . /usr/share/bash-completion/bash_completion
    elif [ -f /etc/bash_completion ]; then
        . /etc/bash_completion
    fi

    # Tool-specific completions
    # <Add bash completion functions here>
fi

# =============================================================================
# Zsh Completions
# =============================================================================

if [ -n "$ZSH_VERSION" ]; then
    # Initialize completion system
    autoload -Uz compinit

    # Cache completions for faster startup
    if [[ -n ${ZDOTDIR}/.zcompdump(#qN.mh+24) ]]; then
        compinit
    else
        compinit -C
    fi

    # Completion styling
    zstyle ':completion:*' menu select
    zstyle ':completion:*' matcher-list 'm:{a-zA-Z}={A-Za-z}'
    zstyle ':completion:*' list-colors "${(s.:.)LS_COLORS}"

    # Tool-specific completions
    # <Add zsh completion functions here>
fi
```

### 3. Add Function-Specific Completions

For functions that need completions, add completion handlers:

**Example: Session selector function**
```bash
# Bash
_<function>_completion() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    # Generate completion options
    COMPREPLY=($(compgen -W "<options>" -- "$cur"))
}
complete -F _<function>_completion <function>

# Zsh
_<function>() {
    local -a options
    options=(<options>)
    _describe '<description>' options
}
compdef _<function> <function>
```

### 4. Common Completion Patterns

| Pattern | Bash | Zsh |
|---------|------|-----|
| Files | `compgen -f` | `_files` |
| Directories | `compgen -d` | `_directories` |
| Commands | `compgen -c` | `_command_names` |
| Custom list | `compgen -W "list"` | `_describe 'desc' list` |

### 5. Tool-Specific Completions

| Component | Completions Needed |
|-----------|-------------------|
| tmux | Session names, window names |
| go | Package names, test patterns |
| python | Virtualenv names |
| node | npm scripts, packages |
| kubernetes | Contexts, namespaces, pods |
| git | Branches, remotes, tags |

### 6. Load Existing Tool Completions

For tools with their own completions, ensure they're loaded:

```bash
# kubectl completions
if command -v kubectl &>/dev/null; then
    if [ -n "$BASH_VERSION" ]; then
        source <(kubectl completion bash)
    elif [ -n "$ZSH_VERSION" ]; then
        source <(kubectl completion zsh)
    fi
fi
```

### 7. Report

Output:
```
Generated Completions: $ARGUMENTS
=================================

Updated: components/$ARGUMENTS/shell/completions.sh

Completions added:
  - <function1>: <completion type>
  - <function2>: <completion type>

External completions loaded:
  - <tool> (native completion)

Shells supported:
  - bash: YES
  - zsh: YES
```
