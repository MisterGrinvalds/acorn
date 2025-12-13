#!/bin/sh
# Secrets loading - silent, secure
# Requires: shell/xdg.sh (for DOTFILES_ROOT)
#
# Secrets directory structure:
#   secrets/
#     *.env       - Environment variable files (KEY=value format)
#     *.sh        - Shell scripts to source
#     .gitignore  - Already configured to ignore all except .gitignore
#
# Security notes:
#   - All output is suppressed
#   - Files are sourced in a subshell to prevent errors from propagating
#   - Only files (not directories) are processed

SECRETS_DIR="${DOTFILES_ROOT}/secrets"

# Exit silently if secrets directory doesn't exist
[ -d "$SECRETS_DIR" ] || return 0

# Load .env files (KEY=value format)
for env_file in "$SECRETS_DIR"/*.env; do
    [ -f "$env_file" ] || continue
    # Read each line, skip comments and empty lines
    while IFS='=' read -r key value; do
        # Skip comments and empty lines
        case "$key" in
            \#*|"") continue ;;
        esac
        # Remove leading/trailing whitespace from key
        key=$(echo "$key" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        # Export the variable (value can contain = signs)
        [ -n "$key" ] && export "$key=$value"
    done < "$env_file" 2>/dev/null
done

# Source .sh files
for sh_file in "$SECRETS_DIR"/*.sh; do
    [ -f "$sh_file" ] || continue
    . "$sh_file" 2>/dev/null
done
