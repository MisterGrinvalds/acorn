package installer

import (
	"runtime"
	"testing"
)

func TestDetectPlatform(t *testing.T) {
	p := DetectPlatform()

	if p.OS == "" {
		t.Error("Expected OS to be detected")
	}
	if p.Arch == "" {
		t.Error("Expected architecture to be detected")
	}

	// Verify OS matches runtime
	if p.OS != runtime.GOOS {
		t.Errorf("Expected OS %s, got %s", runtime.GOOS, p.OS)
	}
	if p.Arch != runtime.GOARCH {
		t.Errorf("Expected arch %s, got %s", runtime.GOARCH, p.Arch)
	}
}

func TestGetMethodKeys(t *testing.T) {
	tests := []struct {
		name     string
		platform Platform
		expected []string
	}{
		{
			name: "darwin",
			platform: Platform{
				OS:   "darwin",
				Arch: "arm64",
			},
			expected: []string{"darwin"},
		},
		{
			name: "linux with distro",
			platform: Platform{
				OS:     "linux",
				Arch:   "amd64",
				Distro: "ubuntu",
			},
			expected: []string{"linux/ubuntu", "linux"},
		},
		{
			name: "linux with distro family",
			platform: Platform{
				OS:           "linux",
				Arch:         "amd64",
				Distro:       "ubuntu",
				DistroFamily: "debian",
			},
			expected: []string{
				"linux/ubuntu",
				"linux/debian",
				"linux",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys := tt.platform.GetMethodKeys()

			if len(keys) != len(tt.expected) {
				t.Errorf("Expected %d keys, got %d", len(tt.expected), len(keys))
				t.Errorf("Expected: %v", tt.expected)
				t.Errorf("Got: %v", keys)
				return
			}

			for i, key := range keys {
				if key != tt.expected[i] {
					t.Errorf("Key %d: expected %s, got %s", i, tt.expected[i], key)
				}
			}
		})
	}
}

func TestLinuxDistroDetection(t *testing.T) {
	// This test verifies the structure of Linux distro detection
	// Actual values depend on the system running the test
	p := DetectPlatform()

	if p.OS == "linux" {
		// On Linux, we expect either distro or package manager to be detected
		if p.Distro == "" && p.PackageManager == "" {
			t.Error("Expected either distro or package manager to be detected on Linux")
		}

		// If distro is detected, family should be too
		if p.Distro != "" && p.DistroFamily == "" {
			// Note: Not all distros have families, so this isn't always an error
			t.Logf("Distro %s detected without family", p.Distro)
		}
	}
}

func TestPlatformString(t *testing.T) {
	tests := []struct {
		name     string
		platform Platform
		want     string
	}{
		{
			name: "darwin with package manager",
			platform: Platform{
				OS:             "darwin",
				Arch:           "arm64",
				PackageManager: "brew",
			},
			want: "darwin arm64 (brew)",
		},
		{
			name: "linux without package manager",
			platform: Platform{
				OS:   "linux",
				Arch: "amd64",
			},
			want: "linux amd64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.platform.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandExists(t *testing.T) {
	// Test with a command that should exist on all systems
	if !commandExists("sh") {
		t.Error("Expected 'sh' command to exist")
	}

	// Test with a command that shouldn't exist
	if commandExists("definitely-not-a-real-command-12345") {
		t.Error("Expected fake command to not exist")
	}
}
