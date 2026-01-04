---
description: "Factor II: Analyze dependency management"
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Glob
---

# 12-Factor: II. Dependencies

Analyze and improve **Dependency** management in this project.

## Factor II: Explicitly declare and isolate dependencies

**Principle**: Never rely on implicit existence of system-wide packages. Declare all dependencies explicitly and isolate them during execution.

**Key Requirements**:
- Dependency declaration manifest (go.mod, package.json, Gemfile)
- Dependency isolation tools (go modules, bundler, virtualenv)
- No reliance on system-wide tools
- Vendor system utilities if needed (ImageMagick, curl)

**Review**:
1. Check dependency declaration completeness
2. Verify isolation mechanisms are in place
3. Identify any implicit system dependencies
4. Ensure new developers can set up with just runtime + dependency manager
5. Review vendor directory usage if applicable

Provide specific improvements for this Go/Cobra project.
