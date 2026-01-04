// Package output provides formatted output utilities for the CLI.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

// Format represents an output format type.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

// ParseFormat parses a format string into a Format type.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "table", "":
		return FormatTable, nil
	case "json":
		return FormatJSON, nil
	case "yaml", "yml":
		return FormatYAML, nil
	default:
		return "", fmt.Errorf("invalid format: %s (must be table, json, or yaml)", s)
	}
}

// Printer provides formatted output capabilities.
type Printer struct {
	writer io.Writer
	format Format
}

// NewPrinter creates a new Printer.
func NewPrinter(w io.Writer, format Format) *Printer {
	return &Printer{
		writer: w,
		format: format,
	}
}

// Print outputs data in the configured format.
func (p *Printer) Print(data interface{}) error {
	switch p.format {
	case FormatJSON:
		return p.printJSON(data)
	case FormatYAML:
		return p.printYAML(data)
	case FormatTable:
		return fmt.Errorf("table format must be implemented per command")
	default:
		return fmt.Errorf("unsupported format: %s", p.format)
	}
}

// printJSON outputs data as JSON.
func (p *Printer) printJSON(data interface{}) error {
	encoder := json.NewEncoder(p.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printYAML outputs data as YAML.
func (p *Printer) printYAML(data interface{}) error {
	encoder := yaml.NewEncoder(p.writer)
	encoder.SetIndent(2)
	return encoder.Encode(data)
}

// Table provides utilities for table output.
type Table struct {
	headers []string
	rows    [][]string
	widths  []int
}

// NewTable creates a new table.
func NewTable(headers ...string) *Table {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	return &Table{
		headers: headers,
		rows:    [][]string{},
		widths:  widths,
	}
}

// AddRow adds a row to the table.
func (t *Table) AddRow(cols ...string) {
	for i, col := range cols {
		if i < len(t.widths) && len(col) > t.widths[i] {
			t.widths[i] = len(col)
		}
	}
	t.rows = append(t.rows, cols)
}

// Render renders the table to a writer.
func (t *Table) Render(w io.Writer) {
	// Print header
	for i, header := range t.headers {
		fmt.Fprintf(w, "%-*s", t.widths[i]+2, header)
	}
	fmt.Fprintln(w)

	// Print separator
	for _, width := range t.widths {
		fmt.Fprintf(w, "%s", strings.Repeat("-", width+2))
	}
	fmt.Fprintln(w)

	// Print rows
	for _, row := range t.rows {
		for i, col := range row {
			if i < len(t.widths) {
				fmt.Fprintf(w, "%-*s", t.widths[i]+2, col)
			}
		}
		fmt.Fprintln(w)
	}
}

// ColorCode represents an ANSI color code.
type ColorCode string

const (
	ColorReset   ColorCode = "\033[0m"
	ColorRed     ColorCode = "\033[31m"
	ColorGreen   ColorCode = "\033[32m"
	ColorYellow  ColorCode = "\033[33m"
	ColorBlue    ColorCode = "\033[34m"
	ColorMagenta ColorCode = "\033[35m"
	ColorCyan    ColorCode = "\033[36m"
	ColorGray    ColorCode = "\033[90m"
)

// Colorize wraps text in ANSI color codes.
func Colorize(text string, color ColorCode) string {
	return string(color) + text + string(ColorReset)
}

// Success returns green-colored text.
func Success(text string) string {
	return Colorize(text, ColorGreen)
}

// Error returns red-colored text.
func Error(text string) string {
	return Colorize(text, ColorRed)
}

// Warning returns yellow-colored text.
func Warning(text string) string {
	return Colorize(text, ColorYellow)
}

// Info returns blue-colored text.
func Info(text string) string {
	return Colorize(text, ColorBlue)
}
