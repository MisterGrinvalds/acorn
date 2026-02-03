// Package scaffold provides the embedded default scaffold configuration.
package scaffold

import (
	_ "embed"

	"github.com/mistergrinvalds/acorn/internal/utils/config"
)

//go:embed default.yaml
var defaultScaffoldYAML []byte

// LoadDefault returns the default scaffold embedded in the binary.
func LoadDefault() (*config.Scaffold, error) {
	return config.ParseScaffoldBytes(defaultScaffoldYAML)
}
