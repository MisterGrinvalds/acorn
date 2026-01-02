---
description: Create and manage Claude Code subagents
argument-hint: [agent-name] [description]
allowed-tools: Read, Write, Edit, Glob
---

# Claude Code Subagents

Help the user create and configure subagents.

## What Are Subagents?

Specialized AI assistants that Claude can delegate tasks to:
- Separate context window (prevents pollution)
- Specific expertise and tools
- Reusable across projects
- Can be shared with teams

## Subagent Locations

| Type | Location | Scope |
|------|----------|-------|
| Project | `.claude/agents/` | Current project (shared) |
| User | `~/.claude/agents/` | All projects (personal) |

## Your Task

Based on the user's request: $ARGUMENTS

1. Create the subagent markdown file
2. Configure appropriate frontmatter
3. Write effective system prompt

## File Format

```markdown
---
name: agent-name
description: When this agent should be used (be specific!)
tools: Read, Grep, Glob, Bash, Edit, Write
model: sonnet
permissionMode: default
skills: skill1, skill2
---

You are a [role] specializing in [domain].

## When Invoked
1. First action
2. Second action
3. Third action

## Your Approach
- Specific guidance
- Best practices
- Constraints

## Output Format
- How to structure responses
- What to include
```

## Frontmatter Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Unique identifier (lowercase, hyphens) |
| `description` | Yes | When to use (be specific for auto-selection) |
| `tools` | No | Allowed tools (inherits all if omitted) |
| `model` | No | `sonnet`, `opus`, `haiku`, or `inherit` |
| `permissionMode` | No | `default`, `acceptEdits`, `bypassPermissions` |
| `skills` | No | Skills to auto-load |

## Built-in Subagents

1. **General-Purpose** - Complex tasks, all tools, Sonnet
2. **Plan** - Research during plan mode, read-only tools
3. **Explore** - Fast codebase search, Haiku, read-only

## Example: Code Reviewer

```markdown
---
name: code-reviewer
description: Expert code review. Use PROACTIVELY after writing or modifying code.
tools: Read, Grep, Glob, Bash
model: sonnet
---

You are a senior code reviewer.

## When Invoked
1. Run git diff to see changes
2. Focus on modified files
3. Begin review immediately

## Review Checklist
- Code clarity and readability
- Proper error handling
- Security vulnerabilities
- Test coverage
- Performance considerations

## Output Format
- **Critical** (must fix)
- **Warnings** (should fix)
- **Suggestions** (nice to have)

Include specific code examples for fixes.
```

## Example: Debugger

```markdown
---
name: debugger
description: Debug errors and failures. Use PROACTIVELY when encountering issues.
tools: Read, Edit, Bash, Grep, Glob
---

You are an expert debugger.

## Process
1. Capture error and stack trace
2. Identify reproduction steps
3. Isolate failure location
4. Implement minimal fix
5. Verify solution

## For Each Issue
- Root cause explanation
- Evidence and diagnosis
- Specific code fix
- Prevention recommendations
```

## Invoking Subagents

### Automatic
Claude selects based on task and description.

### Explicit
```
> Use the code-reviewer subagent to check my changes
> Have the debugger look at this error
```

### Resume Previous
```
> Resume agent abc123 and continue the analysis
```

## Commands

- `/agents` - Manage subagents
- View, create, edit, enable/disable

## Tips

- Use "PROACTIVELY" in description for auto-use
- Limit tools to only what's needed
- Be specific in system prompts
- Check into version control for team sharing
