---
description: Create a Map of Content (MOC) note to organize related notes
argument-hint: <moc-topic> [vault-path]
model: claude-sonnet-4-5-20250929
allowed-tools: Write, Read, Glob, Grep
---

# Create Map of Content (MOC)

Create a Map of Content note to organize and navigate related notes on a specific topic.

**Topic**: $1
**Vault Path**: $2 (optional - defaults to current directory)

## What is a MOC?

A **Map of Content** is a special type of note that:
- Serves as a navigation hub for related notes
- Organizes links to notes on a specific topic
- Provides structure and context
- Acts as an entry point to a knowledge domain
- Reduces cognitive load by grouping related content

**MOC vs Folder**: MOCs use links to organize, not file system hierarchy

## Steps

1. **Analyze Topic**
   - Search vault for notes related to topic
   - Identify common themes and subtopics
   - Find highly connected notes
   - Determine organizational structure

2. **Gather Related Notes**
   - Search by tags related to topic
   - Find notes mentioning topic keywords
   - Check backlinks from seed notes
   - Include both main and peripheral notes

3. **Organize Structure**

   Choose organizational approach:
   - **Hierarchical**: Top-down categorization
   - **Chronological**: Timeline-based
   - **Progression**: Beginner to advanced
   - **Functional**: By use case or application
   - **Custom**: Based on specific needs

4. **Create MOC Note**

Standard MOC structure:
```markdown
---
title: MOC - {{topic}}
type: moc
tags: [moc, {{topic}}]
created: {{date}}
updated: {{date}}
---

# MOC - {{topic}}

> [!info] About this MOC
> This Map of Content organizes all notes related to {{topic}}.
> Use it as a starting point for exploring this topic.

## Overview
[Brief introduction to the topic and what this MOC covers]

## Core Concepts
- [[Fundamental Concept 1]]
- [[Fundamental Concept 2]]
- [[Fundamental Concept 3]]

## Topics

### Subtopic 1
- [[Note 1]] - Brief description
- [[Note 2]] - Brief description
- [[Note 3]] - Brief description

### Subtopic 2
- [[Note 4]] - Brief description
- [[Note 5]] - Brief description

## Resources
- [[Reference 1]]
- [[Reference 2]]

## Related MOCs
- [[MOC - Related Topic 1]]
- [[MOC - Related Topic 2]]

## Unorganized
[Notes that need categorization]
- [[New Note]]

---
**Status**: 🔵 Active | 📊 {{count}} notes
**Last Updated**: {{date}}
```

5. **Link Notes to MOC**
   - Optionally add backlinks from notes to MOC
   - Add "see also: [[MOC - Topic]]" section
   - Update note frontmatter with MOC tag

6. **Add Metadata**
   - Count of included notes
   - Last updated date
   - Status indicator
   - Related tags

## MOC Patterns

### 1. Topic MOC
Organize by subject matter:
```
# MOC - Machine Learning

## Fundamentals
- [[Neural Networks]]
- [[Gradient Descent]]

## Algorithms
### Supervised Learning
- [[Linear Regression]]
- [[Decision Trees]]

### Unsupervised Learning
- [[K-Means Clustering]]
```

### 2. Project MOC
Organize project-related notes:
```
# MOC - Project Name

## Planning
- [[Project Proposal]]
- [[Requirements]]

## Implementation
- [[Architecture Design]]
- [[Development Notes]]

## Meetings
- [[2025-10-15 Kickoff]]
```

### 3. Chronological MOC
Organize by time:
```
# MOC - 2025 Learning Journey

## January
- [[Week 1 Progress]]

## February
- [[New Skills Acquired]]
```

### 4. Index MOC
High-level vault navigation:
```
# MOC - Index

## Areas
- [[MOC - Work]]
- [[MOC - Personal]]

## Projects
- [[MOC - Project A]]

## Resources
- [[MOC - Books]]
```

## Tips for Effective MOCs

1. **Keep it Scannable**: Use clear headings and bullets
2. **Add Context**: Brief descriptions help navigation
3. **Use Callouts**: Highlight important notes or sections
4. **Regular Updates**: Keep MOC current as notes evolve
5. **Link Between MOCs**: Create MOC network for related topics
6. **Status Indicators**: Use emoji or text to show MOC health
7. **Start Simple**: Begin with basic structure, evolve over time

## Output

Display:
1. Path to created MOC
2. Number of notes included
3. Organizational structure used
4. Suggestions for related MOCs
5. Next steps (update note backlinks, create missing notes)

## Notes
- MOCs are living documents - update regularly
- Use "MOC -" prefix in title for easy identification
- Consider creating MOCs when you have 10+ notes on a topic
- MOCs can link to other MOCs (MOC networks)
- Don't over-organize - not every topic needs a MOC
