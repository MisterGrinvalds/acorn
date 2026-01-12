---
description: Add a new API endpoint following project patterns
---

Add a new API endpoint to the Kanban service following existing patterns.

## Required Information

Before proceeding, ensure you know:
1. **Resource name** (e.g., "issues", "comments")
2. **HTTP method** (GET, POST, PUT, DELETE)
3. **Path** (e.g., "/issues", "/issues/{issue_id}")
4. **Authentication requirements** (none, bearer, admin-only)
5. **Request/response structure**

## Implementation Steps

### 1. Check if package exists

Look in `app/domain/` for existing `{resource}app` package.

### 2. If new package needed, create:

```
app/domain/{resource}app/
├── {resource}app.go   # Handlers
├── model.go           # Request/response DTOs
├── route.go           # Route registration
├── filter.go          # Query filters (if list endpoint)
└── order.go           # Ordering options (if list endpoint)
```

### 3. Add handler in `{resource}app.go`

Follow the pattern:
```go
func (a *app) handlerName(ctx context.Context, r *http.Request) web.Encoder {
    // 1. Decode request body (if POST/PUT)
    // 2. Get user/resource from context via mid.GetXXX()
    // 3. Call business layer
    // 4. Return response DTO or error
}
```

### 4. Add route in `route.go`

```go
app.HandlerFunc(http.MethodXXX, version, "/path", api.handler, middlewares...)
```

### 5. Wire up in service (if new package)

Add import and route registration in `api/services/kanban/build/all/all.go`

### 6. Add to Postman collection

Run `/sync-postman` or manually add to `postman/kanban-api.postman_collection.json`

## Patterns to Follow

- Use `web.Decode(r, &dto)` for request body parsing
- Use `errs.New(errs.XXX, err)` for error responses
- Use `mid.GetUserID(ctx)` to get authenticated user
- Return DTO types that implement `web.Encoder` interface
- Business logic belongs in `business/domain/` not handlers
