---
description: End the current session and archive work
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Write, Bash
---

# End Session with Cleanup

End the current session with full cleanup, archiving, and preparation for next session.

## Your Process

This is a comprehensive cleanup combining multiple operations:

1. **Create Session Summary**
   - Document what was accomplished
   - Record key decisions
   - Note incomplete work
   - Write to `SESSION-YYYY-MM-DD-HHMM.md`

2. **Review and Clean TODO List**
   - Mark completed tasks
   - Archive completed items
   - Flag stale items
   - Identify blockers
   - Set priorities for next session

3. **Archive Completed Work**
   - Move completed todos to archive
   - Archive completed plans
   - Create completion summary
   - Clean up temporary files

4. **Organize Files**
   - Ensure `.claude/archive/` structure exists
   - Move session summary to archive if desired
   - Clean up any leftover temp files
   - Verify important work is saved

5. **Prepare for Next Session**
   - Update TODO.md with clear next steps
   - Flag top 3 priority items
   - Note any preparation needed
   - Document any decisions pending user input

6. **Context Cleanup**
   - Identify large/unnecessary files
   - Compress or remove verbose logs
   - Clear any debugging artifacts
   - Verify git status is clean (if applicable)

## Comprehensive Checklist

```markdown
## Session End Checklist

### Documentation
- [ ] Session summary created
- [ ] Key decisions documented
- [ ] Code changes noted
- [ ] Incomplete work recorded

### TODO Management
- [ ] Completed tasks marked
- [ ] Completed tasks archived
- [ ] Stale tasks flagged
- [ ] Priorities set for next session
- [ ] Top 3 next actions identified

### Archive
- [ ] Completed work archived
- [ ] Archive structure organized
- [ ] Old plans moved to archive
- [ ] Temporary files cleaned

### Code/Project State
- [ ] All important changes saved
- [ ] Git status reviewed (if applicable)
- [ ] Build status verified (if applicable)
- [ ] No broken/incomplete features in main files

### Next Session Prep
- [ ] Clear starting point identified
- [ ] Blockers documented
- [ ] Decisions needed flagged
- [ ] Resources/links gathered
```

## Final Report Template

```
🎬 Session End Report
═══════════════════════════════════════

📊 SESSION STATISTICS
────────────────────────────────────────
Duration: [Estimated time]
Tasks completed: X
Files modified: X
Decisions made: X

✅ ACCOMPLISHMENTS
────────────────────────────────────────
1. [Major accomplishment]
2. [Major accomplishment]
3. [Major accomplishment]

🎯 KEY DECISIONS
────────────────────────────────────────
1. [Decision]: [Outcome]
2. [Decision]: [Outcome]

📦 ARCHIVED
────────────────────────────────────────
- Session summary → SESSION-[datetime].md
- Completed todos → archive/completed-todos/[date].md
- Completed plans → archive/YYYY-MM/
- [Other archived items]

📋 TODO STATUS
────────────────────────────────────────
Active tasks: X
High priority: X
Blocked: X
Stale: X (need review)

🎯 NEXT SESSION PRIORITIES
────────────────────────────────────────
1. [Top priority task]
2. [Second priority task]
3. [Third priority task]

⚠️  BLOCKERS / DECISIONS NEEDED
────────────────────────────────────────
- [Blocker 1]: [What's needed to unblock]
- [Decision needed]: [Context]

🧹 CLEANUP PERFORMED
────────────────────────────────────────
- Archived X completed tasks
- Removed X temporary files
- Organized archive structure
- Context reduced by ~X%

💾 IMPORTANT FILES
────────────────────────────────────────
- Session summary: SESSION-[datetime].md
- Active TODO: TODO.md
- Active plans: [list]
- Recent changes: [list]

═══════════════════════════════════════
✨ Session cleaned and ready for next time!
```

**Time to Archive Session Summary?**
Ask user if they want to move SESSION-[datetime].md to archive or keep in root for easy access.

This command performs a full reset - use it at natural breakpoints or when switching major focus areas.
