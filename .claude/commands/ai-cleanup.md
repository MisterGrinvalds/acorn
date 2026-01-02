---
description: Clean up old session files and archives
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Bash
---

# Clean Up Project and Context

Perform comprehensive cleanup of temporary files, stale content, and organize project structure.

## Your Process

1. **Scan for Cleanup Candidates**
   - Temporary files (*.tmp, *.log, etc.)
   - Old context dumps (>30 days)
   - Stale session summaries not archived
   - Duplicate or abandoned files
   - Large files that can be compressed/archived
   - Build artifacts (if not in .gitignore)

2. **Review TODO.md Health**
   - Count total items
   - Identify stale items (>14 days old)
   - Flag items with no recent updates
   - Suggest archiving or removal

3. **Check Archive Organization**
   - Ensure `.claude/archive/YYYY-MM/` structure
   - Move old files to appropriate month
   - Create archive index if useful
   - Verify no duplicates

4. **Clean Git State** (if applicable)
   - Check for untracked files
   - Identify files that should be in .gitignore
   - Note any uncommitted changes
   - Suggest cleanup actions

5. **Organize .claude Directory**
   ```
   .claude/
   ├── agents/           # Should be lean, focused
   ├── commands/         # Should be organized by category
   ├── archive/          # Should be date-organized
   │   ├── YYYY-MM/
   │   └── completed-todos/
   ├── settings.local.json
   └── README.md
   ```

## Report Structure

```
🧹 Cleanup Report
═══════════════════════════════════════

📊 SCAN RESULTS
────────────────────────────────────────
Files scanned: X
Total size: X MB
Issues found: X

🗑️  REMOVED
────────────────────────────────────────
- X temporary files (~X KB)
- X duplicate files (~X KB)
- X build artifacts (~X MB)
Total freed: ~X MB

📦 ARCHIVED
────────────────────────────────────────
- X session summaries → archive/YYYY-MM/
- X context dumps → archive/YYYY-MM/
- X completed plans → archive/YYYY-MM/

🔄 ORGANIZED
────────────────────────────────────────
- TODO.md: X items cleaned, X archived
- Archive structure: Created YYYY-MM folders
- .gitignore: Added X entries

⚠️  NEEDS ATTENTION
────────────────────────────────────────
- [File/issue]: [Recommendation]
- [File/issue]: [Recommendation]

📊 BEFORE → AFTER
────────────────────────────────────────
Total files: X → Y (-Z files)
Total size: X MB → Y MB (-Z MB)
Active context: X files → Y files

✨ RESULT
────────────────────────────────────────
Project is cleaner and more organized!
Active context reduced by ~X%
Archive properly structured
Ready for focused work
```

## When to Run Cleanup

**Regular Cleanup** (every 7-14 days):
- Remove temp files
- Archive completed work
- Review TODO.md

**Deep Cleanup** (monthly):
- Full archive organization
- Consolidate documentation
- Review all files
- Git housekeeping

**Emergency Cleanup** (when context is overwhelming):
- Aggressive archiving
- Remove all non-essential files
- Start fresh with lean context

**Before Major Changes**:
- Clean slate for new work
- Ensure important work is saved
- Archive old context
