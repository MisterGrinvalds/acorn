---
description: Create a new Obsidian note with proper structure and frontmatter
argument-hint: <note-title> [vault-path] [template]
model: claude-sonnet-4-5-20250929
allowed-tools: Write, Read, Glob
---

# Create Obsidian Note

Create a new note in an Obsidian vault with proper markdown formatting and metadata.

**Note Title**: $1
**Vault Path**: $2 (optional - defaults to current directory)
**Template**: $3 (optional - basic, daily, meeting, project)

## Steps

1. **Validate Inputs**
   - Sanitize note title (remove special characters, ensure .md extension)
   - Verify vault path exists
   - Check if note already exists (warn user)

2. **Determine Note Path**
   - If vault has folder structure, ask where to place note
   - Suggest location based on note type or template
   - Handle date-based notes (daily notes in Daily Notes folder)

3. **Generate Frontmatter**

   Standard metadata to include:
   ```yaml
   ---
   title: [Note Title]
   created: [YYYY-MM-DD HH:MM]
   updated: [YYYY-MM-DD HH:MM]
   tags: []
   aliases: []
   ---
   ```

   Template-specific additions:
   - **Daily**: date, day-of-week, weather/mood fields
   - **Meeting**: attendees, agenda, action-items
   - **Project**: status, deadline, stakeholders
   - **Basic**: minimal metadata only

4. **Generate Content Structure**

   **Basic Template**:
   ```markdown
   # [Note Title]

   ## Overview

   ## Content

   ## Related Notes
   -

   ## References
   -
   ```

   **Daily Note Template**:
   ```markdown
   # [Date - Day of Week]

   ## Tasks
   - [ ]

   ## Notes

   ## Journal

   ## Links
   - [[Yesterday]] | [[Tomorrow]]
   ```

   **Meeting Template**:
   ```markdown
   # [Meeting Title]

   **Date**: [Date]
   **Attendees**:
   **Duration**:

   ## Agenda
   1.

   ## Notes

   ## Decisions
   -

   ## Action Items
   - [ ] Task - @person - due date
   ```

5. **Create Automatic Links**
   - Link to related notes if patterns detected in title
   - Add to relevant MOC (Map of Content) if specified
   - Create backlink from related notes if requested

6. **Write File**
   - Create note at determined path
   - Set file timestamps appropriately
   - Confirm creation

## Output

Display:
1. Full path to created note
2. Summary of metadata added
3. Template used
4. Suggested next steps (tags to add, notes to link)
5. Command to open in Obsidian (if applicable)

## Notes
- Use Title Case for note titles by default
- Avoid special characters: /\:*?"<>|
- Consider using YYYY-MM-DD prefix for date-based notes
- Always include basic frontmatter for metadata
