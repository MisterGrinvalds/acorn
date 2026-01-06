---
description: Generate Claude slash commands for a component
argument_hints:
  - tmux
  - go
  - python
  - node
  - kubernetes
---

Generate commands for: $ARGUMENTS

## Instructions

Create standard and tool-specific Claude slash commands for the component.

### 1. Create Command Directory

Create the subdirectory for this component's commands:
```
components/claude/config/commands/$ARGUMENTS/
```

### 2. Generate Standard Commands

Every component gets these two commands:

**explain.md** (or $ARGUMENTS-explain.md):
```yaml
---
description: Explain $ARGUMENTS concepts, commands, and workflows
argument_hints:
  - sessions
  - configuration
  - plugins
  - <tool-specific topics>
---

Explain $ARGUMENTS topic: $ARGUMENTS (if provided, otherwise general overview)

## Instructions

Provide clear, practical explanations for $ARGUMENTS users.

### Topics to Cover

If a specific topic is provided, focus on that. Otherwise, give an overview.

**General Overview:**
- What is $ARGUMENTS and why use it
- Basic concepts and terminology
- Getting started steps

**Configuration:**
- Config file locations
- Key settings and options
- Environment variables

**Common Workflows:**
- Daily usage patterns
- Best practices
- Integration with other tools

### Context

@components/$ARGUMENTS/config.yaml

### Output Format

Use clear headings, code examples, and practical tips.
Reference file locations and available functions where relevant.
```

**coach.md** (or $ARGUMENTS-coach.md):
```yaml
---
description: Interactive coaching session to learn $ARGUMENTS step by step
argument_hints:
  - beginner
  - intermediate
  - advanced
  - <specific-topic>
---

Coach $ARGUMENTS: $ARGUMENTS

## Instructions

Run an interactive coaching session to help the user learn $ARGUMENTS.

### Session Structure

1. **Assess Level**: Ask about current experience
2. **Set Goals**: What do they want to learn/accomplish?
3. **Teach Concepts**: Explain with examples
4. **Practice**: Guide through hands-on exercises
5. **Review**: Summarize key learnings

### Teaching Approach

- Start with fundamentals, build up
- Use practical, real-world examples
- Encourage experimentation
- Provide exercises to try
- Reference available shell functions

### Context

@components/$ARGUMENTS/config.yaml

### Skill Levels

**Beginner:**
- Installation and basic setup
- Core concepts and terminology
- Essential commands

**Intermediate:**
- Configuration customization
- Common workflows
- Troubleshooting

**Advanced:**
- Advanced features
- Scripting and automation
- Performance optimization
```

### 3. Generate Tool-Specific Commands

Based on the component's capabilities, create additional commands:

| Component | Common Task Commands |
|-----------|---------------------|
| tmux | `session-create.md`, `layout.md`, `plugins.md`, `config.md` |
| go | `project-init.md`, `test-coverage.md`, `benchmark.md` |
| python | `venv-create.md`, `test-setup.md`, `uv-migrate.md` |
| node | `clean.md`, `deps-audit.md`, `nvm-setup.md` |
| kubernetes | `deploy.md`, `logs.md`, `context-switch.md`, `resource-audit.md` |
| git | `branch-cleanup.md`, `rebase-guide.md`, `conflict-resolve.md` |
| fzf | `custom-finder.md`, `config.md` |

### 4. Check Existing Commands

Look for existing commands to preserve:
```bash
ls components/claude/config/commands/$ARGUMENTS/*.md 2>/dev/null
```

Migrate any existing commands that should be preserved.

### 5. Report

Output:
```
Generated Commands: $ARGUMENTS
==============================

Location: components/claude/config/commands/$ARGUMENTS/

Standard commands:
  - explain.md → /explain (project:$ARGUMENTS)
  - coach.md → /coach (project:$ARGUMENTS)

Tool-specific commands:
  - <task1>.md → /<task1> (project:$ARGUMENTS)
  - <task2>.md → /<task2> (project:$ARGUMENTS)

Total: N commands

Commands are immediately available via ~/.claude/commands/ symlink.
```
