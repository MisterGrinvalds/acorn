---
description: Explain CloudFlare concepts, services, and workflows
argument-hint: [topic]
allowed-tools: Read, Glob, Grep
---

## Task

Explain the requested topic about CloudFlare. If no specific topic provided, give an overview of CloudFlare services and the dotfiles integration.

## Topics

If `$ARGUMENTS` is provided, explain that specific topic:

- **workers** - CloudFlare Workers serverless functions
- **pages** - CloudFlare Pages static hosting
- **r2** - R2 object storage
- **kv** - KV key-value storage
- **d1** - D1 SQLite databases
- **wrangler** - Wrangler CLI basics
- **config** - wrangler.toml configuration
- **secrets** - Managing worker secrets
- **bindings** - Service bindings (KV, R2, D1)
- **routing** - URL routing and patterns
- **caching** - Edge caching strategies
- **queues** - CloudFlare Queues
- **durable-objects** - Durable Objects for state

## Context

Reference these files for accurate information:
@components/cloudflare/component.yaml
@components/cloudflare/functions.sh
@components/cloudflare/aliases.sh

## Response Format

1. **Concept overview** - What it is and why it matters
2. **Key commands** - Essential wrangler commands
3. **Dotfiles integration** - Available functions and aliases
4. **Examples** - Practical usage examples
5. **Best practices** - Production recommendations
