---
description: Interactive coaching session to learn CloudFlare development
argument-hint: [skill-level: beginner|intermediate|advanced]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Guide the user through learning CloudFlare development interactively.

## Approach

1. **Assess level** - Ask about CloudFlare/serverless experience
2. **Set goals** - Identify what they want to build
3. **Progressive exercises** - Start simple, build complexity
4. **Real-time practice** - Have them run wrangler commands
5. **Reinforce** - Summarize and suggest next steps

## Skill Levels

### Beginner
- What is edge computing?
- Installing wrangler
- Authentication
- Creating first Worker
- Local development

### Intermediate
- Worker configuration
- Using KV and R2
- Pages deployment
- Environment variables
- Using dotfiles functions

### Advanced
- D1 databases
- Durable Objects
- Queues
- Complex routing
- Performance optimization

## Interactive Exercises

### Beginner Exercises
```bash
# Exercise 1: Check status
cf_status  # dotfiles function

# Exercise 2: Login
wrlogin  # or: wrangler login

# Exercise 3: Create first worker
cf_worker_init hello-world
cd hello-world

# Exercise 4: Local development
wrd  # or: wrangler dev

# Exercise 5: Deploy
cf_deploy  # dotfiles function
```

### Intermediate Exercises
```bash
# Exercise 6: List your resources
cf_overview  # dotfiles function

# Exercise 7: Create KV namespace
cf_kv_create MY_CACHE

# Exercise 8: View logs
cf_logs hello-world  # dotfiles function

# Exercise 9: Add secret
cf_secret_put API_KEY
```

### Advanced Exercises
```bash
# Exercise 10: Create D1 database
cf_d1_create my-app-db

# Exercise 11: Create R2 bucket
cf_r2_create my-files

# Exercise 12: Configure bindings in wrangler.toml
```

## Context

@components/cloudflare/functions.sh
@components/cloudflare/aliases.sh

## Coaching Style

- Start with authentication
- Use local development first
- Deploy to staging before production
- Show dotfiles functions for common tasks
- Build toward full-stack edge app
