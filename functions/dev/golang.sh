#!/bin/sh
# Go Development Tools

# Go environment setup
export GOPATH="$HOME/go"
export GOBIN="$GOPATH/bin"
export GO111MODULE=on

# Add Go binary paths to PATH
case ":$PATH:" in
    *":$GOBIN:"*) ;;
    *) PATH="$GOBIN:$PATH" ;;
esac
case ":$PATH:" in
    *":/usr/local/go/bin:"*) ;;
    *) PATH="/usr/local/go/bin:$PATH" ;;
esac

# Go aliases
alias gob='go build'
alias gor='go run'
alias got='go test'
alias gotv='go test -v'
alias gom='go mod'
alias gomi='go mod init'
alias gomt='go mod tidy'
alias gomd='go mod download'
alias gog='go get'
alias gou='go get -u'
alias gof='go fmt ./...'
alias gov='go vet ./...'
alias goi='go install'
alias gover='go version'

# Cobra CLI helpers (for building CLI tools)
cobranew() {
    if [ -z "$1" ]; then
        echo "Usage: cobranew <app-name>"
        return 1
    fi
    cobra-cli init "$1"
    cd "$1" || return 1
    go mod tidy
}

cobradd() {
    if [ -z "$1" ]; then
        echo "Usage: cobradd <command-name>"
        return 1
    fi
    cobra-cli add "$1"
}

# Go project initialization helper
gonew() {
    if [ -z "$1" ]; then
        echo "Usage: gonew <module-name>"
        return 1
    fi
    mkdir -p "$1"
    cd "$1" || return 1
    go mod init "$1"
    cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
EOF
    echo "Go project '$1' initialized!"
}

# Go testing helpers
gotest() {
    if [ -z "$1" ]; then
        go test ./...
    else
        go test ./... -run "$1"
    fi
}

gotestcover() {
    go test ./... -coverprofile=coverage.out
    go tool cover -html=coverage.out -o coverage.html
    echo "Coverage report generated: coverage.html"
}
