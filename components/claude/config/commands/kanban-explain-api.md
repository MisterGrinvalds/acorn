---
description: Explain current API structure and endpoints
---

Analyze and explain the current API structure of the Kanban service.

## Tasks

1. **List all route files** in `app/domain/*/route.go`

2. **For each route package**, show:
   - Package name and purpose
   - All registered endpoints with:
     - HTTP method
     - Full URL path
     - Required authentication/authorization
     - Handler function name

3. **Show service composition**:
   - Which routes are included in the kanban service (`api/services/kanban/build/all/all.go`)
   - Which routes are included in the auth service (`api/services/auth/build/all/all.go`)

4. **Summarize in a table format**:

```
| Method | Path | Auth | Handler | Package |
|--------|------|------|---------|---------|
| GET | /v1/projects | Bearer + Admin | query | projectapp |
```

5. **Note any TODOs or incomplete endpoints** found in route files

## Output Format

Provide a clear, organized summary that a developer can use to understand:
- What endpoints exist
- What authentication they require
- Which service serves them
