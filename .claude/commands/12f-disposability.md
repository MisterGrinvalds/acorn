---
description: "Factor IX: Analyze disposability"
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Grep
---

# 12-Factor: IX. Disposability

Analyze and improve **Disposability** characteristics in this project.

## Factor IX: Maximize robustness with fast startup and graceful shutdown

**Principle**: Processes are disposable—can be started or stopped instantly.

**Key Requirements**:

**Fast Startup** (seconds, not minutes):
- Enables rapid scaling
- Quick deployment of changes
- Easy process migration

**Graceful Shutdown** (on SIGTERM):
- Stop accepting new requests/jobs
- Complete in-flight work
- Exit cleanly
- Return jobs to queue (worker processes)

**Crash Resilience**:
- Handle unexpected termination
- Use robust queueing (Beanstalkd, RabbitMQ)
- Ensure jobs are reentrant
- Make operations idempotent

**Best Practices**:
- Keep HTTP requests short
- Use crash-only design
- Transaction wrapping for consistency
- Queue-based job handling

**Review**:
1. Check startup time (should be < 10 seconds)
2. Verify SIGTERM/SIGINT handling
3. Review graceful shutdown implementation
4. Check job queue integration if applicable
5. Ensure operations are idempotent
6. Test crash recovery

Provide specific disposability improvements for this Go/Cobra project.
