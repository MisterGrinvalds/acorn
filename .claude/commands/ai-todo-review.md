---
description: Review and update TODO.md with current task status
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Edit, Write
---

# Review and Clean TODO List

Review the current TODO.md file, clean up completed/stale items, and reorganize for clarity.

## Your Process

1. **Read Current TODO.md**
   - Load the file
   - Count total tasks
   - Identify status of each task

2. **Categorize Tasks**
   - ✅ **Completed**: Tasks marked as done but not archived
   - 🔴 **Blocked**: Tasks that can't proceed (missing info, dependencies)
   - 🟡 **In Progress**: Currently active tasks
   - 🟢 **Ready**: Can be started immediately
   - ⚪ **Stale**: Not updated recently, unclear status
   - ❌ **Obsolete**: No longer relevant

3. **Analyze Health**
   - Are there too many active tasks? (>20 is problematic)
   - Are there very old tasks? (>7 days may be stale)
   - Are tasks clearly actionable?
   - Are priorities clear?
   - Are there orphaned tasks (no context)?

4. **Clean Up Actions**
   - **Archive completed**: Move to `archive/completed-todos/YYYY-MM-DD-completed.md`
   - **Flag stale**: Add 🟠 stale marker and date last updated
   - **Highlight blocked**: Add 🔴 blocker and reason
   - **Remove obsolete**: Ask user, then delete or archive
   - **Consolidate duplicates**: Merge related tasks
   - **Add missing dates**: Add creation/update timestamps

5. **Reorganize**
   - Group by priority (High/Medium/Low)
   - Group by category/feature area
   - Put blocked items in separate section
   - Clear "Next Actions" section at top

6. **Update TODO.md**
   - Write cleaned version
   - Keep only active/ready/blocked tasks
   - Add summary header with stats
   - Include "Last Reviewed" date

## Report to User

Provide summary:
```
TODO Review Complete:

📊 Statistics:
- Total tasks reviewed: X
- ✅ Archived (completed): X
- 🟢 Ready to start: X
- 🟡 In progress: X
- 🔴 Blocked: X
- ⚪ Stale (needs review): X
- ❌ Removed (obsolete): X

🎯 Recommendations:
- [Top 3 priority tasks to focus on]
- [Any blockers that need resolution]
- [Suggestions for next session]

📁 Files Updated:
- TODO.md (cleaned and reorganized)
- archive/completed-todos/[date]-completed.md (if applicable)
```

Be honest about task health. If the list is unwieldy, say so and suggest strategies.
