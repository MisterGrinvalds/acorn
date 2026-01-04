#!/bin/sh
# components/cloudflare/completions.sh - Tab completions for wrangler
#
# This file is sourced only for INTERACTIVE shells.
# Loads tab completions for wrangler CLI.

# Wrangler doesn't have built-in shell completion generation,
# but we can provide basic completions for common subcommands

if command -v wrangler >/dev/null 2>&1; then
    case "$CURRENT_SHELL" in
        bash)
            # Basic bash completion for wrangler
            _wrangler_completions() {
                local cur="${COMP_WORDS[COMP_CWORD]}"
                local prev="${COMP_WORDS[COMP_CWORD-1]}"

                # Main wrangler subcommands
                local commands="init dev deploy delete tail secret kv r2 d1 pages queues pubsub dispatch-namespace mtls-certificate login logout whoami types generate versions deployments triggers"

                case "$prev" in
                    wrangler|wr)
                        COMPREPLY=($(compgen -W "$commands" -- "$cur"))
                        ;;
                    secret)
                        COMPREPLY=($(compgen -W "put delete list bulk" -- "$cur"))
                        ;;
                    kv)
                        COMPREPLY=($(compgen -W "namespace key bulk" -- "$cur"))
                        ;;
                    namespace)
                        COMPREPLY=($(compgen -W "create delete list" -- "$cur"))
                        ;;
                    r2)
                        COMPREPLY=($(compgen -W "bucket object" -- "$cur"))
                        ;;
                    bucket)
                        COMPREPLY=($(compgen -W "create delete list info" -- "$cur"))
                        ;;
                    d1)
                        COMPREPLY=($(compgen -W "create delete list info execute export migrations backup time-travel" -- "$cur"))
                        ;;
                    pages)
                        COMPREPLY=($(compgen -W "project deployment dev deploy download" -- "$cur"))
                        ;;
                    project)
                        COMPREPLY=($(compgen -W "create list delete" -- "$cur"))
                        ;;
                    deployments)
                        COMPREPLY=($(compgen -W "list view" -- "$cur"))
                        ;;
                    *)
                        COMPREPLY=()
                        ;;
                esac
            }
            complete -F _wrangler_completions wrangler
            complete -F _wrangler_completions wr
            ;;
        zsh)
            # Basic zsh completion for wrangler
            if [ -n "$ZSH_VERSION" ]; then
                _wrangler() {
                    local -a commands
                    commands=(
                        'init:Create a new Worker project'
                        'dev:Start a local development server'
                        'deploy:Deploy your Worker to CloudFlare'
                        'delete:Delete your Worker from CloudFlare'
                        'tail:Stream logs from a deployed Worker'
                        'secret:Manage Worker secrets'
                        'kv:Interact with KV namespaces'
                        'r2:Interact with R2 buckets'
                        'd1:Interact with D1 databases'
                        'pages:Interact with CloudFlare Pages'
                        'queues:Interact with Queues'
                        'login:Login to CloudFlare'
                        'logout:Logout from CloudFlare'
                        'whoami:Show current CloudFlare user'
                        'deployments:Manage deployments'
                    )

                    _arguments -C \
                        '1:command:->command' \
                        '*::arg:->args'

                    case "$state" in
                        command)
                            _describe 'command' commands
                            ;;
                    esac
                }
                compdef _wrangler wrangler
                compdef _wrangler wr
            fi
            ;;
    esac
fi
