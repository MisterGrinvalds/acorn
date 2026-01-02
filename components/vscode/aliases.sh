#!/bin/sh
# components/vscode/aliases.sh - VS Code aliases

# Open current directory
alias c='code .'

# Open with options
alias cg='code --goto'
alias cn='code --new-window'
alias ca='code --add'
alias cr='code --reuse-window'

# Diff mode
alias cdiff='code --diff'

# Extensions
alias cext='code --list-extensions'
alias cexti='code --install-extension'
alias cextu='code --uninstall-extension'
