---
description: Create a structured execution plan for the requested work
argument-hint: <task-description>
model: claude-sonnet-4-5-20250929
allowed-tools: Write, Edit, Read, Glob
---

# Create Execution Plan

Create a structured execution plan for the requested work.

**Task**: $ARGUMENTS

## Your Process

1. **Analyze the Request**
   - Understand the scope and requirements
   - Identify key deliverables
   - Note any ambiguities or unknowns

2. **Break Down the Work**
   - Divide into logical phases
   - Create atomic, actionable tasks
   - Identify dependencies between tasks
   - Estimate complexity (simple/moderate/complex)

3. **Create Plan Document**
   - Write to `PLAN-[feature-name].md`
   - Include: objective, phases, tasks, dependencies, risks
   - Use clear section headers
   - Number tasks for easy reference

4. **Update TODO.md**
   - Add first set of actionable tasks
   - Mark dependencies clearly
   - Set priorities (high/medium/low)
   - Include task references to plan

5. **Present to User**
   - Summarize the plan
   - Highlight critical path
   - Note any decisions needed
   - Suggest starting point

## Plan Structure Template

```markdown
# Plan: [Feature Name]

**Created**: YYYY-MM-DD
**Status**: [Planning/In Progress/Complete]
**Owner**: [If applicable]

## Objective
[What are we building and why?]

## Success Criteria
- [ ] Criterion 1
- [ ] Criterion 2

## Phases

### Phase 1: [Name]
**Goal**: [What this phase achieves]
**Tasks**:
1. [Task 1] - [Complexity: Simple/Moderate/Complex]
2. [Task 2]

**Dependencies**: None / [List]

### Phase 2: [Name]
...

## Risks & Considerations
- Risk 1: [Impact, mitigation]
- Risk 2: [Impact, mitigation]

## Next Steps
1. [Immediate action]
2. [Follow-up action]
```

## Output

Provide:
1. Path to created plan file
2. Summary of plan (2-3 sentences)
3. First 3-5 tasks added to TODO.md
4. Recommended starting point

Be concise but thorough. Focus on actionable tasks, not abstract concepts.
