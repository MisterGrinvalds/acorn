#!/bin/sh
# components/cloudflare/functions.sh - CloudFlare CLI helper functions
#
# This file is sourced only for INTERACTIVE shells.
# Provides helper functions for CloudFlare Workers, Pages, R2, KV, and D1.

# Check CloudFlare CLI status and installation
cf_status() {
    echo "${CLR_BLUE}CloudFlare CLI Status${CLR_RESET}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    if command -v wrangler >/dev/null 2>&1; then
        echo "${CLR_GREEN}✓${CLR_RESET} wrangler installed: $(wrangler --version 2>/dev/null | head -1)"
    else
        echo "${CLR_RED}✗${CLR_RESET} wrangler not found"
        echo "  Install: npm install -g wrangler"
        return 1
    fi

    echo ""
    echo "${CLR_BLUE}Authentication:${CLR_RESET}"
    if wrangler whoami 2>/dev/null | grep -q "You are logged in"; then
        wrangler whoami 2>/dev/null | head -5
    else
        echo "${CLR_YELLOW}Not logged in${CLR_RESET}"
        echo "  Run: wrangler login"
    fi
}

# Show current CloudFlare account
cf_whoami() {
    wrangler whoami
}

# List all Workers
cf_workers() {
    echo "${CLR_BLUE}CloudFlare Workers${CLR_RESET}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    wrangler deployments list 2>/dev/null || echo "No workers found or not authenticated"
}

# List Pages projects
cf_pages() {
    echo "${CLR_BLUE}CloudFlare Pages Projects${CLR_RESET}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    wrangler pages project list 2>/dev/null || echo "No pages projects found or not authenticated"
}

# List R2 buckets
cf_r2_buckets() {
    echo "${CLR_BLUE}CloudFlare R2 Buckets${CLR_RESET}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    wrangler r2 bucket list 2>/dev/null || echo "No R2 buckets found or not authenticated"
}

# List KV namespaces
cf_kv_namespaces() {
    echo "${CLR_BLUE}CloudFlare KV Namespaces${CLR_RESET}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    wrangler kv namespace list 2>/dev/null || echo "No KV namespaces found or not authenticated"
}

# List D1 databases
cf_d1_databases() {
    echo "${CLR_BLUE}CloudFlare D1 Databases${CLR_RESET}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    wrangler d1 list 2>/dev/null || echo "No D1 databases found or not authenticated"
}

# Tail worker logs
# Usage: cf_logs [worker-name]
cf_logs() {
    local worker="${1:-}"
    if [ -z "$worker" ]; then
        echo "Usage: cf_logs <worker-name>"
        echo ""
        echo "Available workers:"
        wrangler deployments list 2>/dev/null
        return 1
    fi
    wrangler tail "$worker"
}

# Deploy current worker
cf_deploy() {
    if [ ! -f "wrangler.toml" ] && [ ! -f "wrangler.json" ]; then
        echo "${CLR_RED}Error:${CLR_RESET} No wrangler.toml or wrangler.json found"
        echo "Initialize with: wrangler init"
        return 1
    fi
    wrangler deploy "$@"
}

# Create a new Worker project
# Usage: cf_worker_init [name]
cf_worker_init() {
    local name="${1:-my-worker}"
    wrangler init "$name"
}

# Create a new Pages project
# Usage: cf_pages_init [name]
cf_pages_init() {
    local name="${1:-my-pages-site}"
    wrangler pages project create "$name"
}

# Create a new R2 bucket
# Usage: cf_r2_create <bucket-name>
cf_r2_create() {
    local bucket="${1:-}"
    if [ -z "$bucket" ]; then
        echo "Usage: cf_r2_create <bucket-name>"
        return 1
    fi
    wrangler r2 bucket create "$bucket"
}

# Create a new KV namespace
# Usage: cf_kv_create <namespace-name>
cf_kv_create() {
    local namespace="${1:-}"
    if [ -z "$namespace" ]; then
        echo "Usage: cf_kv_create <namespace-name>"
        return 1
    fi
    wrangler kv namespace create "$namespace"
}

# Create a new D1 database
# Usage: cf_d1_create <database-name>
cf_d1_create() {
    local database="${1:-}"
    if [ -z "$database" ]; then
        echo "Usage: cf_d1_create <database-name>"
        return 1
    fi
    wrangler d1 create "$database"
}

# Put a secret for current worker
# Usage: cf_secret_put <secret-name>
cf_secret_put() {
    local secret="${1:-}"
    if [ -z "$secret" ]; then
        echo "Usage: cf_secret_put <secret-name>"
        echo "You will be prompted to enter the secret value"
        return 1
    fi
    wrangler secret put "$secret"
}

# List secrets for current worker
cf_secrets() {
    wrangler secret list
}

# Show overview of all CloudFlare resources
cf_overview() {
    echo "${CLR_BLUE}CloudFlare Overview${CLR_RESET}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    cf_whoami 2>/dev/null || true
    echo ""

    echo "${CLR_BLUE}Workers:${CLR_RESET}"
    wrangler deployments list 2>/dev/null | head -10 || echo "  None found"
    echo ""

    echo "${CLR_BLUE}Pages:${CLR_RESET}"
    wrangler pages project list 2>/dev/null | head -10 || echo "  None found"
    echo ""

    echo "${CLR_BLUE}R2 Buckets:${CLR_RESET}"
    wrangler r2 bucket list 2>/dev/null | head -10 || echo "  None found"
    echo ""

    echo "${CLR_BLUE}KV Namespaces:${CLR_RESET}"
    wrangler kv namespace list 2>/dev/null | head -10 || echo "  None found"
    echo ""

    echo "${CLR_BLUE}D1 Databases:${CLR_RESET}"
    wrangler d1 list 2>/dev/null | head -10 || echo "  None found"
}

# Show all CloudFlare functions
cf_help() {
    echo "${CLR_BLUE}CloudFlare Component Functions${CLR_RESET}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "${CLR_GREEN}Status & Info:${CLR_RESET}"
    echo "  cf_status           Check CLI status and auth"
    echo "  cf_whoami           Show current account"
    echo "  cf_overview         Overview of all resources"
    echo ""
    echo "${CLR_GREEN}List Resources:${CLR_RESET}"
    echo "  cf_workers          List Workers"
    echo "  cf_pages            List Pages projects"
    echo "  cf_r2_buckets       List R2 buckets"
    echo "  cf_kv_namespaces    List KV namespaces"
    echo "  cf_d1_databases     List D1 databases"
    echo ""
    echo "${CLR_GREEN}Create Resources:${CLR_RESET}"
    echo "  cf_worker_init      Create new Worker project"
    echo "  cf_pages_init       Create new Pages project"
    echo "  cf_r2_create        Create R2 bucket"
    echo "  cf_kv_create        Create KV namespace"
    echo "  cf_d1_create        Create D1 database"
    echo ""
    echo "${CLR_GREEN}Operations:${CLR_RESET}"
    echo "  cf_deploy           Deploy current worker"
    echo "  cf_logs <worker>    Tail worker logs"
    echo "  cf_secret_put       Add worker secret"
    echo "  cf_secrets          List worker secrets"
    echo ""
    echo "${CLR_GREEN}Aliases:${CLR_RESET}"
    echo "  wr, wrd, wrp, wrr2, wrkv, wrd1"
    echo "  wrlogin, wrlogout, wrwhoami"
}
