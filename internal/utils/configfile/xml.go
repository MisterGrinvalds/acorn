package configfile

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
)

// XMLWriter implements the Writer interface for XML format.
// Used for IntelliJ IDEA keymap and settings files.
type XMLWriter struct{}

func init() {
	Register(&XMLWriter{})
}

// Format returns the format identifier.
func (w *XMLWriter) Format() string {
	return "xml"
}

// Write generates XML content from values.
// Expected structure:
//
//	root: "keymap"  # Root element name
//	attrs:          # Root element attributes
//	  version: "1"
//	  name: "Custom"
//	  parent: "Mac OS X 10.5+"
//	children:       # Child elements
//	  - element: "action"
//	    attrs:
//	      id: "GotoFile"
//	    children:
//	      - element: "keyboard-shortcut"
//	        attrs:
//	          first-keystroke: "meta P"
func (w *XMLWriter) Write(values map[string]interface{}) ([]byte, error) {
	root, ok := values["root"].(string)
	if !ok {
		return nil, fmt.Errorf("xml format requires 'root' element name")
	}

	var b strings.Builder
	b.WriteString(xml.Header)

	// Build root element
	b.WriteString("<")
	b.WriteString(root)

	// Add root attributes
	if attrs, ok := values["attrs"].(map[string]interface{}); ok {
		w.writeAttrs(&b, attrs)
	}

	// Check for children
	children, hasChildren := values["children"]
	if !hasChildren {
		b.WriteString("/>\n")
		return []byte(b.String()), nil
	}

	b.WriteString(">\n")

	// Write children
	if childList, ok := children.([]interface{}); ok {
		for _, child := range childList {
			if childMap, ok := child.(map[string]interface{}); ok {
				w.writeElement(&b, childMap, 1)
			}
		}
	}

	// Close root
	b.WriteString("</")
	b.WriteString(root)
	b.WriteString(">\n")

	return []byte(b.String()), nil
}

// writeElement writes an XML element with optional attributes and children.
func (w *XMLWriter) writeElement(b *strings.Builder, elem map[string]interface{}, indent int) {
	name, ok := elem["element"].(string)
	if !ok {
		return
	}

	indentStr := strings.Repeat("  ", indent)
	b.WriteString(indentStr)
	b.WriteString("<")
	b.WriteString(name)

	// Write attributes
	if attrs, ok := elem["attrs"].(map[string]interface{}); ok {
		w.writeAttrs(b, attrs)
	}

	// Check for children or content
	children, hasChildren := elem["children"]
	content, hasContent := elem["content"].(string)

	if !hasChildren && !hasContent {
		b.WriteString("/>\n")
		return
	}

	b.WriteString(">")

	if hasContent {
		b.WriteString(xmlEscape(content))
	}

	if hasChildren {
		b.WriteString("\n")
		if childList, ok := children.([]interface{}); ok {
			for _, child := range childList {
				if childMap, ok := child.(map[string]interface{}); ok {
					w.writeElement(b, childMap, indent+1)
				}
			}
		}
		b.WriteString(indentStr)
	}

	b.WriteString("</")
	b.WriteString(name)
	b.WriteString(">\n")
}

// writeAttrs writes XML attributes in sorted order.
func (w *XMLWriter) writeAttrs(b *strings.Builder, attrs map[string]interface{}) {
	// Sort keys for deterministic output
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := attrs[k]
		b.WriteString(" ")
		b.WriteString(k)
		b.WriteString("=\"")
		b.WriteString(xmlEscape(fmt.Sprintf("%v", v)))
		b.WriteString("\"")
	}
}

// xmlEscape escapes special XML characters.
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
