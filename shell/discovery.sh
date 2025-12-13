# shell/discovery.sh - Shell and platform detection
# MUST be sourced first - everything else depends on these variables
#
# Exports:
#   CURRENT_SHELL    - bash, zsh, or unknown
#   CURRENT_PLATFORM - darwin, linux, or unknown
#   IS_INTERACTIVE   - true or false
#   IS_LOGIN_SHELL   - true or false

# Detect shell type
if [ -n "$BASH_VERSION" ]; then
    export CURRENT_SHELL="bash"
elif [ -n "$ZSH_VERSION" ]; then
    export CURRENT_SHELL="zsh"
else
    export CURRENT_SHELL="unknown"
fi

# Detect platform
case "$OSTYPE" in
    darwin*)
        export CURRENT_PLATFORM="darwin"
        ;;
    linux*)
        export CURRENT_PLATFORM="linux"
        ;;
    *)
        # Fallback to uname
        case "$(uname -s)" in
            Darwin)
                export CURRENT_PLATFORM="darwin"
                ;;
            Linux)
                export CURRENT_PLATFORM="linux"
                ;;
            *)
                export CURRENT_PLATFORM="unknown"
                ;;
        esac
        ;;
esac

# Detect interactive shell (skip if already set, for testing)
if [ -z "$IS_INTERACTIVE" ]; then
    case "$-" in
        *i*)
            export IS_INTERACTIVE="true"
            ;;
        *)
            export IS_INTERACTIVE="false"
            ;;
    esac
fi

# Detect login shell
if [ "$CURRENT_SHELL" = "bash" ]; then
    if shopt -q login_shell 2>/dev/null; then
        export IS_LOGIN_SHELL="true"
    else
        export IS_LOGIN_SHELL="false"
    fi
elif [ "$CURRENT_SHELL" = "zsh" ]; then
    if [[ -o login ]]; then
        export IS_LOGIN_SHELL="true"
    else
        export IS_LOGIN_SHELL="false"
    fi
else
    export IS_LOGIN_SHELL="unknown"
fi

# Early exit for non-interactive shells (performance optimization)
# Comment this out if you need full config in scripts
if [ "$IS_INTERACTIVE" = "false" ]; then
    return 0 2>/dev/null || exit 0
fi
