---
description: Audit links in Obsidian vault and fix broken connections
argument-hint: [vault-path] [--fix]
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Edit, Glob, Grep
---

# Audit Obsidian Links

Audit all links in an Obsidian vault, identify issues, and optionally fix broken links.

**Vault Path**: $1 (optional - defaults to current directory)
**Auto-Fix**: $2 (optional - use "--fix" to automatically fix issues)

## Audit Types

### 1. Broken Links
Links to notes that don't exist:
- `[[Non-existent Note]]` where file doesn't exist
- Identify link source and target
- Suggest similar note names (fuzzy matching)
- Option to create missing notes
- Option to remove broken links

### 2. Orphaned Notes
Notes with no connections:
- No incoming links (no backlinks)
- No outgoing links
- Completely isolated from graph
- Suggest potential connections
- Recommend tags or MOC inclusion

### 3. Unlinked Mentions
References that could be links:
- Note titles mentioned without `[[]]`
- Case-insensitive matching
- Suggest automatic linking
- Preview context before change

### 4. Malformed Links
Syntax errors in links:
- Single brackets `[Note]`
- Missing closing bracket `[[Note`
- Invalid characters in links
- Spaces vs underscores inconsistency
- Markdown links to internal notes

### 5. Ambiguous Links
Links that could point to multiple notes:
- Multiple notes with same name
- Case sensitivity issues
- Notes in different folders
- Suggest using full path

## Process

1. **Scan Vault**
   - Find all markdown files
   - Build note index with paths
   - Create name-to-path mapping

2. **Parse Links**
   - Extract all `[[wikilinks]]` from each note
   - Parse link components (note, heading, alias)
   - Check `![[ embeds ]]`
   - Identify link locations (file:line)

3. **Validate Links**
   - Check if target note exists
   - Verify heading links point to actual headings
   - Check block references
   - Identify case mismatches

4. **Find Orphans**
   - Build incoming link map
   - Build outgoing link map
   - Identify notes with zero degree (no connections)

5. **Detect Unlinked Mentions**
   - Search for note title occurrences
   - Exclude existing links
   - Filter false positives (common words)
   - Rank by likelihood

6. **Generate Report**
   - Categorize all issues
   - Provide file:line references
   - Include context snippets
   - Suggest fixes

## Fixing Options

When `--fix` flag is provided:

### Auto-Fix (Safe)
- Fix malformed syntax (add missing brackets)
- Standardize link format
- Fix case mismatches to existing notes

### Prompted Fix (Requires Confirmation)
- Create missing notes for broken links
- Add links for unlinked mentions
- Remove broken links
- Merge duplicate notes

### Manual Fix (Report Only)
- Ambiguous links (require user choice)
- Orphaned notes (require strategy decision)
- Complex refactoring

## Report Structure

```markdown
# Link Audit Report
**Generated**: [Date]
**Vault**: [Path]
**Mode**: [Audit Only / Auto-Fix]

## Summary
- Total Links: X
- Broken Links: X (Y%)
- Orphaned Notes: X (Y%)
- Unlinked Mentions: X
- Malformed Links: X

## Critical Issues

### Broken Links (X)
1. `[[Missing Note]]` in `Note Name.md:42`
   - **Context**: "...text around link..."
   - **Suggestion**: Did you mean `[[Similar Note]]`?
   - **Action**: [ ] Create note [ ] Update link [ ] Remove

### Orphaned Notes (X)
1. `Isolated Note.md`
   - **Created**: YYYY-MM-DD
   - **Words**: X
   - **Suggestions**: Link to [[MOC]], Add tags
   - **Action**: [ ] Link [ ] Tag [ ] Archive

## Warnings

### Unlinked Mentions (X)
1. "Reference to Note" in `Source.md:15`
   - **Target**: `Note.md`
   - **Context**: "...mention of Note in text..."
   - **Action**: [ ] Convert to link

### Malformed Links (X)
1. `[Note]` in `File.md:23` (missing bracket)
   - **Fix**: `[[Note]]`
   - **Action**: [ ] Auto-fix

## Statistics by Folder
| Folder | Notes | Links | Broken | Orphaned |
|--------|-------|-------|--------|----------|
| ...    | ...   | ...   | ...    | ...      |

## Recommendations
1. [Prioritized action items]
2. ...
```

## Output

Display:
1. Path to generated audit report
2. Summary statistics
3. Critical issue count
4. Fixes applied (if --fix used)
5. Next steps

## Notes
- Always create backup before using --fix
- Review suggested fixes before applying
- Consider vault size (large vaults may take time)
- Run audit regularly as part of vault maintenance
- Use with version control for safe fixing
