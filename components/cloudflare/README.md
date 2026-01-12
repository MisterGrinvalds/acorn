# CloudFlare Component

Shell integration for CloudFlare CLI (wrangler) providing aliases for Workers, Pages, R2, KV, and D1.

## Installation

```bash
# Install wrangler via npm
./install/install.sh

# Or manually
npm install -g wrangler
```

## Quick Start

```bash
# Login to CloudFlare
wrlogin

# Check current account
wrwhoami
```

## Aliases

### Core Wrangler

| Alias | Command | Description |
|-------|---------|-------------|
| `wr` | `wrangler` | Main wrangler command |
| `wrd` | `wrangler dev` | Start local dev server |
| `wrp` | `wrangler pages` | Pages commands |
| `wrr2` | `wrangler r2` | R2 storage commands |
| `wrkv` | `wrangler kv` | KV store commands |
| `wrd1` | `wrangler d1` | D1 database commands |

### Workers

| Alias | Command | Description |
|-------|---------|-------------|
| `wrlist` | `wrangler deployments list` | List deployments |
| `wrtail` | `wrangler tail` | Tail worker logs |
| `wrpub` | `wrangler deploy` | Deploy worker |

### Pages

| Alias | Command | Description |
|-------|---------|-------------|
| `wrplist` | `wrangler pages project list` | List Pages projects |
| `wrpdeploy` | `wrangler pages deploy` | Deploy to Pages |

### Authentication

| Alias | Command | Description |
|-------|---------|-------------|
| `wrlogin` | `wrangler login` | Login to CloudFlare |
| `wrlogout` | `wrangler logout` | Logout |
| `wrwhoami` | `wrangler whoami` | Show current account |

## Claude Code Integration

This component includes AI assistance for CloudFlare development:

### Agent

- `cloudflare-expert` - Expert guidance on Wrangler CLI, Workers, Pages, R2, KV, and D1

### Commands

| Command | Description |
|---------|-------------|
| `/cf-coach` | Interactive coaching session for CloudFlare development |
| `/cf-explain` | Explain CloudFlare concepts and services |
| `/cf-worker-deploy` | Deploy a CloudFlare Worker |
| `/cf-pages-setup` | Set up a CloudFlare Pages project |
| `/cf-r2-manage` | Manage R2 object storage |
| `/cf-kv-manage` | Manage KV key-value storage |

## Component Structure

```
cloudflare/
├── config.yaml           # Aliases configuration
├── README.md             # This file
└── install/
    └── install.sh        # Wrangler installation script
```

Claude Code AI integration files are in:
- Agent: `ai/agents/cloudflare-expert.md`
- Commands: `ai/commands/cloudflare/`

## Dependencies

- **Required tools**: `wrangler` (via npm)
- **Required components**: `shell`
