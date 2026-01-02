---
description: "Factor VI: Analyze process execution"
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Grep
---

# 12-Factor: VI. Processes

Analyze and improve **Process** architecture in this project.

## Factor VI: Execute the app as one or more stateless processes

**Principle**: Processes are stateless and share-nothing. No persistent local state between requests.

**Key Requirements**:
- Stateless execution
- Share-nothing architecture
- No sticky sessions
- Session state in backing services (Redis, Memcached)
- Asset compilation during build, not runtime
- Memory/filesystem for single-transaction caching only

**Violations to Avoid**:
- Relying on in-memory caching across requests
- Assuming local filesystem persists
- Sticky sessions routing users to same process
- Storing user data in process memory

**Review**:
1. Check for stateful patterns in code
2. Identify session management approach
3. Verify no local state dependencies
4. Review caching strategy
5. Ensure multi-process compatibility

Provide specific recommendations for stateless architecture in this project.
