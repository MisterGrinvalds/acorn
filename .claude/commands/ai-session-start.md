---
description: Initialize a new Claude Code session with goals and context
argument-hint: <session-goal>
model: claude-sonnet-4-5-20250929
allowed-tools: Write, Read, Edit, Glob
---

# Start New Session

Begin a new work session with proper setup and orientation.

## Your Process

1. **Quick Status Check**
   - Review TODO.md
   - Check recent session summaries
   - Identify active work
   - Note any blockers

2. **Session Setup**
   - Confirm focus area for this session
   - Set clear session goals (2-4 specific outcomes)
   - Identify first task to tackle
   - Note any preparation needed

3. **Clean Slate Check**
   - Archive completed work from last session
   - Clean up stale todos
   - Verify important work is saved
   - Check git status if applicable

4. **Present Session Plan**
   - Show current state
   - Propose session goals
   - Recommend starting point
   - Highlight any decisions needed

## Session Start Template

```
🚀 Session Start - YYYY-MM-DD HH:MM
═══════════════════════════════════════

📊 CURRENT STATE
────────────────────────────────────────
Last session: [date/time]
Active tasks: X
In progress: [brief list]

✅ SINCE LAST SESSION
────────────────────────────────────────
- [Any completed work]
- [Any updates]

🎯 PROPOSED SESSION GOALS
────────────────────────────────────────
1. [Goal 1 - specific, achievable]
2. [Goal 2 - specific, achievable]
3. [Goal 3 - specific, achievable]

Success = [Clear definition of "done" for this session]

📋 STARTING POINT
────────────────────────────────────────
First task: [Specific task to begin with]

Context: [Why start here]

Estimated time: [rough estimate]

🔴 BLOCKERS TO RESOLVE
────────────────────────────────────────
- [Blocker]: [How to resolve]
- [Decision needed]: [Options]

🧹 HOUSEKEEPING
────────────────────────────────────────
- [ ] Previous work archived
- [ ] TODO.md reviewed and cleaned
- [ ] Clear focus established
- [ ] Tools/resources ready

═══════════════════════════════════════
Ready to begin? Confirm goals or adjust focus.
```

## Session Goals Best Practices

**Good Session Goals** (Specific, Achievable):
- ✅ "Implement authentication middleware and add tests"
- ✅ "Review and improve all Cobra commands following style guide"
- ✅ "Create Makefile with build, test, and lint targets"

**Poor Session Goals** (Vague, Too Large):
- ❌ "Work on the CLI"
- ❌ "Make it better"
- ❌ "Finish everything"

## Session Types

### Focus Session (2-4 hours)
- Single major task or feature
- Deep work, minimal interruptions
- Goal: Complete or significantly advance one thing

### Sprint Session (30-60 min)
- Small, contained task
- Quick wins
- Goal: Ship something concrete

### Planning Session
- No coding, just planning
- Create execution plans
- Break down work
- Goal: Clear roadmap for next sessions

### Cleanup Session
- Organize, archive, document
- Pay down technical/organizational debt
- Goal: Lean, organized project

### Exploration Session
- Research, spike, experiment
- Learning and discovery
- Goal: Information, not implementation

## Ask User

Before starting work, confirm:
1. "Does this session plan look good?"
2. "Any specific focus or constraints?"
3. "Shall we begin with [first task]?"

## Integration with Other Commands

This command often triggers:
- `/todo-review` if TODO is messy
- `/cleanup` if context is sprawling
- `/status` for quick orientation
- `/plan` if work is undefined

It's the entry point that sets up everything else.
