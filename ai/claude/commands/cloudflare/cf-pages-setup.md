---
description: Set up and deploy a CloudFlare Pages project
argument-hint: <project-name>
allowed-tools: Read, Write, Bash
---

## Task

Help the user set up and deploy a CloudFlare Pages project.

## Quick Setup

Using dotfiles function:
```bash
cf_pages_init my-site
```

## Manual Setup

### Create Pages Project
```bash
wrangler pages project create my-site
```

### Deploy Directory
```bash
# Deploy dist folder
wrangler pages deploy ./dist

# Using alias
wrpdeploy ./dist
```

## Framework Setup

### React / Vite
```bash
npm create vite@latest my-site -- --template react-ts
cd my-site
npm install
npm run build
wrangler pages deploy ./dist
```

### Next.js
```bash
npx create-next-app@latest my-site
cd my-site
npm run build
wrangler pages deploy ./out  # Static export
```

### Astro
```bash
npm create astro@latest my-site
cd my-site
npm run build
wrangler pages deploy ./dist
```

### Plain HTML
```bash
mkdir my-site
cd my-site
echo "<h1>Hello</h1>" > index.html
wrangler pages deploy .
```

## Configuration

### wrangler.toml for Pages
```toml
name = "my-site"
pages_build_output_dir = "./dist"

# Build configuration
[build]
command = "npm run build"

# Environment variables
[env.production.vars]
API_URL = "https://api.example.com"

[env.preview.vars]
API_URL = "https://api-staging.example.com"
```

### _redirects File
```
# Redirect old paths
/old-path /new-path 301

# SPA fallback
/* /index.html 200
```

### _headers File
```
# Security headers
/*
  X-Frame-Options: DENY
  X-Content-Type-Options: nosniff

# Cache static assets
/assets/*
  Cache-Control: public, max-age=31536000
```

## Pages Functions

Add serverless functions:
```
my-site/
├── functions/
│   ├── api/
│   │   └── hello.ts    # /api/hello
│   └── _middleware.ts  # Runs on all routes
├── public/
└── src/
```

### Function Example
```typescript
// functions/api/hello.ts
export const onRequest: PagesFunction = async (context) => {
  return new Response(JSON.stringify({ message: "Hello!" }), {
    headers: { "Content-Type": "application/json" },
  });
};
```

## Deployment Options

### Direct Deploy
```bash
wrangler pages deploy ./dist --project-name=my-site
```

### With Branch
```bash
wrangler pages deploy ./dist --branch=staging
```

### Production Deploy
```bash
wrangler pages deploy ./dist --branch=main
```

## Environment Variables

### Via CLI
```bash
# Set production variable
wrangler pages secret put API_KEY --project-name=my-site

# Set preview variable
wrangler pages secret put API_KEY --project-name=my-site --env=preview
```

### Via Dashboard
- Go to Pages project → Settings → Environment variables

## Custom Domains

```bash
# Add custom domain
wrangler pages project list
# Then configure via dashboard or API
```

## List Projects

```bash
wrangler pages project list
cf_pages  # dotfiles function
wrplist   # alias
```

## Dotfiles Integration

- `cf_pages_init [name]` - Create Pages project
- `cf_pages` - List Pages projects
- `wrp` - wrangler pages (alias)
- `wrplist` - wrangler pages project list (alias)
- `wrpdeploy` - wrangler pages deploy (alias)
