---
description: Create a reusable note template for Obsidian
argument-hint: <template-name> [vault-path]
model: claude-sonnet-4-5-20250929
allowed-tools: Write, Read
---

# Create Obsidian Template

Create a reusable note template with frontmatter, structure, and placeholders.

**Template Name**: $1
**Vault Path**: $2 (optional - defaults to current directory)

## Template Types

### 1. Daily Note
For daily journaling and task tracking:
```markdown
---
date: {{date:YYYY-MM-DD}}
day: {{date:dddd}}
week: {{date:YYYY-[W]ww}}
tags: [daily-note]
---

# {{date:YYYY-MM-DD}} - {{date:dddd}}

## Tasks
- [ ]

## Notes

## Journal

## Grateful For
-

## Tomorrow's Focus
-

---
[[{{date-1d:YYYY-MM-DD}}|Yesterday]] | [[{{date+1d:YYYY-MM-DD}}|Tomorrow]]
```

### 2. Meeting Note
For meeting notes and action items:
```markdown
---
title: {{title}}
date: {{date:YYYY-MM-DD}}
time: {{time}}
type: meeting
attendees: []
tags: [meeting]
---

# {{title}}

**Date**: {{date:YYYY-MM-DD}}
**Time**: {{time}}
**Attendees**:
**Location**:

## Agenda
1.

## Discussion
### Topic 1


## Decisions
-

## Action Items
- [ ] Task - @person - {{date:YYYY-MM-DD}}

## Follow-up
- Next meeting:
- Related: [[]]

## Notes

```

### 3. Project Note
For project planning and tracking:
```markdown
---
title: {{title}}
created: {{date:YYYY-MM-DD}}
status: planning
priority: medium
deadline:
tags: [project]
---

# {{title}}

## Overview
**Status**: 🔵 Planning
**Priority**: Medium
**Deadline**:
**Owner**:

## Objectives
-

## Scope
### In Scope
-

### Out of Scope
-

## Milestones
- [ ] Milestone 1 - Date
- [ ] Milestone 2 - Date

## Tasks
- [ ]

## Resources
- [[]]

## Notes

## Status Updates
### {{date:YYYY-MM-DD}}
-
```

### 4. Book/Article Note
For literature and reference notes:
```markdown
---
title: {{title}}
author:
date-read: {{date:YYYY-MM-DD}}
type: literature
rating:
tags: [book, literature]
---

# {{title}}

**Author**:
**Type**: Book / Article / Paper
**Date Read**: {{date:YYYY-MM-DD}}
**Rating**: ⭐⭐⭐⭐⭐

## Summary
[Brief summary in own words]

## Key Concepts
-

## Quotes
>

## My Thoughts
-

## Related Notes
- [[]]

## References
-
```

### 5. Person Note
For contact and relationship notes:
```markdown
---
title: {{title}}
type: person
tags: [people]
---

# {{title}}

## Contact
- **Email**:
- **Phone**:
- **LinkedIn**:
- **Location**:

## Context
[How I know this person]

## Interactions
### {{date:YYYY-MM-DD}}
-

## Projects
- [[]]

## Notes
-
```

### 6. Basic Note
Minimal template for general notes:
```markdown
---
title: {{title}}
created: {{date:YYYY-MM-DD}}
updated: {{date:YYYY-MM-DD}}
tags: []
---

# {{title}}

## Content


## Related Notes
- [[]]
```

## Steps

1. **Determine Template Type**
   - Ask user which type if not specified
   - Show available template types
   - Explain use case for each

2. **Customize Template**
   - Ask for custom fields to include
   - Determine frontmatter properties
   - Add project-specific sections
   - Include custom placeholders

3. **Create Template File**
   - Save to `Templates/` folder in vault
   - Use `.md` extension
   - Name descriptively (e.g., `Template - Daily Note.md`)

4. **Add Documentation**
   - Include comment at top explaining usage
   - Document placeholders ({{}} syntax)
   - Show example usage
   - Link to related templates

5. **Configure Obsidian**
   - Instructions for setting Templates folder
   - How to insert template (Ctrl/Cmd + T)
   - Configure Templater plugin if needed
   - Set hotkeys for common templates

## Placeholder Syntax

Obsidian supports these placeholders (with Templater plugin):

- `{{title}}` - Note title
- `{{date:FORMAT}}` - Current date in specified format
- `{{date±Xd:FORMAT}}` - Date offset by X days
- `{{time}}` - Current time
- Custom placeholders with Templater

Common date formats:
- `YYYY-MM-DD` - 2025-11-02
- `YYYY-MM-DD HH:mm` - 2025-11-02 14:30
- `dddd, MMMM Do YYYY` - Saturday, November 2nd 2025
- `YYYY-[W]ww` - 2025-W44

## Output

Display:
1. Path to created template file
2. Template type and key features
3. Instructions for using template in Obsidian
4. Suggested hotkey bindings
5. Next steps (create more templates, configure Templater)

## Notes
- Templates require Obsidian's core Templates plugin or Templater plugin
- Placeholders are replaced when template is inserted
- Store all templates in designated Templates folder
- Use descriptive names with "Template -" prefix
- Document custom placeholders in template file
