#!/bin/sh
# components/go/env.sh - Go environment variables

# Go environment setup
export GOPATH="$HOME/go"
export GOBIN="$GOPATH/bin"
export GO111MODULE=on

# Add Go binary paths to PATH
case ":$PATH:" in
    *":$GOBIN:"*) ;;
    *) export PATH="$GOBIN:$PATH" ;;
esac

case ":$PATH:" in
    *":/usr/local/go/bin:"*) ;;
    *) export PATH="/usr/local/go/bin:$PATH" ;;
esac

# macOS Homebrew Go location
if [ -d "/opt/homebrew/opt/go/bin" ]; then
    case ":$PATH:" in
        *":/opt/homebrew/opt/go/bin:"*) ;;
        *) export PATH="/opt/homebrew/opt/go/bin:$PATH" ;;
    esac
fi
