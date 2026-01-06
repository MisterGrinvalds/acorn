---
description: Initialize a new Obsidian vault with best practice structure
argument-hint: <vault-path> [organization-type]
model: claude-sonnet-4-5-20250929
allowed-tools: Write, Bash(mkdir:*), Bash(touch:*)
---

# Initialize Obsidian Vault

Create a new Obsidian vault at the specified path with a well-organized folder structure.

**Vault Path**: $1
**Organization Type**: $2 (optional: minimal, para, zettelkasten - defaults to minimal)

## Steps

1. **Validate Path**
   - Check if the path already exists
   - Ensure parent directory is writable
   - Ask for confirmation if vault directory exists

2. **Create Vault Structure**

   Based on organization type:

   **Minimal** (default):
   ```
   vault/
   ├── .obsidian/          # Created by Obsidian on first open
   ├── Templates/          # Note templates
   ├── Daily Notes/        # Daily journal entries
   └── README.md           # Vault documentation
   ```

   **PARA**:
   ```
   vault/
   ├── .obsidian/
   ├── Templates/
   ├── Projects/           # Active work with deadlines
   ├── Areas/              # Ongoing responsibilities
   ├── Resources/          # Reference materials
   ├── Archive/            # Completed/inactive items
   └── README.md
   ```

   **Zettelkasten**:
   ```
   vault/
   ├── .obsidian/
   ├── Templates/
   ├── Fleeting/           # Quick captures
   ├── Literature/         # Reference notes
   ├── Permanent/          # Processed knowledge
   └── README.md
   ```

3. **Create Welcome Note**
   - Create README.md with:
     - Vault purpose and scope
     - Organizational system explanation
     - Basic usage guidelines
     - Template instructions

4. **Create Essential Templates**
   - Daily note template with date and metadata
   - Basic note template with frontmatter
   - Meeting note template
   - Project template (if PARA)

5. **Create .obsidian Directory**
   - Create basic `.obsidian/` directory structure
   - Add `.gitignore` to exclude workspace settings if desired

6. **Provide Next Steps**
   - Instructions to open vault in Obsidian
   - Recommended initial settings
   - Suggested plugins for the workflow
   - Link to Obsidian documentation

## Output

Display:
1. Vault path and structure created
2. Number of templates added
3. Instructions for opening in Obsidian
4. Next configuration steps
5. Path to README.md for reference

## Notes
- Obsidian will create `.obsidian/` folder automatically on first open
- User can customize the structure after creation
- Templates use YAML frontmatter for metadata
- Consider git initialization for version control
