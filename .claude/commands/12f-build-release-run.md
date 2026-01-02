---
description: "Factor V: Analyze build, release, and run stages"
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Bash
---

# 12-Factor: V. Build, Release, Run

Analyze and improve **Build, Release, Run** separation in this project.

## Factor V: Strictly separate build and run stages

**The Three Stages**:

**Build**: Code → Executable Bundle
- Fetch dependencies
- Compile binaries and assets
- Uses specific git commit

**Release**: Build + Config → Deployment
- Combines build with deployment config
- Immutable with unique release ID
- Ready for immediate execution

**Run**: Execute in Production
- Launch processes against selected release
- No code changes possible at runtime
- Minimal complexity for reliability

**Key Requirements**:
- Cannot modify code at runtime
- Every release has unique ID
- Releases are immutable
- Rollback capability
- Clear separation of stages

**Review**:
1. Check build process (Makefile, CI/CD)
2. Verify release versioning strategy
3. Ensure runtime doesn't allow code changes
4. Review deployment pipeline
5. Check rollback mechanisms

Provide specific recommendations for this Go/Cobra project's build pipeline.
