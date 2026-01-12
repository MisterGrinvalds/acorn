package configfile

import (
	"strings"
	"testing"
)

func TestTextListWriterFormat(t *testing.T) {
	w := &TextListWriter{}
	if w.Format() != "textlist" {
		t.Errorf("Format() = %q, want %q", w.Format(), "textlist")
	}
}

func TestTextListWriterSimpleItems(t *testing.T) {
	w := &TextListWriter{}

	values := map[string]any{
		"items": []any{
			"item-one",
			"item-two",
			"item-three",
		},
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	output := string(content)
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d: %s", len(lines), output)
	}

	if lines[0] != "item-one" {
		t.Errorf("Line 1 = %q, want %q", lines[0], "item-one")
	}
}

func TestTextListWriterItemsWithComments(t *testing.T) {
	w := &TextListWriter{}

	values := map[string]any{
		"items": []any{
			map[string]any{
				"id":      "catppuccin-theme",
				"comment": "Catppuccin color scheme",
			},
			map[string]any{
				"id":      "org.rust.lang",
				"comment": "Rust",
			},
		},
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	output := string(content)

	if !strings.Contains(output, "catppuccin-theme") {
		t.Error("Missing catppuccin-theme")
	}

	if !strings.Contains(output, "# Catppuccin color scheme") {
		t.Error("Missing comment for catppuccin-theme")
	}

	if !strings.Contains(output, "org.rust.lang") {
		t.Error("Missing org.rust.lang")
	}
}

func TestTextListWriterWithHeader(t *testing.T) {
	w := &TextListWriter{}

	values := map[string]any{
		"header": "# Plugin List\n# For IntelliJ IDEA\n",
		"items": []any{
			"plugin-one",
		},
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	output := string(content)

	if !strings.HasPrefix(output, "# Plugin List") {
		t.Errorf("Should start with header, got: %s", output)
	}

	if !strings.Contains(output, "# For IntelliJ IDEA") {
		t.Error("Missing second header line")
	}

	if !strings.Contains(output, "plugin-one") {
		t.Error("Missing plugin item")
	}
}

func TestTextListWriterWithSections(t *testing.T) {
	w := &TextListWriter{}

	values := map[string]any{
		"header": "# My Plugins\n",
		"sections": []any{
			map[string]any{
				"name": "Theme",
				"items": []any{
					map[string]any{
						"id":      "catppuccin-theme",
						"comment": "Color scheme",
					},
				},
			},
			map[string]any{
				"name": "Languages",
				"items": []any{
					map[string]any{
						"id":      "org.rust.lang",
						"comment": "Rust",
					},
					map[string]any{
						"id":      "org.jetbrains.plugins.go",
						"comment": "Go",
					},
				},
			},
		},
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	output := string(content)

	if !strings.Contains(output, "# My Plugins") {
		t.Error("Missing header")
	}

	if !strings.Contains(output, "# Theme") {
		t.Error("Missing Theme section")
	}

	if !strings.Contains(output, "# Languages") {
		t.Error("Missing Languages section")
	}

	if !strings.Contains(output, "catppuccin-theme") {
		t.Error("Missing catppuccin-theme")
	}

	if !strings.Contains(output, "org.rust.lang") {
		t.Error("Missing org.rust.lang")
	}
}

func TestTextListWriterEmptyValues(t *testing.T) {
	w := &TextListWriter{}

	values := map[string]any{}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("Expected empty output, got: %s", string(content))
	}
}

func TestTextListWriterItemWithoutComment(t *testing.T) {
	w := &TextListWriter{}

	values := map[string]any{
		"items": []any{
			map[string]any{
				"id": "simple-plugin",
			},
		},
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	output := string(content)

	if !strings.Contains(output, "simple-plugin") {
		t.Error("Missing plugin id")
	}

	if strings.Contains(output, "#") {
		t.Error("Should not have comment marker for item without comment")
	}
}
