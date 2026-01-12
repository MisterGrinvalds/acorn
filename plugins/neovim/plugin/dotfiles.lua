-- plugin/dotfiles.lua
-- Auto-loaded plugin entry point
--
-- This file is sourced automatically when the plugin directory
-- is added to Neovim's runtimepath.
--
-- Note: The actual setup is deferred until require("dotfiles").setup()
-- is called in the user's config. This allows for lazy loading.

-- Guard against double-loading
if vim.g.loaded_dotfiles then
  return
end
vim.g.loaded_dotfiles = true

-- Create a global command that works even before setup()
-- This is useful for quick installation without needing config
vim.api.nvim_create_user_command("DotfilesSetup", function()
  require("dotfiles").setup({
    auto_setup = {
      enabled = false,  -- Don't auto-prompt since user explicitly called this
    },
  })
  require("dotfiles").install()
end, {
  desc = "Quick setup and run dotfiles installer",
})
