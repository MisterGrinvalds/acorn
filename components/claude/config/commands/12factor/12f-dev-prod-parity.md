---
description: "Factor X: Analyze dev/prod parity"
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Glob
---

# 12-Factor: X. Dev/Prod Parity

Analyze and improve **Dev/Prod Parity** in this project.

## Factor X: Keep development, staging, and production as similar as possible

**Principle**: Minimize gaps between environments to enable continuous deployment.

**The Three Gaps**:

**Time Gap**:
- Traditional: Weeks/months to production
- 12-Factor: Hours/minutes to production
- Goal: Deploy frequently

**Personnel Gap**:
- Traditional: Developers write, ops deploy
- 12-Factor: Developers deploy
- Goal: Developers close to production

**Tools Gap**:
- Traditional: SQLite dev, PostgreSQL prod
- 12-Factor: Same services everywhere
- Goal: Maximum parity

**Backing Services Consistency**:
- Same database type/version
- Same cache implementation
- Same message queue
- No "lightweight" alternatives in dev

**Modern Tools Enable Parity**:
- Docker for environment consistency
- Vagrant for local VMs
- Homebrew/apt for package parity
- Infrastructure as Code

**Review**:
1. Compare dev vs. production environments
2. Identify backing service discrepancies
3. Check deployment frequency (time gap)
4. Verify who deploys (personnel gap)
5. Review containerization usage
6. Assess environment setup automation

Provide specific recommendations for achieving dev/prod parity in this project.
