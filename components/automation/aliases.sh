#!/bin/sh
# components/automation/aliases.sh - Automation framework aliases

# Main CLI
alias auto='$AUTO_CLI'

# Module shortcuts
alias autodev='auto dev'
alias autok8s='auto k8s'
alias autogithub='auto github'
alias autosystem='auto system'
alias autoconfig='auto config'
alias autocloud='auto cloud'

# Quick actions
alias ainit='auto dev init'
alias adeploy='auto k8s deploy'
alias apr='auto github pr create'
alias abackup='auto system backup'
alias acleanup='auto system cleanup'
