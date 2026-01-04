---
description: "Factor III: Analyze configuration management"
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Grep
---

# 12-Factor: III. Config

Analyze and improve **Configuration** management in this project.

## Factor III: Store config in the environment

**Principle**: Strict separation of config from code. Configuration includes anything that varies between deployments.

**Key Requirements**:
- Use environment variables for configuration
- No constants in code (database URLs, credentials, API keys)
- Avoid grouped environments (dev, test, prod configs)
- Granular, orthogonal env vars
- Test: Could you open-source without exposing credentials?

**Configuration Hierarchy**:
1. Command-line flags (highest priority)
2. Environment variables
3. Configuration files
4. Default values (lowest priority)

**Review**:
1. Identify hardcoded configuration
2. Check for credentials in code/version control
3. Review Viper/config integration (Cobra best practice)
4. Ensure proper env var usage
5. Verify config flexibility across deployments

Provide specific recommendations for this Cobra CLI project.
