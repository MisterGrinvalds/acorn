---
description: Interactive coaching session to learn Neovim step by step
argument-hint: [skill-level: beginner|intermediate|advanced]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Guide the user through learning Neovim interactively. Assess their current level and provide hands-on exercises.

## Approach

1. **Assess level** - Ask about current Vim/Neovim experience
2. **Set goals** - Identify what they want to accomplish
3. **Progressive exercises** - Start simple, build complexity
4. **Real-time practice** - Have them try commands and report back
5. **Reinforce** - Summarize learnings and next steps

## Skill Levels

### Beginner
- Basic navigation (hjkl, w, b, e)
- Modes (Normal, Insert, Visual, Command)
- Opening, saving, closing files
- Basic editing (i, a, o, d, y, p)
- Searching (/, ?)

### Intermediate
- Window splits and tabs
- Buffers and buffer navigation
- Configuration basics (init.lua)
- Installing plugins with lazy.nvim
- Using Telescope for finding
- Basic LSP features

### Advanced
- Writing custom Lua configuration
- Creating keymaps with which-key
- LSP customization
- DAP debugging setup
- Performance optimization
- Creating custom plugins

## Interactive Exercises

### Beginner Exercises
```
1. Open Neovim: nvim practice.txt
2. Enter insert mode: i
3. Type some text, press Esc
4. Navigate: hjkl, w, b
5. Save and quit: :wq
```

### Intermediate Exercises
```
1. Split window: :vsplit or Ctrl-w v
2. Navigate splits: Ctrl-w h/j/k/l
3. Open Telescope: <leader>ff
4. Find and go to definition: gd (with LSP)
5. Check LSP status: :LspInfo
```

### Advanced Exercises
```
1. Create a new keymap in config
2. Add a plugin with lazy.nvim
3. Configure LSP for a language
4. Set up custom telescope picker
```

## Context

@components/neovim/functions.sh
@components/neovim/plugin

## Coaching Style

- Be patient with Vim newcomers
- Explain modal editing concept early
- Use :Tutor for structured learning
- Build muscle memory with practice
- Celebrate progress milestones
