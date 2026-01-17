package configfile

import (
	"strings"
)

// TextListWriter implements the Writer interface for simple text list format.
// Used for plugin lists, extension lists, etc.
type TextListWriter struct{}

func init() {
	Register(&TextListWriter{})
}

// Format returns the format identifier.
func (w *TextListWriter) Format() string {
	return "textlist"
}

// Write generates a text list from values.
// Expected structure:
//
//	header: |
//	  # Header comments
//	  # Multiple lines supported
//	sections:
//	  - name: "Section Name"
//	    items:
//	      - id: "item-id"
//	        comment: "Optional description"
//	      - id: "another-id"
//
// Or simpler flat structure:
//
//	header: "# Header"
//	items:
//	  - id: "item-id"
//	    comment: "description"
func (w *TextListWriter) Write(values map[string]interface{}) ([]byte, error) {
	var b strings.Builder

	// Write header if present
	if header, ok := values["header"].(string); ok && header != "" {
		b.WriteString(header)
		if !strings.HasSuffix(header, "\n") {
			b.WriteString("\n")
		}
	}

	// Check for sections (grouped items)
	if sections, ok := values["sections"].([]interface{}); ok {
		for i, section := range sections {
			if sectionMap, ok := section.(map[string]interface{}); ok {
				w.writeSection(&b, sectionMap, i > 0)
			}
		}
		return []byte(b.String()), nil
	}

	// Check for flat items list
	if items, ok := values["items"].([]interface{}); ok {
		for _, item := range items {
			w.writeItem(&b, item)
		}
	}

	return []byte(b.String()), nil
}

// writeSection writes a section with name and items.
func (w *TextListWriter) writeSection(b *strings.Builder, section map[string]interface{}, addBlankLine bool) {
	if addBlankLine {
		b.WriteString("\n")
	}

	// Write section name as comment
	if name, ok := section["name"].(string); ok && name != "" {
		b.WriteString("# ")
		b.WriteString(name)
		b.WriteString("\n")
	}

	// Write items
	if items, ok := section["items"].([]interface{}); ok {
		for _, item := range items {
			w.writeItem(b, item)
		}
	}
}

// writeItem writes a single item line.
func (w *TextListWriter) writeItem(b *strings.Builder, item interface{}) {
	switch v := item.(type) {
	case string:
		// Simple string item
		b.WriteString(v)
		b.WriteString("\n")
	case map[string]interface{}:
		// Item with id and optional comment
		id, _ := v["id"].(string)
		comment, _ := v["comment"].(string)

		if id == "" {
			return
		}

		b.WriteString(id)
		if comment != "" {
			// Pad to align comments (use 36 chars for id column)
			padding := 36 - len(id)
			if padding < 1 {
				padding = 1
			}
			b.WriteString(strings.Repeat(" ", padding))
			b.WriteString("# ")
			b.WriteString(comment)
		}
		b.WriteString("\n")
	}
}
