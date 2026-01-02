---
description: Add a new target to the Makefile
argument-hint: <target-name> <description>
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Edit
---

# Add Makefile Target

Add a new target to the existing Makefile with proper formatting and documentation.

## Task

1. Verify a Makefile exists in the project root
2. Ask the user for:
   - Target name (validate kebab-case format)
   - Target description (for help system)
   - Target category/section (Development, Building, Testing, etc.)
   - Command(s) to execute
   - Dependencies (other targets this depends on)

3. Add the target to the appropriate section:
   - Insert after the section header (##@)
   - Include .PHONY declaration
   - Add help comment (##)
   - Format commands with @ prefix for silent execution
   - Add echo statements for user feedback

4. Verify the syntax is valid
5. Show the added target and update the help output

## Example Output

```makefile
##@ Development

.PHONY: watch
watch: ## Watch for changes and rebuild
	@echo "Watching for changes..."
	@air
```

## Requirements

- Follow existing Makefile style and formatting
- Add .PHONY declaration
- Include descriptive help comment
- Use @ for silent commands
- Add informative echo messages
