---
description: Create the expert agent definition for a component
argument_hints:
  - tmux
  - go
  - python
  - node
  - kubernetes
---

Generate agent for: $ARGUMENTS

## Instructions

Create the expert agent definition for the specified component.

### 1. Research the Component

Read the component's configuration to understand:
- What the tool does
- Available functions and aliases
- Configuration locations
- Common workflows

Check: `components/$ARGUMENTS/config.yaml`

### 2. Create Agent File

Create `components/claude/config/agents/$ARGUMENTS-expert.md`:

```yaml
---
name: $ARGUMENTS-expert
description: Expert in $ARGUMENTS <brief description of tool>
tools: Read, Write, Edit, Glob, Grep, Bash
model: sonnet
---
```

### 3. Write Agent Content

Include these sections:

**Introduction:**
```markdown
You are a **$ARGUMENTS Expert** specializing in <domain>.
```

**Core Competencies:**
- List 5-7 key areas of expertise
- Based on the tool's capabilities
- Include configuration, usage, troubleshooting

**Key Concepts:**
- Include a hierarchy or structure diagram if applicable
- List configuration file locations
- Document important environment variables

**Available Shell Functions:**
- Group by purpose (setup, daily use, management)
- Include function name and brief description
- Reference from component config.yaml

**Key Aliases:**
- List the most useful aliases
- Format: `alias` - description

**Available Commands:**
- List slash commands in `components/claude/config/commands/$ARGUMENTS/`
- Format: `/command-name` - description (project:$ARGUMENTS)

**Best Practices:**
- Include 3-5 sections with practical advice
- Cover common workflows
- Mention integration with other tools

**Component Files:**
- Reference: `components/$ARGUMENTS/config.yaml`
- Reference: `components/$ARGUMENTS/config/<tool-configs>`

**Your Approach:**
```markdown
When providing $ARGUMENTS guidance:
1. Assess current setup and user's needs
2. Recommend appropriate configuration
3. Implement with clear examples
4. Explain commands and options
5. Reference available shell functions from the dotfiles

Always reference file locations when discussing code.
```

### 4. Check for Existing Agent

Look for existing agent to migrate or reference:
- `components/claude/config/agents/$ARGUMENTS-expert.md` (current location)

If exists, preserve good content and update structure.

### 5. Report

Output:
```
Generated Agent: $ARGUMENTS-expert
==================================

Location: components/claude/config/agents/$ARGUMENTS-expert.md

Sections included:
  - Core Competencies (N items)
  - Key Concepts
  - Available Shell Functions (N items)
  - Key Aliases (N items)
  - Available Commands (N items)
  - Best Practices (N sections)
  - Component Files
  - Approach

The agent is immediately available via ~/.claude/agents/ symlink.
```
