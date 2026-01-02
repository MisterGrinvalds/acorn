#!/bin/sh
# components/huggingface/env.sh - Hugging Face environment variables

# XDG-compliant cache location for Hugging Face models
export HF_HOME="${XDG_CACHE_HOME:-$HOME/.cache}/huggingface"
export TRANSFORMERS_CACHE="$HF_HOME/transformers"
