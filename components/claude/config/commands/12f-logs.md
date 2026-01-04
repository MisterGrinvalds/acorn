---
description: "Factor XI: Analyze logging"
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Grep
---

# 12-Factor: XI. Logs

Analyze and improve **Logging** approach in this project.

## Factor XI: Treat logs as event streams

**Principle**: Apps write to stdout; execution environment handles routing and storage.

**Key Requirements**:
- Write unbuffered to `stdout`
- Never manage log files
- Never route/store logs in app code
- Treat as time-ordered event stream
- Aggregate from all processes and services

**Event Stream Approach**:
- Logs are continuous flows, not files
- No fixed beginning or end
- Environment captures and routes
- Apps don't see final destination

**Development vs. Production**:

**Development**:
- View stream in terminal foreground
- Real-time visibility

**Production/Staging**:
- Environment captures from all processes
- Routes to archival destinations
- Splunk, Logplex, Fluentd, Hadoop/Hive
- Long-term storage and analysis

**Capabilities Enabled**:
- Find past events
- Large-scale graphing (requests/minute)
- Active alerting (errors/minute threshold)
- Historical trend analysis

**Review**:
1. Check current logging implementation
2. Verify stdout usage (not log files)
3. Review log aggregation strategy
4. Ensure structured logging
5. Check log levels and formatting
6. Verify no buffering issues

Provide specific logging improvements for this Go/Cobra project.
