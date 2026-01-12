---
description: Deploy a CloudFlare Worker
argument-hint: [--env staging|production]
allowed-tools: Read, Bash
---

## Task

Help the user deploy a CloudFlare Worker.

## Quick Deploy

Using dotfiles function:
```bash
cf_deploy
# Deploys current worker from wrangler.toml
```

## Deployment Commands

### Basic Deploy
```bash
wrangler deploy
wrpub  # alias
```

### Deploy to Environment
```bash
wrangler deploy --env staging
wrangler deploy --env production
```

### Deploy with Dry Run
```bash
wrangler deploy --dry-run
# Shows what would be deployed
```

## Pre-Deploy Checklist

### 1. Verify Authentication
```bash
cf_status  # or: wrangler whoami
```

### 2. Check Configuration
```bash
# Ensure wrangler.toml exists
cat wrangler.toml
```

### 3. Test Locally
```bash
wrangler dev
wrd  # alias
# Test at http://localhost:8787
```

### 4. Build (if needed)
```bash
npm run build  # If custom build step
```

## wrangler.toml Example

```toml
name = "my-worker"
main = "src/index.ts"
compatibility_date = "2024-01-01"

# Production environment
[env.production]
name = "my-worker-prod"
route = "example.com/*"

# Staging environment
[env.staging]
name = "my-worker-staging"
route = "staging.example.com/*"

# Environment variables
[vars]
ENVIRONMENT = "development"

[env.production.vars]
ENVIRONMENT = "production"
```

## Managing Secrets

```bash
# Add secret (prompted for value)
wrangler secret put API_KEY
cf_secret_put API_KEY  # dotfiles function

# Add secret for specific environment
wrangler secret put API_KEY --env production

# List secrets
wrangler secret list
cf_secrets  # dotfiles function

# Delete secret
wrangler secret delete API_KEY
```

## Post-Deploy Verification

### Check Deployment
```bash
# List deployments
wrangler deployments list
wrlist  # alias
cf_workers  # dotfiles function

# View specific deployment
wrangler deployments view
```

### View Logs
```bash
wrangler tail my-worker
cf_logs my-worker  # dotfiles function
```

### Rollback
```bash
# View deployment history
wrangler deployments list

# Rollback to previous
wrangler rollback
```

## CI/CD Integration

### GitHub Actions
```yaml
name: Deploy Worker
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: cloudflare/wrangler-action@v3
        with:
          apiToken: ${{ secrets.CF_API_TOKEN }}
```

## Dotfiles Integration

- `cf_deploy` - Deploy current worker
- `cf_status` - Check auth status
- `cf_logs <worker>` - Tail logs
- `cf_secret_put <name>` - Add secret
- `cf_secrets` - List secrets
- `wrpub` - wrangler deploy (alias)
- `wrd` - wrangler dev (alias)
