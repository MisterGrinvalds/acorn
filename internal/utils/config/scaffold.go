package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Scaffold defines the complete router tree configuration.
// It determines which groups exist, their metadata, and which
// components belong to each group.
type Scaffold struct {
	// Name identifies this scaffold configuration.
	Name string `yaml:"name"`

	// Description is a human-readable description of this scaffold.
	Description string `yaml:"description,omitempty"`

	// Groups defines the command groups and their components.
	// Map key is the group command name (e.g., "devops").
	// Order is preserved via GroupOrder.
	Groups map[string]ScaffoldGroup `yaml:"groups"`

	// GroupOrder defines the order groups appear in help output.
	// Derived from YAML key order during unmarshaling.
	GroupOrder []string `yaml:"-"`

	// ShellOrder controls shell script generation ordering.
	// Bootstrap components always load first. Optional components follow.
	// If Optional is nil (omitted from YAML), all grouped components are
	// auto-included in group declaration order.
	ShellOrder ShellOrder `yaml:"shell_order"`
}

// ShellOrder defines the two-tier shell script loading order.
type ShellOrder struct {
	// Bootstrap lists components that MUST load first (shell detection,
	// XDG setup, theme, core). These are always included and always
	// take precedence over optional ordering.
	Bootstrap []string `yaml:"bootstrap"`

	// Optional lists additional components for shell script generation.
	// If nil (key omitted from YAML), all components found in scaffold
	// groups are auto-included in group declaration order.
	// If empty list (optional: []), no optional shell scripts are generated.
	// If populated, only the listed components are included, in this order.
	Optional *[]string `yaml:"optional,omitempty"`
}

// ScaffoldGroup defines a command group in the CLI router tree.
type ScaffoldGroup struct {
	// Description is shown in CLI help for this group.
	Description string `yaml:"description"`

	// Aliases are alternative names for the group command.
	Aliases []string `yaml:"aliases,omitempty"`

	// Components lists the component names to include in this group.
	Components []string `yaml:"components"`
}

// LoadScaffold loads the scaffold configuration.
// It searches for .sapling/scaffold.yaml using the same upward
// search as other sapling config. Returns nil if no scaffold is found.
func LoadScaffold() (*Scaffold, error) {
	path, err := findScaffoldFile()
	if err != nil {
		return nil, nil // No scaffold found, not an error
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read scaffold file %s: %w", path, err)
	}

	scaffold, err := parseScaffold(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse scaffold file %s: %w", path, err)
	}

	return scaffold, nil
}

// findScaffoldFile locates the scaffold.yaml file.
// Searches: $SAPLING_DIR/scaffold.yaml, then walks up from CWD looking for
// .sapling/scaffold.yaml.
func findScaffoldFile() (string, error) {
	if dir := os.Getenv("SAPLING_DIR"); dir != "" {
		path := filepath.Join(dir, "scaffold.yaml")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	dir := cwd
	for {
		candidate := filepath.Join(dir, ".sapling", "scaffold.yaml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("scaffold.yaml not found")
}

// parseScaffold parses scaffold YAML while preserving group order.
func parseScaffold(data []byte) (*Scaffold, error) {
	scaffold := &Scaffold{}
	if err := yaml.Unmarshal(data, scaffold); err != nil {
		return nil, err
	}

	// Extract group order from YAML by parsing the raw node tree.
	// yaml.v3 preserves key order in mapping nodes.
	var raw yaml.Node
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	scaffold.GroupOrder = extractGroupOrder(&raw)

	return scaffold, nil
}

// extractGroupOrder pulls group key order from the YAML AST.
func extractGroupOrder(root *yaml.Node) []string {
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil
	}

	doc := root.Content[0]
	if doc.Kind != yaml.MappingNode {
		return nil
	}

	// Find the "groups" key
	for i := 0; i < len(doc.Content)-1; i += 2 {
		if doc.Content[i].Value == "groups" {
			groupsNode := doc.Content[i+1]
			if groupsNode.Kind != yaml.MappingNode {
				return nil
			}
			var order []string
			for j := 0; j < len(groupsNode.Content)-1; j += 2 {
				order = append(order, groupsNode.Content[j].Value)
			}
			return order
		}
	}

	return nil
}

// ParseScaffoldBytes parses scaffold YAML from raw bytes.
// Exported so the embedded default scaffold can use it.
func ParseScaffoldBytes(data []byte) (*Scaffold, error) {
	return parseScaffold(data)
}

// AllComponents returns a flat list of all component names from the scaffold,
// in group order then component order within each group.
func (s *Scaffold) AllComponents() []string {
	var all []string
	seen := make(map[string]bool)

	for _, groupName := range s.GroupOrder {
		group, ok := s.Groups[groupName]
		if !ok {
			continue
		}
		for _, comp := range group.Components {
			if !seen[comp] {
				all = append(all, comp)
				seen[comp] = true
			}
		}
	}

	return all
}

// ResolveShellOrder computes the final shell component loading order.
//
// The resolution strategy is:
//  1. shell_order.bootstrap always comes first, exactly as specified.
//  2. If shell_order.optional is explicitly set, those components follow
//     in the specified order (user takes full control).
//  3. If shell_order.optional is omitted (nil), all components from the
//     scaffold groups are auto-included in group declaration order,
//     skipping anything already in bootstrap.
func (s *Scaffold) ResolveShellOrder() []string {
	// Start with bootstrap — always first
	bootstrapSet := make(map[string]bool, len(s.ShellOrder.Bootstrap))
	var result []string
	for _, name := range s.ShellOrder.Bootstrap {
		result = append(result, name)
		bootstrapSet[name] = true
	}

	if s.ShellOrder.Optional != nil {
		// User explicitly specified optional components — use as-is
		for _, name := range *s.ShellOrder.Optional {
			if !bootstrapSet[name] {
				result = append(result, name)
			}
		}
	} else {
		// Auto-derive from groups in declaration order
		for _, groupName := range s.GroupOrder {
			group, ok := s.Groups[groupName]
			if !ok {
				continue
			}
			for _, comp := range group.Components {
				if !bootstrapSet[comp] {
					result = append(result, comp)
					bootstrapSet[comp] = true // dedup across groups
				}
			}
		}
	}

	return result
}

// ComponentGroup returns the group name a component belongs to,
// or empty string if not found.
func (s *Scaffold) ComponentGroup(component string) string {
	for groupName, group := range s.Groups {
		for _, comp := range group.Components {
			if comp == component {
				return groupName
			}
		}
	}
	return ""
}
