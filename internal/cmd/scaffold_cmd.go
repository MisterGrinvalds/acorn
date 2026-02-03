package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/components"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

// scaffoldCmd is the root command for scaffold management.
var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Manage scaffold configuration",
	Long: `View and manage the scaffold that controls acorn's router tree.

The scaffold defines which command groups exist, which components
belong to each group, and the shell script generation order.

Source: .sapling/scaffold.yaml (or embedded default)

Examples:
  acorn scaffold show            # Show full scaffold config
  acorn scaffold groups          # List groups and their components
  acorn scaffold components      # List all registered components
  acorn scaffold unmapped        # Show components not in any group
  acorn scaffold shell-order     # Show resolved shell loading order`,
	Aliases: []string{"sc"},
}

// scaffoldShowCmd prints the full resolved scaffold.
var scaffoldShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current scaffold configuration",
	Long: `Display the full scaffold configuration including name, description,
groups, components, and resolved shell order.

Examples:
  acorn scaffold show
  acorn scaffold show -o json
  acorn scaffold show -o yaml`,
	RunE: runScaffoldShow,
}

// scaffoldGroupsCmd lists all groups.
var scaffoldGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List all groups and their components",
	Long: `Display each group defined in the scaffold along with its aliases
and the components it contains.

Examples:
  acorn scaffold groups
  acorn scaffold groups -o json`,
	RunE: runScaffoldGroups,
}

// scaffoldComponentsCmd lists all registered components.
var scaffoldComponentsCmd = &cobra.Command{
	Use:   "components",
	Short: "List all registered components and their group",
	Long: `Display every component registered in the binary, showing which
scaffold group (if any) it belongs to.

Examples:
  acorn scaffold components
  acorn scaffold components -o json`,
	Aliases: []string{"comps"},
	RunE:    runScaffoldComponents,
}

// scaffoldUnmappedCmd shows components not in any group.
var scaffoldUnmappedCmd = &cobra.Command{
	Use:   "unmapped",
	Short: "Show components not in any scaffold group",
	Long: `List components that are compiled into the binary but are not
assigned to any group in the current scaffold. These components
will not appear in the CLI and will not get shell scripts unless
explicitly added to shell_order.optional.

Examples:
  acorn scaffold unmapped
  acorn scaffold unmapped -o json`,
	RunE: runScaffoldUnmapped,
}

// scaffoldShellOrderCmd shows the resolved shell order.
var scaffoldShellOrderCmd = &cobra.Command{
	Use:   "shell-order",
	Short: "Show the resolved shell script loading order",
	Long: `Display the final shell component loading order after resolution.

Shows bootstrap components (always first) and optional components
(either user-specified or auto-derived from groups).

Examples:
  acorn scaffold shell-order
  acorn scaffold shell-order -o json`,
	Aliases: []string{"so"},
	RunE:    runScaffoldShellOrder,
}

func init() {
	rootCmd.AddCommand(scaffoldCmd)
	scaffoldCmd.AddCommand(scaffoldShowCmd)
	scaffoldCmd.AddCommand(scaffoldGroupsCmd)
	scaffoldCmd.AddCommand(scaffoldComponentsCmd)
	scaffoldCmd.AddCommand(scaffoldUnmappedCmd)
	scaffoldCmd.AddCommand(scaffoldShellOrderCmd)
}

func runScaffoldShow(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	s, err := loadScaffold()
	if err != nil {
		return fmt.Errorf("failed to load scaffold: %w", err)
	}
	if s == nil {
		return fmt.Errorf("no scaffold configuration found")
	}

	shellOrder := s.ResolveShellOrder()

	type groupInfo struct {
		Description string   `json:"description" yaml:"description"`
		Aliases     []string `json:"aliases,omitempty" yaml:"aliases,omitempty"`
		Components  []string `json:"components" yaml:"components"`
	}

	groups := make(map[string]groupInfo, len(s.Groups))
	for name, g := range s.Groups {
		groups[name] = groupInfo{
			Description: g.Description,
			Aliases:     g.Aliases,
			Components:  g.Components,
		}
	}

	result := map[string]any{
		"name":        s.Name,
		"description": s.Description,
		"group_order": s.GroupOrder,
		"groups":      groups,
		"shell_order": map[string]any{
			"bootstrap": s.ShellOrder.Bootstrap,
			"optional":  s.ShellOrder.Optional,
			"resolved":  shellOrder,
		},
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(result)
	}

	fmt.Fprintf(os.Stdout, "%s %s\n", output.Info("Scaffold:"), s.Name)
	if s.Description != "" {
		fmt.Fprintf(os.Stdout, "%s %s\n", output.Info("Description:"), s.Description)
	}
	fmt.Fprintf(os.Stdout, "\n%s\n", output.Info("Groups:"))
	for _, gn := range s.GroupOrder {
		g := s.Groups[gn]
		aliases := ""
		if len(g.Aliases) > 0 {
			aliases = fmt.Sprintf(" (%s)", strings.Join(g.Aliases, ", "))
		}
		fmt.Fprintf(os.Stdout, "  %s%s — %s\n", output.Info(gn), aliases, g.Description)
		for _, c := range g.Components {
			fmt.Fprintf(os.Stdout, "    - %s\n", c)
		}
	}

	fmt.Fprintf(os.Stdout, "\n%s\n", output.Info("Shell Order (resolved):"))
	fmt.Fprintf(os.Stdout, "  Bootstrap: %s\n", strings.Join(s.ShellOrder.Bootstrap, ", "))
	if s.ShellOrder.Optional != nil {
		fmt.Fprintf(os.Stdout, "  Optional:  %s (explicit)\n", strings.Join(*s.ShellOrder.Optional, ", "))
	} else {
		fmt.Fprintf(os.Stdout, "  Optional:  (auto-derived from groups)\n")
	}
	return nil
}

func runScaffoldGroups(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	s, err := loadScaffold()
	if err != nil {
		return fmt.Errorf("failed to load scaffold: %w", err)
	}
	if s == nil {
		return fmt.Errorf("no scaffold configuration found")
	}

	type groupEntry struct {
		Name        string   `json:"name" yaml:"name"`
		Description string   `json:"description" yaml:"description"`
		Aliases     []string `json:"aliases,omitempty" yaml:"aliases,omitempty"`
		Components  []string `json:"components" yaml:"components"`
	}

	var entries []groupEntry
	for _, gn := range s.GroupOrder {
		g := s.Groups[gn]
		entries = append(entries, groupEntry{
			Name:        gn,
			Description: g.Description,
			Aliases:     g.Aliases,
			Components:  g.Components,
		})
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"groups": entries})
	}

	for _, e := range entries {
		aliases := ""
		if len(e.Aliases) > 0 {
			aliases = fmt.Sprintf(" (%s)", strings.Join(e.Aliases, ", "))
		}
		fmt.Fprintf(os.Stdout, "%s%s\n", output.Info(e.Name), aliases)
		fmt.Fprintf(os.Stdout, "  %s\n", e.Description)
		for _, c := range e.Components {
			reg := components.Get(c)
			status := output.Success("●")
			if reg == nil {
				status = output.Warning("○") // registered in scaffold but not compiled in
			}
			fmt.Fprintf(os.Stdout, "  %s %s\n", status, c)
		}
		fmt.Fprintln(os.Stdout)
	}
	return nil
}

func runScaffoldComponents(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	s, err := loadScaffold()
	if err != nil {
		return fmt.Errorf("failed to load scaffold: %w", err)
	}
	if s == nil {
		return fmt.Errorf("no scaffold configuration found")
	}

	type compEntry struct {
		Name  string `json:"name" yaml:"name"`
		Group string `json:"group" yaml:"group"`
	}

	all := components.All()
	names := make([]string, 0, len(all))
	for name := range all {
		names = append(names, name)
	}
	sort.Strings(names)

	var entries []compEntry
	for _, name := range names {
		group := s.ComponentGroup(name)
		entries = append(entries, compEntry{Name: name, Group: group})
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"components": entries})
	}

	fmt.Fprintf(os.Stdout, "%-25s %s\n", output.Info("COMPONENT"), output.Info("GROUP"))
	for _, e := range entries {
		group := e.Group
		if group == "" {
			group = output.Warning("(unmapped)")
		} else {
			group = output.Info(group)
		}
		fmt.Fprintf(os.Stdout, "%-25s %s\n", e.Name, group)
	}
	return nil
}

func runScaffoldUnmapped(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	s, err := loadScaffold()
	if err != nil {
		return fmt.Errorf("failed to load scaffold: %w", err)
	}
	if s == nil {
		return fmt.Errorf("no scaffold configuration found")
	}

	// Build set of all grouped components
	grouped := make(map[string]bool)
	for _, g := range s.Groups {
		for _, c := range g.Components {
			grouped[c] = true
		}
	}

	all := components.All()
	var unmapped []string
	for name := range all {
		if !grouped[name] {
			unmapped = append(unmapped, name)
		}
	}
	sort.Strings(unmapped)

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"unmapped": unmapped, "count": len(unmapped)})
	}

	if len(unmapped) == 0 {
		fmt.Fprintf(os.Stdout, "%s All registered components are mapped to a group\n", output.Success("✓"))
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s %d component(s) not in any scaffold group:\n\n",
		output.Warning("!"), len(unmapped))
	for _, name := range unmapped {
		fmt.Fprintf(os.Stdout, "  %s %s\n", output.Warning("○"), name)
	}
	fmt.Fprintf(os.Stdout, "\nThese components are compiled in but won't appear in the CLI\n")
	fmt.Fprintf(os.Stdout, "or get shell scripts unless added to a group or shell_order.optional.\n")
	return nil
}

func runScaffoldShellOrder(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	s, err := loadScaffold()
	if err != nil {
		return fmt.Errorf("failed to load scaffold: %w", err)
	}
	if s == nil {
		return fmt.Errorf("no scaffold configuration found")
	}

	resolved := s.ResolveShellOrder()
	autoDerive := s.ShellOrder.Optional == nil

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{
			"bootstrap":    s.ShellOrder.Bootstrap,
			"optional":     s.ShellOrder.Optional,
			"auto_derived": autoDerive,
			"resolved":     resolved,
		})
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Bootstrap (always first):"))
	for i, name := range s.ShellOrder.Bootstrap {
		fmt.Fprintf(os.Stdout, "  %2d. %s\n", i+1, output.Info(name))
	}

	fmt.Fprintln(os.Stdout)
	if autoDerive {
		fmt.Fprintf(os.Stdout, "%s (auto-derived from groups)\n", output.Info("Optional:"))
	} else {
		fmt.Fprintf(os.Stdout, "%s (explicitly configured)\n", output.Info("Optional:"))
	}

	offset := len(s.ShellOrder.Bootstrap)
	for i, name := range resolved[offset:] {
		fmt.Fprintf(os.Stdout, "  %2d. %s\n", offset+i+1, name)
	}

	fmt.Fprintf(os.Stdout, "\n  %s %d total components\n", output.Info("→"), len(resolved))
	return nil
}
