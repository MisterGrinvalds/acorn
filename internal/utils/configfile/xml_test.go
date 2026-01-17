package configfile

import (
	"strings"
	"testing"
)

func TestXMLWriterFormat(t *testing.T) {
	w := &XMLWriter{}
	if w.Format() != "xml" {
		t.Errorf("Format() = %q, want %q", w.Format(), "xml")
	}
}

func TestXMLWriterSimpleElement(t *testing.T) {
	w := &XMLWriter{}

	values := map[string]interface{}{
		"root": "config",
		"attrs": map[string]interface{}{
			"version": "1",
		},
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	output := string(content)

	if !strings.Contains(output, "<?xml version=") {
		t.Error("Missing XML declaration")
	}

	if !strings.Contains(output, `<config version="1"/>`) {
		t.Errorf("Expected self-closing element, got:\n%s", output)
	}
}

func TestXMLWriterWithChildren(t *testing.T) {
	w := &XMLWriter{}

	values := map[string]interface{}{
		"root": "keymap",
		"attrs": map[string]interface{}{
			"name":    "Custom",
			"version": "1",
		},
		"children": []interface{}{
			map[string]interface{}{
				"element": "action",
				"attrs": map[string]interface{}{
					"id": "GotoFile",
				},
				"children": []interface{}{
					map[string]interface{}{
						"element": "keyboard-shortcut",
						"attrs": map[string]interface{}{
							"first-keystroke": "meta P",
						},
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

	if !strings.Contains(output, "<keymap") {
		t.Error("Missing keymap element")
	}

	if !strings.Contains(output, `name="Custom"`) {
		t.Error("Missing name attribute")
	}

	if !strings.Contains(output, `<action id="GotoFile">`) {
		t.Error("Missing action element")
	}

	if !strings.Contains(output, `<keyboard-shortcut first-keystroke="meta P"/>`) {
		t.Error("Missing keyboard-shortcut element")
	}

	if !strings.Contains(output, "</action>") {
		t.Error("Missing action closing tag")
	}

	if !strings.Contains(output, "</keymap>") {
		t.Error("Missing keymap closing tag")
	}
}

func TestXMLWriterApplicationConfig(t *testing.T) {
	w := &XMLWriter{}

	values := map[string]interface{}{
		"root": "application",
		"children": []interface{}{
			map[string]interface{}{
				"element": "component",
				"attrs": map[string]interface{}{
					"name": "EditorSettings",
				},
				"children": []interface{}{
					map[string]interface{}{
						"element": "option",
						"attrs": map[string]interface{}{
							"name":  "LINE_NUMBERS_SHOWN",
							"value": "true",
						},
					},
					map[string]interface{}{
						"element": "option",
						"attrs": map[string]interface{}{
							"name":  "RIGHT_MARGIN",
							"value": "120",
						},
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

	if !strings.Contains(output, "<application>") {
		t.Error("Missing application element")
	}

	if !strings.Contains(output, `<component name="EditorSettings">`) {
		t.Error("Missing component element")
	}

	if !strings.Contains(output, `<option name="LINE_NUMBERS_SHOWN" value="true"/>`) {
		t.Error("Missing LINE_NUMBERS_SHOWN option")
	}
}

func TestXMLWriterEscaping(t *testing.T) {
	w := &XMLWriter{}

	values := map[string]interface{}{
		"root": "test",
		"attrs": map[string]interface{}{
			"value": `<>&"'`,
		},
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	output := string(content)

	if !strings.Contains(output, "&lt;") {
		t.Error("< not escaped")
	}
	if !strings.Contains(output, "&gt;") {
		t.Error("> not escaped")
	}
	if !strings.Contains(output, "&amp;") {
		t.Error("& not escaped")
	}
	if !strings.Contains(output, "&quot;") {
		t.Error("\" not escaped")
	}
}

func TestXMLWriterMissingRoot(t *testing.T) {
	w := &XMLWriter{}

	values := map[string]interface{}{
		"attrs": map[string]interface{}{
			"version": "1",
		},
	}

	_, err := w.Write(values)
	if err == nil {
		t.Error("Expected error for missing root element")
	}
}

func TestXMLWriterElementWithContent(t *testing.T) {
	w := &XMLWriter{}

	values := map[string]interface{}{
		"root": "config",
		"children": []interface{}{
			map[string]interface{}{
				"element": "description",
				"content": "This is a test",
			},
		},
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	output := string(content)

	if !strings.Contains(output, "<description>This is a test</description>") {
		t.Errorf("Expected element with content, got:\n%s", output)
	}
}
