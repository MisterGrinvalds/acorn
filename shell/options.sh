#!/bin/sh
# Shell options configuration
# Requires: shell/discovery.sh

case "$CURRENT_SHELL" in
    bash)
        # History options
        shopt -s histappend         # Append to history, don't overwrite
        shopt -s cmdhist            # Save multi-line commands as single entry

        # Directory navigation
        shopt -s autocd 2>/dev/null # cd by typing directory name (bash 4+)
        shopt -s cdspell            # Autocorrect typos in cd
        shopt -s dirspell 2>/dev/null  # Autocorrect directory names (bash 4+)

        # Globbing
        shopt -s globstar 2>/dev/null  # ** matches recursively (bash 4+)
        shopt -s nocaseglob         # Case-insensitive globbing
        shopt -s extglob            # Extended pattern matching

        # Misc
        shopt -s checkwinsize       # Update LINES/COLUMNS after each command
        shopt -s no_empty_cmd_completion  # Don't complete on empty line

        # History control
        export HISTCONTROL="erasedups:ignoredups:ignorespace"
        export HISTSIZE=10000
        export HISTFILESIZE=20000
        export HISTIGNORE="ls:ll:cd:pwd:bg:fg:history:clear:exit"
        ;;

    zsh)
        # History options
        setopt APPEND_HISTORY       # Append to history, don't overwrite
        setopt HIST_IGNORE_DUPS     # Don't store duplicates
        setopt HIST_IGNORE_SPACE    # Don't store commands starting with space
        setopt HIST_EXPIRE_DUPS_FIRST  # Remove duplicates first when trimming
        setopt HIST_FIND_NO_DUPS    # Don't show duplicates in search
        setopt SHARE_HISTORY        # Share history between sessions
        setopt EXTENDED_HISTORY     # Save timestamp

        # Directory navigation
        setopt AUTO_CD              # cd by typing directory name
        setopt AUTO_PUSHD           # Push directories onto stack
        setopt PUSHD_IGNORE_DUPS    # Don't push duplicates
        setopt PUSHD_SILENT         # Don't print directory stack
        setopt CDABLE_VARS          # cd to named directories

        # Globbing
        setopt EXTENDED_GLOB        # Extended pattern matching
        setopt NO_CASE_GLOB         # Case-insensitive globbing
        setopt GLOB_DOTS            # Include dotfiles in glob

        # Misc
        setopt CORRECT              # Spelling correction for commands
        setopt NO_BEEP              # Don't beep on errors
        setopt INTERACTIVE_COMMENTS # Allow comments in interactive shell

        # History settings
        export HISTSIZE=10000
        export SAVEHIST=20000
        ;;
esac
