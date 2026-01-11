package installer

import (
	"testing"

	"github.com/mistergrinvalds/acorn/internal/componentconfig"
)

func TestGetExecutor(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		wantErr bool
	}{
		{
			name:    "brew",
			method:  "brew",
			wantErr: false,
		},
		{
			name:    "apt",
			method:  "apt",
			wantErr: false,
		},
		{
			name:    "npm",
			method:  "npm",
			wantErr: false,
		},
		{
			name:    "go",
			method:  "go",
			wantErr: false,
		},
		{
			name:    "curl",
			method:  "curl",
			wantErr: false,
		},
		{
			name:    "unsupported",
			method:  "unsupported",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor, err := GetExecutor(tt.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExecutor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if executor == nil {
					t.Error("Expected non-nil executor")
				}
				if executor.Type() != tt.method {
					t.Errorf("Type() = %s, want %s", executor.Type(), tt.method)
				}
			}
		})
	}
}

func TestBrewExecutor(t *testing.T) {
	e := &BrewExecutor{}

	if e.Type() != "brew" {
		t.Errorf("Type() = %s, want brew", e.Type())
	}

	// Test availability check
	_ = e.Available()
}

func TestAptExecutor(t *testing.T) {
	e := &AptExecutor{}

	if e.Type() != "apt" {
		t.Errorf("Type() = %s, want apt", e.Type())
	}

	// Test availability check
	_ = e.Available()
}

func TestNpmExecutor(t *testing.T) {
	e := &NpmExecutor{}

	if e.Type() != "npm" {
		t.Errorf("Type() = %s, want npm", e.Type())
	}

	// Test availability check
	_ = e.Available()
}

func TestGoExecutor(t *testing.T) {
	e := &GoExecutor{}

	if e.Type() != "go" {
		t.Errorf("Type() = %s, want go", e.Type())
	}

	// Test availability check
	_ = e.Available()
}

func TestCurlExecutor(t *testing.T) {
	e := &CurlExecutor{}

	if e.Type() != "curl" {
		t.Errorf("Type() = %s, want curl", e.Type())
	}

	// Test availability check
	_ = e.Available()
}

func TestExecutorAvailability(t *testing.T) {
	executors := []MethodExecutor{
		&BrewExecutor{},
		&AptExecutor{},
		&NpmExecutor{},
		&GoExecutor{},
		&CurlExecutor{},
	}

	for _, e := range executors {
		t.Run(e.Type(), func(t *testing.T) {
			// Just verify Available() doesn't panic
			_ = e.Available()
		})
	}
}

func TestSelectMethodWithPlatform(t *testing.T) {
	methods := map[string]componentconfig.InstallMethod{
		"darwin":       {Type: "brew", Package: "node"},
		"linux/ubuntu": {Type: "apt", Package: "nodejs"},
		"linux":        {Type: "brew", Package: "node"},
	}

	tests := []struct {
		name     string
		platform *Platform
		wantType string
		wantOk   bool
	}{
		{
			name:     "darwin prefers brew",
			platform: &Platform{OS: "darwin", Arch: "arm64"},
			wantType: "brew",
			wantOk:   true,
		},
		{
			name:     "ubuntu prefers apt",
			platform: &Platform{OS: "linux", Distro: "ubuntu", Arch: "amd64"},
			wantType: "apt",
			wantOk:   true,
		},
		{
			name:     "fedora falls back to linux",
			platform: &Platform{OS: "linux", Distro: "fedora", Arch: "amd64"},
			wantType: "brew",
			wantOk:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method, ok := SelectMethod(methods, tt.platform)
			if ok != tt.wantOk {
				t.Errorf("ok = %v, want %v", ok, tt.wantOk)
			}
			if tt.wantOk && method.Type != tt.wantType {
				t.Errorf("type = %s, want %s", method.Type, tt.wantType)
			}
		})
	}
}
