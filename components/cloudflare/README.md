# CloudFlare Component

Shell integration for CloudFlare CLI (wrangler) providing aliases and functions for Workers, Pages, R2, KV, and D1.

## Prerequisites

Install wrangler CLI:

```bash
npm install -g wrangler
# or use npx wrangler
```

## Quick Start

```bash
# Login to CloudFlare
wrlogin

# Check status
cf_status

# See all available functions
cf_help
```

## Aliases

| Alias | Command | Description |
|-------|---------|-------------|
| `wr` | `wrangler` | Main wrangler command |
| `wrd` | `wrangler dev` | Start local dev server |
| `wrp` | `wrangler pages` | Pages commands |
| `wrr2` | `wrangler r2` | R2 storage commands |
| `wrkv` | `wrangler kv` | KV store commands |
| `wrd1` | `wrangler d1` | D1 database commands |
| `wrlist` | `wrangler deployments list` | List deployments |
| `wrtail` | `wrangler tail` | Tail worker logs |
| `wrpub` | `wrangler deploy` | Deploy worker |
| `wrlogin` | `wrangler login` | Login to CloudFlare |
| `wrlogout` | `wrangler logout` | Logout |
| `wrwhoami` | `wrangler whoami` | Show current account |

## Functions

### Status & Info

| Function | Description |
|----------|-------------|
| `cf_status` | Check CLI installation and auth status |
| `cf_whoami` | Show current CloudFlare account |
| `cf_overview` | Overview of all CloudFlare resources |
| `cf_help` | Show all available functions |

### List Resources

| Function | Description |
|----------|-------------|
| `cf_workers` | List all Workers |
| `cf_pages` | List Pages projects |
| `cf_r2_buckets` | List R2 buckets |
| `cf_kv_namespaces` | List KV namespaces |
| `cf_d1_databases` | List D1 databases |

### Create Resources

| Function | Description |
|----------|-------------|
| `cf_worker_init [name]` | Create new Worker project |
| `cf_pages_init [name]` | Create new Pages project |
| `cf_r2_create <name>` | Create R2 bucket |
| `cf_kv_create <name>` | Create KV namespace |
| `cf_d1_create <name>` | Create D1 database |

### Operations

| Function | Description |
|----------|-------------|
| `cf_deploy` | Deploy current worker |
| `cf_logs <worker>` | Tail worker logs |
| `cf_secret_put <name>` | Add worker secret |
| `cf_secrets` | List worker secrets |

## XDG Compliance

Wrangler configuration is stored in XDG-compliant locations:
- Config: `~/.config/wrangler/`

## Dependencies

- **Required tools**: `wrangler`
- **Required components**: `shell`
