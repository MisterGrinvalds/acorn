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

### 1. Generate Standard Commands

Every component gets these two commands in `components/$ARGUMENTS/ai/claude/commands/`:

**$ARGUMENTS-explain.md:**
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

### Output Format

Use clear headings, code examples, and practical tips.
Reference file locations and available functions where relevant.
```

**$ARGUMENTS-coach.md:**
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

### 2. Generate Tool-Specific Commands

Based on the component's capabilities, create additional commands:

| Component | Common Task Commands |
|-----------|---------------------|
| tmux | `tmux-session-create`, `tmux-layout`, `tmux-plugins`, `tmux-config` |
| go | `go-project-init`, `go-test-coverage`, `go-benchmark`, `cobra-router` |
| python | `python-venv-create`, `python-test-setup`, `python-uv-migrate` |
| node | `node-clean`, `node-deps-audit` |
| kubernetes | `k8s-deploy`, `k8s-logs`, `k8s-context-switch`, `k8s-resource-audit` |
| git | `git-branch-cleanup`, `git-rebase-guide`, `git-conflict-resolve` |
| fzf | `fzf-custom-finder`, `fzf-config` |

### 3. Check Existing Commands

Look for existing commands to preserve:
```bash
# New location
ls components/$ARGUMENTS/ai/claude/commands/$ARGUMENTS-*.md 2>/dev/null

# Legacy location
ls components/claude/config/commands/$ARGUMENTS-*.md 2>/dev/null
```

Copy any existing commands that aren't being replaced.

### 4. Report

Output:
```
Generated Commands: $ARGUMENTS
==============================

Standard commands:
  - components/$ARGUMENTS/ai/claude/commands/$ARGUMENTS-explain.md
  - components/$ARGUMENTS/ai/claude/commands/$ARGUMENTS-coach.md

Tool-specific commands:
  - components/$ARGUMENTS/ai/claude/commands/$ARGUMENTS-<task1>.md
  - components/$ARGUMENTS/ai/claude/commands/$ARGUMENTS-<task2>.md

Migrated from legacy:
  - <any copied commands>

Total: N commands

To inject: acorn ai inject $ARGUMENTS
```
