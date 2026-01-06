---
description: Analyze Obsidian vault health, structure, and linking patterns
argument-hint: [vault-path]
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Glob, Grep, Bash
---

# Analyze Obsidian Vault

Perform a comprehensive analysis of an Obsidian vault's health, structure, and linking patterns.

**Vault Path**: $1 (optional - defaults to current directory)

## Analysis Categories

### 1. Vault Statistics
Gather and report:
- Total number of notes
- Total word count across vault
- Average note size
- Date range (oldest to newest note)
- Last modified notes (top 10)
- Folder distribution

### 2. Linking Analysis
Analyze link structure:
- Total internal links (Wikilinks)
- Total backlinks
- Most linked notes (top 20)
- Least linked notes
- **Orphaned notes** (no incoming or outgoing links)
- **Hub notes** (many incoming links)
- Average links per note
- Broken links (links to non-existent notes)

### 3. Content Analysis
Review content patterns:
- Tag usage and frequency
- Most common tags (top 20)
- Notes without tags
- Notes without frontmatter
- Duplicate note titles
- Empty or stub notes (< 100 words)
- Very large notes (> 5000 words)

### 4. Organizational Analysis
Assess structure:
- Folder depth and hierarchy
- Files in root vs folders
- Folder with most notes
- Empty folders
- Naming consistency
- Date-based note patterns

### 5. Health Checks
Identify issues:
- Broken links and references
- Missing backlinks (unlinked mentions)
- Isolated clusters (groups of notes disconnected from main graph)
- Naming conflicts
- Special characters in filenames
- Inconsistent metadata

## Process

1. **Scan Vault**
   - Recursively find all .md files
   - Exclude .obsidian directory
   - Build index of all notes

2. **Parse Notes**
   - Extract frontmatter from each note
   - Identify all [[wikilinks]]
   - Count words and metadata
   - Note last modified times

3. **Build Link Graph**
   - Create adjacency map of all links
   - Calculate link metrics
   - Identify orphans and hubs
   - Find broken links

4. **Generate Report**
   Create comprehensive report with:
   - Executive summary
   - Key metrics and visualizations
   - Issue highlights
   - Recommendations
   - Action items

## Report Structure

```markdown
# Vault Analysis Report
**Generated**: [Date]
**Vault**: [Path]

## Executive Summary
- [Brief overview of vault health]

## Statistics
- Total Notes: X
- Total Words: X
- Average Note Size: X words
- Date Range: YYYY-MM-DD to YYYY-MM-DD

## Link Analysis
- Total Links: X
- Orphaned Notes: X
- Broken Links: X
- Top Hub Notes:
  1. [[Note]] (X links)
  2. ...

## Issues Detected
### Critical
- [ ] X broken links
- [ ] X orphaned notes

### Warning
- [ ] X empty notes
- [ ] X duplicate titles

### Info
- [ ] X notes without tags
- [ ] X notes without frontmatter

## Recommendations
1. [Specific actionable recommendation]
2. ...

## Detailed Findings
[Comprehensive breakdown by category]
```

## Output

Display:
1. Path to generated analysis report
2. Key metrics summary
3. Critical issues count
4. Top 3 recommendations
5. Next steps

## Notes
- Large vaults (1000+ notes) may take time to analyze
- Consider generating report as markdown file in vault
- Use results to prioritize vault maintenance
- Re-run analysis periodically to track improvements
