# Claude Aggregate

Scan repositories for Claude Code agents, commands, and subagents, then aggregate them into the central dotfiles location.

## Usage

```
/claude-aggregate [directory]
```

- `directory`: Optional. Directory containing repos to scan (default: ~/Repos)

## Task

1. **Scan the target directory** for repositories containing `.claude/agents/`, `.claude/commands/`, or `.claude/subagents/`

2. **For each repository found**:
   - List all `.md` files in `.claude/agents/`, `.claude/commands/`, `.claude/subagents/`
   - Skip files that are session-specific (SESSION_CONTEXT.md, session-notes.md, etc.)
   - Skip files that are identical to existing files in the central location

3. **Handle naming conflicts**:
   - If a file with the same name already exists, compare contents
   - If contents differ, prefix the new file with the repo name (e.g., `myrepo-agent-name.md`)
   - Report conflicts to the user

4. **Copy new files** to `$DOTFILES_ROOT/components/claude/config/`:
   - Agents → `config/agents/`
   - Commands → `config/commands/`
   - Subagents → `config/subagents/`

5. **Generate a summary report**:
   - Number of repos scanned
   - New agents/commands/subagents added
   - Conflicts detected and how they were resolved
   - Files skipped (duplicates or session files)

## Example Output

```
Scanning ~/Repos for Claude Code configurations...

Repos scanned: 12
- ~/Repos/project-a: 2 agents, 3 commands
- ~/Repos/project-b: 1 command (skipped: duplicate of existing)
- ~/Repos/project-c: 1 agent (renamed: project-c-api-expert.md)

Summary:
  New agents: 3
  New commands: 3
  New subagents: 0
  Skipped (duplicates): 2
  Renamed (conflicts): 1

Files added:
  agents/api-expert.md (from project-a)
  agents/db-expert.md (from project-a)
  agents/project-c-api-expert.md (from project-c, renamed)
  commands/deploy.md (from project-a)
  commands/test-e2e.md (from project-a)
  commands/lint-fix.md (from project-a)
```

## Implementation Notes

- Use `find` or `fd` to locate `.claude` directories
- Compare file contents with `diff` or by reading both files
- Preserve file permissions when copying
- Skip any files in `.gitignore` patterns
- The central location is `$DOTFILES_ROOT/components/claude/config/`

## Arguments

$ARGUMENTS
