---
description: "Factor I: Analyze codebase and version control"
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Bash
---

# 12-Factor: I. Codebase

Analyze and improve **Codebase** management in this project.

## Factor I: One codebase tracked in revision control, many deploys

**Principle**: A twelve-factor app is always tracked in version control with a one-to-one mapping between codebase and app.

**Key Requirements**:
- Single repository per application
- Multiple deployments from one codebase (dev, staging, production)
- Shared code extracted into libraries, not duplicated
- Same codebase across all environments (different versions may run)

**Review**:
1. Verify proper version control setup
2. Check for code duplication that should be libraries
3. Ensure deployment strategy aligns with one-codebase principle
4. Identify any multi-codebase issues (indicates distributed system)

Provide specific recommendations for this project.
