-- lua/dotfiles/installer.lua
-- Installer interface for running dotfiles install commands

local config = require("dotfiles.config")

local M = {}

-- Notify helper
local function notify(msg, level)
  if config.options.notify and config.options.notify.enabled then
    vim.notify(msg, level or vim.log.levels.INFO, { title = "Dotfiles" })
  end
end

-- Run a shell command and return output
local function run_cmd(cmd, opts)
  opts = opts or {}
  local output = {}

  if opts.async then
    -- Async execution with terminal
    vim.cmd("split | terminal " .. cmd)
    vim.cmd("startinsert")
    return nil
  else
    -- Sync execution
    local handle = io.popen(cmd .. " 2>&1")
    if handle then
      for line in handle:lines() do
        table.insert(output, line)
      end
      handle:close()
    end
    return output
  end
end

-- Check installer availability
function M.check()
  local installer = config.get_installer_path()
  if not installer then
    notify("Dotfiles installer not found. Set dotfiles_root in config.", vim.log.levels.ERROR)
    return false
  end

  local root = config.options.dotfiles_root
  notify("Dotfiles root: " .. root, vim.log.levels.INFO)
  notify("Installer: " .. installer, vim.log.levels.INFO)
  notify("Installed: " .. tostring(config.is_installed()), vim.log.levels.INFO)
  return true
end

-- Install dotfiles (minimal - just dotfiles + configs)
function M.install_minimal()
  local installer = config.get_installer_path()
  if not installer then
    notify("Installer not found", vim.log.levels.ERROR)
    return false
  end

  notify("Running minimal install (dotfiles + app configs)...", vim.log.levels.INFO)

  -- Run dotfiles-only install, then link configs
  local cmd = string.format(
    "cd %s && ./install.sh --dotfiles && ./install.sh --yes-to-all --skip-gui 2>&1 | head -50",
    vim.fn.shellescape(config.options.dotfiles_root)
  )

  -- Use terminal for better output
  vim.cmd("botright split | resize 15 | terminal bash -c " .. vim.fn.shellescape(cmd))
  vim.cmd("startinsert")

  return true
end

-- Install with specific options
function M.install(opts)
  opts = opts or {}
  local installer = config.get_installer_path()
  if not installer then
    notify("Installer not found", vim.log.levels.ERROR)
    return false
  end

  local args = config.build_install_args(opts)
  local args_str = table.concat(args, " ")

  notify("Running installer with: " .. (args_str ~= "" and args_str or "(interactive)"), vim.log.levels.INFO)

  local cmd = string.format(
    "cd %s && ./install.sh %s",
    vim.fn.shellescape(config.options.dotfiles_root),
    args_str
  )

  -- Always use terminal for interactive installer
  vim.cmd("botright split | resize 20 | terminal bash -c " .. vim.fn.shellescape(cmd))
  vim.cmd("startinsert")

  return true
end

-- Install dotfiles only
function M.install_dotfiles()
  return M.install({ dotfiles_only = true })
end

-- Install with --auto flag (non-interactive, full install)
function M.install_auto()
  return M.install({ auto = true })
end

-- Install with --yes-to-all flag
function M.install_yes()
  return M.install({ yes_to_all = true, skip_gui = true })
end

-- Interactive component installer
function M.install_components()
  return M.install({ components = true })
end

-- Install specific component
function M.install_component(name)
  if not name or name == "" then
    notify("Component name required", vim.log.levels.ERROR)
    return false
  end
  return M.install({ component = name })
end

-- List components and their status
function M.list_components()
  local installer = config.get_installer_path()
  if not installer then
    notify("Installer not found", vim.log.levels.ERROR)
    return false
  end

  local cmd = string.format(
    "cd %s && ./install.sh --list-components",
    vim.fn.shellescape(config.options.dotfiles_root)
  )

  vim.cmd("botright split | resize 25 | terminal bash -c " .. vim.fn.shellescape(cmd))
  vim.cmd("startinsert")

  return true
end

-- Update dotfiles (git pull + reload)
function M.update()
  local root = config.options.dotfiles_root
  if not root then
    notify("Dotfiles root not set", vim.log.levels.ERROR)
    return false
  end

  notify("Updating dotfiles...", vim.log.levels.INFO)

  local cmd = string.format("cd %s && git pull", vim.fn.shellescape(root))
  vim.cmd("botright split | resize 10 | terminal bash -c " .. vim.fn.shellescape(cmd))
  vim.cmd("startinsert")

  return true
end

-- Show installer help
function M.help()
  local installer = config.get_installer_path()
  if not installer then
    notify("Installer not found", vim.log.levels.ERROR)
    return false
  end

  local cmd = string.format(
    "cd %s && ./install.sh --help",
    vim.fn.shellescape(config.options.dotfiles_root)
  )

  vim.cmd("botright split | resize 30 | terminal bash -c " .. vim.fn.shellescape(cmd))

  return true
end

-- Auto-setup hook (called on startup if configured)
function M.auto_setup()
  if not config.options.auto_setup.enabled then
    return
  end

  -- Check if already installed
  if config.options.auto_setup.check_installed and config.is_installed() then
    if not config.options.auto_setup.silent then
      notify("Dotfiles already installed", vim.log.levels.DEBUG)
    end
    return
  end

  -- Check if installer exists
  local installer = config.get_installer_path()
  if not installer then
    if not config.options.auto_setup.silent then
      notify("Dotfiles installer not found - skipping auto-setup", vim.log.levels.WARN)
    end
    return
  end

  -- Prompt user before auto-install
  vim.ui.select(
    { "Yes, install now", "No, skip", "Never ask again" },
    {
      prompt = "Dotfiles not installed. Install now?",
    },
    function(choice)
      if choice == "Yes, install now" then
        if config.options.auto_setup.minimal then
          M.install_minimal()
        else
          M.install()
        end
      elseif choice == "Never ask again" then
        -- Could persist this preference
        notify("Auto-setup disabled for this session", vim.log.levels.INFO)
      end
    end
  )
end

return M
