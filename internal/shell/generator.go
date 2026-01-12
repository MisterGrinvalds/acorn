package shell

import (
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/componentconfig"
	"github.com/mistergrinvalds/acorn/internal/configfile"
)

// Generator creates shell scripts from component configuration.
type Generator struct {
	platform string
	dryRun   bool
}

// NewGenerator creates a new Generator for the current platform.
func NewGenerator() *Generator {
	return &Generator{
		platform: runtime.GOOS,
		dryRun:   false,
	}
}

// NewGeneratorWithDryRun creates a new Generator with dry run mode.
func NewGeneratorWithDryRun(dryRun bool) *Generator {
	return &Generator{
		platform: runtime.GOOS,
		dryRun:   dryRun,
	}
}

// Generate creates a complete shell script from a BaseConfig.
func (g *Generator) Generate(cfg *componentconfig.BaseConfig) string {
	var b strings.Builder

	// Generate env section
	if len(cfg.Env) > 0 {
		g.generateEnv(&b, cfg.Env)
	}

	// Generate path additions
	if len(cfg.Paths) > 0 {
		g.generatePaths(&b, cfg.Paths)
	}

	// Generate aliases
	if len(cfg.Aliases) > 0 {
		g.generateAliases(&b, cfg.Aliases)
	}

	// Generate wrapper functions
	if len(cfg.Wrappers) > 0 {
		g.generateWrappers(&b, cfg.Wrappers)
	}

	// Include raw shell functions
	if len(cfg.ShellFunctions) > 0 {
		g.generateShellFunctions(&b, cfg.ShellFunctions)
	}

	return b.String()
}

// generateEnv generates environment variable exports.
func (g *Generator) generateEnv(b *strings.Builder, env map[string]string) {
	// Sort keys for deterministic output
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		b.WriteString(fmt.Sprintf("export %s=\"%s\"\n", k, v))
	}
	b.WriteString("\n")
}

// generatePaths generates PATH additions.
func (g *Generator) generatePaths(b *strings.Builder, paths []componentconfig.PathEntry) {
	for _, p := range paths {
		// Check platform condition
		if p.Condition != "" && p.Condition != g.platform {
			continue
		}

		// Add to PATH if not already present
		b.WriteString(fmt.Sprintf("case \":$PATH:\" in\n"))
		b.WriteString(fmt.Sprintf("    *\":%s:\"*) ;;\n", p.Path))
		b.WriteString(fmt.Sprintf("    *) export PATH=\"%s:$PATH\" ;;\n", p.Path))
		b.WriteString("esac\n")
	}
	b.WriteString("\n")
}

// generateAliases generates shell aliases.
func (g *Generator) generateAliases(b *strings.Builder, aliases map[string]string) {
	// Sort keys for deterministic output
	keys := make([]string, 0, len(aliases))
	for k := range aliases {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		cmd := aliases[name]
		b.WriteString(fmt.Sprintf("alias %s='%s'\n", name, cmd))
	}
	b.WriteString("\n")
}

// generateWrappers generates wrapper functions that call acorn commands.
func (g *Generator) generateWrappers(b *strings.Builder, wrappers []componentconfig.Wrapper) {
	for _, w := range wrappers {
		g.generateWrapper(b, w)
	}
}

// generateWrapper generates a single wrapper function.
func (g *Generator) generateWrapper(b *strings.Builder, w componentconfig.Wrapper) {
	b.WriteString(fmt.Sprintf("# %s\n", w.Name))
	b.WriteString(fmt.Sprintf("%s() {\n", w.Name))

	// Add usage check if requires arg
	if w.RequiresArg || w.Usage != "" {
		b.WriteString("    if [ -z \"$1\" ]; then\n")
		if w.Usage != "" {
			b.WriteString(fmt.Sprintf("        echo \"Usage: %s\"\n", w.Usage))
		} else {
			b.WriteString(fmt.Sprintf("        echo \"Usage: %s <arg>\"\n", w.Name))
		}
		b.WriteString("        return 1\n")
		b.WriteString("    fi\n")
	}

	// Build the command
	if w.DefaultArg != "" {
		b.WriteString(fmt.Sprintf("    %s \"${1:-%s}\" \"${@:2}\"\n", w.Command, w.DefaultArg))
	} else {
		b.WriteString(fmt.Sprintf("    %s \"$@\"\n", w.Command))
	}

	// Handle post action
	if w.PostAction == "cd" {
		b.WriteString("    && cd \"$1\"\n")
	}

	b.WriteString("}\n\n")
}

// generateShellFunctions generates raw shell functions.
func (g *Generator) generateShellFunctions(b *strings.Builder, funcs map[string]string) {
	// Sort keys for deterministic output
	keys := make([]string, 0, len(funcs))
	for k := range funcs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		body := funcs[name]
		b.WriteString(fmt.Sprintf("# %s\n", name))
		b.WriteString(fmt.Sprintf("%s() {\n", name))

		// Indent the body
		lines := strings.Split(strings.TrimSpace(body), "\n")
		for _, line := range lines {
			if line == "" {
				b.WriteString("\n")
			} else {
				b.WriteString(fmt.Sprintf("    %s\n", line))
			}
		}

		b.WriteString("}\n\n")
	}
}

// GenerateComponent creates a Component from a BaseConfig.
// This bridges the new config system with the existing shell.Component type.
func (g *Generator) GenerateComponent(cfg *componentconfig.BaseConfig) *Component {
	return &Component{
		Name:        cfg.Name,
		Description: cfg.Description,
		Env:         g.generateEnvString(cfg),
		Aliases:     g.generateAliasesString(cfg.Aliases),
		Functions:   g.generateFunctionsString(cfg),
	}
}

// generateEnvString generates just the env section as a string.
func (g *Generator) generateEnvString(cfg *componentconfig.BaseConfig) string {
	var b strings.Builder

	// Generate env exports
	if len(cfg.Env) > 0 {
		b.WriteString("# " + cfg.Name + " environment setup\n")
		keys := make([]string, 0, len(cfg.Env))
		for k := range cfg.Env {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := cfg.Env[k]
			b.WriteString(fmt.Sprintf("export %s=\"%s\"\n", k, v))
		}
		b.WriteString("\n")
	}

	// Generate PATH additions
	if len(cfg.Paths) > 0 {
		b.WriteString("# Add paths to PATH\n")
		for _, p := range cfg.Paths {
			// Check platform condition
			if p.Condition != "" && p.Condition != g.platform {
				continue
			}

			// Conditional path check (e.g., for Homebrew on macOS)
			if p.Condition != "" {
				b.WriteString(fmt.Sprintf("if [ -d \"%s\" ]; then\n", p.Path))
				b.WriteString(fmt.Sprintf("    case \":$PATH:\" in\n"))
				b.WriteString(fmt.Sprintf("        *\":%s:\"*) ;;\n", p.Path))
				b.WriteString(fmt.Sprintf("        *) export PATH=\"%s:$PATH\" ;;\n", p.Path))
				b.WriteString("    esac\n")
				b.WriteString("fi\n")
			} else {
				b.WriteString(fmt.Sprintf("case \":$PATH:\" in\n"))
				b.WriteString(fmt.Sprintf("    *\":%s:\"*) ;;\n", p.Path))
				b.WriteString(fmt.Sprintf("    *) export PATH=\"%s:$PATH\" ;;\n", p.Path))
				b.WriteString("esac\n")
			}
		}
	}

	return b.String()
}

// generateAliasesString generates just the aliases section as a string.
func (g *Generator) generateAliasesString(aliases map[string]string) string {
	if len(aliases) == 0 {
		return ""
	}

	var b strings.Builder
	keys := make([]string, 0, len(aliases))
	for k := range aliases {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		cmd := aliases[name]
		b.WriteString(fmt.Sprintf("alias %s='%s'\n", name, cmd))
	}
	return b.String()
}

// generateFunctionsString generates wrappers and shell functions as a string.
// Functions prefixed with __ are considered init functions and will be called
// automatically after definition.
func (g *Generator) generateFunctionsString(cfg *componentconfig.BaseConfig) string {
	var b strings.Builder
	var initFunctions []string // Track init functions to call

	// Generate wrappers
	for _, w := range cfg.Wrappers {
		g.generateWrapper(&b, w)
	}

	// Generate shell functions
	funcs := cfg.GetShellFunctions()
	keys := make([]string, 0, len(funcs))
	for k := range funcs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		body := funcs[name]
		b.WriteString(fmt.Sprintf("# %s\n", name))
		b.WriteString(fmt.Sprintf("%s() {\n", name))
		lines := strings.Split(strings.TrimSpace(body), "\n")
		for _, line := range lines {
			if line == "" {
				b.WriteString("\n")
			} else {
				b.WriteString(fmt.Sprintf("    %s\n", line))
			}
		}
		b.WriteString("}\n\n")

		// Track init functions (prefixed with __)
		if strings.HasPrefix(name, "__") {
			initFunctions = append(initFunctions, name)
		}
	}

	// Call init functions in order
	if len(initFunctions) > 0 {
		b.WriteString("# Call init functions\n")
		for _, name := range initFunctions {
			b.WriteString(fmt.Sprintf("%s\n", name))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// GenerateConfigFiles generates all config files for a component.
// Returns the list of generated files and any error.
func (g *Generator) GenerateConfigFiles(cfg *componentconfig.BaseConfig) ([]*configfile.GeneratedFile, error) {
	if len(cfg.Files) == 0 {
		return nil, nil
	}

	manager := configfile.NewManager(g.dryRun)
	return manager.GenerateFiles(cfg.Files)
}
