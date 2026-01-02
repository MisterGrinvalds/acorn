---
description: Archive completed plans and clean up session files
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Write, Bash
---

# Archive Completed Work

Archive completed work to reduce active context while preserving historical record.

## Your Process

1. **Identify Completed Work**
   - Scan TODO.md for completed tasks
   - Check for completed plan files (PLAN-*.md with status: Complete)
   - Look for session notes or temporary files
   - Identify old context dumps or research notes

2. **Create Archive Structure**
   ```
   .claude/archive/
   ├── YYYY-MM/
   │   ├── completed-work-YYYY-MM-DD.md
   │   └── decisions-YYYY-MM.md
   └── completed-todos/
       └── YYYY-MM-DD-completed.md
   ```

3. **Archive Completed TODOs**
   - Extract completed tasks from TODO.md
   - Create `archive/completed-todos/YYYY-MM-DD-completed.md`
   - Format with completion date and context
   - Remove from active TODO.md

4. **Archive Completed Plans**
   - Move PLAN-*.md files marked complete to archive
   - Rename to `archive/YYYY-MM/plan-[name]-completed.md`
   - Keep reference link in active docs if needed

5. **Create Completion Summary**
   - Write `archive/YYYY-MM/completed-work-YYYY-MM-DD.md`
   - Include: what was completed, key decisions, outcomes
   - Add links to relevant code/commits
   - Note any follow-up items

6. **Clean Active Directory**
   - Remove archived files from root
   - Update TODO.md
   - Update any active plans that referenced archived work
   - Clean up temporary files

7. **Update Documentation**
   - Update TODO.md header with last archive date
   - Add archive index entry (if maintaining one)
   - Note reduced context size

## Archive Summary Template

```markdown
# Completed Work - YYYY-MM-DD

## Summary
[2-3 sentence overview of what was accomplished]

## Completed Tasks
- [Task 1] - [Brief outcome]
- [Task 2] - [Brief outcome]

## Key Decisions
1. **[Decision topic]**: [What was decided and why]
2. **[Decision topic]**: [What was decided and why]

## Outcomes
- [Deliverable 1]: [Location/status]
- [Deliverable 2]: [Location/status]

## Metrics (if applicable)
- Files modified: X
- Tests added: X
- Documentation pages: X

## Follow-up Items
- [ ] [Item moved to TODO.md]
- [ ] [Item moved to TODO.md]

## References
- Code: [commit hash, file paths]
- Plans: [links to plan files]
- Related work: [links]
```

## Report to User

```
✅ Archive Complete

📦 Archived:
- X completed tasks → archive/completed-todos/[date]
- X plan files → archive/YYYY-MM/
- [Other files]

📝 Created:
- archive/YYYY-MM/completed-work-[date].md

🧹 Cleaned:
- Removed X completed items from TODO.md
- Removed X temporary files
- Active context reduced by ~X%

📋 Remaining Active Work:
- X tasks in TODO.md
- X active plans

🎯 Next Focus:
[Top 2-3 priority items from remaining work]
```

Keep archives detailed but concise. They should be useful for future reference without being overwhelming.
