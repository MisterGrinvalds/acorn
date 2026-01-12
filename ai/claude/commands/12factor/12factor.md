---
description: Analyze project against all 12-factor principles
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Glob, Grep
---

# 12-Factor Methodology Expert

Analyze this project against the 12-factor methodology for building modern SaaS applications.

Review the codebase and identify:
1. Which of the 12 factors are currently being followed
2. Which factors are being violated or could be improved
3. Specific recommendations for better alignment with 12-factor principles

The 12 factors are:
1. **Codebase**: One codebase in version control, many deploys
2. **Dependencies**: Explicitly declare and isolate dependencies
3. **Config**: Store config in the environment
4. **Backing Services**: Treat backing services as attached resources
5. **Build, Release, Run**: Strictly separate build and run stages
6. **Processes**: Execute as stateless processes
7. **Port Binding**: Export services via port binding
8. **Concurrency**: Scale out via the process model
9. **Disposability**: Fast startup and graceful shutdown
10. **Dev/Prod Parity**: Keep environments as similar as possible
11. **Logs**: Treat logs as event streams
12. **Admin Processes**: Run admin tasks as one-off processes

Focus on providing practical, actionable improvements for this specific codebase.
