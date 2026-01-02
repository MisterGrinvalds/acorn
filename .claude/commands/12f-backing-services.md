---
description: "Factor IV: Analyze backing services"
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Grep
---

# 12-Factor: IV. Backing Services

Analyze and improve **Backing Services** integration in this project.

## Factor IV: Treat backing services as attached resources

**Principle**: Make no distinction between local and third-party services. All accessed via configuration, not code.

**Backing Services Include**:
- Databases (MySQL, PostgreSQL, MongoDB)
- Messaging systems (RabbitMQ, Kafka)
- Caching (Redis, Memcached)
- Email services (SMTP, Postmark)
- Third-party APIs (AWS, Stripe, Twilio)

**Key Requirements**:
- Access via URL + credentials in config
- Swap services without code changes
- Each distinct service is a separate resource
- Loosely coupled attachments

**Review**:
1. Identify all backing services used
2. Check if services are configurable
3. Verify no hardcoded service locations
4. Ensure services can be swapped via config only
5. Review resource abstraction patterns

Provide specific improvements for service integration in this project.
