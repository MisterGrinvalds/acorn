-- lua/dotfiles/config.lua
-- Configuration module for dotfiles installer plugin

local M = {}

-- Default configuration matching install.sh options
M.defaults = {
  -- Path to dotfiles repository (auto-detected if nil)
  dotfiles_root = nil,

  -- Installation options (maps to install.sh flags)
  install = {
    dotfiles = true,          -- Install shell bootstrap files
    app_configs = true,       -- Link application configs (git, ssh)
    package_manager = false,  -- Install package manager tools
    dev_tools = false,        -- Install development tools
    cloud_tools = false,      -- Install cloud/k8s tools
    skip_gui = true,          -- Skip GUI apps (VS Code, Docker) - default true for remote
  },

  -- Component-based installation
  components = {
    enabled = false,          -- Use component-based installation
    categories = {},          -- Categories to install: "core", "dev", "cloud", "ai", "database"
    list = {},                -- Specific components to install
  },

  -- Auto-setup on Neovim start
  auto_setup = {
    enabled = false,          -- Run setup automatically on start
    check_installed = true,   -- Only run if not already installed
    silent = false,           -- Suppress output during auto-setup
    minimal = true,           -- Use minimal install (dotfiles + configs only)
  },

  -- Notifications
  notify = {
    enabled = true,           -- Show notifications
    level = vim.log.levels.INFO,
  },
}

-- Current configuration (populated by setup())
M.options = {}

-- Detect DOTFILES_ROOT from environment or common locations
local function detect_dotfiles_root()
  -- Check environment variable first
  local env_root = vim.env.DOTFILES_ROOT
  if env_root and vim.fn.isdirectory(env_root) == 1 then
    return env_root
  end

  -- Check common locations
  local common_paths = {
    vim.fn.expand("~/.config/dotfiles"),
    vim.fn.expand("~/.dotfiles"),
    vim.fn.expand("~/dotfiles"),
    vim.fn.expand("~/Repos/personal/bash-profile"),
  }

  for _, path in ipairs(common_paths) do
    if vim.fn.isdirectory(path) == 1 and vim.fn.filereadable(path .. "/install.sh") == 1 then
      return path
    end
  end

  return nil
end

-- Check if dotfiles are already installed
function M.is_installed()
  local bashrc = vim.fn.expand("~/.bashrc")
  local zshrc = vim.fn.expand("~/.zshrc")

  -- Check if either rc file exists and contains DOTFILES_ROOT
  for _, rc_file in ipairs({ bashrc, zshrc }) do
    if vim.fn.filereadable(rc_file) == 1 then
      local content = vim.fn.readfile(rc_file)
      for _, line in ipairs(content) do
        if line:match("DOTFILES_ROOT") then
          return true
        end
      end
    end
  end

  return false
end

-- Get install.sh path
function M.get_installer_path()
  local root = M.options.dotfiles_root
  if not root then
    return nil
  end
  local installer = root .. "/install.sh"
  if vim.fn.filereadable(installer) == 1 then
    return installer
  end
  return nil
end

-- Build command line args from options
function M.build_install_args(opts)
  opts = opts or {}
  local args = {}

  -- Mode flags
  if opts.auto then
    table.insert(args, "--auto")
  elseif opts.yes_to_all then
    table.insert(args, "--yes-to-all")
  end

  -- Specific install targets
  if opts.dotfiles_only then
    table.insert(args, "--dotfiles")
  elseif opts.dev_tools_only then
    table.insert(args, "--dev-tools")
  elseif opts.cloud_tools_only then
    table.insert(args, "--cloud-tools")
  end

  -- Skip GUI apps
  if opts.skip_gui or M.options.install.skip_gui then
    table.insert(args, "--skip-gui")
  end

  -- Component-based
  if opts.components then
    table.insert(args, "--components")
  elseif opts.component then
    table.insert(args, "--component")
    table.insert(args, opts.component)
  end

  return args
end

-- Setup function - merges user config with defaults
function M.setup(opts)
  opts = opts or {}
  M.options = vim.tbl_deep_extend("force", {}, M.defaults, opts)

  -- Auto-detect dotfiles root if not specified
  if not M.options.dotfiles_root then
    M.options.dotfiles_root = detect_dotfiles_root()
  end

  return M.options
end

return M
