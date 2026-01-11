package installer

import (
	"bytes"
	"context"
	"testing"

	"github.com/mistergrinvalds/acorn/internal/componentconfig"
)

func TestNewInstaller(t *testing.T) {
	i := NewInstaller()

	if i.loader == nil {
		t.Error("Expected loader to be initialized")
	}
	if i.platform == nil {
		t.Error("Expected platform to be initialized")
	}
	if i.stdout == nil {
		t.Error("Expected stdout to be initialized")
	}
}

func TestInstallerOptions(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	i := NewInstaller(
		WithDryRun(true),
		WithVerbose(true),
		WithOutput(stdout, stderr),
	)

	if !i.dryRun {
		t.Error("Expected dry-run to be enabled")
	}
	if !i.verbose {
		t.Error("Expected verbose to be enabled")
	}
	if i.stdout != stdout {
		t.Error("Expected custom stdout")
	}
	if i.stderr != stderr {
		t.Error("Expected custom stderr")
	}
}

func TestPlan(t *testing.T) {
	tests := []struct {
		name      string
		component string
		wantErr   bool
	}{
		{
			name:      "valid component with install config",
			component: "cloudflare",
			wantErr:   false,
		},
		{
			name:      "nonexistent component",
			component: "nonexistent",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewInstaller()
			plan, err := i.Plan(context.Background(), tt.component)

			if (err != nil) != tt.wantErr {
				t.Errorf("Plan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if plan == nil {
					t.Error("Expected plan to be non-nil")
				}
				if plan.Component != tt.component {
					t.Errorf("Expected component %s, got %s", tt.component, plan.Component)
				}
			}
		})
	}
}

func TestInstallDryRun(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	i := NewInstaller(
		WithDryRun(true),
		WithOutput(stdout, stderr),
	)

	result, err := i.Install(context.Background(), "cloudflare")
	if err != nil {
		t.Fatalf("Install() error = %v", err)
	}

	if !result.Success {
		t.Error("Expected dry-run install to succeed")
	}
	if !result.DryRun {
		t.Error("Expected result to indicate dry-run")
	}

	// In dry-run mode, all tools should be skipped
	for _, tr := range result.Tools {
		if !tr.Skipped {
			t.Errorf("Tool %s should be skipped in dry-run", tr.Name)
		}
	}
}

func TestLoadInstallConfig(t *testing.T) {
	i := NewInstaller()

	cfg, err := i.loadInstallConfig("cloudflare")
	if err != nil {
		t.Fatalf("loadInstallConfig() error = %v", err)
	}

	if len(cfg.Tools) == 0 {
		t.Error("Expected cloudflare to have install tools")
	}

	// Verify wrangler tool exists
	var found bool
	for _, tool := range cfg.Tools {
		if tool.Name == "wrangler" {
			found = true
			if tool.Check == "" {
				t.Error("Expected wrangler to have check command")
			}
			if len(tool.Methods) == 0 {
				t.Error("Expected wrangler to have install methods")
			}
			break
		}
	}
	if !found {
		t.Error("Expected to find wrangler tool in cloudflare config")
	}
}

func TestInstallTool_AlreadyInstalled(t *testing.T) {
	i := NewInstaller()

	tool := PlannedTool{
		Name:             "test-tool",
		AlreadyInstalled: true,
		Version:          "1.0.0",
	}

	result := i.installTool(context.Background(), tool)

	if !result.Success {
		t.Error("Expected success for already installed tool")
	}
	if !result.Skipped {
		t.Error("Expected tool to be skipped")
	}
	if result.SkipReason != "already installed" {
		t.Errorf("Expected skip reason 'already installed', got %s", result.SkipReason)
	}
	if result.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", result.Version)
	}
}

func TestInstallTool_DryRun(t *testing.T) {
	i := NewInstaller(WithDryRun(true))

	tool := PlannedTool{
		Name: "test-tool",
		Method: componentconfig.InstallMethod{
			Type:    "npm",
			Package: "test-pkg",
			Global:  true,
		},
	}

	result := i.installTool(context.Background(), tool)

	if !result.Success {
		t.Error("Expected success for dry-run")
	}
	if !result.Skipped {
		t.Error("Expected tool to be skipped in dry-run")
	}
	if result.SkipReason != "dry run" {
		t.Errorf("Expected skip reason 'dry run', got %s", result.SkipReason)
	}
}

func TestInstallTool_UnsupportedMethod(t *testing.T) {
	i := NewInstaller()

	tool := PlannedTool{
		Name: "test-tool",
		Method: componentconfig.InstallMethod{
			Type: "unsupported-method",
		},
	}

	result := i.installTool(context.Background(), tool)

	if result.Success {
		t.Error("Expected failure for unsupported method")
	}
	if result.Error == nil {
		t.Error("Expected error for unsupported method")
	}
}
