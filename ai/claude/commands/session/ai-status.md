---
description: Display current session status, progress, and next steps
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Glob
---

# Quick Status Check

Get a quick overview of project status, active work, and next steps.

## Your Process

1. **Read Key Files**
   - TODO.md
   - Recent PLAN-*.md files
   - Latest SESSION-*.md if exists
   - Git status (if applicable)

2. **Analyze Current State**
   - What work is in progress?
   - What's blocked?
   - What's ready to start?
   - What was recently completed?

3. **Provide Quick Summary**
   - Current focus area
   - Active task count
   - Top priorities
   - Any blockers
   - Recommended next action

## Status Report Format

```
🔍 Project Status - YYYY-MM-DD HH:MM
═══════════════════════════════════════

📌 CURRENT FOCUS
────────────────────────────────────────
[Main feature/area being worked on]

📊 TODO SUMMARY
────────────────────────────────────────
Total active tasks: X
├─ 🔴 Blocked: X
├─ 🟡 In progress: X
├─ 🟢 Ready: X
└─ ⚪ Stale/needs review: X

🎯 TOP 3 PRIORITIES
────────────────────────────────────────
1. [Priority task 1] - [Status]
2. [Priority task 2] - [Status]
3. [Priority task 3] - [Status]

✅ RECENTLY COMPLETED
────────────────────────────────────────
- [Recent completion 1]
- [Recent completion 2]

🔴 BLOCKERS
────────────────────────────────────────
- [Blocker 1]: [What's needed]
- [Blocker 2]: [What's needed]

📁 ACTIVE PLANS
────────────────────────────────────────
- PLAN-[name].md - [Status/progress]
- PLAN-[name].md - [Status/progress]

💾 RECENT CHANGES
────────────────────────────────────────
Modified files: [count or list top 5]
Last significant change: [description]

🎯 RECOMMENDED NEXT ACTION
────────────────────────────────────────
[Specific actionable next step]

Context: [Why this is the recommended next step]

═══════════════════════════════════════
```

## Quick Health Check

Also provide a quick health assessment:

```
🏥 HEALTH CHECK
────────────────────────────────────────
TODO List: [✅ Healthy | ⚠️  Getting large | 🔴 Unwieldy]
Context: [✅ Lean | ⚠️  Moderate | 🔴 Sprawling]
Archive: [✅ Organized | ⚠️  Needs attention | 🔴 Chaotic]
Git State: [✅ Clean | ⚠️  Uncommitted changes | 🔴 Messy]

💡 Recommendations:
- [Suggestion 1]
- [Suggestion 2]
```

## When to Use /status

**Session Start**:
- Orient yourself to current work
- See what's highest priority
- Check for any blockers

**Mid-Session Check**:
- Verify you're on track
- See if priorities have shifted
- Quick reminder of goals

**Before Asking for Help**:
- Understand current state
- Have context for questions
- Know what's blocking progress

**When Feeling Lost**:
- Re-establish focus
- See the big picture
- Get recommended next step

## Speed Matters

This command should be FAST - under 30 seconds.
- Read only essential files
- Don't analyze deeply
- Surface critical info quickly
- Give clear next action

Think of it as a dashboard, not a detailed report.
