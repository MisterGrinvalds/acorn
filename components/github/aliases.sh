#!/bin/sh
# components/github/aliases.sh - GitHub CLI aliases

# Pull Requests
alias ghpr='gh pr create'
alias ghprs='gh pr status'
alias ghprv='gh pr view'
alias ghprc='gh pr checkout'
alias ghprm='gh pr merge'
alias ghprl='gh pr list'

# Issues
alias ghissue='gh issue create'
alias ghissues='gh issue list'
alias ghissuev='gh issue view'

# Repository
alias ghrepo='gh repo view'
alias ghrepoc='gh repo clone'
alias ghrepof='gh repo fork'

# Actions/Workflow
alias ghrun='gh run list'
alias ghrunv='gh run view'
alias ghrunw='gh run watch'

# Browse
alias ghweb='gh browse'
