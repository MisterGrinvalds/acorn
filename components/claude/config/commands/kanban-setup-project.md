---
description: Initialize Go project structure for kanban API
---

Initialize a clean Go project structure for the kanban API with the following:

1. Create `go.mod` with module name `github.com/mistergrinvalds/kanban`
2. Create directory structure:
   - `cmd/api/` - Main application entry point
   - `internal/handler/` - HTTP handlers
   - `internal/service/` - Business logic
   - `internal/repository/` - Data access layer
   - `internal/model/` - Domain models
   - `internal/middleware/` - HTTP middleware
   - `migrations/` - Database migrations
   - `tests/` - Integration tests
3. Create `.env.example` with configuration template
4. Create `Makefile` with common tasks (run, test, migrate, build)
5. Update `.gitignore` for Go projects
6. Create basic `main.go` with graceful shutdown

Use idiomatic Go patterns and include comments explaining the architecture.
