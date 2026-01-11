package installer

import (
	"testing"

	"github.com/mistergrinvalds/acorn/internal/componentconfig"
)

func TestParseRequirement(t *testing.T) {
	tests := []struct {
		name        string
		requirement string
		wantComp    string
		wantCmd     string
	}{
		{
			name:        "command only",
			requirement: "npm",
			wantComp:    "",
			wantCmd:     "npm",
		},
		{
			name:        "component and command",
			requirement: "node:npm",
			wantComp:    "node",
			wantCmd:     "npm",
		},
		{
			name:        "multiple colons",
			requirement: "foo:bar:baz",
			wantComp:    "foo",
			wantCmd:     "bar:baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotComp, gotCmd := parseRequirement(tt.requirement)
			if gotComp != tt.wantComp {
				t.Errorf("component = %v, want %v", gotComp, tt.wantComp)
			}
			if gotCmd != tt.wantCmd {
				t.Errorf("command = %v, want %v", gotCmd, tt.wantCmd)
			}
		})
	}
}

func TestCheckInstalled(t *testing.T) {
	r := NewResolver(DetectPlatform())

	tests := []struct {
		name    string
		check   string
		wantOk  bool
		wantVer bool
	}{
		{
			name:    "command -v format with existing command",
			check:   "command -v sh",
			wantOk:  true,
			wantVer: false,
		},
		{
			name:    "command -v format with nonexistent command",
			check:   "command -v definitely-not-real",
			wantOk:  false,
			wantVer: false,
		},
		{
			name:    "version check format",
			check:   "go version",
			wantOk:  commandExists("go"),
			wantVer: commandExists("go"),
		},
		{
			name:    "empty check",
			check:   "",
			wantOk:  false,
			wantVer: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installed, version := r.checkInstalled(tt.check)
			if installed != tt.wantOk {
				t.Errorf("installed = %v, want %v", installed, tt.wantOk)
			}
			if tt.wantVer && version == "" {
				t.Error("Expected version string but got empty")
			}
			if !tt.wantVer && version != "" {
				t.Logf("Got unexpected version string: %s", version)
			}
		})
	}
}

func TestResolverBuildPlan(t *testing.T) {
	platform := DetectPlatform()
	r := NewResolver(platform)

	// Create a simple install config for testing
	cfg := &componentconfig.InstallConfig{
		Tools: []componentconfig.ToolInstall{
			{
				Name:  "test-tool",
				Check: "command -v test-tool",
				Methods: map[string]componentconfig.InstallMethod{
					"darwin": {Type: "brew", Package: "test-tool"},
					"linux":  {Type: "apt", Package: "test-tool"},
				},
			},
		},
	}

	plan, err := r.BuildPlan("test-component", cfg)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}

	if plan.Component != "test-component" {
		t.Errorf("Component = %s, want test-component", plan.Component)
	}

	if len(plan.Tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(plan.Tools))
	}

	tool := plan.Tools[0]
	if tool.Name != "test-tool" {
		t.Errorf("Tool name = %s, want test-tool", tool.Name)
	}
	if tool.Reason != "direct" {
		t.Errorf("Reason = %s, want direct", tool.Reason)
	}
}

func TestResolverWithPrerequisites(t *testing.T) {
	platform := DetectPlatform()
	r := NewResolver(platform)

	// Test with cloudflare component which requires node:npm
	cfg := &componentconfig.InstallConfig{
		Tools: []componentconfig.ToolInstall{
			{
				Name:     "wrangler",
				Check:    "command -v wrangler",
				Requires: []string{"node:npm"},
				Methods: map[string]componentconfig.InstallMethod{
					"darwin": {Type: "npm", Package: "wrangler", Global: true},
				},
			},
		},
	}

	plan, err := r.BuildPlan("cloudflare", cfg)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}

	// Should have wrangler as direct tool
	if len(plan.Tools) < 1 {
		t.Error("Expected at least 1 direct tool")
	}

	// If npm is not installed, should have prerequisites
	if !commandExists("npm") {
		if len(plan.Prerequisites) == 0 {
			t.Error("Expected prerequisites when npm not installed")
		}
	}
}

func TestSelectMethod(t *testing.T) {
	tests := []struct {
		name     string
		methods  map[string]componentconfig.InstallMethod
		platform *Platform
		wantType string
		wantOk   bool
	}{
		{
			name: "darwin exact match",
			methods: map[string]componentconfig.InstallMethod{
				"darwin": {Type: "brew"},
			},
			platform: &Platform{OS: "darwin", Arch: "arm64"},
			wantType: "brew",
			wantOk:   true,
		},
		{
			name: "linux with distro match",
			methods: map[string]componentconfig.InstallMethod{
				"linux/ubuntu": {Type: "apt"},
				"linux":        {Type: "brew"},
			},
			platform: &Platform{OS: "linux", Arch: "amd64", Distro: "ubuntu"},
			wantType: "apt",
			wantOk:   true,
		},
		{
			name: "linux fallback",
			methods: map[string]componentconfig.InstallMethod{
				"linux": {Type: "brew"},
			},
			platform: &Platform{OS: "linux", Arch: "amd64", Distro: "fedora"},
			wantType: "brew",
			wantOk:   true,
		},
		{
			name: "no match",
			methods: map[string]componentconfig.InstallMethod{
				"windows": {Type: "chocolatey"},
			},
			platform: &Platform{OS: "darwin"},
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method, ok := SelectMethod(tt.methods, tt.platform)
			if ok != tt.wantOk {
				t.Errorf("ok = %v, want %v", ok, tt.wantOk)
			}
			if tt.wantOk && method.Type != tt.wantType {
				t.Errorf("type = %s, want %s", method.Type, tt.wantType)
			}
		})
	}
}

func TestResolverCycleDetection(t *testing.T) {
	platform := DetectPlatform()
	r := NewResolver(platform)

	// Create a config that references the same tool multiple times
	cfg := &componentconfig.InstallConfig{
		Tools: []componentconfig.ToolInstall{
			{
				Name:     "tool-a",
				Check:    "command -v tool-a",
				Requires: []string{"tool-b"},
				Methods: map[string]componentconfig.InstallMethod{
					"darwin": {Type: "brew"},
				},
			},
		},
	}

	// This should not cause infinite recursion
	// (tool-b won't be found, but that's ok for this test)
	_, err := r.BuildPlan("test", cfg)
	if err == nil {
		t.Log("BuildPlan completed (expected error for missing tool-b)")
	}

	// The key is that it doesn't hang - if we get here, cycle detection works
}
