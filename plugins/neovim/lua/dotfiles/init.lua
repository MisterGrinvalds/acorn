-- lua/dotfiles/init.lua
-- Neovim plugin for dotfiles installer
--
-- A minimal plugin to install and manage the bash-profile dotfiles
-- from within Neovim - perfect for setting up new/foreign environments.
--
-- Usage:
--   require("dotfiles").setup({
--     dotfiles_root = "~/path/to/dotfiles",  -- optional, auto-detected
--     auto_setup = { enabled = true },       -- prompt on start if not installed
--   })
--
-- Commands:
--   :Dotfiles           - Interactive installer
--   :DotfilesMinimal    - Quick install (dotfiles + configs only)
--   :DotfilesComponents - Component-based installer menu
--   :DotfilesComponent <name> - Install specific component
--   :DotfilesUpdate     - Git pull dotfiles repo
--   :DotfilesStatus     - Show installation status
--   :DotfilesHelp       - Show installer help

local config = require("dotfiles.config")
local installer = require("dotfiles.installer")

local M = {}

-- Setup function - call this from your init.lua
function M.setup(opts)
  config.setup(opts)

  -- Register commands
  M.create_commands()

  -- Setup auto-setup if enabled
  if config.options.auto_setup.enabled then
    vim.api.nvim_create_autocmd("VimEnter", {
      callback = function()
        -- Defer to let Neovim fully load
        vim.defer_fn(function()
          installer.auto_setup()
        end, 100)
      end,
      once = true,
    })
  end
end

-- Create user commands
function M.create_commands()
  -- Main interactive installer
  vim.api.nvim_create_user_command("Dotfiles", function(opts)
    if opts.args == "" then
      installer.install()
    elseif opts.args == "minimal" then
      installer.install_minimal()
    elseif opts.args == "auto" then
      installer.install_auto()
    elseif opts.args == "yes" then
      installer.install_yes()
    else
      vim.notify("Unknown option: " .. opts.args, vim.log.levels.ERROR)
    end
  end, {
    nargs = "?",
    complete = function()
      return { "minimal", "auto", "yes" }
    end,
    desc = "Run dotfiles installer",
  })

  -- Minimal install (dotfiles + configs)
  vim.api.nvim_create_user_command("DotfilesMinimal", function()
    installer.install_minimal()
  end, {
    desc = "Install dotfiles (minimal: shell configs only)",
  })

  -- Component-based installer
  vim.api.nvim_create_user_command("DotfilesComponents", function()
    installer.install_components()
  end, {
    desc = "Interactive component-based installer",
  })

  -- Install specific component
  vim.api.nvim_create_user_command("DotfilesComponent", function(opts)
    installer.install_component(opts.args)
  end, {
    nargs = 1,
    complete = function()
      -- List available components
      local root = config.options.dotfiles_root
      if not root then
        return {}
      end
      local components = {}
      local comp_dir = root .. "/components"
      local handle = vim.loop.fs_scandir(comp_dir)
      if handle then
        while true do
          local name, type = vim.loop.fs_scandir_next(handle)
          if not name then
            break
          end
          if type == "directory" and name ~= "_template" then
            table.insert(components, name)
          end
        end
      end
      return components
    end,
    desc = "Install tools for a specific component",
  })

  -- List components
  vim.api.nvim_create_user_command("DotfilesListComponents", function()
    installer.list_components()
  end, {
    desc = "List all components and their status",
  })

  -- Update (git pull)
  vim.api.nvim_create_user_command("DotfilesUpdate", function()
    installer.update()
  end, {
    desc = "Update dotfiles repository (git pull)",
  })

  -- Status check
  vim.api.nvim_create_user_command("DotfilesStatus", function()
    installer.check()
  end, {
    desc = "Check dotfiles installation status",
  })

  -- Help
  vim.api.nvim_create_user_command("DotfilesHelp", function()
    installer.help()
  end, {
    desc = "Show dotfiles installer help",
  })
end

-- Expose submodules
M.config = config
M.installer = installer

-- Convenience functions (can be called directly)
M.install = installer.install
M.install_minimal = installer.install_minimal
M.install_components = installer.install_components
M.update = installer.update
M.is_installed = config.is_installed

return M
