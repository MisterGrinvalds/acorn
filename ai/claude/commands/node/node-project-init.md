---
description: Initialize a new Node.js/TypeScript project
argument-hint: <project-name>
allowed-tools: Read, Write, Bash
---

## Task

Help the user initialize a new Node.js project with TypeScript and modern tooling.

## Quick Start

Using dotfiles function:
```bash
node_init my-project
# Creates TypeScript project with pnpm
```

## Manual Setup

### Step 1: Create Directory
```bash
mkdir my-project
cd my-project
```

### Step 2: Initialize Package
```bash
# With pnpm (recommended)
pnpm init

# With npm
npm init -y
```

### Step 3: Add TypeScript
```bash
# With pnpm
pnpm add -D typescript @types/node tsx

# With npm
npm install -D typescript @types/node tsx
```

### Step 4: Configure TypeScript

Create `tsconfig.json`:
```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "outDir": "./dist",
    "rootDir": "./src",
    "declaration": true,
    "sourceMap": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
```

### Step 5: Create Structure
```bash
mkdir src tests
echo 'console.log("Hello, TypeScript!");' > src/index.ts
```

### Step 6: Add Scripts

Update `package.json`:
```json
{
  "type": "module",
  "scripts": {
    "dev": "tsx watch src/index.ts",
    "build": "tsc",
    "start": "node dist/index.js",
    "test": "vitest",
    "lint": "eslint src/",
    "typecheck": "tsc --noEmit"
  }
}
```

## Project Types

### CLI Application
```bash
pnpm add commander chalk
pnpm add -D @types/node tsx
```

### Express API
```bash
pnpm add express
pnpm add -D @types/express tsx
```

### Fastify API
```bash
pnpm add fastify
pnpm add -D tsx
```

## Recommended Dev Dependencies

```bash
# Linting & Formatting
pnpm add -D eslint prettier eslint-config-prettier

# Testing
pnpm add -D vitest

# Build
pnpm add -D tsup  # For libraries
```

## .nvmrc File

Pin Node version:
```bash
node --version > .nvmrc
# Creates file with current version
```

## .gitignore

```gitignore
node_modules/
dist/
.env
*.log
.DS_Store
coverage/
```

## Verification

```bash
# Run in dev mode
pnpm dev
# or: prd (alias)

# Build
pnpm build
# or: prb (alias)

# Type check
pnpm typecheck
```
