---
description: "Factor XII: Analyze admin processes"
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Grep
---

# 12-Factor: XII. Admin Processes

Analyze and improve **Admin Process** management in this project.

## Factor XII: Run admin/management tasks as one-off processes

**Principle**: One-off admin tasks run in identical environment as regular app processes.

**Common Admin Tasks**:
- Database migrations (rake db:migrate, go run migrations)
- Interactive console/REPL
- Running one-time scripts
- Data fixes and imports
- Running scheduled tasks manually

**Key Requirements**:
- Same environment as long-running processes
- Same codebase and version
- Same configuration
- Same dependency isolation
- Ship admin code with app code

**Best Practices**:

**Identical Environment**:
- Use same release version
- Access same backing services
- Use same dependency tools
- Apply same isolation mechanisms

**Local Development**:
- Direct shell commands
- Run from app directory

**Production**:
- SSH or remote command execution
- Platform-provided mechanisms
- Containerized execution

**Language Features**:
- Built-in REPL (go run, python, irb)
- Easy one-off script execution
- Interactive debugging capabilities

**Review**:
1. Identify existing admin tasks
2. Check if they run in same environment
3. Verify codebase synchronization
4. Review migration strategy
5. Check console/REPL availability
6. Ensure dependency consistency for admin tasks

Provide specific recommendations for admin processes in this Go/Cobra project.
