---
description: Dump all relevant context for debugging
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Glob
---

# Dump Current Context to File

Create a comprehensive snapshot of current session context for reference or before major transitions.

## Your Process

1. **Capture Current State**
   - Current working directory
   - Active TODO items
   - In-progress work
   - Recent file modifications
   - Key decisions or context from session
   - Outstanding questions

2. **Analyze Project State**
   - List active files and their purposes
   - Identify important code sections
   - Note configuration state
   - Record environment details

3. **Create Context Document**
   - Write to `CONTEXT-YYYY-MM-DD.md`
   - Use structured format (see template)
   - Include file paths and line numbers
   - Add relevant code snippets if needed

4. **Make it Useful**
   - Focus on what someone needs to resume work
   - Include "why" not just "what"
   - Note any non-obvious state
   - Highlight key relationships between files

## Context Dump Template

```markdown
# Context Dump - YYYY-MM-DD HH:MM

## Current Focus
[What is the main work happening right now?]

## Project State
**Repository**: [path]
**Branch**: [if git]
**Working Directory**: [path]

## Active Work 🚧

### In Progress
1. **[Task/Feature]**
   - Status: [% complete or stage]
   - Files: [paths]
   - Next step: [what needs to happen]
   - Context: [any important background]

2. **[Task/Feature]**
   - ...

### Recently Completed
1. [Task] - [Outcome]
2. [Task] - [Outcome]

## File State 📁

### Modified (Uncommitted)
- `path/to/file.go` - [What changed and why]
- `path/to/file.go` - [What changed and why]

### Key Files in Focus
| File | Purpose | Status | Notes |
|------|---------|--------|-------|
| cmd/root.go | CLI root | Modified | Added config flag |
| pkg/auth.go | Auth logic | In progress | Implementing OAuth |

### New Files Created
- `path/to/file.go` - [Purpose]

## Outstanding Questions ❓
1. [Question] - [Context about why this matters]
2. [Question] - [Context]

## Key Decisions Made 🎯
1. **[Topic]**: [Decision] - [Rationale]
2. **[Topic]**: [Decision] - [Rationale]

## Dependencies & Blockers 🔴
- [Dependency]: [Status, ETA, or blocker reason]
- [Blocker]: [What's blocking and what's needed]

## Active TODO Items (Snapshot)
```
[Copy of current TODO.md or top 10 items]
```

## Environment Details
- Go version: [version]
- Key dependencies: [list major ones]
- Configuration: [any important config state]

## Code Context 💻

### Important Code Sections
**File: cmd/serve.go:45-67**
```go
// Snippet of important code with context
```
Purpose: [Why this code matters]

**File: pkg/database.go:123**
```go
// Another important snippet
```
Purpose: [Why this code matters]

## Mental Model 🧠
[Explain the current understanding of the system, key relationships, architecture decisions]

## Next Steps (When Resuming)
1. [First thing to do]
2. [Second thing to do]
3. [Third thing to do]

## References
- Related docs: [links]
- External resources: [links]
- Previous context dumps: [links]

## Session History
- Started: [when this work began]
- Last major change: [what and when]
- Total time invested: [estimate if known]
```

## When to Use Context Dumps

**Use before**:
- Switching to different feature/workstream
- Taking break for >24 hours
- Major refactoring or architectural change
- Handing off work to another developer
- Context feeling overwhelming (>100K tokens)

**Use after**:
- Completing major milestone
- Making significant architectural decisions
- Complex debugging session (capture findings)
- Learning something important about the system

## Report to User

```
💾 Context Dump Created

📄 File: CONTEXT-YYYY-MM-DD.md

📊 Captured:
- X active tasks
- X modified files
- X key decisions
- X outstanding questions

🎯 Focus: [Main work area]

✅ Use this context dump to:
- Resume work after a break
- Onboard someone to this work
- Remember why decisions were made
- See project state at this point in time

💡 Tip: Reference this dump when you return to this work
```

Context dumps are like save points in a game - they let you resume exactly where you left off.
