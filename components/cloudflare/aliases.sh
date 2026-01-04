#!/bin/sh
# components/cloudflare/aliases.sh - CloudFlare CLI aliases
#
# This file is sourced only for INTERACTIVE shells.
# Provides short aliases for common wrangler commands.

# Core wrangler aliases
alias wr='wrangler'
alias wrd='wrangler dev'
alias wrp='wrangler pages'
alias wrr2='wrangler r2'
alias wrkv='wrangler kv'

# Workers management
alias wrlist='wrangler deployments list'
alias wrtail='wrangler tail'
alias wrpub='wrangler deploy'

# Pages management
alias wrplist='wrangler pages project list'
alias wrpdeploy='wrangler pages deploy'

# R2 storage
alias wrr2list='wrangler r2 bucket list'

# KV storage
alias wrkvlist='wrangler kv namespace list'

# D1 database
alias wrd1='wrangler d1'
alias wrd1list='wrangler d1 list'

# Secrets management
alias wrsecret='wrangler secret'
alias wrsecrets='wrangler secret list'

# Login/logout
alias wrlogin='wrangler login'
alias wrlogout='wrangler logout'
alias wrwhoami='wrangler whoami'
