---
description: "Factor VII: Analyze port binding"
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Grep
---

# 12-Factor: VII. Port Binding

Analyze and improve **Port Binding** implementation in this project.

## Factor VII: Export services via port binding

**Principle**: The app is completely self-contained and exports HTTP (or other protocols) by binding to a port.

**Key Requirements**:
- No reliance on runtime webserver injection
- Embed webserver library in app (Tornado, Thin, Jetty)
- Listen on port, respond to requests
- Routing layer forwards public traffic to port-bound processes
- Works for HTTP, XMPP, Redis protocol, etc.

**Self-Contained Approach**:
- Webserver is a library, not external container
- Declared as dependency
- Operates entirely in user space
- In production: routing maps hostname to port

**Service Composition**:
- One app can be backing service for another
- Share port URL as config resource

**Review**:
1. Check if app is self-contained (no Apache/nginx dependency)
2. Verify port binding implementation
3. Review webserver library integration
4. Ensure port configuration flexibility
5. Check routing/reverse proxy setup for production

Provide recommendations for proper port binding in this Go/Cobra project.
