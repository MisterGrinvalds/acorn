---
description: Optimize Makefile for better performance
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Edit
---

# Optimize Makefile

Review and optimize the Makefile for performance, maintainability, and best practices.

## Task

1. Analyze the current Makefile for:
   - **Performance**: Parallel execution opportunities, unnecessary rebuilds
   - **Organization**: Logical grouping, section ordering
   - **Maintainability**: Variable usage, DRY principles
   - **Best practices**: .PHONY usage, error handling, cross-platform support
   - **Documentation**: Help completeness, clarity

2. Identify optimization opportunities:
   - Consolidate duplicate commands into variables
   - Extract repeated patterns into reusable functions
   - Add dependency chains for common workflows
   - Optimize slow targets (caching, conditional execution)
   - Remove unused targets or variables

3. Suggest improvements:
   - Add composite targets (e.g., `pre-commit: fmt vet lint test`)
   - Implement conditional logic for platform-specific code
   - Add tool detection and installation helpers
   - Include version management patterns
   - Add CI/CD integration targets

4. Show before/after comparisons for major changes

5. Apply approved optimizations

## Optimization Categories

**Performance**:
- Parallel target execution
- Conditional rebuilds
- Caching strategies

**Code Quality**:
- Variable extraction
- Function definitions
- Dependency chains

**Developer Experience**:
- Composite workflow targets
- Better error messages
- Tool installation automation

## Output Format

```
Makefile Optimization Report
============================

Performance Improvements:
• Extract repeated 'go build' flags to GO_BUILD_FLAGS variable
• Add dependency chain: ci -> deps lint test build

Code Quality:
• Consolidate 3 similar build targets using pattern rules
• Extract platform detection to variables

Developer Experience:
• Add 'pre-commit' composite target
• Add 'install-tools' for developer setup

Estimated impact: 30% fewer lines, 2x faster common workflows

Would you like me to apply these optimizations? (y/n)
```
