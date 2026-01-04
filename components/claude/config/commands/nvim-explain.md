---
description: Explain Neovim concepts, configuration, and features
argument-hint: [topic]
allowed-tools: Read, Glob, Grep
---

## Task

Explain the requested topic about Neovim. If no specific topic provided, give an overview of Neovim and the dotfiles integration.

## Topics

If `$ARGUMENTS` is provided, explain that specific topic:

- **config** - Configuration structure and init.lua
- **lua** - Lua basics for Neovim configuration
- **options** - vim.opt and editor settings
- **keymaps** - Key mapping with vim.keymap.set
- **lazy** - lazy.nvim plugin manager
- **lsp** - Language Server Protocol integration
- **mason** - Mason LSP/DAP installer
- **treesitter** - Treesitter syntax highlighting
- **telescope** - Telescope fuzzy finder
- **completion** - nvim-cmp and completion setup
- **dap** - Debug Adapter Protocol
- **git** - Git integration plugins
- **themes** - Catppuccin and colorscheme setup
- **xdg** - XDG directory compliance
- **dotfiles-plugin** - The dotfiles.nvim integration

## Context

Reference these files for accurate information:
@components/neovim/component.yaml
@components/neovim/functions.sh
@components/neovim/aliases.sh
@components/neovim/plugin

## Response Format

1. **Concept overview** - What it is and why it matters
2. **Configuration** - How to set it up
3. **Key commands** - Useful Neovim commands
4. **Dotfiles integration** - Available functions and plugin
5. **Examples** - Practical Lua code examples
