---
description: Generate a summary of the current session
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Glob
---

# Create Session Summary

Create a comprehensive summary of the current session's work, decisions, and outcomes.

## Your Process

1. **Review Session Activity**
   - Check TODO.md for completed tasks
   - Review modified files (use git status/diff if available)
   - Identify key decisions made
   - Note any problems encountered and solutions

2. **Gather Context**
   - What was the session goal?
   - What was actually accomplished?
   - What decisions were made?
   - What remains unfinished?
   - What should happen next?

3. **Create Summary Document**
   - Write to `SESSION-YYYY-MM-DD-HHMM.md`
   - Use clear structure (see template below)
   - Be specific with file paths and line numbers
   - Include both wins and challenges

4. **Update TODO.md**
   - Add "Last Session" reference at top
   - Note any new tasks discovered
   - Update task statuses

## Session Summary Template

```markdown
# Session Summary - YYYY-MM-DD HH:MM

## Goal
[What was the intended focus of this session?]

## Accomplished ✅
1. **[Task/Feature]**
   - File(s): path/to/file.go:123
   - Description: [What was done]
   - Outcome: [Result]

2. **[Task/Feature]**
   - ...

## Decisions Made 🎯
1. **[Decision topic]**
   - **Decision**: [What was decided]
   - **Rationale**: [Why]
   - **Alternatives considered**: [If any]
   - **Impact**: [What this affects]

2. **[Decision topic]**
   - ...

## Challenges & Solutions 🔧
1. **[Problem]**
   - Solution: [How it was resolved]
   - Learning: [Takeaway]

2. **[Problem]**
   - ...

## Code Changes 📝
- Modified: [list of files]
- Created: [list of files]
- Deleted: [list of files]
- Lines changed: ~X

## Testing 🧪
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed
- [ ] Edge cases considered

## Documentation 📚
- [ ] Code comments added
- [ ] README updated
- [ ] API docs updated
- [ ] Examples created

## Incomplete Work ⏳
1. [Task] - [Why incomplete, what's needed]
2. [Task] - [Why incomplete, what's needed]

## Next Session 🎯
**Priority tasks**:
1. [Top priority]
2. [Second priority]
3. [Third priority]

**Preparation needed**:
- [Any research, setup, or decisions needed]

## References
- Related commits: [hashes]
- Related issues: [numbers]
- Documentation: [links]
- External resources: [links]

## Notes
[Any other observations, ideas, or context]
```

## Report to User

Provide brief summary:
```
📊 Session Summary Created

✅ Completed:
- [X tasks completed]
- [Key accomplishment 1]
- [Key accomplishment 2]

🎯 Key Decisions:
- [Decision 1]
- [Decision 2]

⏳ Remaining Work:
- [X tasks in TODO.md]
- Next priority: [Top task]

💾 Saved to: SESSION-[date-time].md
```

The summary should be detailed enough to be useful in 3 months, but concise enough to read in 2 minutes.
