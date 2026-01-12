---
description: "Factor VIII: Analyze concurrency model"
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Grep
---

# 12-Factor: VIII. Concurrency

Analyze and improve **Concurrency** approach in this project.

## Factor VIII: Scale out via the process model

**Principle**: Processes are first-class citizens. Scale horizontally by running multiple processes.

**Key Requirements**:
- Different process types for different workloads (web, worker, cron)
- Process formation (types + quantities)
- Horizontal scaling across machines
- Share-nothing architecture enables simple scaling
- Never daemonize or write PID files
- Use OS process manager (systemd, cloud platform, Foreman)

**Process Model**:
- Web processes handle HTTP requests
- Worker processes handle background jobs
- Cron processes handle scheduled tasks
- Each type scales independently

**Implementation**:
- Let OS manage processes (output streams, restarts, shutdowns)
- Internal multiplexing OK (threading, async/await, goroutines)
- Design for adding concurrency easily

**Review**:
1. Identify process types needed (web, worker, etc.)
2. Check process management approach
3. Verify horizontal scaling capability
4. Review goroutine usage (Go-specific)
5. Ensure no daemonization or PID files
6. Check process manager integration

Provide specific concurrency recommendations for this Go/Cobra project.
