---
description: Install tools for a component or update installer with component
---

Install a component's required tools or update the installer to include a component.

Component name: $ARGUMENTS

## Instructions

### 1. Validate Component Exists

Check that `components/$ARGUMENTS/component.yaml` exists and is valid.

### 2. Read Component Setup Info

Extract from `component.yaml`:
- `requires.tools` - CLI tools that must be installed
- `setup.brew` - Homebrew packages (macOS)
- `setup.apt` - APT packages (Linux)
- `setup.post_install` - Post-installation script

### 3. Check Current Tool Status

For each tool in `requires.tools`:
- Check if already installed with `command -v`
- Report which tools are missing

### 4. Installation Options

Offer the user these options:
1. **Install now** - Run brew/apt install for missing tools
2. **Update installer** - Add an install function to `install.sh`
3. **Generate script** - Create a standalone installation script
4. **Skip** - Just show what would be installed

### 5. If "Install now" Selected

```bash
# macOS
brew install <packages from setup.brew>

# Linux
sudo apt-get install -y <packages from setup.apt>
```

### 6. If "Update installer" Selected

Add a function to `install.sh`:
```bash
# Install <component> tools
install_<component>_tools() {
    log_info "Installing <component> tools..."

    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install <brew packages>
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        sudo apt-get install -y <apt packages>
    fi

    log_success "<component> tools installed!"
}
```

Then add to the appropriate section (dev-tools, cloud-tools, ai-tools, etc.) or create a new component-specific install option.

### 7. Report Summary

```
Component: <name>
Category: <category>
Required tools: <list>
Status: <installed/missing>

Installation options:
- brew: <packages>
- apt: <packages>
- post_install: <script if any>

Installer: <updated/not updated>
```

## Example Output

```
Component: cloudflare
Category: cloud
Description: CloudFlare CLI (wrangler) integration

Required Tools:
  - wrangler: NOT INSTALLED

Setup Configuration:
  - brew: [] (npm-based, use: npm install -g wrangler)
  - apt: []

Recommended Installation:
  npm install -g wrangler

Would you like to:
1. Install wrangler now (npm install -g wrangler)
2. Add to installer under cloud tools
3. Skip
```
