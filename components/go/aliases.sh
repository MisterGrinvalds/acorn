#!/bin/sh
# components/go/aliases.sh - Go aliases

# Build and run
alias gob='go build'
alias gor='go run'
alias goi='go install'

# Testing
alias got='go test'
alias gotv='go test -v'
alias gotc='go test -cover'

# Module management
alias gom='go mod'
alias gomi='go mod init'
alias gomt='go mod tidy'
alias gomd='go mod download'

# Dependencies
alias gog='go get'
alias gou='go get -u'

# Code quality
alias gof='go fmt ./...'
alias gov='go vet ./...'

# Info
alias gover='go version'
alias goenv='go env'
