---
description: Show detailed help for a component's functions and aliases
---

Show help for component: $ARGUMENTS

## Instructions

Generate comprehensive help documentation for the specified component.

### 1. Read Component Metadata

Parse `components/$ARGUMENTS/component.yaml` to extract:
- Name, version, description
- Category
- Required tools and components
- What it provides (aliases, functions, completions)

### 2. Extract Aliases

Read `components/$ARGUMENTS/aliases.sh` and list all aliases with their expansions.

### 3. Extract Functions

Read `components/$ARGUMENTS/functions.sh` and for each function:
- Function name and signature
- Any comment block above the function (description)
- Parameters if documented

### 4. Extract Environment Variables

Read `components/$ARGUMENTS/env.sh` and list:
- All exported variables
- Their default values
- Any comments describing them

### 5. Generate Help Output

## Output Format

```
================================================================================
$ARGUMENTS - <description>
================================================================================

Version: <version>
Category: <category>

DEPENDENCIES
------------
Tools:      <tool1>, <tool2> (optional: <tool3>)
Components: <component1>, <component2>

ALIASES
-------
<alias>     <expansion>
            <description if available>

FUNCTIONS
---------
<function>()
    <description>
    Usage: <example>

ENVIRONMENT VARIABLES
---------------------
<VAR_NAME>
    Default: <value>
    <description>

FILES
-----
<list any config files or data files the component uses>

EXAMPLES
--------
<usage examples from README.md if available>
```

If the component doesn't exist, report:
```
Error: Component '$ARGUMENTS' not found.

Available components:
<list all component directories>
```
