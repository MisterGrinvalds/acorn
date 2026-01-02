---
description: Add a new section to the Makefile
argument-hint: <section-name>
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Edit
---

# Add Makefile Section

Add a new organized section to the Makefile with related targets.

## Task

1. Verify a Makefile exists
2. Ask the user for:
   - Section name (e.g., "Docker", "Deployment", "Database")
   - Initial targets to include in the section

3. Create the section with:
   - Section header using ##@ format
   - Targets with proper formatting
   - .PHONY declarations
   - Help comments for each target
   - Logical ordering of related targets

4. Insert the section in an appropriate location:
   - After General/Development sections for workflow sections
   - Before Utilities section for tooling sections
   - At the end for specialized sections

5. Update help output to show the new section

## Common Section Templates

**Docker Section**:
- docker-build: Build Docker image
- docker-run: Run container locally
- docker-push: Push to registry
- docker-clean: Remove images

**Deployment Section**:
- deploy-dev: Deploy to development
- deploy-staging: Deploy to staging
- deploy-prod: Deploy to production
- rollback: Rollback deployment

**Database Section**:
- db-migrate: Run migrations
- db-rollback: Rollback migrations
- db-seed: Seed database
- db-reset: Reset database

## Requirements

- Use ##@ for section headers
- Group related targets logically
- Follow naming conventions (verb-noun pattern)
- Include all necessary .PHONY declarations
- Add helpful documentation
