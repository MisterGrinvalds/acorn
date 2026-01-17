// Package version provides build-time version information for the acorn CLI.
package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// Build-time variables set via ldflags
var (
	// Version is the semantic version (e.g., v1.0.0)
	Version = "dev"
	// Commit is the git commit SHA
	Commit = "unknown"
	// Date is the build date in RFC3339 format
	Date = "unknown"
	// BuiltBy indicates who/what built this binary
	BuiltBy = "unknown"
)

// Info holds complete version information
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	BuiltBy   string `json:"built_by"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// Get returns the current version information
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		BuiltBy:   BuiltBy,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// String returns a human-readable version string
func (i Info) String() string {
	return fmt.Sprintf("acorn %s (%s) built %s by %s\n%s %s/%s",
		i.Version, i.Commit[:min(7, len(i.Commit))], i.Date, i.BuiltBy,
		i.GoVersion, i.OS, i.Arch)
}

// Short returns just the version number
func (i Info) Short() string {
	return i.Version
}

// GetModuleVersion attempts to get version from module info (for go install)
func GetModuleVersion() string {
	if Version != "dev" {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}
	return Version
}
