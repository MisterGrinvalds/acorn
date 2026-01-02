#!/bin/sh
# components/python/env.sh - Python environment variables

# Virtual environment location (optional central storage)
export ENVS_LOCATION="${ENVS_LOCATION:-$HOME/.virtualenvs}"

# Python startup file for interactive sessions
export PYTHONSTARTUP="${XDG_CONFIG_HOME:-$HOME/.config}/python/pythonrc"

# IPython configuration directory
export IPYTHONDIR="${XDG_CONFIG_HOME:-$HOME/.config}/ipython"

# UV cache location (XDG compliant)
export UV_CACHE_DIR="${XDG_CACHE_HOME:-$HOME/.cache}/uv"
