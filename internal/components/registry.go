// Package components provides the component registry for acorn.
// Components self-register via init() and the scaffold determines
// which components are active and how they're grouped in the CLI.
package components

import (
	"sync"

	"github.com/spf13/cobra"
)

// Registration holds the metadata and command factory for a component.
type Registration struct {
	// Name is the unique identifier for this component (e.g., "docker").
	Name string

	// Description is a short description shown in CLI help.
	Description string

	// Aliases are alternative names for the command (e.g., ["dk", "d"]).
	Aliases []string

	// RegisterCmd returns the root cobra.Command for this component.
	// Called once when building the router tree from the scaffold.
	RegisterCmd func() *cobra.Command
}

var (
	mu       sync.RWMutex
	registry = make(map[string]*Registration)
)

// Register adds a component to the global registry.
// Typically called from init() in each component's cmd file.
func Register(r *Registration) {
	mu.Lock()
	defer mu.Unlock()
	registry[r.Name] = r
}

// Get returns a component registration by name, or nil if not found.
func Get(name string) *Registration {
	mu.RLock()
	defer mu.RUnlock()
	return registry[name]
}

// All returns a copy of all registered components.
func All() map[string]*Registration {
	mu.RLock()
	defer mu.RUnlock()
	result := make(map[string]*Registration, len(registry))
	for k, v := range registry {
		result[k] = v
	}
	return result
}

// Names returns the names of all registered components.
func Names() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
