#!/bin/sh
# core/theme.sh - Catppuccin Mocha color theme
# https://github.com/catppuccin/catppuccin
#
# Provides color variables for use throughout the shell configuration.
# Supports both 256-color and truecolor (24-bit) terminals.

# =============================================================================
# Truecolor Detection
# =============================================================================

_has_truecolor() {
    case "$COLORTERM" in
        truecolor|24bit) return 0 ;;
    esac
    case "$TERM" in
        *-truecolor|*-24bit|xterm-256color|screen-256color) return 0 ;;
    esac
    return 1
}

# RGB escape sequence helpers
_rgb() {
    printf '\033[38;2;%d;%d;%dm' "$1" "$2" "$3"
}

_rgb_bg() {
    printf '\033[48;2;%d;%d;%dm' "$1" "$2" "$3"
}

# =============================================================================
# Catppuccin Mocha Palette
# =============================================================================

if _has_truecolor; then
    # Truecolor (24-bit RGB) - exact Catppuccin Mocha colors
    THEME_ROSEWATER=$(_rgb 245 224 220)
    THEME_FLAMINGO=$(_rgb 242 205 205)
    THEME_PINK=$(_rgb 245 194 231)
    THEME_MAUVE=$(_rgb 203 166 247)
    THEME_RED=$(_rgb 243 139 168)
    THEME_MAROON=$(_rgb 235 160 172)
    THEME_PEACH=$(_rgb 250 179 135)
    THEME_YELLOW=$(_rgb 249 226 175)
    THEME_GREEN=$(_rgb 166 227 161)
    THEME_TEAL=$(_rgb 148 226 213)
    THEME_SKY=$(_rgb 137 220 235)
    THEME_SAPPHIRE=$(_rgb 116 199 236)
    THEME_BLUE=$(_rgb 137 180 250)
    THEME_LAVENDER=$(_rgb 180 190 254)
    THEME_TEXT=$(_rgb 205 214 244)
    THEME_SUBTEXT1=$(_rgb 186 194 222)
    THEME_SUBTEXT0=$(_rgb 166 173 200)
    THEME_OVERLAY2=$(_rgb 147 153 178)
    THEME_OVERLAY1=$(_rgb 127 132 156)
    THEME_OVERLAY0=$(_rgb 108 112 134)
    THEME_SURFACE2=$(_rgb 88 91 112)
    THEME_SURFACE1=$(_rgb 69 71 90)
    THEME_SURFACE0=$(_rgb 49 50 68)
    THEME_BASE=$(_rgb 30 30 46)
    THEME_MANTLE=$(_rgb 24 24 37)
    THEME_CRUST=$(_rgb 17 17 27)
else
    # 256-color fallback
    THEME_ROSEWATER='\033[38;5;224m'
    THEME_FLAMINGO='\033[38;5;224m'
    THEME_PINK='\033[38;5;218m'
    THEME_MAUVE='\033[38;5;183m'
    THEME_RED='\033[38;5;211m'
    THEME_MAROON='\033[38;5;217m'
    THEME_PEACH='\033[38;5;223m'
    THEME_YELLOW='\033[38;5;223m'
    THEME_GREEN='\033[38;5;157m'
    THEME_TEAL='\033[38;5;158m'
    THEME_SKY='\033[38;5;117m'
    THEME_SAPPHIRE='\033[38;5;117m'
    THEME_BLUE='\033[38;5;117m'
    THEME_LAVENDER='\033[38;5;183m'
    THEME_TEXT='\033[38;5;189m'
    THEME_SUBTEXT1='\033[38;5;189m'
    THEME_SUBTEXT0='\033[38;5;146m'
    THEME_OVERLAY2='\033[38;5;146m'
    THEME_OVERLAY1='\033[38;5;102m'
    THEME_OVERLAY0='\033[38;5;102m'
    THEME_SURFACE2='\033[38;5;60m'
    THEME_SURFACE1='\033[38;5;60m'
    THEME_SURFACE0='\033[38;5;238m'
    THEME_BASE='\033[38;5;235m'
    THEME_MANTLE='\033[38;5;234m'
    THEME_CRUST='\033[38;5;233m'
fi

# Reset and formatting codes
THEME_RESET='\033[0m'
THEME_BOLD='\033[1m'

# =============================================================================
# Semantic Color Aliases
# =============================================================================

THEME_GIT_CLEAN="$THEME_GREEN"
THEME_GIT_DIRTY="$THEME_PEACH"
THEME_GIT_AHEAD="$THEME_RED"
THEME_GIT_UNKNOWN="$THEME_MAUVE"
THEME_PATH="$THEME_BLUE"
THEME_USER="$THEME_TEAL"
THEME_HOST="$THEME_SAPPHIRE"
THEME_PROMPT="$THEME_TEXT"
THEME_WARNING="$THEME_YELLOW"
THEME_ERROR="$THEME_RED"
THEME_SUCCESS="$THEME_GREEN"
THEME_INFO="$THEME_BLUE"

# =============================================================================
# LS Colors
# =============================================================================

case "$CURRENT_PLATFORM" in
    darwin)
        export LSCOLORS="GxFxCxDxBxegedabagaced"
        ;;
    linux)
        export LS_COLORS="di=1;34:ln=1;35:so=1;32:pi=33:ex=1;31:bd=34;46:cd=34;43:su=30;41:sg=30;46:tw=30;42:ow=30;43"
        ;;
esac

# =============================================================================
# Export All Theme Variables
# =============================================================================

export THEME_ROSEWATER THEME_FLAMINGO THEME_PINK THEME_MAUVE
export THEME_RED THEME_MAROON THEME_PEACH THEME_YELLOW
export THEME_GREEN THEME_TEAL THEME_SKY THEME_SAPPHIRE
export THEME_BLUE THEME_LAVENDER THEME_TEXT THEME_SUBTEXT1
export THEME_SUBTEXT0 THEME_OVERLAY2 THEME_OVERLAY1 THEME_OVERLAY0
export THEME_SURFACE2 THEME_SURFACE1 THEME_SURFACE0
export THEME_BASE THEME_MANTLE THEME_CRUST
export THEME_RESET THEME_BOLD
export THEME_GIT_CLEAN THEME_GIT_DIRTY THEME_GIT_AHEAD THEME_GIT_UNKNOWN
export THEME_PATH THEME_USER THEME_HOST THEME_PROMPT
export THEME_WARNING THEME_ERROR THEME_SUCCESS THEME_INFO
