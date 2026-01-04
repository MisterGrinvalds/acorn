---
description: Sync Postman collection with current API routes
---

Analyze all route files and update the Postman collection to match current API endpoints.

## Steps

1. **Scan all route files** in `app/domain/*/route.go`

2. **For each route file**, extract:
   - HTTP method (GET, POST, PUT, DELETE, PATCH)
   - URL path (e.g., `/projects`, `/projects/{project_id}`)
   - Handler name (for description)
   - Middleware (to determine auth requirements)

3. **Parse existing Postman collection** at `postman/kanban-api.postman_collection.json`

4. **Compare and update**:
   - Add new endpoints that don't exist in Postman
   - Flag removed endpoints (don't auto-delete, warn user)
   - Update paths if they changed

5. **Organize by folder** based on route package:
   - `checkapp` -> Health folder
   - `authapp` -> Auth folder
   - `projectapp` -> Projects folder
   - New packages -> New folders

6. **For each endpoint, ensure**:
   - Correct HTTP method
   - URL uses `{{base_url}}` or `{{auth_url}}` variable
   - Path parameters use `:param` Postman syntax
   - Auth inherits from collection (Bearer) unless it's a no-auth endpoint
   - Request body template for POST/PUT/PATCH

7. **Model reference**: Check `app/domain/*/model.go` for request body structure

## Output

- Update `postman/kanban-api.postman_collection.json`
- Report what was added/changed/flagged for removal

## Example Route Parsing

```go
// From route.go:
app.HandlerFunc(http.MethodGet, version, "/projects/{project_id}", api.queryByID, authen, ruleAuthorizeProject)

// Becomes Postman request:
{
  "name": "Get Project by ID",
  "request": {
    "method": "GET",
    "url": {
      "raw": "{{base_url}}/v1/projects/:project_id",
      "host": ["{{base_url}}"],
      "path": ["v1", "projects", ":project_id"],
      "variable": [{"key": "project_id", "value": "", "description": "UUID of the project"}]
    }
  }
}
```

## Notes

- Auth service routes use `{{auth_url}}`
- Kanban service routes use `{{base_url}}`
- `HandlerFuncNoMid` routes should have `"auth": {"type": "noauth"}`
- Keep existing request body examples if they exist
