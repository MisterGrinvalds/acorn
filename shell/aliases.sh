#!/bin/sh
# Shell-portable aliases
# Requires: shell/discovery.sh

# Shell reload
alias resource='source ~/.bashrc 2>/dev/null || source ~/.zshrc 2>/dev/null'

# Git shortcuts
alias g='git'
alias gs='git status'
alias ga='git add'
alias gc='git commit'
alias gp='git push'
alias gl='git pull'
alias gd='git diff'
alias gco='git checkout'
alias gb='git branch'
alias glog='git log --oneline --graph --decorate'

# Navigation
alias ..='cd ..'
alias ...='cd ../..'
alias ....='cd ../../..'
alias ll='ls -alh'
alias llr='ls -alhr'
alias lls='ls -alhS'
alias llsr='ls -alhSr'
alias lld='ls -alht'
alias lldr='ls -alhtr'
alias mkdir='mkdir -pv'

# Python development
alias pip='python -m pip'
alias pip3='python3 -m pip'
alias py='python'
alias py3='python3'
alias ipy='ipython'
alias ptest='python -m pytest'
alias ptestv='python -m pytest -v'
alias black='python -m black'
alias isort='python -m isort'
alias flake8='python -m flake8'

# FastAPI development
alias uvdev='uvicorn main:app --reload'
alias uvprod='uvicorn main:app --host 0.0.0.0 --port 8000'

# Tmux enhanced
alias ta='tmux attach-session -t'
alias tn='tmux new-session -s'
alias tk='tmux kill-session -t'
alias tko='tmux kill-session -a'
alias ti='tmux info'
alias ts='tmux list-sessions'
alias tks='tmux kill-server'
alias td='tmux detach'

# Tmux project sessions
alias twork='tmux new-session -s work -d'
alias tdev='tmux new-session -s dev -d'
alias tk8s='tmux new-session -s k8s -d'

# Platform-specific aliases
case "$CURRENT_PLATFORM" in
    darwin)
        alias getsshkey='pbcopy < ~/.ssh/id_rsa.pub'
        alias perm="stat -f '%Lp'"
        alias lldc='ls -alhtU'       # List by date created (macOS only)
        alias lldcr='ls -alhtUr'
        ;;
    linux)
        alias getsshkey='xclip -selection clipboard < ~/.ssh/id_rsa.pub'
        alias perm='stat -c "%a"'
        # xclip shortcuts
        alias c='xclip'
        alias cs='xclip -selection clipboard'
        alias v='xclip -o'
        alias vs='xclip -o -selection clipboard'
        ;;
esac

# Tree alternative using ls
alias tree='ls -R | grep ":$" | sed -e "s/:$//" -e "s/[^-][^\/]*\//--/g" -e "s/^/   /" -e "s/-/|/"'

# =============================================================================
# Kubernetes
# =============================================================================
alias k='kubectl'
alias kgp='kubectl get pods'
alias kgpa='kubectl get pods --all-namespaces'
alias kgs='kubectl get svc'
alias kgn='kubectl get nodes'
alias kga='kubectl get all'
alias kgaa='kubectl get all --all-namespaces'
alias kd='kubectl describe'
alias kdp='kubectl describe pod'
alias kds='kubectl describe svc'
alias kdn='kubectl describe node'
alias kl='kubectl logs'
alias klf='kubectl logs -f'
alias kex='kubectl exec -it'
alias kaf='kubectl apply -f'
alias kdf='kubectl delete -f'
alias kctx='kubectl config use-context'
alias kns='kubectl config set-context --current --namespace'
alias ktop='kubectl top pods'
alias ktopn='kubectl top nodes'

# Helm
alias hls='helm list'
alias hlsa='helm list --all-namespaces'
alias hi='helm install'
alias hu='helm upgrade'
alias hd='helm delete'
alias hs='helm status'
alias hh='helm history'
alias hr='helm rollback'
alias ht='helm template'

# k9s (TUI)
alias k9='k9s'
alias k9a='k9s --all-namespaces'

# ArgoCD
alias argocd-login='argocd login --grpc-web'
alias argocd-apps='argocd app list'
alias argocd-sync='argocd app sync'

# Kind (Kubernetes in Docker)
alias kind-clusters='kind get clusters'
alias kind-nodes='kind get nodes'

# =============================================================================
# Terraform
# =============================================================================
alias tf='terraform'
alias tfi='terraform init'
alias tfp='terraform plan'
alias tfa='terraform apply'
alias tfaa='terraform apply -auto-approve'
alias tfd='terraform destroy'
alias tfda='terraform destroy -auto-approve'
alias tff='terraform fmt'
alias tfv='terraform validate'
alias tfs='terraform state'
alias tfsl='terraform state list'
alias tfo='terraform output'
alias tfw='terraform workspace'
alias tfwl='terraform workspace list'
alias tfws='terraform workspace select'

# =============================================================================
# Cloud CLIs
# =============================================================================
# AWS
alias awsw='aws sts get-caller-identity'
alias awsr='aws configure list'
alias awsp='aws configure list-profiles'
alias s3ls='aws s3 ls'
alias s3cp='aws s3 cp'
alias s3sync='aws s3 sync'

# Azure
alias azw='az account show'
alias azl='az login'
alias azs='az account set --subscription'
alias azls='az account list --output table'

# DigitalOcean
alias dow='doctl account get'
alias dols='doctl compute droplet list'
alias dok8s='doctl kubernetes cluster list'

# Vault (HashiCorp)
alias vst='vault status'
alias vlogin='vault login'
alias vread='vault read'
alias vwrite='vault write'
alias vlist='vault list'

# Cloudflare Tunnel
alias cft='cloudflared tunnel'
alias cftls='cloudflared tunnel list'
alias cftr='cloudflared tunnel run'

# =============================================================================
# Data Processing
# =============================================================================
alias jqc='jq -C'
alias jqr='jq -r'
alias yqc='yq -C'
alias yqr='yq -r'

# =============================================================================
# Lazygit
# =============================================================================
alias lg='lazygit'

# =============================================================================
# Shellcheck
# =============================================================================
alias sc='shellcheck'
alias scf='shellcheck -f diff'
