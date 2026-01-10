# Plan: Add Installation Support to config.yaml

**Status:** Active
**Created:** 2026-01-06
**Branch:** feat/claude-tool-agents-commands

## Goal

Replace shell-based install scripts with declarative `install:` sections in config.yaml, enabling `acorn <component> install` commands.

## User Requirements

- **Separate CLI command**: `acorn <component> install`
- **Platform-based method selection**: Auto-select brew (macOS), apt (Linux), etc.
- **Recursive prerequisite installation**: Auto-install missing prerequisites

## Schema Design

Add `install:` section to config.yaml:

```yaml
install:
  tools:
    - name: wrangler
      check: "command -v wrangler"        # Verification command
      methods:
        darwin:
          type: npm
          package: wrangler
          global: true
        linux:
          type: npm
          package: wrangler
          global: true
      requires:
        - node:npm                        # component:command format
      post_install:
        message: "Run 'wrangler login' to authenticate"
```

### Supported Installation Types

| Type | Example | Use Case |
|------|---------|----------|
| `brew` | `brew install node` | macOS packages |
| `apt` | `apt install nodejs` | Debian/Ubuntu |
| `npm` | `npm install -g wrangler` | Node packages |
| `curl` | Script installers | NVM, UV |
| `go` | `go install ...` | Go tools |

## Implementation Phases

### Phase 1: Schema Extension
**Files:**
- `internal/componentconfig/schema.go` - Add types

```go
type InstallConfig struct {
    Tools []ToolInstall `yaml:"tools,omitempty"`
}

type ToolInstall struct {
    Name        string                   `yaml:"name"`
    Check       string                   `yaml:"check"`
    Methods     map[string]InstallMethod `yaml:"methods"`
    Requires    []string                 `yaml:"requires,omitempty"`
    PostInstall PostInstallConfig        `yaml:"post_install,omitempty"`
}

type InstallMethod struct {
    Type    string `yaml:"type"`              // brew, apt, npm, curl, go
    Package string `yaml:"package,omitempty"` // Package name if different
    Global  bool   `yaml:"global,omitempty"`  // For npm/pip
    URL     string `yaml:"url,omitempty"`     // For curl type
}

type PostInstallConfig struct {
    Message string `yaml:"message,omitempty"`
}
```

### Phase 2: Installer Package
**Create:** `internal/installer/`

```
internal/installer/
├── installer.go      # Main Installer, Plan(), Install()
├── platform.go       # DetectPlatform(), GetMethodKey()
├── resolver.go       # Prerequisite resolution
├── methods.go        # BrewExecutor, NpmExecutor, AptExecutor
├── types.go          # InstallPlan, InstallResult, Platform
└── installer_test.go
```

**Key Functions:**

```go
// installer.go
func NewInstaller(opts ...Option) *Installer
func (i *Installer) Plan(ctx context.Context, component string) (*InstallPlan, error)
func (i *Installer) Install(ctx context.Context, component string) (*InstallResult, error)

// platform.go
func DetectPlatform() *Platform  // Returns OS, distro, package manager

// resolver.go - Recursive prerequisite resolution
func (r *Resolver) BuildPlan(ctx, component, cfg, platform) (*InstallPlan, error)
```

### Phase 3: CLI Integration
**Files:**
- `internal/cmd/cloudflare.go` - Add install subcommand (pilot)

```go
var cfInstallCmd = &cobra.Command{
    Use:   "install",
    Short: "Install CloudFlare CLI tools",
    RunE:  runCfInstall,
}

func runCfInstall(cmd *cobra.Command, args []string) error {
    inst := installer.NewInstaller(
        installer.WithDryRun(cfDryRun),
        installer.WithVerbose(cfVerbose),
    )
    return inst.Install(cmd.Context(), "cloudflare")
}
```

### Phase 4: Config Migration
**Files:**
- `internal/componentconfig/config/cloudflare/config.yaml` - Add install section
- `internal/componentconfig/config/node/config.yaml` - Add install section (for npm prerequisite)
- `components/cloudflare/install/install.sh` - Delete after migration

**Example cloudflare config.yaml addition:**

```yaml
install:
  tools:
    - name: wrangler
      check: "command -v wrangler"
      methods:
        darwin:
          type: npm
          package: wrangler
          global: true
        linux:
          type: npm
          package: wrangler
          global: true
      requires:
        - node:npm
      post_install:
        message: "Run 'wrangler login' to authenticate with CloudFlare"
```

## Prerequisite Resolution Flow

1. Parse `requires` list (format: `command` or `component:command`)
2. Check if command exists via `tools.CommandExists()`
3. If missing and component specified, load that component's install config
4. Recursively resolve that tool's prerequisites
5. Build topological install order: prerequisites first

**Example:**
```
cloudflare install → wrangler needs npm → node component → npm needs node → brew install node
```

## Files to Modify

| File | Change |
|------|--------|
| `internal/componentconfig/schema.go` | Add InstallConfig, ToolInstall, InstallMethod types |
| `internal/cmd/cloudflare.go` | Add `cfInstallCmd` subcommand |
| `internal/componentconfig/config/cloudflare/config.yaml` | Add install section |
| `internal/componentconfig/config/node/config.yaml` | Add install section |

## Files to Create

| File | Purpose |
|------|---------|
| `internal/installer/installer.go` | Main installer with Plan/Install |
| `internal/installer/platform.go` | Platform detection |
| `internal/installer/resolver.go` | Prerequisite resolution |
| `internal/installer/methods.go` | Package manager executors |
| `internal/installer/types.go` | Type definitions |
| `internal/installer/installer_test.go` | Tests |

## Files to Delete

| File | Reason |
|------|--------|
| `components/cloudflare/install/install.sh` | Replaced by config.yaml |

## CLI Usage After Implementation

```bash
# Install cloudflare tools (auto-installs npm via node if needed)
acorn cf install

# Show what would be installed
acorn cf install --dry-run

# Verbose output
acorn cf install -v
```

## Reuses Existing Code

- `internal/tools/checker.go` → `CommandExists()` for checking prerequisites
- `internal/componentconfig/loader.go` → Load install configs
- Existing `cfDryRun`, `cfVerbose` flags in cloudflare.go
