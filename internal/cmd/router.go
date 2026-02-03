package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components"
	"github.com/mistergrinvalds/acorn/internal/scaffold"
	"github.com/mistergrinvalds/acorn/internal/utils/config"
	"github.com/spf13/cobra"
)

// buildRouter loads the scaffold and dynamically builds the Cobra command tree.
// Called during initialization to wire components into group commands.
func buildRouter() {
	s, err := loadScaffold()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load scaffold: %v\n", err)
		return
	}
	if s == nil {
		return
	}

	buildRouterFromScaffold(rootCmd, s)
}

// loadScaffold tries to load scaffold from .sapling first, then falls back
// to the embedded default.
func loadScaffold() (*config.Scaffold, error) {
	// Try .sapling/scaffold.yaml first
	s, err := config.LoadScaffold()
	if err != nil {
		return nil, err
	}
	if s != nil {
		return s, nil
	}

	// Fall back to embedded default
	return scaffold.LoadDefault()
}

// buildRouterFromScaffold creates group commands and wires registered
// components into them based on the scaffold configuration.
func buildRouterFromScaffold(root *cobra.Command, s *config.Scaffold) {
	for _, groupName := range s.GroupOrder {
		group, ok := s.Groups[groupName]
		if !ok {
			continue
		}

		groupCmd := &cobra.Command{
			Use:   groupName,
			Short: group.Description,
			Long:  fmt.Sprintf("Commands for %s.", group.Description),
		}
		if len(group.Aliases) > 0 {
			groupCmd.Aliases = group.Aliases
		}

		hasComponents := false
		for _, compName := range group.Components {
			reg := components.Get(compName)
			if reg != nil {
				groupCmd.AddCommand(reg.RegisterCmd())
				hasComponents = true
			}
		}

		// Only add the group if it has at least one component
		if hasComponents {
			root.AddCommand(groupCmd)
		}
	}
}
