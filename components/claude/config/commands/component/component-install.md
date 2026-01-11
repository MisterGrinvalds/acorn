---
description: Install tools for a component using declarative config
argument_hints:
  - cloudflare
  - node
  - go
  - python
  - tmux
---

Install a component's required tools via `acorn <component> install`.

Component name: $ARGUMENTS

## Instructions

### 1. Check Component Has Install Config

Read `internal/componentconfig/config/$ARGUMENTS/config.yaml` and verify it has an `install:` section.

If no `install:` section exists:
```
Component '$ARGUMENTS' does not have installation configuration.

To add installation config, run:
  /component:gen-install $ARGUMENTS
```

### 2. Show Installation Plan (Dry Run)

Run the install command in dry-run mode:

```bash
acorn $ARGUMENTS install --dry-run
```

This shows:
- Platform detected (darwin/linux)
- Tools to install
- Prerequisites resolved
- Installation methods selected

### 3. Confirm and Install

If the user confirms, run the actual installation:

```bash
acorn $ARGUMENTS install
```

For verbose output:

```bash
acorn $ARGUMENTS install --verbose
```

### 4. Report Results

```
Installation: $ARGUMENTS
========================

Platform: <darwin|linux>
Package Manager: <brew|apt>

Tools Installed:
  ✓ <tool1> - <description>
  ✓ <tool2> - <description>

Prerequisites Resolved:
  ✓ <prereq> (from <component>)

Post-Install Messages:
  - <any messages from post_install.message>

Status: Complete
```

## How Installation Works

The `internal/installer/` package handles installation:

1. **Platform Detection** - Identifies OS, distro, and package manager
2. **Prerequisite Resolution** - Recursively resolves `requires` dependencies
3. **Method Selection** - Chooses install method based on platform
4. **Execution** - Runs the appropriate install command
5. **Verification** - Runs `check` command to verify success

## Supported Install Types

| Type | Command | Example |
|------|---------|---------|
| brew | `brew install <package>` | go, tmux, fzf |
| apt | `sudo apt-get install -y <package>` | golang-go, tmux |
| npm | `npm install -g <package>` | wrangler, pnpm |
| pip | `pip install <package>` | uv |
| go | `go install <package>@latest` | cobra-cli |
| curl | `curl -fsSL <url> \| sh` | nvm, rustup |

## Adding Install Command to Components

If a component doesn't have the `install` subcommand in its CLI, you need to add it.

Example from `internal/cmd/cloudflare.go`:

```go
// Install subcommand
installCmd := &cobra.Command{
    Use:   "install",
    Short: "Install cloudflare tools",
    Long:  "Install required tools for the cloudflare component",
    RunE: func(cmd *cobra.Command, args []string) error {
        dryRun, _ := cmd.Flags().GetBool("dry-run")
        verbose, _ := cmd.Flags().GetBool("verbose")

        cfg, err := componentconfig.LoadComponent("cloudflare")
        if err != nil {
            return fmt.Errorf("failed to load config: %w", err)
        }

        inst := installer.New(verbose)
        plan, err := inst.Plan(cfg)
        if err != nil {
            return fmt.Errorf("failed to plan installation: %w", err)
        }

        if dryRun {
            plan.Print()
            return nil
        }

        result := inst.Install(plan)
        result.Print()
        return nil
    },
}
installCmd.Flags().Bool("dry-run", false, "Show what would be installed")
installCmd.Flags().Bool("verbose", false, "Show detailed output")
```

## Example Workflow

```bash
# Check what would be installed
$ acorn cf install --dry-run

Installation Plan: cloudflare
=============================
Platform: darwin (brew)

Tools to install:
  1. wrangler (npm) - CloudFlare Workers CLI
     Requires: node:npm

Prerequisites:
  - node (already installed)
  - npm (already installed)

# Actually install
$ acorn cf install

Installing cloudflare tools...
  ✓ wrangler installed via npm

Post-install:
  Run 'wrangler login' to authenticate with CloudFlare

Installation complete!
```
