---
description: Check if a feature exists before building it
---

# Feature Existence Check

Before building a new feature, verify if it already exists or if similar patterns are available.

## Search Process

Given the user's feature request (e.g., "Issues", "Comments", "Boards"):

### 1. Search Domain Packages
```bash
# App layer
ls -la app/domain/

# Business layer
ls -la business/domain/
```

### 2. Search for Entity Names
```bash
# Case-insensitive search for the entity
grep -ri "issue\|issues" app/domain/ business/domain/ --include="*.go" -l
grep -ri "comment\|comments" app/domain/ business/domain/ --include="*.go" -l
```

### 3. Check Database Schema
```bash
# Look for existing tables
grep -i "CREATE TABLE" business/sdk/migrate/sql/migrate.sql
```

### 4. Check API Design Doc
Read `.claude/commands/design-api.md` to see if the feature was planned.

## Report Format

Provide a clear report:

### If Feature Exists:
```
✅ FOUND: {Feature} already exists!

Location:
- App layer: app/domain/{feature}app/
- Business layer: business/domain/{feature}bus/
- Database: {table_name} table in migrate.sql

Endpoints:
- GET /v1/{feature}s
- POST /v1/{feature}s
- ...

Would you like me to explain how it works or help you extend it?
```

### If Similar Feature Exists:
```
⚠️ SIMILAR: {Feature} doesn't exist, but {SimilarFeature} does.

You can use {SimilarFeature} as a reference pattern:
- app/domain/{similar}app/ - for handler patterns
- business/domain/{similar}bus/ - for business logic

Want to start building {Feature} using {SimilarFeature} as a template?
Use /coach to walk through it step by step.
```

### If Feature is Planned:
```
📋 PLANNED: {Feature} is in the API design but not implemented yet.

From design-api.md:
- GET /api/v1/{feature}s - List all
- POST /api/v1/{feature}s - Create new
- ...

Ready to implement? Use /coach to build it together.
```

### If Feature is New:
```
🆕 NEW: {Feature} doesn't exist and wasn't planned.

Before building, consider:
1. How does {Feature} relate to existing entities?
2. What operations are needed?
3. Who can access it?

Ready to design and build? Use /coach to start from scratch.
```

## Next Steps

Based on the check result, suggest:
- `/coach` - For guided implementation
- `/explain-api` - To understand existing patterns
- `/add-endpoint` - For quick reference
